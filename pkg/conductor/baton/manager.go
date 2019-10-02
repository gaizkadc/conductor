/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package baton

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/conductor/internal/entities"
	"github.com/nalej/conductor/internal/persistence/app_cluster"
	"github.com/nalej/conductor/internal/structures"
	"github.com/nalej/conductor/pkg/conductor/observer"
	"github.com/nalej/conductor/pkg/conductor/plandesigner"
	"github.com/nalej/conductor/pkg/conductor/requirementscollector"
	"github.com/nalej/conductor/pkg/conductor/scorer"
	"github.com/nalej/conductor/pkg/utils"

	//"github.com/nalej/deployment-manager/pkg/network"
	"github.com/nalej/derrors"
	pbAppClusterApi "github.com/nalej/grpc-app-cluster-api-go"
	pbApplication "github.com/nalej/grpc-application-go"
	pbConductor "github.com/nalej/grpc-conductor-go"
	pbDeploymentManager "github.com/nalej/grpc-deployment-manager-go"
	pbNetwork "github.com/nalej/grpc-network-go"
	pbOrganization "github.com/nalej/grpc-organization-go"
	pbCoordinator "github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/nalej-bus/pkg/queue/network/ops"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

// Time to wait between checks in the queue in milliseconds.
const (
	CheckSleepTime = 2000
	// Timeout in seconds for queries to the application clusters.
	ConductorAppTimeout = 60
	// Maximum number of retries per request
	ConductorMaxDeploymentRetries = 3
	// Time to wait between retries in seconds
	ConductorSleepBetweenRetries = 25
	// Time to wait to receive a terminating status when draining clusters in seconds
	ConductorDrainClusterAppTimeout = time.Second * 60
	// Timeout when sending messages to the queue
	ConductorQueueTimeout = time.Second * 5
	// Initial address to use during the definition of VSA
	ConductorBaseVSA = "172.16.0.1"
	// Initial address to use during the definition of VSA
	ConductorOutboundVSA = "172.18.0.1"
	// Suffix to add to the outbound FQDNs
	OutboundSuffix = "-OUT"
)

type Manager struct {
	// Connections helper
	ConnHelper *utils.ConnectionsHelper
	// ScorerMethod
	ScorerMethod scorer.Scorer
	// Requirements collector
	ReqCollector requirementscollector.RequirementsCollector
	// Plan designer
	Designer plandesigner.PlanDesigner
	// queue for incoming requests
	Queue structures.RequestsQueue
	// Pending plans
	PendingPlans *structures.PendingPlans
	// Application client
	AppClient pbApplication.ApplicationsClient
	// Networking manager client
	NetClient pbNetwork.NetworksClient
	// DNS manager client
	DNSClient pbNetwork.DNSClient
	// UnifiedLogging client
	UnifiedLoggingClient pbCoordinator.CoordinatorClient
	// Information about applications deployed in clusters
	AppClusterDB *app_cluster.AppClusterDB
	// NetworkOps queue producer
	NetworkOpsProducer *ops.NetworkOpsProducer
}

func NewManager(connHelper *utils.ConnectionsHelper, queue structures.RequestsQueue, scorer scorer.Scorer,
	reqColl requirementscollector.RequirementsCollector, designer plandesigner.PlanDesigner,
	pendingPlans *structures.PendingPlans, appClusterDB *app_cluster.AppClusterDB, networkOpsProducer *ops.NetworkOpsProducer) *Manager {
	// initialize clients
	pool := connHelper.GetSystemModelClients()
	if pool != nil && len(pool.GetConnections()) == 0 {
		log.Panic().Msg("system model clients were not started")
		return nil
	}
	conn := pool.GetConnections()[0]
	// Create associated clients
	appClient := pbApplication.NewApplicationsClient(conn)
	// Network client
	netPool := connHelper.GetNetworkingClients()
	if netPool != nil && len(netPool.GetConnections()) == 0 {
		log.Panic().Msg("networking client was not started")
		return nil
	}
	// UnifiedLogging client
	ulPool := connHelper.GetUnifiedLoggingClients()
	if ulPool != nil && len(ulPool.GetConnections()) == 0 {
		log.Panic().Msg("unified logging client was not started")
		return nil
	}

	netClient := pbNetwork.NewNetworksClient(netPool.GetConnections()[0])
	dnsClient := pbNetwork.NewDNSClient(netPool.GetConnections()[0])
	ulClient := pbCoordinator.NewCoordinatorClient(ulPool.GetConnections()[0])
	return &Manager{ConnHelper: connHelper, Queue: queue, ScorerMethod: scorer, ReqCollector: reqColl,
		Designer: designer, AppClient: appClient, PendingPlans: pendingPlans, NetClient: netClient,
		DNSClient: dnsClient, UnifiedLoggingClient: ulClient, AppClusterDB: appClusterDB,
		NetworkOpsProducer: networkOpsProducer}
}

// Check iteratively if there is anything to be processed in the queue.
func (c *Manager) Run() {
	sleep := time.Tick(time.Millisecond * CheckSleepTime)
	for {
		select {
		case <-sleep:
			//TODO revisit this solution because it could lead to intensive active queue checking
			forNextIteration := make([]*entities.DeploymentRequest, 0)
			for c.Queue.AvailableRequests() {
				log.Info().Int("queued requests", c.Queue.Len()).Msg("there are pending deployment requests")
				next := c.Queue.NextRequest()
				readyToProcess := true
				// Check if there was any error with this deployment
				if next.NumRetries > 0 {
					// this app had a retry, check if enough time passed since the last check
					elapsedTime := time.Now().Unix() - next.TimeRetry.Unix()
					if elapsedTime < ConductorSleepBetweenRetries {
						log.Debug().Str("requestId", next.RequestId).Msg("not enough time elapsed before retry")
						readyToProcess = false
					}
				}

				if readyToProcess {
					c.processQueuedRequest(next)
				} else {
					// queue it for later
					forNextIteration = append(forNextIteration, next)
				}
			}
			// Add again to the queue the non-processed entries
			if len(forNextIteration) > 0 {
				log.Info().Int("pending", len(forNextIteration)).Msg("some deployments were excluded in this round")
			}
			for _, toAdd := range forNextIteration {
				c.Queue.PushRequest(toAdd)
			}
		}
	}
}

// Process a queued deployment request.
func (c *Manager) processQueuedRequest(req *entities.DeploymentRequest) {
	err := c.ProcessDeploymentRequest(req)
	if err != nil {

		// Update this deployment request
		req.NumRetries = req.NumRetries + 1

		// Prepare connections to update status information
		smConn := c.ConnHelper.SMClients.GetConnections()[0]
		client := pbApplication.NewApplicationsClient(smConn)

		var updateRequest pbApplication.UpdateAppStatusRequest

		if req.NumRetries >= ConductorMaxDeploymentRetries {
			log.Error().Str("requestId", req.RequestId).Msg("exceeded number of retries")
			// Consider this deployment to be failed
			// Update instance value to ERROR
			updateRequest = pbApplication.UpdateAppStatusRequest{
				AppInstanceId:  req.InstanceId,
				OrganizationId: req.OrganizationId,
				Status:         pbApplication.ApplicationStatus_DEPLOYMENT_ERROR,
				Info:           fmt.Sprintf("Exceeded number of retries. Latest known error: [%v]", err.Error()),
			}
		} else {
			log.Error().Err(err).Str("requestId", req.RequestId).Msg("enqueue again after errors")
			currentTime := time.Now()
			req.TimeRetry = &currentTime
			c.Queue.PushRequest(req)

			updateRequest = pbApplication.UpdateAppStatusRequest{
				AppInstanceId:  req.InstanceId,
				OrganizationId: req.OrganizationId,
				Status:         pbApplication.ApplicationStatus_DEPLOYMENT_ERROR,
				Info:           err.Error(),
			}

		}

		_, errUpdate := client.UpdateAppStatus(context.Background(), &updateRequest)
		if errUpdate != nil {
			log.Error().Interface("request", updateRequest).Msg("error updating application instance status")
		}
		// call the Rollback
		c.Rollback(req.OrganizationId, req.InstanceId, make([]string, 0))
	}
}

// Push a request into the queue.
func (c *Manager) PushRequest(req *pbConductor.DeploymentRequest) error {
	log.Debug().Interface("request", req).Msg("received deployment request")
	// Get ParameterizedDescriptor
	desc, err := c.AppClient.GetParametrizedDescriptor(context.Background(), req.AppInstanceId)
	if err != nil {
		log.Error().Err(err).Msg("error getting application descriptor")
		return err
	}

	toEnqueue := entities.DeploymentRequest{
		RequestId:      req.RequestId,
		InstanceId:     req.AppInstanceId.AppInstanceId,
		OrganizationId: req.AppInstanceId.OrganizationId,
		ApplicationId:  desc.AppDescriptorId,
		NumRetries:     0,
		TimeRetry:      nil,
		AppInstanceId:  req.AppInstanceId.AppInstanceId,
	}
	err = c.Queue.PushRequest(&toEnqueue)
	if err != nil {
		return err
	}

	if err != nil {
		log.Error().Err(err).Msgf("problems updating application %s", req.AppInstanceId.AppInstanceId)
		return err
	}

	return nil
}

func (c *Manager) ProcessDeploymentRequest(req *entities.DeploymentRequest) derrors.Error {
	if req == nil {
		err := derrors.NewFailedPreconditionError("the queue was unexpectedly empty")
		log.Error().Err(err)
		return err
	}

	retrievedAppInstance, err := c.AppClient.GetAppInstance(context.Background(),
		&pbApplication.AppInstanceId{OrganizationId: req.OrganizationId, AppInstanceId: req.InstanceId})
	if err != nil {
		err := derrors.NewGenericError("impossible to obtain application descriptor")
		log.Error().Err(err).Msg("impossible to retrieve app instance")
		return err
	}

	appInstance := entities.NewAppInstanceFromGRPC(retrievedAppInstance)

	// 1) collect requirements for the application descriptor
	// Get the application descriptor (parametrized)
	appDescriptor, err := c.AppClient.GetParametrizedDescriptor(context.Background(),
		&pbApplication.AppInstanceId{AppInstanceId: appInstance.AppInstanceId, OrganizationId: appInstance.OrganizationId})
	if err != nil {
		err := derrors.NewNotFoundError("impossible to find application descriptor", err)
		log.Error().Err(err).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId).
			Msg("application descriptor not found when processing deployment request")
		return err
	}

	foundRequirements, err := c.ReqCollector.FindRequirements(appDescriptor, appInstance.AppInstanceId)
	if err != nil {
		err := derrors.NewGenericError("impossible to find requirements for application")
		log.Error().Err(err).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
		return err
	}

	// 2) score requirements
	scoreResult, err := c.ScorerMethod.ScoreRequirements(req.OrganizationId, foundRequirements)

	if err != nil {
		err := derrors.NewGenericError("error scoring request")
		log.Error().Err(err).Str("requestId", req.RequestId).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
		return err
	}

	log.Info().Msgf("conductor maximum score for %s has score %v from %d potential candidates",
		req.RequestId, scoreResult.DeploymentsScore, scoreResult.NumEvaluatedClusters)

	// 3) design plan
	// Elaborate deployment plan for the application
	plan, err := c.Designer.DesignPlan(appInstance, *scoreResult, *req, nil, nil)

	if err != nil {
		log.Error().Err(err).Str("requestId", req.RequestId).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
		return derrors.AsError(err, fmt.Sprintf("plan design failed for descriptor %s", err.Error()))
	}

	// 4) Create the virtual service addresses
	vsa, err := c.createVSA(entities.NewParametrizedDescriptorFromGRPC(appDescriptor), appInstance.AppInstanceId)
	if err != nil {
		err := derrors.NewGenericError("impossible to create VAS", err)
		log.Error().Err(err).Str("appDescriptorId", appDescriptor.AppDescriptorId).Msg("impossible to create VAS")
		return err
	}

	// 5) Create ZT-network with Network manager
	// we use the app instance id as the network id
	ztNetworkId, zt_err := c.CreateZTNetwork(retrievedAppInstance.AppInstanceId, req.OrganizationId,
		retrievedAppInstance.AppInstanceId, vsa)
	if zt_err != nil {
		err := derrors.NewGenericError("impossible to create zt network before deployment", zt_err)
		log.Error().Err(zt_err).Str("requestId", req.RequestId).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
		return err
	}

	// 6) deploy fragments
	// Tell deployment managers to execute plans
	err_deploy := c.DeployPlan(plan, ztNetworkId, req.NumRetries)
	if err_deploy != nil {
		err := derrors.NewGenericError("error deploying plan request", err_deploy)
		log.Error().Err(err_deploy).Str("requestId", req.RequestId).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
		return err
	}
	return nil
}

// Reschedule all the current app instances if needed.
func (c *Manager) scheduleRunningApps(organizationId string) {
	log.Debug().Str("organizationId", organizationId).Msg("schedule running apps for potential replanning")
	// get the list of available application instances
	ctx, cancel := context.WithTimeout(context.Background(), ConductorQueueTimeout)
	instances, err := c.AppClient.ListAppInstances(ctx, &pbOrganization.OrganizationId{OrganizationId: organizationId})
	defer cancel()
	if err != nil {
		log.Error().Err(err).Msg("impossible to retrieve application instances")
		return
	}

	c.ConnHelper.UpdateClusterConnections(organizationId)

	for i, appInstance := range instances.Instances {
		log.Debug().Msgf("Check if instance %d out of %d instances has to be scheduled", i+1, len(instances.Instances))

		if err != nil {
			log.Error().Err(err).Msg("impossible to get app descriptor")
		}

		// summary of service groups with multiple replica set
		replicated := make(map[string]*pbApplication.ServiceGroupInstance, 0)
		numReplicas := make(map[string]int, 0)
		for _, group := range appInstance.Groups {
			if group.Specs != nil && group.Specs.MultiClusterReplica {
				replicated[group.ServiceGroupId] = group
				if _, found := numReplicas[group.ServiceGroupId]; !found {
					numReplicas[group.ServiceGroupId] = 1
				} else {
					numReplicas[group.ServiceGroupId]++
				}
			}
		}

		log.Debug().Msgf("found %d multi replica cluster groups", len(replicated))

		// check if these groups are instantiated as many times as it can be
		// store the list of service groups for scheduling
		toReschedule := make([]string, 0)
		for _, group := range replicated {
			expectedReplicas := 0
			log.Debug().Str("groupId", group.ServiceGroupId).Msg("check group id")
			for _, cluster := range c.ConnHelper.ClusterReference {
				if c.clusterCanDeployGroup(cluster.Labels, group.Specs.DeploymentSelectors) {
					expectedReplicas++
				}
			}
			// check if expected and current do not match
			if expectedReplicas > numReplicas[group.ServiceGroupId] {
				toReschedule = append(toReschedule, group.ServiceGroupId)
			}
		}

		// Reschedule....
		if len(toReschedule) > 0 {
			conv := entities.NewAppInstanceFromGRPC(appInstance)
			c.scheduleServiceGroups(toReschedule, conv)
		} else {
			log.Debug().Str("appInstanceId", appInstance.AppInstanceId).Msg("no groups to reschedule for this app")
		}
	}
}

// Schedule the deployment of a set of service groups.
// params:
//  serviceGroupIds list of service group ids to be scheduled
//  appInstanceId
//  appDescriptorId
//  organizationId
func (c *Manager) scheduleServiceGroups(serviceGroupIds []string, appInstance entities.AppInstance) {
	log.Debug().Msgf("schedule %d services from app %s", len(serviceGroupIds), appInstance.AppInstanceId)

	// get application descriptor
	ctx, cancel := context.WithTimeout(context.Background(), ConductorQueueTimeout)
	defer cancel()
	appDescriptor, err := c.AppClient.GetParametrizedDescriptor(ctx,
		&pbApplication.AppInstanceId{OrganizationId: appInstance.OrganizationId, AppInstanceId: appInstance.AppInstanceId})
	if err != nil {
		log.Error().Err(err).Msg("scheduleServiceGroups aborted due to descriptor impossible to retrieve")
		return
	}

	foundRequirements, err := c.ReqCollector.FindRequirementsForGroups(serviceGroupIds, appInstance.AppInstanceId, appDescriptor)
	if err != nil {
		err := derrors.NewGenericError("impossible to find requirements for application")
		log.Error().Err(err).Str("appDescriptorId", appInstance.AppDescriptorId)
		return
	}

	// 2) score requirements
	scoreResult, err := c.ScorerMethod.ScoreRequirements(appInstance.OrganizationId, foundRequirements)

	if err != nil {
		err := derrors.NewGenericError("error scoring request")
		log.Error().Err(err).Str("appDescriptorId", appInstance.AppDescriptorId)
		return
	}

	log.Info().Msgf("conductor maximum score has score %v from %d potential candidates",
		scoreResult.DeploymentsScore, scoreResult.NumEvaluatedClusters)

	// 3) design plan
	// Elaborate deployment plan for the application
	req := entities.DeploymentRequest{
		OrganizationId: appInstance.OrganizationId,
		AppInstanceId:  appInstance.AppInstanceId,
		ApplicationId:  appInstance.AppDescriptorId,
		RequestId:      uuid.New().String(),
		InstanceId:     uuid.New().String(),
	}

	// design a plan for the service groups contained into the deployment fragment
	// build a summary of the groups running in the cluster
	allocatedGroupsPerClusters := make(map[string][]string, 0)
	for clusterId := range c.ConnHelper.ClusterReference {
		// get the list of the deployment fragments for the same service group of the
		// target fragment in the cluster
		clusterFragments, err := c.AppClusterDB.GetFragmentsApp(clusterId, appInstance.AppInstanceId)
		if err != nil {
			log.Error().Msg("error when getting deployment fragment from cluster")
			continue
		}
		if clusterFragments != nil {
			entry, found := allocatedGroupsPerClusters[clusterId]
			if !found {
				allocatedGroupsPerClusters[clusterId] = make([]string, 0)
				entry = allocatedGroupsPerClusters[clusterId]
			}
			// TODO this check assumes all stages in a deployment fragment to be in the same service group
			for _, cf := range clusterFragments {
				entry = append(entry, cf.Stages[0].Services[0].ServiceGroupId)
			}
			allocatedGroupsPerClusters[clusterId] = entry
		}
	}

	log.Debug().Interface("allocatedGroupsPerCluster", allocatedGroupsPerClusters).
		Interface("serviceGroupIds", serviceGroupIds).
		Msg("design a plan for deployment fragments")
	plan, err := c.Designer.DesignPlan(appInstance, *scoreResult, req, serviceGroupIds, allocatedGroupsPerClusters)

	if err != nil {
		log.Error().Err(err).Str("appDescriptorId", appInstance.AppDescriptorId)
		return
	}

	// 4) deploy fragments
	// Tell deployment managers to execute plans

	// We need the ZT network to deploy the plan.
	ctxNet, cancelNet := context.WithTimeout(context.Background(), ConductorQueueTimeout)
	defer cancelNet()
	network, err := c.AppClient.GetAppZtNetwork(ctxNet,
		&pbApplication.GetAppZtNetworkRequest{
			AppInstanceId: appInstance.AppInstanceId, OrganizationId: appInstance.OrganizationId})

	if err != nil {
		log.Error().Err(err).Msg("service groups could not be deployed. The network id was not found")
		return
	}

	err_deploy := c.DeployPlan(plan, network.NetworkId, 0)
	if err_deploy != nil {
		//err := derrors.NewGenericError("error deploying plan request", err_deploy)
		log.Error().Err(err_deploy).Str("requestId", req.RequestId).Str("appDescriptorId", appInstance.AppDescriptorId)
		return
	}

}

// Check if a cluster can deploy a service group
func (c *Manager) clusterCanDeployGroup(cluster map[string]string, group map[string]string) bool {
	if group == nil {
		return true
	}

	for expectedKey, expectedValue := range group {
		foundValue, foundKey := cluster[expectedKey]
		if !foundKey {
			return false
		}
		if foundValue != expectedValue {
			return true
		}
	}
	// everything was correct
	return true
}

// Analyze the best deployment options for a single deployment fragment.
func (c *Manager) ProcessDeploymentFragment(fragment *entities.DeploymentFragment) derrors.Error {
	if fragment == nil {
		err := derrors.NewFailedPreconditionError("the requested fragment was nil")
		log.Error().Err(err)
		return err
	}

	// Get the ServiceGroup structure

	retrievedAppInstance, err := c.AppClient.GetAppInstance(context.Background(),
		&pbApplication.AppInstanceId{OrganizationId: fragment.OrganizationId, AppInstanceId: fragment.AppInstanceId})
	if err != nil {
		err := derrors.NewGenericError("impossible to obtain application descriptor")
		log.Error().Err(err).Msg("impossible to retrieve app instance")
		return err
	}

	appInstance := entities.NewAppInstanceFromGRPC(retrievedAppInstance)

	// 1) collect requirements for the application descriptor
	// Get the application descriptor (parametrized)
	appDescriptor, err := c.AppClient.GetParametrizedDescriptor(context.Background(),
		&pbApplication.AppInstanceId{AppInstanceId: fragment.AppInstanceId, OrganizationId: fragment.OrganizationId})
	if err != nil {
		err := derrors.NewNotFoundError("impossible to find application descriptor", err)
		log.Error().Err(err).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId).
			Msg("application descriptor not found when processing deployment request")
		return err
	}

	log.Debug().Interface("re-scheduled fragment", fragment).Msg("this is the fragment to deploy")

	// Find the service groups in the deployment fragment
	serviceGroupsMap := make(map[string]bool, 0)
	serviceGroupIds := make([]string, 0)
	for _, sg := range fragment.Stages {
		for _, s := range sg.Services {
			if _, found := serviceGroupsMap[s.ServiceGroupId]; !found {
				log.Debug().Str("groupName", s.Name).Str("serviceGroupId", s.ServiceGroupId).Msg("services contained in running service group")
				serviceGroupsMap[s.ServiceGroupId] = true
				serviceGroupIds = append(serviceGroupIds, s.ServiceGroupId)
			} else {
				log.Debug().Str("groupName", s.Name).Str("serviceGroupId", s.ServiceGroupId).Msg("excluded from the list")
			}
		}
	}
	log.Debug().Interface("serviceGroupIds", serviceGroupIds).
		Msgf("a total of %d service groups have to be rescheduled", len(serviceGroupIds))

	foundRequirements, err := c.ReqCollector.FindRequirementsForGroups(serviceGroupIds, appInstance.AppInstanceId, appDescriptor)
	if err != nil {
		err := derrors.NewGenericError("impossible to find requirements for application")
		log.Error().Err(err).Str("appDescriptorId", fragment.AppDescriptorId)
		return err
	}

	// 2) score requirements
	scoreResult, err := c.ScorerMethod.ScoreRequirements(fragment.OrganizationId, foundRequirements)

	if err != nil {
		err := derrors.NewGenericError("error scoring request")
		log.Error().Err(err).Str("appDescriptorId", fragment.AppDescriptorId)
		return err
	}

	log.Info().Msgf("conductor maximum score has score %v from %d potential candidates",
		scoreResult.DeploymentsScore, scoreResult.NumEvaluatedClusters)

	// 3) design plan
	// Elaborate deployment plan for the application
	req := entities.DeploymentRequest{
		OrganizationId: fragment.OrganizationId,
		AppInstanceId:  fragment.AppInstanceId,
		ApplicationId:  fragment.AppDescriptorId,
		RequestId:      uuid.New().String(),
		InstanceId:     fragment.FragmentId,
	}

	// design a plan for the service groups contained into the deployment fragment
	// build a summary of the groups running in the cluster
	allocatedGroupsPerClusters := make(map[string][]string, 0)
	for clusterId := range c.ConnHelper.ClusterReference {
		// get the list of the deployment fragments for the same service group of the
		// target fragment in the cluster
		clusterFragments, err := c.AppClusterDB.GetFragmentsApp(clusterId, fragment.AppInstanceId)
		if err != nil {
			log.Error().Msg("error when getting deployment fragment from cluster")
			continue
		}
		if clusterFragments != nil {
			entry, found := allocatedGroupsPerClusters[clusterId]
			if !found {
				allocatedGroupsPerClusters[clusterId] = make([]string, 0)
				entry = allocatedGroupsPerClusters[clusterId]
			}
			// TODO this check assumes all stages in a deployment fragment to be in the same service group
			for _, cf := range clusterFragments {
				entry = append(entry, cf.Stages[0].Services[0].ServiceGroupId)
			}
			allocatedGroupsPerClusters[clusterId] = entry
		}
	}

	log.Debug().Interface("allocatedGroupsPerCluster", allocatedGroupsPerClusters).
		Interface("serviceGroupIds", serviceGroupIds).
		Msg("design a plan for deployment fragments")
	plan, err := c.Designer.DesignPlan(appInstance, *scoreResult, req, serviceGroupIds, allocatedGroupsPerClusters)

	if err != nil {
		log.Error().Err(err).Str("appDescriptorId", fragment.AppDescriptorId)
		return derrors.AsError(err, fmt.Sprintf("plan design failed for descriptor %s", err.Error()))
	}

	// 4) deploy fragments
	// Tell deployment managers to execute plans
	err_deploy := c.DeployPlan(plan, fragment.ZtNetworkID, 0)
	if err_deploy != nil {
		err := derrors.NewGenericError("error deploying plan request", err_deploy)
		log.Error().Err(err_deploy).Str("requestId", req.RequestId).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
		return err
	}
	return nil
}

// Create a new zero tier network and return the corresponding network id.
// params:
//  name of the network
//  organizationId for this network
// returns:
//  networkId or error otherwise
func (c *Manager) CreateZTNetwork(name string, organizationId string, appInstanceId string, vsa map[string]string) (string, error) {
	request := pbNetwork.AddNetworkRequest{
		Name:           name,
		OrganizationId: organizationId,
		AppInstanceId:  appInstanceId,
		Vsa:            vsa}

	log.Debug().Interface("addNetworkRequest", request).Msgf("create a network request")

	ztNetworkId, err := c.NetClient.AddNetwork(context.Background(), &request)

	if err != nil {
		log.Error().Err(err).Msgf("there was a problem when creating network for name: %s with org: %s", name, organizationId)
		return "", err
	}
	return ztNetworkId.NetworkId, err
}

// For a given collection of plans, tell the corresponding deployment managers to run the deployment.
// params:
//  plan to be deployed
//  ztNetworkId identifier for the zt network to be created
//  numRetry number of retry of this plan
// returns:
//  error if any
func (c *Manager) DeployPlan(plan *entities.DeploymentPlan, ztNetworkId string, numRetry int32) error {
	// Add this plan to the list of pending entries
	c.PendingPlans.AddPendingPlan(plan)

	// retrieve the application instance and update the cluster ids
	appInstance, err := c.AppClient.GetAppInstance(context.Background(),
		&pbApplication.AppInstanceId{OrganizationId: plan.OrganizationId, AppInstanceId: plan.AppInstanceId})
	if err != nil {
		return derrors.NewInternalError("impossible to find application instance before deployment", err)
	}

	// before sending the deployment plan we check that all the involved clusters are ready
	for _, fragment := range plan.Fragments {
		targetCluster, found := c.ConnHelper.ClusterReference[fragment.ClusterId]
		if !found {
			msg := fmt.Sprintf("unknown target address for cluster with id %s", fragment.ClusterId)
			err := errors.New(msg)
			log.Error().Msgf(msg)
			return err
		}
		if targetCluster.Cordon {
			msg := fmt.Sprintf("the cluster %s with address %s is cordoned", fragment.ClusterId, targetCluster.Hostname)
			err := errors.New(msg)
			log.Error().Str("clusterId", fragment.ClusterId).Msg(msg)
			return err
		}
	}

	// time to deploy
	for fragmentIndex, fragment := range plan.Fragments {
		log.Debug().Interface("fragment", fragment).Msg("fragment to be deployed")
		log.Info().Str("deploymentId", fragment.DeploymentId).
			Msgf("start fragment %s deployment with %d out of %d fragments", fragment.DeploymentId, fragmentIndex+1, len(plan.Fragments))

		targetCluster, found := c.ConnHelper.ClusterReference[fragment.ClusterId]
		if !found {
			msg := fmt.Sprintf("unknown target address for cluster with id %s", fragment.ClusterId)
			err := errors.New(msg)
			log.Error().Msgf(msg)
			return err
		}

		clusterAddress := fmt.Sprintf("%s:%d", targetCluster.Hostname, utils.APP_CLUSTER_API_PORT)
		log.Debug().Str("clusterAddress", clusterAddress).Msg("Deploying plan")
		conn, err := c.ConnHelper.GetClusterClients().GetConnection(clusterAddress)

		if err != nil {
			log.Error().Err(err).Msgf("problem creating connection with %s", clusterAddress)
			return err
		}

		// build a request
		request := pbDeploymentManager.DeploymentFragmentRequest{
			RequestId:      uuid.New().String(),
			Fragment:       fragment.ToGRPC(),
			ZtNetworkId:    ztNetworkId,
			RollbackPolicy: pbDeploymentManager.RollbackPolicy_NONE,
			NumRetry:       numRetry,
		}

		client := pbAppClusterApi.NewDeploymentManagerClient(conn)

		log.Debug().Interface("deploymentFragmentRequest", request).
			Msg("deployment fragment request")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*ConductorAppTimeout)
		defer cancel()
		response, err := client.Execute(ctx, &request)

		log.Debug().Interface("deploymentFragmentResponse", response).Interface("deploymentFragmentError", err).
			Msg("finished fragment deployment")

		if err != nil {
			// TODO define how to proceed in case of error
			log.Error().Err(err).Str("deploymentId", fragment.DeploymentId).Msg("problem deploying fragment")
			return err
		}

		// update the db of fragments deployed on that cluster
		// update the value of the ztNetworkId in the local entity
		fragment.ZtNetworkID = ztNetworkId
		err = c.AppClusterDB.AddDeploymentFragment(&fragment)
		if err != nil {
			log.Error().Err(err).Msg("there was a problem when storing information about a deployment fragment")
		}

		// update the corresponding cluster id on the resources
		for _, stage := range fragment.Stages {
			for _, serv := range stage.Services {
				// update the cluster for this service
				for _, groups := range appInstance.Groups {
					for _, instServ := range groups.ServiceInstances {
						if instServ.ServiceInstanceId == serv.ServiceInstanceId {
							instServ.DeployedOnClusterId = fragment.ClusterId
						}
					}
				}
			}
		}
	}

	// update the instance
	log.Debug().Str("appInstanceId", appInstance.AppInstanceId).Msg("update app instance after deployment")
	_, err = c.AppClient.UpdateAppInstance(context.Background(), appInstance)
	if err != nil {
		log.Error().Err(err).Msg("impossible to update application instance after deployment was done")
		return err
	}

	return nil
}

// Undeploy
func (c *Manager) Undeploy(request *entities.UndeployRequest) error {
	/*
	   err := c.PendingPlans.RemovePendingPlanByApp(request.AppInstanceId)
	   if err != nil {
	       log.Error().Err(err).Msg("impossible to remove a pending pending plan")
	   }
	*/
	return c.HardUndeploy(request.OrganizationId, request.AppInstanceId)
}

// Undeploy function that maintains the application instance in the system.
func (c *Manager) SoftUndeploy(organizationId string, appInstanceId string) error {
	// find application instance
	appInstance, err := c.AppClient.GetAppInstance(context.Background(),
		&pbApplication.AppInstanceId{OrganizationId: organizationId, AppInstanceId: appInstanceId})
	if err != nil {
		log.Error().Err(err).Str("appInstanceID", appInstanceId).Msgf("impossible to obtain application descriptor")
		return err
	}

	clusterFound := make(map[string]bool, 0)
	clusterIds := make([]string, 0)
	for _, g := range appInstance.Groups {
		for _, svc := range g.ServiceInstances {
			if svc.DeployedOnClusterId != "" {
				if _, found := clusterFound[svc.DeployedOnClusterId]; !found {
					clusterFound[svc.DeployedOnClusterId] = true
					clusterIds = append(clusterIds, svc.DeployedOnClusterId)
				}
			}
		}
	}

	// call Rollback
	c.Rollback(organizationId, appInstanceId, clusterIds)

	// create observer and the array of entries to be observed
	toObserve := make([]observer.ObservableDeploymentFragment, 0)
	for _, clusterId := range clusterIds {
		fragments, err := c.AppClusterDB.GetFragmentsApp(clusterId, appInstanceId)
		if err != nil {
			log.Error().Err(err).Msg("error when getting fragments from cluster")
			continue
		}
		for _, fr := range fragments {
			toObserve = append(toObserve, observer.ObservableDeploymentFragment{ClusterId: clusterId,
				FragmentId: fr.FragmentId, AppInstanceId: appInstanceId})
		}
	}

	observer := observer.NewDeploymentFragmentsObserver(toObserve, c.AppClusterDB)
	// Run an observer in a separated thread to send the schedule to the queue when is terminating
	go observer.Observe(ConductorDrainClusterAppTimeout, entities.FRAGMENT_TERMINATING,
		func(d *entities.DeploymentFragment) derrors.Error {
			return c.AppClusterDB.DeleteDeploymentFragment(d.ClusterId, d.FragmentId)
		}, nil)

	// terminate execution
	return c.undeployClustersInstance(appInstance.OrganizationId, appInstance.AppInstanceId, clusterIds)
}

// Undeploy function that removes the application instance
func (c *Manager) HardUndeploy(organizationId string, appInstanceId string) error {
	// find application instance
	appInstance, err := c.AppClient.GetAppInstance(context.Background(),
		&pbApplication.AppInstanceId{OrganizationId: organizationId, AppInstanceId: appInstanceId})
	if err != nil {
		log.Error().Err(err).Str("appInstanceID", appInstanceId).Msgf("impossible to obtain application descriptor")
		return err
	}

	// find in what clusters this app is running
	clusterFound := make(map[string]bool, 0)
	clusterIds := make([]string, 0)
	// set of services
	services := make(map[string]bool, 1)
	for _, g := range appInstance.Groups {
		for _, svc := range g.ServiceInstances {
			// add this service to the set
			services[svc.Name] = true
			if svc.DeployedOnClusterId != "" {
				if _, found := clusterFound[svc.DeployedOnClusterId]; !found {
					clusterFound[svc.DeployedOnClusterId] = true
					clusterIds = append(clusterIds, svc.DeployedOnClusterId)
				}
			}
		}
	}

	instID := &pbApplication.AppInstanceId{
		OrganizationId: organizationId,
		AppInstanceId:  appInstanceId,
	}

	// TODO: see if this method could be moved to application manager ( app-manager is responsible for creating it)
	_, err = c.AppClient.RemoveParametrizedDescriptor(context.Background(), instID)
	if err != nil {
		log.Error().Err(err).Str("app_instance_id", appInstanceId).Msg("could not remove parametrized descriptor from system model")
	}

	// Remove from the associated request from the queue
	removed := c.Queue.Remove(appInstanceId)
	if !removed {
		log.Info().Interface("appInstanceId", appInstanceId).Msg("no request was found in the queue for this deployed app")
	}

	// Remove any DNS entry
	for serviceName := range services {
		req := pbNetwork.DeleteDNSEntryRequest{
			OrganizationId: organizationId,
			ServiceName:    utils.GetVSAName(serviceName, organizationId, appInstanceId),
		}
		ctx, cancel := context.WithTimeout(context.Background(), ConductorQueueTimeout)
		errQueue := c.NetworkOpsProducer.Send(ctx, &req)
		cancel()
		if errQueue != nil {
			log.Error().Err(errQueue).Interface("request", req).Msg("faile sending delete dns entry request to the queue")
		}
	}

	// create observer and the array of entries to be observed
	toObserve := make([]observer.ObservableDeploymentFragment, 0)
	for _, clusterId := range clusterIds {
		fragments, err := c.AppClusterDB.GetFragmentsApp(clusterId, appInstanceId)
		if err != nil {
			log.Error().Err(err).Msg("error when getting fragments from cluster")
			continue
		}
		for _, fr := range fragments {
			toObserve = append(toObserve, observer.ObservableDeploymentFragment{ClusterId: clusterId,
				FragmentId: fr.FragmentId, AppInstanceId: instID.AppInstanceId})
		}
	}

	observer := observer.NewDeploymentFragmentsObserver(toObserve, c.AppClusterDB)
	// Run an observer in a separated thread to send the schedule to the queue when is terminating
	go observer.Observe(
		ConductorDrainClusterAppTimeout,
		entities.FRAGMENT_TERMINATING,
		// function to execute in every deployment fragment
		func(d *entities.DeploymentFragment) derrors.Error {
			log.Debug().Str("deploymentFragmentId", d.DeploymentId).Msg("deployment fragment terminated has" +
				" to be removed")
			err := c.AppClusterDB.DeleteDeploymentFragment(d.ClusterId, d.FragmentId)
			if err != nil {
				return err
			}
			if !c.PendingPlans.PlanHasPendingFragments(d.DeploymentId) {
				log.Debug().Msg("pending plan has no fragments, remove it")
				c.PendingPlans.RemovePendingPlan(d.DeploymentId)
			}
			return nil
		},
		// finally remove the application instance
		func() {

			// call Rollback
			c.Rollback(organizationId, appInstanceId, clusterIds)

			// Remove any entry from the system model
			instID := &pbApplication.AppInstanceId{
				OrganizationId: organizationId,
				AppInstanceId:  appInstanceId,
			}
			ctx, cancel := context.WithTimeout(context.Background(), ConductorQueueTimeout)
			defer cancel()
			log.Debug().Interface("appInstanceId", instID).Msg("remove application instance")
			_, err = c.AppClient.RemoveAppInstance(ctx, instID)
			if err != nil {
				log.Error().Err(err).Str("app_instance_id", appInstanceId).Msg("could not remove instance from system model")
			}
		})

	// terminate execution
	return c.undeployClustersInstance(appInstance.OrganizationId, appInstance.AppInstanceId, clusterIds)
}

// Private function to communicate the application clusters to remove a running instance.
// params:
//  organizationId
//  appInstanceId
//  targetClusters list with the cluster ids where this function must run. If no list is set, then all clusters are informed.
// return:
//  error if any
func (c *Manager) undeployClustersInstance(organizationId string, appInstanceId string, targetClusters []string) error {
	// Send undeploy requests to application clusters
	log.Debug().Str("app_instance_id", appInstanceId).Msg("undeploy app instance with id")

	err := c.ConnHelper.UpdateClusterConnections(organizationId)
	if err != nil {
		log.Error().Err(err).Str("organizationID", organizationId).Msg("error updating connections for organization")
		return err
	}
	if len(c.ConnHelper.ClusterReference) == 0 {
		log.Error().Msgf("no clusters found for organization %s", organizationId)
		return nil
	}
	log.Debug().Interface("number", len(c.ConnHelper.ClusterReference)).Msg("Known clusters")

	log.Debug().Int("number of cluster to send undeploy", len(targetClusters)).Msg("send undeploy to clusters")
	if len(targetClusters) == 0 {
		log.Error().Msg("no clusters found to send undeploy notification")
	}

	for _, clusterId := range targetClusters {

		clusterEntry, found := c.ConnHelper.ClusterReference[clusterId]
		if !found {
			log.Error().Str("clusterId", clusterId).Str("clusterHost", clusterEntry.Hostname).Msg("unknown clusterHost for the clusterId")
			return errors.New(fmt.Sprintf("unknown host for cluster id %s", clusterId))
		}

		log.Debug().Str("clusterId", clusterId).Str("clusterHost", clusterEntry.Hostname).Msg("conductor query deployment-manager cluster")

		clusterAddress := fmt.Sprintf("%s:%d", clusterEntry.Hostname, utils.APP_CLUSTER_API_PORT)
		conn, err := c.ConnHelper.GetClusterClients().GetConnection(clusterAddress)
		if err != nil {
			log.Error().Err(err).Str("clusterHost", clusterEntry.Hostname).Msg("impossible to get connection for the host")
			return err
		}

		dmClient := pbAppClusterApi.NewDeploymentManagerClient(conn)

		undeployRequest := pbDeploymentManager.UndeployRequest{
			OrganizationId: organizationId,
			AppInstanceId:  appInstanceId,
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*ConductorAppTimeout)
		_, err = dmClient.Undeploy(ctx, &undeployRequest)

		if err != nil {
			log.Error().Str("app_instance_id", appInstanceId).Msg("could not undeploy app")
			return err
		}
		cancel()
	}
	return nil
}

// Private function to undeploy a fragment from a cluster.
// params:
//  organizationId
//  fragmentId
//  targetCluster
// return:
//  error if any
func (c *Manager) undeployFragment(organizationId string, appInstanceId string, fragmentId string, targetCluster string) error {
	// Send undeploy requests to application clusters
	log.Debug().Str("app_instance_id", appInstanceId).Msg("undeploy fragment from cluster")

	// Unauthorize the members of this fragment
	fragment, err := c.AppClusterDB.GetFragmentsApp(targetCluster, appInstanceId)
	if err != nil {
		log.Error().Msg("impossible to get deployment fragment to unauthorize")
	}

	var targetFragment *entities.DeploymentFragment = nil
	for _, f := range fragment {
		// find the cluster we are looking for
		if f.FragmentId == fragmentId {
			targetFragment = &f
			break
		}
	}

	if targetFragment == nil {
		// this is extremely weird to occur
		log.Error().Msg("a deployment fragment could not be found. We cannot unauthorize fragment entries")
		return derrors.NewFailedPreconditionError("a deployment fragment could not be found. We cannot unauthorize fragment entries")
	}

	// unauthorize every service contained in the fragment
	for _, stage := range targetFragment.Stages {
		for _, serv := range stage.Services {
			// Unauthorize this entry
			ctxFrg, cancelFrg := context.WithTimeout(context.Background(), ConductorQueueTimeout)
			req := pbNetwork.DisauthorizeMemberRequest{
				AppInstanceId: appInstanceId, OrganizationId: organizationId,
				ServiceGroupInstanceId: serv.ServiceGroupInstanceId, ServiceApplicationInstanceId: serv.ServiceInstanceId}
			errPostMsg := c.NetworkOpsProducer.Send(ctxFrg, &req)
			cancelFrg()
			if errPostMsg != nil {
				log.Error().Err(errPostMsg).Interface("request", req).
					Msg("there was an error posting an undeploy fragment request")
			}
		}
	}

	errUpdate := c.ConnHelper.UpdateClusterConnections(organizationId)
	if errUpdate != nil {
		log.Error().Err(err).Str("organizationID", organizationId).Msg("error updating connections for organization")
		return err
	}
	if len(c.ConnHelper.ClusterReference) == 0 {
		log.Error().Msgf("no clusters found for organization %s", organizationId)
		return nil
	}

	clusterEntry, found := c.ConnHelper.ClusterReference[targetCluster]
	if !found {
		log.Error().Str("clusterId", targetCluster).Str("clusterHost", clusterEntry.Hostname).Msg("unknown clusterHost for the clusterId")
		return errors.New(fmt.Sprintf("unknown host for cluster id %s", targetCluster))
	}

	log.Debug().Str("targetCluster", targetCluster).Str("clusterHost", clusterEntry.Hostname).Msg("conductor query deployment-manager cluster")

	clusterAddress := fmt.Sprintf("%s:%d", clusterEntry.Hostname, utils.APP_CLUSTER_API_PORT)
	conn, errUpdate := c.ConnHelper.GetClusterClients().GetConnection(clusterAddress)
	if err != nil {
		log.Error().Err(err).Str("clusterHost", clusterEntry.Hostname).Msg("impossible to get connection for the host")
		return err
	}

	dmClient := pbAppClusterApi.NewDeploymentManagerClient(conn)

	undeployFragmentRequest := pbDeploymentManager.UndeployFragmentRequest{
		OrganizationId:       organizationId,
		DeploymentFragmentId: fragmentId,
		AppInstanceId:        appInstanceId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ConductorAppTimeout)
	defer cancel()
	_, errUpdate = dmClient.UndeployFragment(ctx, &undeployFragmentRequest)

	if err != nil {
		log.Error().Str("app_instance_id", appInstanceId).Msg("could not undeploy app")
		return err
	}

	return nil
}

func (c *Manager) expireUnifiedLogging(organizationId string, appInstanceId string) {

	// 4) NP-916 Link unified logging expire with the undeploy operation
	log.Info().Str("organizationID", organizationId).Str("instanceID", appInstanceId).
		Msg("Expire logging")
	_, err := c.UnifiedLoggingClient.Expire(context.Background(), &pbCoordinator.ExpirationRequest{
		OrganizationId: organizationId,
		AppInstanceId:  appInstanceId,
	})
	if err != nil {
		log.Warn().Str("error", conversions.ToDerror(err).DebugReport()).Msg("Error expiring unified Logging")
	}
	log.Info().Str("organizationID", organizationId).Str("instanceID", appInstanceId).
		Msg("Expire logging finish")

}

// Run operations to remove additional information generated in the system after instantiating an application.
func (c *Manager) Rollback(organizationId string, appInstanceId string, clusterIds []string) error {
	// Remove any related pending plan
	log.Debug().Str("appInstanceId", appInstanceId).Msg("Rollback app instance")

	// 1) Remove from the list of pendings plans
	c.PendingPlans.RemovePendingPlanByApp(appInstanceId)

	// 2) Delete zt network
	c.unauthorizeEntries(organizationId, appInstanceId, clusterIds)

	req := pbNetwork.DeleteNetworkRequest{AppInstanceId: appInstanceId, OrganizationId: organizationId}
	log.Debug().Interface("deleteNetworkRequest", req).Msg("delete zt network")
	_, err := c.NetClient.DeleteNetwork(context.Background(), &req)
	if err != nil {
		// TODO decide what to do here
		log.Error().Err(err).Str("appInstanceId", appInstanceId).Str("organizationId", organizationId).
			Msg("impossible to delete zerotier network")
	}

	// NP-1031. Improve undeploy performance
	go c.expireUnifiedLogging(organizationId, appInstanceId)

	// 3) Remove app entry points
	_, err = c.AppClient.RemoveAppEndpoints(context.Background(), &pbApplication.RemoveAppEndpointRequest{
		OrganizationId: organizationId,
		AppInstanceId:  appInstanceId,
	})
	if err != nil {
		log.Error().Err(err).Str("app_instance_id", appInstanceId).Msg("could not remove app endpoint  from system model")
	}

	// 4) Remove any service group instance
	_, err = c.AppClient.RemoveServiceGroupInstances(context.Background(), &pbApplication.RemoveServiceGroupInstancesRequest{
		OrganizationId: organizationId,
		AppInstanceId:  appInstanceId,
	})
	if err != nil {
		log.Error().Err(err).Str("app_instance_id", appInstanceId).Msg("could not remove service group instances")
	}

	return nil
}

// Unauthorize all the network entries for a certain application
// params:
//  organizationId
//  appInstanceId
//  clusterIds
func (c *Manager) unauthorizeEntries(organizationId string, appInstanceId string, clusterIds []string) {
	for _, clusterId := range clusterIds {
		fragments, err := c.AppClusterDB.GetFragmentsApp(clusterId, appInstanceId)
		if err != nil {
			log.Error().Err(err).Msg("impossible to find fragments to unauthorize entries")
			continue
		}
		// unauthorize every entry
		for _, f := range fragments {
			for _, ds := range f.Stages {
				for _, serv := range ds.Services {
					unauthorizeReq := pbNetwork.DisauthorizeMemberRequest{
						OrganizationId:               serv.OrganizationId,
						AppInstanceId:                serv.AppInstanceId,
						ServiceGroupInstanceId:       serv.ServiceGroupInstanceId,
						ServiceApplicationInstanceId: serv.ServiceInstanceId,
					}
					ctx, cancel := context.WithTimeout(context.Background(), ConductorQueueTimeout)
					err := c.NetworkOpsProducer.Send(ctx, &unauthorizeReq)
					cancel()
					if err != nil {
						log.Error().Err(err).Msg("problem sending unauthorize request to the queue")
					}
				}
			}
		}
	}
}

// Create the virtual application addresses for a given application descriptor
// params:
//  appDescriptor requiring the VAS entries
//  appInstanceId to work with
// return:
//  map with the VSA list
//  error if the operation failed
func (c *Manager) createVSA(appDescriptor entities.AppDescriptor, appInstanceId string) (map[string]string, derrors.Error) {
	currentIp := net.ParseIP(ConductorBaseVSA).To4()
	// store the generated vsa
	vsa := make(map[string]string, 0)
	servicesByName := make(map[string]entities.Service)

	for _, sg := range appDescriptor.Groups {
		for _, serv := range sg.Services {
			servicesByName[serv.Name] = serv
			fqdn := utils.GetVSAName(serv.Name, appDescriptor.OrganizationId, appInstanceId)
			dnsRequest := pbNetwork.AddDNSEntryRequest{
				OrganizationId: serv.OrganizationId,
				ServiceName:    serv.Name,
				Fqdn:           fqdn,
				Ip:             currentIp.String(),
				Tags: []string{
					fmt.Sprintf("appInstanceId:%s", appInstanceId),
					fmt.Sprintf("organizationId:%s", appDescriptor.OrganizationId),
					fmt.Sprintf("descriptorId:%s", appDescriptor.AppDescriptorId),
					fmt.Sprintf("serviceGroupId:%s", sg.ServiceGroupId),
					fmt.Sprintf("serviceId:%s", serv.ServiceId),
				},
			}
			ctx, cancel := context.WithTimeout(context.Background(), ConductorQueueTimeout)
			err := c.NetworkOpsProducer.Send(ctx, &dnsRequest)
			cancel()
			if err != nil {
				log.Error().Err(err).Interface("request", dnsRequest).Msg("impossible to send a dns entry request")
				return nil, err
			}
			vsa[fqdn] = currentIp.String()
			// Increase the IP
			currentIp = utils.NextIP(currentIp, 1)
		}
	}

	currentOutboundIp := net.ParseIP(ConductorOutboundVSA).To4()
	for _, securityRule := range appDescriptor.Rules {
		// TODO beware to not overflow the maximum length for a DNS name (253 chars)
		fqdn := utils.GetVSAName(securityRule.TargetServiceName+OutboundSuffix, appDescriptor.OrganizationId, appInstanceId)
		targetService := servicesByName[securityRule.TargetServiceName]
		dnsRequest := pbNetwork.AddDNSEntryRequest{
			OrganizationId: securityRule.OrganizationId,
			ServiceName:    targetService.Name,
			AppInstanceId:  appInstanceId,
			Fqdn:           fqdn,
			Ip:             currentOutboundIp.String(),
			Tags: []string{
				fmt.Sprintf("appInstanceId:%s", appInstanceId),
				fmt.Sprintf("organizationId:%s", appDescriptor.OrganizationId),
				fmt.Sprintf("descriptorId:%s", appDescriptor.AppDescriptorId),
				fmt.Sprintf("serviceGroupId:%s", targetService.ServiceGroupId),
				fmt.Sprintf("serviceId:%s", targetService.ServiceId),
				fmt.Sprintf("securityRule:%s", securityRule.Name),
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), ConductorQueueTimeout)
		err := c.NetworkOpsProducer.Send(ctx, &dnsRequest)
		cancel()
		if err != nil {
			log.Error().Err(err).Interface("request", dnsRequest).Msg("impossible to send a dns entry request")
			return nil, err
		}
		vsa[fqdn] = currentOutboundIp.String()
		// Increase the IP
		currentOutboundIp = utils.NextIP(currentOutboundIp, 1)
	}

	return vsa, nil
}

// Drain a cluster if and only if it is already cordoned, removed all the running applications and schedule the removed
// fragments.
func (c *Manager) DrainCluster(drainRequest *pbConductor.DrainClusterRequest) {
	// get all the apps deployed on that cluster

	fragmentIds, err := c.AppClusterDB.GetFragmentsInCluster(drainRequest.ClusterId.ClusterId)
	if err != nil {
		log.Error().Err(err).Str("clusterId", drainRequest.ClusterId.ClusterId).
			Msg("impossible to obtain fragments running in cluster")
		return
	}

	if len(fragmentIds) == 0 {
		log.Info().Str("clusterId", drainRequest.ClusterId.ClusterId).
			Msg("nothing to do for drain. Target cluster has no running deployments")
		return
	}

	// entries to schedule again
	toReschedule := make([]observer.ObservableDeploymentFragment, len(fragmentIds))

	for i, fragId := range fragmentIds {
		toReschedule[i] = observer.ObservableDeploymentFragment{ClusterId: drainRequest.ClusterId.ClusterId,
			FragmentId: fragId.FragmentId, AppInstanceId: fragId.AppInstanceId}
	}

	// reschedule removed deployment fragments
	log.Info().Str("clusterId", drainRequest.ClusterId.ClusterId).
		Int("numFragmentsToReschedule", len(toReschedule)).
		Msg("schedule drained operations to be scheduled again...")

	observer := observer.NewDeploymentFragmentsObserver(toReschedule, c.AppClusterDB)
	// Run an observer in a separated thread to send the schedule to the queue when is terminating
	go observer.ObserveOrganizationLevel(ConductorDrainClusterAppTimeout, entities.FRAGMENT_TERMINATING,
		drainRequest.ClusterId.OrganizationId, c.scheduleDeploymentFragment, nil)

	log.Info().Str("clusterId", drainRequest.ClusterId.ClusterId).Msg("schedule drained operations to be scheduled again done")

	// Drain the whole cluster
	log.Info().Str("clusterId", drainRequest.ClusterId.ClusterId).Msg("start cluster drain operation...")
	for _, fragment := range toReschedule {
		c.undeployFragment(drainRequest.ClusterId.OrganizationId, fragment.AppInstanceId, fragment.FragmentId, drainRequest.ClusterId.ClusterId)
	}
	log.Info().Str("clusterId", drainRequest.ClusterId.ClusterId).Msg("cluster drain operation complete")
}

// This function schedules a existing deployment fragment to be deployed again an updates the corresponding db status.
// params:
//  d deployment fragment to be deployed again
// return:
//  error if any
func (c *Manager) scheduleDeploymentFragment(d *entities.DeploymentFragment) derrors.Error {
	log.Debug().Str("deploymentFragmentId", d.DeploymentId).Msg("deployment fragment to be re-scheduled")
	// Update the application instance removing the affected service group
	ctx, cancel := context.WithTimeout(context.Background(), ConductorAppTimeout*time.Second)
	defer cancel()
	appInstance, errGet := c.AppClient.GetAppInstance(ctx, &pbApplication.AppInstanceId{OrganizationId: d.OrganizationId, AppInstanceId: d.AppInstanceId})
	if errGet != nil {
		log.Error().Err(errGet).Msg("impossible to get application instance when scheduling deployment fragment")
		return derrors.NewInternalError("impossible to get application instance when scheduling deployment fragment", errGet)
	}

	// Delete the associated service group instance
	serviceGroups := appInstance.Groups
	indexToDelete := 0
	for i, sg := range serviceGroups {
		if sg.ServiceGroupInstanceId == d.Stages[0].Services[0].ServiceGroupInstanceId {
			indexToDelete = i
			break
		}
	}
	serviceGroups = append(serviceGroups[:indexToDelete], serviceGroups[indexToDelete+1:]...)
	appInstance.Groups = serviceGroups

	ctx2, cancel2 := context.WithTimeout(context.Background(), ConductorAppTimeout*time.Second)
	defer cancel2()
	_, errUpdate := c.AppClient.UpdateAppInstance(ctx2, appInstance)
	if errUpdate != nil {
		log.Error().Err(errUpdate).Msg("impossible to update application when scheduling deployment fragment")
		return derrors.NewInternalError("impossible to update application when scheduling deployment fragment", errUpdate)
	}

	err := c.AppClusterDB.DeleteDeploymentFragment(d.ClusterId, d.FragmentId)
	if err != nil {
		log.Error().Err(err).Msg("impossible to update local version of deployment fragment")
		return err
	}

	// remove the pending fragment
	c.PendingPlans.RemoveFragment(d.FragmentId)

	deploymentErr := c.ProcessDeploymentFragment(d)
	if deploymentErr != nil {
		// if the deployment of a fragment fails we can consider this application instance to be in an error state
		log.Error().Str("appInstanceId", d.AppInstanceId).Str("fragmentId", d.FragmentId).
			Msg("fragment deployment failed")

		// get all the clusters for this application
		errUndeploy := c.SoftUndeploy(d.OrganizationId, d.AppInstanceId)
		if errUndeploy != nil {
			log.Error().Err(errUndeploy).Msg("problem during soft undeploy of application")
			return deploymentErr
		}

		statusReq := pbApplication.UpdateAppStatusRequest{AppInstanceId: d.AppInstanceId, OrganizationId: d.OrganizationId,
			Info:   fmt.Sprintf("impossible to find candidates for deployment of required service group replica %s", d.Stages[0].Services[0].ServiceGroupId),
			Status: pbApplication.ApplicationStatus_PLANNING_ERROR}
		ctxErr, cancelErr := context.WithTimeout(context.Background(), ConductorAppTimeout*time.Second)
		defer cancelErr()
		_, updateErr := c.AppClient.UpdateAppStatus(ctxErr, &statusReq)
		if updateErr != nil {
			log.Error().Str("appInstanceId", d.AppInstanceId).Msg("error updating status")
		}
	}

	return deploymentErr
}

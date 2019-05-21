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
    "github.com/nalej/derrors"
    pbAppClusterApi "github.com/nalej/grpc-app-cluster-api-go"
    pbApplication "github.com/nalej/grpc-application-go"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbDeploymentManager "github.com/nalej/grpc-deployment-manager-go"
    pbNetwork "github.com/nalej/grpc-network-go"
    pbCoordinator "github.com/nalej/grpc-unified-logging-go"
    "github.com/nalej/grpc-utils/pkg/conversions"
    "github.com/rs/zerolog/log"
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
}

func NewManager(connHelper *utils.ConnectionsHelper, queue structures.RequestsQueue, scorer scorer.Scorer,
    reqColl requirementscollector.RequirementsCollector, designer plandesigner.PlanDesigner,
    pendingPlans *structures.PendingPlans, appClusterDB *app_cluster.AppClusterDB) *Manager {
    // initialize clients
    pool := connHelper.GetSystemModelClients()
    if pool!=nil && len(pool.GetConnections())==0{
        log.Panic().Msg("system model clients were not started")
        return nil
    }
    conn := pool.GetConnections()[0]
    // Create associated clients
    appClient := pbApplication.NewApplicationsClient(conn)
    // Network client
    netPool := connHelper.GetNetworkingClients()
    if netPool != nil && len(netPool.GetConnections())==0{
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
        Designer: designer, AppClient:appClient, PendingPlans: pendingPlans, NetClient: netClient,
        DNSClient: dnsClient, UnifiedLoggingClient:ulClient, AppClusterDB: appClusterDB}
}

// Check iteratively if there is anything to be processed in the queue.
func (c *Manager) Run() {
	sleep := time.Tick(time.Millisecond * CheckSleepTime)
	for {
		select {
		case <-sleep:
		    //TODO revisit this solution because it could lead to intensive active queue checking
		    forNextIteration := make([]*entities.DeploymentRequest,0)
		    for c.Queue.AvailableRequests() {
                log.Info().Int("queued requests", c.Queue.Len()).Msg("there are pending deployment requests")
			    next := c.Queue.NextRequest()
                readyToProcess := true
			    // Check if there was any error with this deployment
			    if next.NumRetries > 0 {
			        // this app had a retry, check if enough time passed since the last check
			        elapsedTime := time.Now().Unix()-next.TimeRetry.Unix()
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
                log.Info().Int("pending",len(forNextIteration)).Msg("some deployments were excluded in this round")
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
                AppInstanceId: req.InstanceId,
                OrganizationId: req.OrganizationId,
                Status: pbApplication.ApplicationStatus_DEPLOYMENT_ERROR,
                Info: fmt.Sprintf("Exceeded number of retries. Latest known error: [%v]", err.Error()),
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
        // TODO review how to proceed with zt network id
        c.Rollback(req.OrganizationId, req.InstanceId)
    }
}

// Push a request into the queue.
func(c *Manager) PushRequest(req *pbConductor.DeploymentRequest)  error{
    log.Debug().Interface("request",req).Msg("received deployment request")
    // Get ParameterizedDescriptor
    desc, err := c.AppClient.GetParametrizedDescriptor(context.Background(), req.AppInstanceId)
    if err!= nil {
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
        log.Error().Err(err).Msgf("problems updating application %s",req.AppInstanceId.AppInstanceId)
        return  err
    }

    return nil
}

func(c *Manager) ProcessDeploymentRequest(req *entities.DeploymentRequest) derrors.Error{
    if req == nil {
        err := derrors.NewFailedPreconditionError("the queue was unexpectedly empty")
        log.Error().Err(err)
        return err
    }

    // TODO get all the data from the system model
    // Get the ServiceGroup structure

    retrievedAppInstance, err  := c.AppClient.GetAppInstance(context.Background(),
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
    scoreResult, err := c.ScorerMethod.ScoreRequirements (req.OrganizationId,foundRequirements)

    if err != nil {
        err := derrors.NewGenericError("error scoring request")
        log.Error().Err(err).Str("requestId",req.RequestId).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
        return err
    }

    log.Info().Msgf("conductor maximum score for %s has score %v from %d potential candidates",
        req.RequestId, scoreResult.DeploymentsScore, scoreResult.NumEvaluatedClusters)


    // 3) design plan
    // Elaborate deployment plan for the application
    plan, err := c.Designer.DesignPlan(appInstance, *scoreResult, *req, nil,nil)

    if err != nil{
        log.Error().Err(err).Str("requestId",req.RequestId).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
        return derrors.AsError(err,fmt.Sprintf("plan design failed for descriptor %s",err.Error()))
    }

    // 4) Create ZT-network with Network manager
    // we use the app instance id as the network id
    ztNetworkId, zt_err := c.CreateZTNetwork(retrievedAppInstance.AppInstanceId, req.OrganizationId, retrievedAppInstance.AppInstanceId)
    if zt_err != nil {
        err := derrors.NewGenericError("impossible to create zt network before deployment", zt_err)
        log.Error().Err(zt_err).Str("requestId",req.RequestId).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
        return err
    }

    // 5) deploy fragments
    // Tell deployment managers to execute plans
    err_deploy := c.DeployPlan(plan, ztNetworkId, req.NumRetries)
    if err_deploy != nil {
        err := derrors.NewGenericError("error deploying plan request", err_deploy)
        log.Error().Err(err_deploy).Str("requestId",req.RequestId).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
        return err
    }
    return nil
}

// Analyze the best deployment options for a single deployment fragment.
func(c *Manager) ProcessDeploymentFragment(fragment *entities.DeploymentFragment) derrors.Error{
    if fragment == nil {
        err := derrors.NewFailedPreconditionError("the requested fragment was nil")
        log.Error().Err(err)
        return err
    }

    // Get the ServiceGroup structure

    retrievedAppInstance, err  := c.AppClient.GetAppInstance(context.Background(),
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

    log.Debug().Interface("re-scheduled fragment",fragment).Msg("this is the fragment to deploy")

    // Find the service groups in the deployment fragment
    serviceGroupsMap := make(map[string]bool,0)
    serviceGroupIds := make([]string, 0)
    for _, sg := range fragment.Stages{
        for _, s := range sg.Services {
            if _, found := serviceGroupsMap[s.ServiceGroupId]; !found {
                log.Debug().Str("groupName",s.Name).Str("serviceGroupId",s.ServiceGroupId).Msg("services contained in running service group")
                serviceGroupsMap[s.ServiceGroupId]=true
                serviceGroupIds = append(serviceGroupIds, s.ServiceGroupId)
            } else {
                log.Debug().Str("groupName",s.Name).Str("serviceGroupId",s.ServiceGroupId).Msg("excluded from the list")
            }
        }
    }
    log.Debug().Interface("serviceGroupIds", serviceGroupIds).
        Msgf("a total of %d service groups have to be rescheduled",len(serviceGroupIds))


    foundRequirements, err := c.ReqCollector.FindRequirementsForGroups(serviceGroupIds, appInstance.AppInstanceId, appDescriptor)
    if err != nil {
        err := derrors.NewGenericError("impossible to find requirements for application")
        log.Error().Err(err).Str("appDescriptorId", fragment.AppDescriptorId)
        return err
    }

    // 2) score requirements
    scoreResult, err := c.ScorerMethod.ScoreRequirements (fragment.OrganizationId,foundRequirements)

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
        AppInstanceId: fragment.AppInstanceId,
        ApplicationId: fragment.AppDescriptorId,
        RequestId: uuid.New().String(),
        InstanceId: fragment.FragmentId,
    }

    // design a plan for the service groups contained into the deployment fragment
    // build a summary of the groups running in the cluster
    allocatedGroupsPerClusters := make(map[string][]string,0)
    for clusterId, _ := range c.ConnHelper.ClusterReference {
        df, err := c.AppClusterDB.GetDeploymentFragment(clusterId,fragment.FragmentId)
        if err != nil {
            log.Error().Msg("error when getting deployment fragment from cluster")
            continue
        }
        if df != nil {
            entry, found := allocatedGroupsPerClusters[clusterId]
            if !found {
                allocatedGroupsPerClusters[clusterId] = make([]string,0)
                entry = allocatedGroupsPerClusters[clusterId]
            }
            // every stage corresponds to one service group, any service has the service group id
            for _, s := range df.Stages {
                entry = append(entry, s.Services[0].ServiceGroupId)
            }
            allocatedGroupsPerClusters[clusterId] = entry
        }
    }

    log.Debug().Interface("allocatedGroupsPerCluster",allocatedGroupsPerClusters).
        Interface("serviceGroupIds", serviceGroupIds).
        Msg("design a plan for deployment fragments")
    plan, err := c.Designer.DesignPlan(appInstance, *scoreResult, req, serviceGroupIds, allocatedGroupsPerClusters)

    if err != nil{
        log.Error().Err(err).Str("appDescriptorId", fragment.AppDescriptorId)
        return derrors.AsError(err,fmt.Sprintf("plan design failed for descriptor %s",err.Error()))
    }

    // 4) deploy fragments
    // Tell deployment managers to execute plans
    err_deploy := c.DeployPlan(plan, fragment.ZtNetworkID, 0)
    if err_deploy != nil {
        err := derrors.NewGenericError("error deploying plan request", err_deploy)
        log.Error().Err(err_deploy).Str("requestId",req.RequestId).Str("appDescriptorId", retrievedAppInstance.AppDescriptorId)
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
func (c *Manager) CreateZTNetwork(name string, organizationId string, appInstanceId string) (string, error){
    request := pbNetwork.AddNetworkRequest{ Name: name, OrganizationId: organizationId, AppInstanceId: appInstanceId }

    log.Debug().Interface("addNetworkRequest",request).Msgf("create a network request")

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
    for _, fragment := range plan.Fragments{
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
        log.Info().Str("deploymentId",fragment.DeploymentId).
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

        ctx, cancel := context.WithTimeout(context.Background(), time.Second * ConductorAppTimeout)
        defer cancel()
        response, err := client.Execute(ctx, &request)


        log.Debug().Interface("deploymentFragmentResponse", response).Interface("deploymentFragmentError",err).
            Msg("finished fragment deployment")

        if err != nil {
            // TODO define how to proceed in case of error
            log.Error().Err(err).Str("deploymentId",fragment.DeploymentId).Msg("problem deploying fragment")
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
func (c* Manager) Undeploy (request *entities.UndeployRequest) error {
    return c.HardUndeploy(request.OrganizationId,request.AppInstanceId)
}

// Undeploy function that maintains the application instance in the system.
func(c *Manager) SoftUndeploy(organizationId string, appInstanceId string) error {
    // find application instance
    appInstance, err  := c.AppClient.GetAppInstance(context.Background(),
        &pbApplication.AppInstanceId{OrganizationId: organizationId, AppInstanceId: appInstanceId})
    if err != nil {
        log.Error().Err(err).Str("appInstanceID",appInstanceId).Msgf("impossible to obtain application descriptor")
        return err
    }

    // call Rollback
    c.Rollback(organizationId, appInstanceId)

    // find in what clusters this app is running
    clusterFound := make(map[string]bool, 0)
    clusterIds := make([]string,0)
    for _, g := range appInstance.Groups {
        for _, svc := range g.ServiceInstances {
            if svc.DeployedOnClusterId != "" {
                clusterFound[svc.DeployedOnClusterId] = true
                clusterIds = append(clusterIds, svc.DeployedOnClusterId)
            }
        }
    }


    // create observer and the array of entries to be observed
    toObserve := make([]observer.ObservableDeploymentFragment, 0)
    for _, clusterId := range clusterIds {
        fragments, err := c.AppClusterDB.GetFragmentsInCluster(clusterId)
        if err != nil {
            log.Error().Err(err).Msg("error when getting fragments from cluster")
            continue
        }
        for _, fr := range fragments {
            toObserve = append(toObserve, observer.ObservableDeploymentFragment{ClusterId: clusterId,
                FragmentId: fr.FragmentId, AppInstanceId: appInstance.AppInstanceId})
        }
    }

    observer := observer.NewDeploymentFragmentsObserver(toObserve, c.AppClusterDB)
    // Run an observer in a separated thread to send the schedule to the queue when is terminating
    go observer.Observe(ConductorDrainClusterAppTimeout,entities.FRAGMENT_TERMINATING,
        func(d *entities.DeploymentFragment) derrors.Error {
            return c.AppClusterDB.DeleteDeploymentFragment(d.ClusterId,d.FragmentId)
        })

    // terminate execution
    return c.undeployClustersInstance(appInstance.OrganizationId, appInstance.AppInstanceId,clusterIds)
}

// Undeploy function that removes the application instance
func(c *Manager) HardUndeploy(organizationId string, appInstanceId string) error {
    // find application instance
    appInstance, err  := c.AppClient.GetAppInstance(context.Background(),
        &pbApplication.AppInstanceId{OrganizationId: organizationId, AppInstanceId: appInstanceId})
    if err != nil {
        log.Error().Err(err).Str("appInstanceID",appInstanceId).Msgf("impossible to obtain application descriptor")
        return err
    }

    // call Rollback
    c.Rollback(organizationId, appInstanceId)

    // Remove any entry from the system model
    instID := &pbApplication.AppInstanceId{
        OrganizationId: organizationId,
        AppInstanceId: appInstanceId,
    }
    _, err = c.AppClient.RemoveAppInstance(context.Background(), instID)
    if err != nil{
        log.Error().Err(err).Str("app_instance_id", appInstanceId).Msg("could not remove instance from system model")
    }

    // TODO: see if this method could be moved to application manager ( app-manager is responsible for creating it)
    _, err = c.AppClient.RemoveParametrizedDescriptor(context.Background(), instID)
    if err != nil{
        log.Error().Err(err).Str("app_instance_id", appInstanceId).Msg("could not remove parametrized descriptor from system model")
    }

    // Remove from the associated request from the queue
    removed := c.Queue.Remove(appInstanceId)
    if !removed {
        log.Info().Interface("appInstanceId", appInstanceId).Msg("no request was found in the queue for this deployed app")
    }

    // find in what clusters this app is running
    clusterFound := make(map[string]bool, 0)
    clusterIds := make([]string,0)
    for _, g := range appInstance.Groups {
        for _, svc := range g.ServiceInstances {
            if svc.DeployedOnClusterId != "" {
                clusterFound[svc.DeployedOnClusterId] = true
                clusterIds = append(clusterIds, svc.DeployedOnClusterId)
            }
        }
    }


    // create observer and the array of entries to be observed
    toObserve := make([]observer.ObservableDeploymentFragment, 0)
    for _, clusterId := range clusterIds {
        fragments, err := c.AppClusterDB.GetFragmentsInCluster(clusterId)
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
    go observer.Observe(ConductorDrainClusterAppTimeout,entities.FRAGMENT_TERMINATING,
        func(d *entities.DeploymentFragment) derrors.Error {
            return c.AppClusterDB.DeleteDeploymentFragment(d.ClusterId,d.FragmentId)
        })

    // terminate execution
    return c.undeployClustersInstance(appInstance.OrganizationId, appInstance.AppInstanceId,clusterIds)
}

// Private function to communicate the application clusters to remove a running instance.
// params:
//  organizationId
//  appInstanceId
//  targetClusters list with the cluster ids where this function must run. If no list is set, then all clusters are informed.
// return:
//  error if any
func(c *Manager) undeployClustersInstance(organizationId string, appInstanceId string, targetClusters []string) error {
    // Send undeploy requests to application clusters
    log.Debug().Str("app_instance_id", appInstanceId).Msg("undeploy app instance with id")

    err := c.ConnHelper.UpdateClusterConnections(organizationId)
    if err != nil {
        log.Error().Err(err).Str("organizationID",organizationId).Msg("error updating connections for organization")
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

    for _, clusterId := range targetClusters{

        clusterEntry, found := c.ConnHelper.ClusterReference[clusterId]
        if !found {
            log.Error().Str("clusterId",clusterId).Str("clusterHost",clusterEntry.Hostname).Msg("unknown clusterHost for the clusterId")
            return errors.New(fmt.Sprintf("unknown host for cluster id %s", clusterId))
        }

        log.Debug().Str("clusterId", clusterId).Str("clusterHost", clusterEntry.Hostname).Msg("conductor query deployment-manager cluster")


        clusterAddress := fmt.Sprintf("%s:%d",clusterEntry.Hostname,utils.APP_CLUSTER_API_PORT)
        conn, err := c.ConnHelper.GetClusterClients().GetConnection(clusterAddress)
        if err != nil {
            log.Error().Err(err).Str("clusterHost", clusterEntry.Hostname).Msg("impossible to get connection for the host")
            return err
        }

        dmClient := pbAppClusterApi.NewDeploymentManagerClient(conn)

        undeployRequest := pbDeploymentManager.UndeployRequest{
            OrganizationId: organizationId,
            AppInstanceId: appInstanceId,
        }
        ctx, cancel := context.WithTimeout(context.Background(), time.Second * ConductorAppTimeout)
        _, err = dmClient.Undeploy(ctx, &undeployRequest)

        if err != nil {
            log.Error().Str("app_instance_id", appInstanceId).Msg("could not undeploy app")
            return err
        }
        cancel()
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
func (c *Manager) Rollback(organizationId string, appInstanceId string) error {
    // Remove any related pending plan
    log.Debug().Str("appInstanceId",appInstanceId).Msg("Rollback app instance")

    // 1) Remove from the list of pendings plans
    c.PendingPlans.RemovePendingPlanByApp(appInstanceId)


    // 2) Delete zt network
    req := pbNetwork.DeleteNetworkRequest{AppInstanceId: appInstanceId, OrganizationId: organizationId}
    log.Debug().Interface("deleteNetworkRequest", req).Msg("delete zt network")
    _, err := c.NetClient.DeleteNetwork(context.Background(), &req)
    if err != nil {
        // TODO decide what to do here
        log.Error().Err(err).Str("appInstanceId",appInstanceId).Str("organizationId", organizationId).
            Msg("impossible to delete zerotier network")
    }

    // 3) Remove associated DNS entries if any
    log.Debug().Msgf("remove DNS entries for %s in %s",appInstanceId, organizationId)
    deleteReq := pbNetwork.DeleteDNSEntryRequest{
        OrganizationId: organizationId,
        AppInstanceId: appInstanceId,
    }

    _, err = c.DNSClient.DeleteDNSEntry(context.Background(), &deleteReq)
    if err != nil {
        // TODO decide what to do here
        log.Error().Str("appInstanceId",appInstanceId).Err(err).
            Msgf("error removing dns entries for appInstance %s", deleteReq.OrganizationId)
    }

    // NP-1031. Improve undeploy performance
    go c.expireUnifiedLogging(organizationId, appInstanceId)
      
    // 5) Remove app entry points  
    _, err = c.AppClient.RemoveAppEndpoints(context.Background(), &pbApplication.RemoveAppEndpointRequest{
        OrganizationId: organizationId,
        AppInstanceId: appInstanceId,
    })
    if err != nil{
        log.Error().Err(err).Str("app_instance_id", appInstanceId).Msg("could not remove app endpoint  from system model")
    }

    // 6) Remove any service group instance
    _, err = c.AppClient.RemoveServiceGroupInstances(context.Background(), &pbApplication.RemoveServiceGroupInstancesRequest{
        OrganizationId: organizationId,
        AppInstanceId: appInstanceId,
    })
    if err != nil {
        log.Error().Err(err).Str("app_instance_id", appInstanceId).Msg("could not remove service group instances")
    }

    return nil
}


// Drain a cluster if and only if it is already cordoned, removed all the running applications and schedule the removed
// fragments.
func (c *Manager) DrainCluster(drainRequest *pbConductor.DrainClusterRequest) {
    // get all the apps deployed on that cluster

    fragmentIds, err := c.AppClusterDB.GetFragmentsInCluster(drainRequest.ClusterId.ClusterId)
    if err != nil {
        log.Error().Err(err).Str("clusterId",drainRequest.ClusterId.ClusterId).
            Msg("impossible to obtain fragments running in cluster")
        return
    }

    if len(fragmentIds) == 0 {
        log.Info().Str("clusterId",drainRequest.ClusterId.ClusterId).
            Msg("nothing to do for drain. Target cluster has no running deployments")
        return
    }


    // entries to schedule again
    toReschedule := make([]observer.ObservableDeploymentFragment,len(fragmentIds))

    for i, fragId := range fragmentIds {
        toReschedule[i] = observer.ObservableDeploymentFragment{ClusterId: drainRequest.ClusterId.ClusterId,
            FragmentId: fragId.FragmentId, AppInstanceId: fragId.AppInstanceId}
    }

    // reschedule removed deployment fragments
    log.Info().Str("clusterId",drainRequest.ClusterId.ClusterId).
        Int("numFragmentsToReschedule", len(toReschedule)).
        Msg("schedule drained operations to be scheduled again...")

    observer := observer.NewDeploymentFragmentsObserver(toReschedule, c.AppClusterDB)
    // Run an observer in a separated thread to send the schedule to the queue when is terminating
    go observer.Observe(ConductorDrainClusterAppTimeout,entities.FRAGMENT_TERMINATING, c.scheduleDeploymentFragment)

    log.Info().Str("clusterId",drainRequest.ClusterId.ClusterId).Msg("schedule drained operations to be scheduled again done")

    // Drain the whole cluster
    log.Info().Str("clusterId",drainRequest.ClusterId.ClusterId).Msg("start cluster drain operation...")
    for _, fragment := range toReschedule {
        c.undeployClustersInstance(drainRequest.ClusterId.OrganizationId, fragment.AppInstanceId, []string{drainRequest.ClusterId.ClusterId})
    }
    log.Info().Str("clusterId",drainRequest.ClusterId.ClusterId).Msg("cluster drain operation complete")
}

// This function schedules a existing deployment fragment to be deployed again an updates the corresponding db status.
// params:
//  d deployment fragment to be deployed again
// return:
//  error if any
func (c *Manager) scheduleDeploymentFragment(d *entities.DeploymentFragment) derrors.Error {
    log.Debug().Str("deploymentFragmentId", d.DeploymentId).Msg("deployment fragment to be re-scheduled")
    // Update the application instance removing the affected service group
    ctx, cancel :=  context.WithTimeout(context.Background(), ConductorAppTimeout * time.Second)
    defer cancel()
    appInstance, errGet := c.AppClient.GetAppInstance(ctx,&pbApplication.AppInstanceId{OrganizationId: d.OrganizationId, AppInstanceId: d.AppInstanceId})
    if errGet!= nil {
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


    ctx2, cancel2 :=  context.WithTimeout(context.Background(), ConductorAppTimeout * time.Second)
    defer cancel2()
    _, errUpdate := c.AppClient.UpdateAppInstance(ctx2,appInstance)
    if errUpdate != nil {
        log.Error().Err(errUpdate).Msg("impossible to update application when scheduling deployment fragment")
        return derrors.NewInternalError("impossible to update application when scheduling deployment fragment", errUpdate)
    }


    err := c.AppClusterDB.DeleteDeploymentFragment(d.ClusterId, d.FragmentId)
    if err != nil {
        log.Error().Err(err).Msg("impossible to update local version of deployment fragment")
        return err
    }

    return c.ProcessDeploymentFragment(d)
}
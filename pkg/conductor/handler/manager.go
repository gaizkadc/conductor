/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package handler

import (
    "context"
    "errors"
    "fmt"
    "github.com/google/uuid"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/conductor/monitor"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/conductor/requirementscollector"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/conductor/pkg/utils"
    pbAppClusterApi "github.com/nalej/grpc-app-cluster-api-go"
    pbApplication "github.com/nalej/grpc-application-go"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbDeploymentManager "github.com/nalej/grpc-deployment-manager-go"
    pbNetwork "github.com/nalej/grpc-network-go"
    "github.com/rs/zerolog/log"
    "time"
)

// Time to wait between checks in the queue in milliseconds.
const CheckSleepTime = 2000

type Manager struct {
    // Connections helper
    ConnHelper *utils.ConnectionsHelper
    // ScorerMethod
    ScorerMethod scorer.Scorer
    // Requirements collector
    ReqCollector requirementscollector.RequirementsCollector
    // Plan designer
    Designer plandesigner.PlanDesigner
    // Queue for incoming requests
    Queue RequestsQueue
    // Monitoring service
    Monitor monitor.Manager
    // Application client
    AppClient pbApplication.ApplicationsClient
    // Networking manager client
    NetClient pbNetwork.NetworksClient
    // DNS manager client
    DNSClient pbNetwork.DNSClient
}

func NewManager(connHelper *utils.ConnectionsHelper, queue RequestsQueue, scorer scorer.Scorer,
    reqColl requirementscollector.RequirementsCollector, designer plandesigner.PlanDesigner, monitor monitor.Manager) *Manager {
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
    netClient := pbNetwork.NewNetworksClient(netPool.GetConnections()[0])
    dnsClient := pbNetwork.NewDNSClient(netPool.GetConnections()[0])
    return &Manager{ConnHelper: connHelper, Queue: queue, ScorerMethod: scorer, ReqCollector: reqColl,
        Designer: designer, AppClient:appClient, Monitor: monitor, NetClient: netClient, DNSClient: dnsClient}
}

// Check iteratively if there is anything to be processed in the queue.
func (c *Manager) Run() {
	sleep := time.Tick(time.Millisecond * CheckSleepTime)
	for {
		select {
		case <-sleep:
			for c.Queue.AvailableRequests() {
				c.ProcessDeploymentRequest()
			}
		}
	}
}

// Push a request into the queue.

func(c *Manager) PushRequest(req *pbConductor.DeploymentRequest) (*entities.DeploymentRequest, error){
    log.Debug().Msgf("push request %s", req.RequestId)
    desc, err := c.AppClient.GetAppDescriptor(context.Background(), req.AppId)
    if err!= nil {
        log.Error().Err(err).Msg("error getting application descriptor")
        return nil,err
    }
    // Create new application instance
    addReq := pbConductor.AddAppInstanceRequest{
        OrganizationId: desc.OrganizationId,
        AppDescriptorId: desc.AppDescriptorId,
        Name: req.Name,
        Description: req.Description,
    }
    // Add instance, by default this is created with queue status
    instance,err := c.AppClient.AddAppInstance(context.Background(),&addReq)
    if err != nil {
        log.Error().Err(err).Msg("error adding application instance")
        return nil,err
    }

    toEnqueue := entities.DeploymentRequest{
        RequestId:      req.RequestId,
        InstanceId:     instance.AppInstanceId,
        OrganizationId: req.AppId.OrganizationId,
        ApplicationId:  req.AppId.AppDescriptorId,
    }
    err = c.Queue.PushRequest(&toEnqueue)
    if err != nil {
        return &toEnqueue,err
    }

    if err != nil {
        log.Error().Err(err).Msgf("problems updating application %s",req.AppId)
        return &toEnqueue, err
    }

    return &toEnqueue, nil
}

func(c *Manager) ProcessDeploymentRequest(){
    req := c.Queue.NextRequest()
    if req == nil {
        log.Error().Msg("the queue was unexpectedly empty")
        return
    }

    // TODO get all the data from the system model
    // Get the ServiceGroup structure

    appInstance, err  := c.AppClient.GetAppInstance(context.Background(),
        &pbApplication.AppInstanceId{OrganizationId: req.OrganizationId, AppInstanceId: req.InstanceId})
    if err != nil {
        log.Error().Err(err).Msgf("impossible to obtain application descriptor %s", appInstance.AppDescriptorId)
        return
    }

    // 1) collect requirements for the application descriptor
    foundRequirements, err := c.ReqCollector.FindRequirements(appInstance)
    if err != nil {
        log.Error().Err(err).Msgf("impossible to find requirements for application %s", appInstance.AppDescriptorId)
        return
    }

    // 2) score requirements
    scoreResult, err := c.ScorerMethod.ScoreRequirements (req.OrganizationId,foundRequirements)

    if err != nil {
        log.Error().Err(err).Msgf("error scoring request %s", req.RequestId)
        return
    }

    log.Info().Msgf("conductor maximum score for %s has score %v from %d potential candidates",
        req.RequestId, scoreResult.Scoring, scoreResult.TotalEvaluated)


    // 3) design plan
    // TODO elaborate plan, modify system model accordingly
    // Elaborate deployment plan for the application
    plan, err := c.Designer.DesignPlan(appInstance, scoreResult)

    if err != nil{
        log.Error().Err(err).Msgf("error designing plan for request %s",req.RequestId)
        return
    }

    // 4) Create ZT-network with Network manager
    ztNetworkId, err := c.CreateZTNetwork(appInstance.AppInstanceId, req.OrganizationId)
    if err != nil {
        log.Error().Err(err).Msg("impossible to create zt network before deployment")
        return
    }

    // 5) deploy fragments
    // Tell deployment managers to execute plans
    err = c.DeployPlan(plan, ztNetworkId)
    if err != nil {
        log.Error().Err(err).Msgf("error deploying plan request %s", req.RequestId)
        // Run a rollback
        c.rollback(plan, ztNetworkId)
        return
    }
}

// Create a new zero tier network and return the corresponding network id.
// params:
//  name of the network
//  organizationId for this network
// returns:
//  networkId or error otherwise
func (c *Manager) CreateZTNetwork(name string, organizationId string) (string, error){
    log.Debug().Msgf("create zt network with name %s in organization %s",name, organizationId)
    request := pbNetwork.AddNetworkRequest{ Name: name, OrganizationId: organizationId }

    ztNetworkId, err := c.NetClient.AddNetwork(context.Background(), &request)

    if err != nil {
        log.Error().Err(err).Msgf("there was a problem when creating network for name: %s with org: %s", name, organizationId)
        return "", err
    }
    return ztNetworkId.NetworkId, err
}

// For a given collection of plans, tell the corresponding deployment managers to run the deployment.
func (c *Manager) DeployPlan(plan *entities.DeploymentPlan, ztNetworkId string) error {
    // Start monitoring this fragment
    c.Monitor.AddPlanToMonitor(plan)

    for fragmentIndex, fragment := range plan.Fragments {
        log.Info().Msgf("start fragment %s deployment with %d out of %d fragments", fragment.DeploymentId, fragmentIndex, len(plan.Fragments))

        targetHostname, found := c.ConnHelper.ClusterReference[fragment.ClusterId]
        if !found {
            msg := fmt.Sprintf("unknown target address for cluster with id %s", fragment.ClusterId)
            err := errors.New(msg)
            log.Error().Msgf(msg)
            return err
        }

        clusterAddress := fmt.Sprintf("%s:%d", targetHostname, utils.APP_CLUSTER_API_PORT)

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
        }

        client := pbAppClusterApi.NewDeploymentManagerClient(conn)
        _, err = client.Execute(context.Background(), &request)

        if err != nil {
            // TODO define how to proceed in case of error
            log.Error().Err(err).Msgf("problem deploying fragment %s", fragment.DeploymentId)
            return err
        }
    }



    return nil
}

// Undeploy
func (c* Manager) Undeploy (request *entities.UndeployRequest) error {

    log.Debug().Msgf("remove DNS entries for %s in %s",request.AppInstanceId,request.OrganizationId)
    deleteReq := pbNetwork.DeleteDNSEntryRequest{
        OrganizationId: request.OrganizationId,
        AppInstanceId: request.AppInstanceId,
    }

    _, err := c.DNSClient.DeleteDNSEntry(context.Background(), &deleteReq)
    if err != nil {
        log.Error().Err(err).Msgf("error removing dns entries for appInstance %s", deleteReq.OrganizationId)
    }


    log.Debug().Msgf("undeploy app instance with id %s",request.AppInstanceId)

    err = c.ConnHelper.UpdateClusterConnections(request.OrganizationId)
    if err != nil {
        log.Error().Err(err).Msgf("error updating connections for organization %s", request.OrganizationId)
        return err
    }
    if len(c.ConnHelper.ClusterReference) == 0 {
        log.Error().Msgf("no clusters found for organization %s", request.OrganizationId)
        return nil
    }

    log.Debug().Msgf("There are %d known clusters",len(c.ConnHelper.ClusterReference))

    for clusterId, clusterHost := range c.ConnHelper.ClusterReference {
        log.Debug().Msgf("conductor query deployment-manager cluster %s at %s", clusterId, clusterHost)

        clusterAddress := fmt.Sprintf("%s:%d",clusterHost,utils.APP_CLUSTER_API_PORT)
        conn, err := c.ConnHelper.GetClusterClients().GetConnection(clusterAddress)
        if err != nil {
            log.Error().Err(err).Msgf("impossible to get connection for %s",clusterHost)
            return err
        }

        dmClient := pbDeploymentManager.NewDeploymentManagerClient(conn)

        undeployRequest := pbDeploymentManager.UndeployRequest{
            OrganizationId: request.OrganizationId,
            AppInstanceId: request.AppInstanceId,
        }
        _, err = dmClient.Undeploy(context.Background(), &undeployRequest)
        if err != nil {
            log.Error().Msgf("could not undeploy app %s", request.AppInstanceId)
            return err
        }
    }

    return nil
}

// Return the system to the status before instantiating the given deployment plan and zt network id.
func (c *Manager) rollback (plan *entities.DeploymentPlan, ztNetworkId string) error {
    // Delete zt network
    req := pbNetwork.DeleteNetworkRequest{NetworkId: ztNetworkId, OrganizationId: plan.OrganizationId}
    _, err := c.NetClient.DeleteNetwork(context.Background(), &req)
    if err != nil {
        // TODO decide what to do here
        log.Error().Msgf("impossible to delete zerotier network %s", ztNetworkId)
    }

    // Remove associated DNS entries if any
    log.Debug().Msgf("remove DNS entries for %s in %s",plan.AppInstanceId,plan.OrganizationId)
    deleteReq := pbNetwork.DeleteDNSEntryRequest{
        OrganizationId: plan.OrganizationId,
        AppInstanceId: plan.AppInstanceId,
    }

    _, err = c.DNSClient.DeleteDNSEntry(context.Background(), &deleteReq)
    if err != nil {
        // TODO decide what to do here
        log.Error().Err(err).Msgf("error removing dns entries for appInstance %s", deleteReq.OrganizationId)
    }

    return nil
}







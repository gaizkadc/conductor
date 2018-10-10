/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package handler

import (
    "github.com/rs/zerolog/log"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
    pbDeploymentManager "github.com/nalej/grpc-deployment-manager-go"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "time"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/conductor"
    "context"
    "github.com/nalej/conductor/pkg/conductor/requirementscollector"
    "github.com/nalej/conductor/internal/entities"
    "github.com/google/uuid"
    "github.com/nalej/conductor/pkg/conductor/monitor"
)

// Time to wait between checks in the queue in milliseconds.
const CheckSleepTime = 2000

type Manager struct {
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
    appClient pbApplication.ApplicationsClient
}


func NewManager(queue RequestsQueue, scorer scorer.Scorer, reqColl requirementscollector.RequirementsCollector,
    designer plandesigner.PlanDesigner, monitor monitor.Manager) *Manager {
    // initialize clients
    pool := conductor.GetSystemModelClients()
    if pool!=nil && len(pool.GetConnections())==0{
        log.Panic().Msg("system model clients were not started")
        return nil
    }
    conn := pool.GetConnections()[0]
    appClient := pbApplication.NewApplicationsClient(conn)
    return &Manager{Queue: queue, ScorerMethod: scorer, ReqCollector: reqColl, Designer: designer, appClient:appClient,
    Monitor: monitor}
}

// Check iteratively if there is anything to be processed in the queue.
func (c *Manager) Run() {
    sleep := time.Tick(time.Millisecond * CheckSleepTime)
    for{
        select {
        case <- sleep:
            for c.Queue.AvailableRequests() {
                c.ProcessDeploymentRequest()
            }
        }
    }
}

// Push a request into the queue.
func(c *Manager) PushRequest(req *pbConductor.DeploymentRequest) (*entities.DeploymentRequest, error){
    log.Debug().Msgf("push request %s", req.RequestId)
    desc, err := c.appClient.GetAppDescriptor(context.Background(), req.AppId)
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
    instance,err := c.appClient.AddAppInstance(context.Background(),&addReq)
    if err != nil {
        log.Error().Err(err).Msg("error adding application instance")
        return nil,err
    }

    toEnqueue := entities.DeploymentRequest{
        RequestID: req.RequestId,
        InstanceID: instance.AppInstanceId,
        OrganizationID: req.AppId.OrganizationId,
        ApplicationID: req.AppId.AppDescriptorId,
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

    appInstance, err  := c.appClient.GetAppInstance(context.Background(),
        &pbApplication.AppInstanceId{OrganizationId: req.OrganizationID, AppInstanceId: req.InstanceID})
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
    scoreResult, err := c.ScorerMethod.ScoreRequirements (foundRequirements)

    if err != nil {
        log.Error().Err(err).Msgf("error scoring request %s", req.RequestID)
        return
    }

    log.Info().Msgf("conductor maximum score for %s is for cluster %s among %d possible",
        scoreResult.RequestID, scoreResult.ClusterID, scoreResult.TotalEvaluated)


    // 3) design plan
    // TODO elaborate plan, modify system model accordingly
    // Elaborate deployment plan for the application
    plan, err := c.Designer.DesignPlan(appInstance, scoreResult)

    if err != nil{
        log.Error().Err(err).Msgf("error designing plan for request %s",req.RequestID)
        return
    }

    // 4) deploy fragments
    // Tell deployment managers to execute plans
    err = c.DeployPlan(plan)
    if err != nil {
        log.Error().Err(err).Msgf("error deploying plan request %s", req.RequestID)
        return
    }
}


// For a given collection of plans, tell the corresponding deployment managers to run the deployment.
func (c *Manager) DeployPlan(plan *pbConductor.DeploymentPlan) error {
    // Start monitoring this fragment
    c.Monitor.AddPlanToMonitor(plan)

    for fragmentIndex, fragment := range plan.Fragments {
        log.Info().Msgf("start fragment %s deployment with %d out of %d fragments", fragment.DeploymentId, fragmentIndex, len(plan.Fragments))
        // TODO get cluster IP address from system model
        conductor.GetDMClients().AddConnection("127.0.0.1:5200")
        clusterIP := "127.0.0.1:5200"
        conn,err := conductor.GetDMClients().GetConnection(clusterIP)
        if err!=nil{
            log.Error().Err(err).Msgf("problem creating connection with %s",clusterIP)
            // TODO define what to do in this case. Run rollback?
            return err
        }

        // build a request
        request := pbDeploymentManager.DeploymentFragmentRequest{RequestId: uuid.New().String(),Fragment: fragment}
        client := pbDeploymentManager.NewDeploymentManagerClient(conn)
        _, err = client.Execute(context.Background(),&request)

        if err!=nil {
            // TODO define how to proceed in case of error
            log.Error().Err(err).Msgf("problem deploying fragment %s", fragment.DeploymentId)
            return err
        }

        // TODO define how to modify the system model according to the response
    }

    return nil
}








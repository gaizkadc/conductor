/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 */

package handler

import (
    "github.com/phf/go-queue/queue"
    "github.com/rs/zerolog/log"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
    pbDeploymentManager "github.com/nalej/grpc-deployment-manager-go"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "sync"
    "time"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/conductor"
    "context"
    "github.com/nalej/conductor/pkg/conductor/requirementscollector"
)

// Time to wait between checks in the queue in milliseconds.
const CheckSleepTime = 2000

type Manager struct {
    // queue for incoming messages
    queue *queue.Queue
    // ScorerMethod
    ScorerMethod scorer.Scorer
    // Requirements collector
    ReqCollector requirementscollector.RequirementsCollector
    // Plan designer
    Designer plandesigner.PlanDesigner
    // Mutex for queue operations
    mux sync.RWMutex
}

func NewManager(queue *queue.Queue, scorer scorer.Scorer, reqColl requirementscollector.RequirementsCollector,
    designer plandesigner.PlanDesigner, port uint32) *Manager {
    // instantiate a server
    return &Manager{queue: queue, ScorerMethod: scorer, ReqCollector: reqColl, Designer: designer}
}

// Check iteratively if there is anything to be processed in the queue.
func (c *Manager) Run() {
    sleep := time.Tick(time.Millisecond * CheckSleepTime)
    for{
        select {
        case <- sleep:
            for c.AvailableRequests() {
                c.ProcessDeploymentRequest()
            }
        }
    }
}

func(c *Manager) ProcessDeploymentRequest(){
    req := c.NextRequest()
    if req == nil {
        log.Error().Msg("the queue was unexpectedly empty")
        return
    }

    // TODO get all the data from the system model

    // Get the ServiceGroup structure
    // This is hardcoded for testing purposes
    appDescriptor := pbApplication.AppDescriptor{
        Name:"app_descriptor_test",
        Description: "app_descriptor_test description",
        AppDescriptorId: "app_descriptor_id",
        OrganizationId: "organization_test",
        EnvironmentVariables: map[string]string{"var1":"var1_value", "var2":"var2_value"},
        Labels: map[string]string{"label1":"label1_value", "label2":"label2_value"},
    }

    // 1) collect requirements for the application descriptor
    foundRequirements, err := c.ReqCollector.FindRequirements(&appDescriptor)
    if err != nil {
        log.Error().Err(err).Msgf("impossible to find requirements for application %s", appDescriptor.AppDescriptorId)
        return
    }

    // 2) score requirements
    scoreResult, err := c.ScorerMethod.ScoreRequirements (foundRequirements)

    if err != nil {
        log.Error().Err(err).Msgf("error scoring request %s", req.RequestId)
        return
    }

    log.Info().Msgf("conductor maximum score for %s is for cluster %s among %d possible",
        scoreResult.RequestID, scoreResult.ClusterID, scoreResult.TotalEvaluated)


    // 3) design plan
    // TODO elaborate plan, modify system model accordingly
    // Elaborate deployment plan for the application
    plan, err := c.Designer.DesignPlan(&appDescriptor, scoreResult)

    if err != nil{
        log.Error().Err(err).Msgf("error designing plan for request %s",req.RequestId)
        return
    }

    // 4) deploy fragments
    // Tell deployment managers to execute plans
    err = c.DeployPlan(plan)
    if err != nil {
        log.Error().Err(err).Msgf("error deploying plan request %s", req.RequestId)
        return
    }
}


// For a given collection of plans, tell the corresponding deployment managers to run the deployment.
func (c *Manager) DeployPlan(plan *pbConductor.DeploymentPlan) error {
    for fragmentIndex, fragment := range plan.Fragments {
        log.Info().Msgf("start fragment %s deployment with %d out of %d fragments", fragment.DeploymentId, fragmentIndex, len(plan.Fragments))
        // TODO get cluster IP address from system model
        conductor.GetDMClients().AddConnection("127.0.0.1:5002")
        clusterIP := "127.0.0.1:5002"
        conn,err := conductor.GetDMClients().GetConnection(clusterIP)
        if err!=nil{
            log.Error().Err(err).Msgf("problem creating connection with %s",clusterIP)
            // TODO define what to do in this case. Run rollback?
            return err
        }

        // build a request
        request := pbDeploymentManager.DeployFragmentRequest{}
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


// Thread-safe method to access queued requests
func(c *Manager) NextRequest() *pbConductor.DeploymentRequest {
    c.mux.Lock()
    toReturn := c.queue.PopFront().(*pbConductor.DeploymentRequest)
    defer c.mux.Unlock()
    return toReturn
}

// Thread-safe function to find whether there are more requests available or not.
func(c *Manager) AvailableRequests() bool {
    c.mux.RLock()
    available := c.queue.Len()!=0
    defer c.mux.RUnlock()
    return available
}

// Push a new request to the que for later processing.
//  params:
//   req entry to be enqueued
func (c *Manager) PushRequest(req *pbConductor.DeploymentRequest) error {
    c.mux.Lock()
    c.queue.PushBack(req)
    defer c.mux.Unlock()
    return nil
}








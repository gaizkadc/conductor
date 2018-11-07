/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nalej/conductor/internal/entities"
	"github.com/nalej/conductor/pkg/conductor"
	"github.com/nalej/conductor/pkg/conductor/monitor"
	"github.com/nalej/conductor/pkg/conductor/plandesigner"
	"github.com/nalej/conductor/pkg/conductor/requirementscollector"
	"github.com/nalej/conductor/pkg/conductor/scorer"
	"github.com/nalej/conductor/pkg/utils"
	pbApplication "github.com/nalej/grpc-application-go"
	pbConductor "github.com/nalej/grpc-conductor-go"
	pbDeploymentManager "github.com/nalej/grpc-deployment-manager-go"
	"github.com/rs/zerolog/log"
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
	if pool != nil && len(pool.GetConnections()) == 0 {
		log.Panic().Msg("system model clients were not started")
		return nil
	}
	conn := pool.GetConnections()[0]
	// Create associated clients
	appClient := pbApplication.NewApplicationsClient(conn)
	return &Manager{Queue: queue, ScorerMethod: scorer, ReqCollector: reqColl, Designer: designer, appClient: appClient,
		Monitor: monitor}
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
func (c *Manager) PushRequest(req *pbConductor.DeploymentRequest) (*entities.DeploymentRequest, error) {
	log.Debug().Msgf("push request %s", req.RequestId)
	desc, err := c.appClient.GetAppDescriptor(context.Background(), req.AppId)
	if err != nil {
		log.Error().Err(err).Msg("error getting application descriptor")
		return nil, err
	}
	// Create new application instance
	addReq := pbConductor.AddAppInstanceRequest{
		OrganizationId:  desc.OrganizationId,
		AppDescriptorId: desc.AppDescriptorId,
		Name:            req.Name,
		Description:     req.Description,
	}
	// Add instance, by default this is created with queue status
	instance, err := c.appClient.AddAppInstance(context.Background(), &addReq)
	if err != nil {
		log.Error().Err(err).Msg("error adding application instance")
		return nil, err
	}

	toEnqueue := entities.DeploymentRequest{
		RequestId:      req.RequestId,
		InstanceId:     instance.AppInstanceId,
		OrganizationId: req.AppId.OrganizationId,
		ApplicationId:  req.AppId.AppDescriptorId,
	}
	err = c.Queue.PushRequest(&toEnqueue)
	if err != nil {
		return &toEnqueue, err
	}

	if err != nil {
		log.Error().Err(err).Msgf("problems updating application %s", req.AppId)
		return &toEnqueue, err
	}

	return &toEnqueue, nil
}

func (c *Manager) ProcessDeploymentRequest() {
	req := c.Queue.NextRequest()
	if req == nil {
		log.Error().Msg("the queue was unexpectedly empty")
		return
	}

	// TODO get all the data from the system model
	// Get the ServiceGroup structure

	appInstance, err := c.appClient.GetAppInstance(context.Background(),
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
	scoreResult, err := c.ScorerMethod.ScoreRequirements(req.OrganizationId, foundRequirements)

	if err != nil {
		log.Error().Err(err).Msgf("error scoring request %s", req.RequestId)
		return
	}

	log.Info().Msgf("conductor maximum score for %s is for cluster %s among %d possible",
		scoreResult.RequestId, scoreResult.ClusterId, scoreResult.TotalEvaluated)

	// 3) design plan
	// TODO elaborate plan, modify system model accordingly
	// Elaborate deployment plan for the application
	plan, err := c.Designer.DesignPlan(appInstance, scoreResult)

	if err != nil {
		log.Error().Err(err).Msgf("error designing plan for request %s", req.RequestId)
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
func (c *Manager) DeployPlan(plan *entities.DeploymentPlan) error {
	// Start monitoring this fragment
	c.Monitor.AddPlanToMonitor(plan)

	for fragmentIndex, fragment := range plan.Fragments {
		log.Info().Msgf("start fragment %s deployment with %d out of %d fragments", fragment.DeploymentId, fragmentIndex, len(plan.Fragments))

		targetHostname, found := conductor.ClusterReference[fragment.ClusterId]
		if !found {
			msg := fmt.Sprintf("unknown target address for cluster with id %s", fragment.ClusterId)
			err := errors.New(msg)
			log.Error().Msgf(msg)
			return err
		}

		dmAddress := fmt.Sprintf("%s:%d", targetHostname, utils.DEPLOYMENT_MANAGER_PORT)

		conn, err := conductor.GetDMClients().GetConnection(dmAddress)

		if err != nil {
			log.Error().Err(err).Msgf("problem creating connection with %s", dmAddress)
			// TODO define what to do in this case.
			return err
		}

		// build a request
		request := pbDeploymentManager.DeploymentFragmentRequest{
			RequestId:      uuid.New().String(),
			Fragment:       fragment.ToGRPC(),
			RollbackPolicy: pbDeploymentManager.RollbackPolicy_NONE}
		client := pbDeploymentManager.NewDeploymentManagerClient(conn)
		_, err = client.Execute(context.Background(), &request)

		if err != nil {
			// TODO define how to proceed in case of error
			log.Error().Err(err).Msgf("problem deploying fragment %s", fragment.DeploymentId)
			return err
		}
	}

	return nil
}

/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

// Business logic for the conductor monitor service.

package monitor

import (
    "github.com/nalej/conductor/internal/structures"
    "github.com/nalej/conductor/pkg/conductor/baton"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
    "github.com/rs/zerolog/log"
    "errors"
    "fmt"
    "github.com/nalej/conductor/internal/entities"
    "context"
    "github.com/nalej/conductor/pkg/utils"
    "time"
)

type Manager struct {
    // structure controlling pending plans to be monitored
    pendingPlans *structures.PendingPlans
    ConnHelper *utils.ConnectionsHelper
    AppClient pbApplication.ApplicationsClient
    // Queue of deployment requests
    queue structures.RequestsQueue
    // Access to baton to operate with deployments
    manager *baton.Manager
}

func NewManager(connHelper *utils.ConnectionsHelper, queue structures.RequestsQueue, pendingPlans *structures.PendingPlans,
    manager *baton.Manager) *Manager {
    // initialize clients
    pool := connHelper.GetSystemModelClients()
    if pool!=nil && len(pool.GetConnections())==0{
        log.Panic().Msg("system model clients were not started")
        return nil
    }
    conn := pool.GetConnections()[0]
    appClient := pbApplication.NewApplicationsClient(conn)
    return &Manager{ConnHelper: connHelper, AppClient: appClient, pendingPlans: pendingPlans,
        queue: queue, manager: manager}
}

// Add a plan to be monitored.
func (m *Manager) AddPlanToMonitor(plan *entities.DeploymentPlan) {
    m.pendingPlans.AddPendingPlan(plan)
}

func(m *Manager) UpdateFragmentStatus(request *pbConductor.DeploymentFragmentUpdateRequest) error {
    log.Debug().Interface("request", request).Str("status", request.Status.String()).Msg("monitor received fragment update status")

    // Check if we are monitoring the fragment
    found := m.pendingPlans.MonitoredFragment(request.FragmentId)
    if !found {
        err := errors.New(fmt.Sprintf("fragment %s is not monitored", request.FragmentId))
        return err
    }

    var newStatus *pbApplication.UpdateAppStatusRequest
    newStatus = nil

    //failedDeployment := false

    if entities.DeploymentStatusToGRPC[request.Status] == entities.FRAGMENT_DONE {
        log.Info().Str("fragmentId", request.FragmentId).Msgf("deployment fragment was done")
        //m.pendingPlans.RemoveFragment(request.FragmentId)
        newStatus = &pbApplication.UpdateAppStatusRequest{
            OrganizationId: request.OrganizationId,
            AppInstanceId: request.AppInstanceId,
            Status: pbApplication.ApplicationStatus_RUNNING,
        }
    }

    if entities.DeploymentStatusToGRPC[request.Status] == entities.FRAGMENT_DEPLOYING {
        log.Info().Str("fragmentId", request.FragmentId).Msgf("deployment fragment is being deployed")
        newStatus = &pbApplication.UpdateAppStatusRequest{
            OrganizationId: request.OrganizationId,
            AppInstanceId: request.AppInstanceId,
            Status: pbApplication.ApplicationStatus_DEPLOYING,
        }
    }

    if entities.DeploymentStatusToGRPC[request.Status] == entities.FRAGMENT_ERROR {
        log.Info().Str("deploymentId", request.DeploymentId).Msg("deployment fragment failed")
        newStatus = m.processFailedFragment(request)
        //failedDeployment = true
    }


    // If no more fragments are pending... we stop monitoring the deployment plan
    /*
    if !failedDeployment && !m.pendingPlans.PlanHasPendingFragments(request.DeploymentId) {
        log.Info().Str("deploymentId", request.DeploymentId).Msg("deployment plan was done")
        // time to delete this plan
        // m.pendingPlans.RemovePendingPlan(request.DeploymentId)
        // update the application status in the system model
        newStatus = &pbApplication.UpdateAppStatusRequest{
            OrganizationId: request.OrganizationId,
            AppInstanceId: request.AppInstanceId,
            Status: pbApplication.ApplicationStatus_RUNNING,
        }
    }
    */

    if newStatus != nil {
        _, err := m.AppClient.UpdateAppStatus(context.Background(), newStatus)
        if err != nil {
            log.Error().Err(err).Interface("request", newStatus).Msg("impossible to update app status")
            return err
        }
        log.Debug().Str("instanceId", request.AppInstanceId).Str("status", newStatus.Status.String()).Msg("set instance to new status")
    }

    log.Debug().Interface("request", request).Msg("finished processing update fragment")

    return nil
}

func(m *Manager) UpdateServicesStatus(request *pbConductor.DeploymentServiceUpdateRequest) error {
    log.Debug().Interface("request", request).Msg("monitor received deployment service update")
        for _, update := range request.List {
        updateService := pbApplication.UpdateServiceStatusRequest{
            OrganizationId: update.OrganizationId,
            ServiceId: update.ServiceInstanceId,
            AppInstanceId: update.ApplicationInstanceId,
            Status: update.Status,
            DeployedOnClusterId: request.ClusterId,
            Endpoints: update.Endpoints,
        }
        _, err := m.AppClient.UpdateServiceStatus(context.Background(), &updateService)
        if err != nil {
            log.Error().Err(err).Interface("request", updateService).Msg("impossible to update service status")
            return err
        }
    }

    return nil
}

// Execute the corresponding operations for a failed fragment.
//  params:
//   deploymentId deployment id
//   appInstanceId application id
//  return:
//   update request status
func(m *Manager) processFailedFragment(request *pbConductor.DeploymentFragmentUpdateRequest) *pbApplication.UpdateAppStatusRequest {
    // get deployment request associated with this plan
    plan, isThere := m.pendingPlans.Pending[request.DeploymentId]
    if !isThere {
        log.Error().Str("deploymentId",request.DeploymentId).Msg("no pending plan for the deployment id")
        return nil
    }

    toReturn := &pbApplication.UpdateAppStatusRequest{
        OrganizationId: request.OrganizationId,
        AppInstanceId: request.AppInstanceId,
    }
    // How many times have we tried to deploy this?
    if plan.DeploymentRequest.NumRetries < baton.ConductorMaxDeploymentRetries -1{
        // there is room for one more attempt
        plan.DeploymentRequest.NumRetries = plan.DeploymentRequest.NumRetries + 1
        t := time.Now()
        plan.DeploymentRequest.TimeRetry = &t
        log.Info().Interface("fragmentUpdate",request).Int32("numRetries",plan.DeploymentRequest.NumRetries).
            Msg("fragment deployment failed. Enqueue deployment for another retry")
        // Push this into the queue
        if err:=m.queue.PushRequest(plan.DeploymentRequest); err!= nil {
            log.Error().Err(err).Interface("deploymentRequest",plan.DeploymentRequest).Msg("impossible to" +
                "enqueue the deployment request")
            toReturn.Status = pbApplication.ApplicationStatus_ERROR
        } else {
            toReturn.Status = pbApplication.ApplicationStatus_QUEUED
        }
    } else {
        // no more retries for this request
        log.Info().Interface("fragmentUpdate",request).Int32("numRetries",plan.DeploymentRequest.NumRetries).
            Msg("exceeded number of retries")
        toReturn.Status = pbApplication.ApplicationStatus_ERROR
    }

    // Undeploy the application
    undeployRequest := &entities.UndeployRequest{AppInstanceId: request.AppInstanceId,
        OrganizationId: request.OrganizationId}
    err := m.manager.Undeploy(undeployRequest)
    if err != nil {
        log.Error().Err(err).Interface("fragmentUpdate",request).Msg("error undeploying failed application")
    }

    // Remove references to this pending plan
    m.pendingPlans.RemovePendingPlan(request.DeploymentId)

    return toReturn
}



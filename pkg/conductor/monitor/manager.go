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


    if entities.DeploymentStatusToGRPC[request.Status] == entities.FRAGMENT_DONE {
        log.Info().Str("fragmentId", request.FragmentId).Msgf("deployment fragment was done")
        newStatus = &pbApplication.UpdateAppStatusRequest{
            OrganizationId: request.OrganizationId,
            AppInstanceId: request.AppInstanceId,
            Status: pbApplication.ApplicationStatus_RUNNING,
        }
        // This fragment is no longer pending
        m.pendingPlans.SetFragmentNoPending(request.FragmentId)
    }

    if entities.DeploymentStatusToGRPC[request.Status] == entities.FRAGMENT_DEPLOYING {
        log.Info().Str("fragmentId", request.FragmentId).Msgf("deployment fragment is being deployed")
        newStatus = &pbApplication.UpdateAppStatusRequest{
            OrganizationId: request.OrganizationId,
            AppInstanceId: request.AppInstanceId,
            Status: pbApplication.ApplicationStatus_DEPLOYING,
            Info: request.Info,
        }
        // This fragment is pending
        m.pendingPlans.SetFragmentPending(request.FragmentId)
    }

    if entities.DeploymentStatusToGRPC[request.Status] == entities.FRAGMENT_ERROR {
        log.Info().Str("deploymentId", request.DeploymentId).Msg("deployment fragment failed")
        newStatus = &pbApplication.UpdateAppStatusRequest{
            OrganizationId: request.OrganizationId,
            AppInstanceId: request.AppInstanceId,
            Status: pbApplication.ApplicationStatus_DEPLOYMENT_ERROR,
            Info: request.Info,
        }
        // This fragment is pending
        newStatus = m.processFailedFragment(request)
    }


    if newStatus != nil {
        log.Debug().Msg("update app status")
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

    // Create a map of service group metadata to avoid duplicated queries
    groupMetadata := make(map[string]*pbApplication.InstanceMetadata,0)

    for _, update := range request.List {
        updateService := pbApplication.UpdateServiceStatusRequest{
            OrganizationId: update.OrganizationId,
            AppInstanceId: update.ApplicationInstanceId,
            Status: update.Status,
            DeployedOnClusterId: request.ClusterId,
            ServiceGroupInstanceId: update.ServiceGroupInstanceId,
            ServiceInstanceId: update.ServiceInstanceId,
            Endpoints: update.Endpoints,
        }
        _, err := m.AppClient.UpdateServiceStatus(context.Background(), &updateService)
        if err != nil {
            log.Error().Err(err).Interface("request", updateService).Msg("impossible to update service status")
            return err
        }

        // TODO: Improve the update. Reduce the number of calls.
        // udpate application endpoints
        for _, endpoint := range update.Endpoints{
            _, err := m.AppClient.AddAppEndpoint(context.Background(), &pbApplication.AppEndpoint{
                OrganizationId:         update.OrganizationId,
                AppInstanceId:          update.ApplicationInstanceId,
                ServiceGroupInstanceId: update.ServiceGroupInstanceId,
                ServiceInstanceId:      update.ServiceInstanceId,
                //Port:
                //Protocol:
                EndpointInstance: endpoint,
            })
            if err != nil {
                log.Error().Err(err).Interface("endpoint", endpoint).Msg("impossible to add application endpoint")
                return err
            }
        }

        // Update the corresponding service group instance
        // Get the service metadata just in case we don't queried it yet
        meta, found := groupMetadata[update.ServiceGroupInstanceId]
        if !found {
            meta, err = m.AppClient.GetServiceGroupInstanceMetadata(context.Background(), &pbApplication.GetServiceGroupInstanceMetadataRequest{
                AppInstanceId: update.ApplicationInstanceId,
                OrganizationId: update.OrganizationId,
                ServiceGroupInstanceId: update.ServiceGroupInstanceId,
            })
            if err != nil {
                log.Error().Err(err).Interface("request", updateService).Msg("service group instance metadata not found")
                return err
            }
            groupMetadata[update.ServiceGroupInstanceId] = meta
        }

        // update the status of the service
        meta.Status[updateService.ServiceInstanceId] = updateService.Status
    }

    for _,v := range groupMetadata {
        log.Debug().Str("serviceGroupInstanceId", v.MonitoredInstanceId).Msg("update service group instance metadata")
        _, err := m.AppClient.UpdateServiceGroupInstanceMetadata(context.Background(),v)
        if err != nil {
            log.Error().Err(err).Msg("impossible to update serviceGroupInstanceId metadata")
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
        // Fill it below
        // Info:
        // Status:
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
            toReturn.Info = "impossible to queue application after failed deployment"

        } else {
            toReturn.Status = pbApplication.ApplicationStatus_QUEUED
            toReturn.Info = "app queued after failed deployment"
        }
    } else {
        // no more retries for this request
        log.Info().Interface("fragmentUpdate",request).Int32("numRetries",plan.DeploymentRequest.NumRetries).
            Msg("exceeded number of retries")
        toReturn.Status = pbApplication.ApplicationStatus_ERROR
        toReturn.Info = "exceeded number of retries"
    }

    // rollback
    // TODO check how to proceed with remaining zt networks
    m.manager.Rollback(request.OrganizationId, request.AppInstanceId, "")

    // Undeploy the application
    undeployRequest := &entities.UndeployRequest{AppInstanceId: request.AppInstanceId,
        OrganizationId: request.OrganizationId}
    err := m.manager.Undeploy(undeployRequest)
    if err != nil {
        log.Error().Err(err).Interface("fragmentUpdate",request).Msg("error undeploying failed application")
    }

    return toReturn
}



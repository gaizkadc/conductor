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

    if entities.DeploymentStatusToGRPC[request.Status] == entities.FRAGMENT_DONE {
        log.Info().Str("fragmentId", request.FragmentId).Msgf("deployment fragment was done")
        // This fragment is no longer pending
        m.pendingPlans.SetFragmentNoPending(request.FragmentId)
    }

    if entities.DeploymentStatusToGRPC[request.Status] == entities.FRAGMENT_DEPLOYING {
        log.Info().Str("fragmentId", request.FragmentId).Msgf("deployment fragment is being deployed")
        // This fragment is pending
        m.pendingPlans.SetFragmentPending(request.FragmentId)
    }

    if entities.DeploymentStatusToGRPC[request.Status] == entities.FRAGMENT_ERROR {
        log.Info().Str("deploymentId", request.DeploymentId).Msg("deployment fragment failed")
        // This fragment is pending
        m.processFailedFragment(request)
    }

    log.Debug().Interface("request", request).Msg("finished processing update fragment")

    return nil
}


func(m *Manager) UpdateServicesStatus(request *pbConductor.DeploymentServiceUpdateRequest) error {

    log.Debug().Interface("request", request).Msg("monitor received deployment service update")

    instancesToUpdate := make(map[string]*pbApplication.AppInstance, 0)

    // collect the required application instances
    for _, update := range request.List {
        _, found := instancesToUpdate[update.ApplicationInstanceId]
        if !found {
            retrievedInstance, err := m.AppClient.GetAppInstance(context.Background(),
                &pbApplication.AppInstanceId{OrganizationId: update.OrganizationId, AppInstanceId: update.ApplicationInstanceId})
            if err != nil {
                log.Error().Err(err).Msg("impossible to retrieve application instance to update service status")
            } else {
                instancesToUpdate[retrievedInstance.AppInstanceId] = retrievedInstance
            }
        }
    }

    // update applications using current status info
    for _, update := range request.List {
        instance, found := instancesToUpdate[update.ApplicationInstanceId]
        if found {
            // update the status of this application instance according to this service
            m.updateAppInstanceServiceStatus(instance, update)
        }
    }

    // update application instances in the system model
    for appInstanceId, appInstance := range instancesToUpdate {
        _, err := m.AppClient.UpdateAppInstance(context.Background(), appInstance)
        if err != nil {
            log.Error().Err(err).Str("appInstanceId", appInstanceId).Msg("impossible to update application instance")
        }
    }

    return nil

}

func (m *Manager) updateAppInstanceServiceStatus(instance *pbApplication.AppInstance, update *pbConductor.ServiceUpdate) {
    // Update status
    var targetService *pbApplication.ServiceInstance
    var targetGroup *pbApplication.ServiceGroupInstance
    for _, g := range instance.Groups {
        if g.ServiceGroupInstanceId == update.ServiceGroupInstanceId {
            targetGroup = g
            for _, serv := range g.ServiceInstances {
                if serv.ServiceInstanceId == update.ServiceInstanceId {
                    targetService = serv
                }
            }
        }
    }

    if targetService == nil {
        log.Error().Str("serviceInstanceId", update.ServiceInstanceId).
            Str("applicationInstanceId", update.ApplicationInstanceId).
            Msg("impossible to find service in application instance to be updated")
        return
    }
    // update service status
    if targetService.Status != update.Status {
        log.Debug().Str("serviceInstance", targetService.ServiceInstanceId).
            Msg(fmt.Sprintf("update service instance status from %s ---> %s", targetService.Status, update.Status))
        targetService.Status = update.Status
    }
    targetService.Info = update.Info

    // Update metadata entry in the group
    groupFinalStatus := pbApplication.ServiceStatus_SERVICE_RUNNING
    targetGroup.Metadata.Status[update.ServiceInstanceId] = targetService.Status
    targetGroup.Metadata.Info[update.ServiceInstanceId] = targetService.Info
    // decide the final status for this group
    for _, status := range targetGroup.Metadata.Status {
        if status == pbApplication.ServiceStatus_SERVICE_ERROR {
            groupFinalStatus = pbApplication.ServiceStatus_SERVICE_ERROR
            break
        }
        if groupFinalStatus > status {
            groupFinalStatus = status
        }
    }
    if targetGroup.Status != groupFinalStatus {
        log.Debug().Str("groupInstance", targetGroup.ServiceGroupInstanceId).
            Msg(fmt.Sprintf("update service group status from %s ---> %s", targetGroup.Status, groupFinalStatus))
        targetGroup.Status = groupFinalStatus
    }



    // decide the final status for this instance
    groupsSummary := pbApplication.ServiceStatus_SERVICE_RUNNING
    for _, g := range instance.Groups {
        if g.Status == pbApplication.ServiceStatus_SERVICE_ERROR {
            groupsSummary = pbApplication.ServiceStatus_SERVICE_ERROR
            break
        }
        if g.Status < groupsSummary {
            groupsSummary = g.Status
        }
    }

    var finalAppStatus pbApplication.ApplicationStatus
    switch groupsSummary {
    case pbApplication.ServiceStatus_SERVICE_SCHEDULED:
        finalAppStatus = pbApplication.ApplicationStatus_SCHEDULED
        break
    case pbApplication.ServiceStatus_SERVICE_WAITING:
        finalAppStatus = pbApplication.ApplicationStatus_PLANNING
        break
    case pbApplication.ServiceStatus_SERVICE_DEPLOYING:
        finalAppStatus = pbApplication.ApplicationStatus_DEPLOYING
        break
    case pbApplication.ServiceStatus_SERVICE_RUNNING:
        finalAppStatus = pbApplication.ApplicationStatus_RUNNING
        break
    case pbApplication.ServiceStatus_SERVICE_ERROR:
        finalAppStatus = pbApplication.ApplicationStatus_ERROR
        break
    }

    if instance.Status != finalAppStatus {
        log.Debug().Str("appInstanceId", instance.AppInstanceId).
            Msg(fmt.Sprintf("update app instance status from %s ---> %s", instance.Status, finalAppStatus))
        instance.Status = finalAppStatus
    }

    // Update endpoints
    if update.Endpoints != nil && len(update.Endpoints) > 0 {
        targetService.Endpoints = update.Endpoints
    }

    // Update metadata entry for the service groups with the service service group id
    availableReplicas := int32(0)
    unavailableReplicas := int32(0)
    for _, g := range instance.Groups {
        if g.ServiceGroupId == update.ServiceGroupId {
            if g.Status == pbApplication.ServiceStatus_SERVICE_RUNNING {
                availableReplicas++
            } else {
                unavailableReplicas++
            }
        }
    }

    // Update groups
    for _, g := range instance.Groups {
        if g.ServiceGroupId == update.ServiceGroupId {
            g.Metadata.AvailableReplicas = availableReplicas
            g.Metadata.UnavailableReplicas = unavailableReplicas
        }
    }
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

    // Undeploy the application
    err := m.manager.SoftUndeploy(request.OrganizationId, request.AppInstanceId)
    if err != nil {
        log.Error().Err(err).Interface("fragmentUpdate",request).Msg("error undeploying failed application")
    }

    return toReturn
}



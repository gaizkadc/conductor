/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

// Business logic for the conductor monitor service.

package monitor

import (
    "github.com/nalej/conductor/internal/structures"
    "github.com/nalej/conductor/pkg/conductor/baton"
    "github.com/nalej/derrors"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
    "github.com/nalej/grpc-utils/pkg/conversions"
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

    var finalStatus entities.DeploymentFragmentStatus

    switch entities.DeploymentStatusToGRPC[request.Status] {
    case entities.FRAGMENT_DONE:
        log.Info().Str("fragmentId", request.FragmentId).Msgf("deployment fragment was done")
        // This fragment is no longer pending
        m.pendingPlans.SetFragmentNoPending(request.FragmentId)
        finalStatus = entities.FRAGMENT_DONE

    case entities.FRAGMENT_DEPLOYING:
        log.Info().Str("fragmentId", request.FragmentId).Msgf("deployment fragment is being deployed")
        // This fragment is pending
        m.pendingPlans.SetFragmentPending(request.FragmentId)
        finalStatus = entities.FRAGMENT_DEPLOYING

    case entities.FRAGMENT_ERROR:
        finalStatus = entities.FRAGMENT_ERROR
        log.Info().Str("deploymentId", request.DeploymentId).Msg("deployment fragment failed")
        // This fragment is pending
        newStatus := m.processFailedFragment(request)
        _, err := m.AppClient.UpdateAppStatus(context.Background(), newStatus)
        if err != nil {
            log.Error().Err(err).Msg("problem found when update app status after failed fragment")
        }
        finalStatus = entities.FRAGMENT_DEPLOYING

    case entities.FRAGMENT_TERMINATING:
        log.Info().Str("deploymentId", request.FragmentId).Msg("deployment fragment terminating")
        finalStatus = entities.FRAGMENT_TERMINATING

    case entities.FRAGMENT_RETRYING:
        log.Info().Str("deploymentId", request.DeploymentId).Msg("deployment fragment retrying")
        m.pendingPlans.SetFragmentPending(request.FragmentId)
        finalStatus = entities.FRAGMENT_RETRYING

    case entities.FRAGMENT_WAITING:
        log.Info().Str("deploymentId", request.DeploymentId).Msg("deployment fragment waiting")
        m.pendingPlans.SetFragmentPending(request.FragmentId)
        finalStatus = entities.FRAGMENT_WAITING

    default:
        log.Info().Msg("received a non processable status in an update fragment update")
        return nil
    }

    log.Debug().Interface("finalStatus",finalStatus).Msg("update deployment fragment status")

    // Update the view of this deployment fragment in the DB
    df, err := m.manager.AppClusterDB.GetDeploymentFragment(request.ClusterId, request.FragmentId)
    if err != nil {
        e := derrors.NewInternalError("impossible to get deployment fragment status from database", err)
        return e
    }
    df.Status = finalStatus
    err = m.manager.AppClusterDB.AddDeploymentFragment(df)
    if err != nil {
        e := derrors.NewInternalError("impossible to update deployment fragment status in database", err)
        return e
    }

    log.Debug().Interface("request", request).Msg("finished processing update fragment")

    return nil
}


func(m *Manager) UpdateServicesStatus(request *pbConductor.DeploymentServiceUpdateRequest) error {

    log.Debug().Interface("updateRequest", request).Msg("monitor received deployment service update")

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
            // update appEndpoint in the system model
            m.updateAppEndpointInstances(update)
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

// updateAppEndpointInstances calls system-model to update global_fqdn
func (m *Manager) updateAppEndpointInstances(update *pbConductor.ServiceUpdate) {
    if update.Endpoints != nil && len(update.Endpoints) > 0 {
        for _, endpoint := range update.Endpoints {
            _, err := m.AppClient.AddAppEndpoint(context.Background(), &pbApplication.AddAppEndpointRequest{
                OrganizationId: update.OrganizationId,
                AppInstanceId: update.ApplicationInstanceId,
                ServiceGroupInstanceId: update.ServiceGroupInstanceId,
                ServiceInstanceId: update.ServiceInstanceId,
                ServiceName: update.ServiceName,
                EndpointInstance: endpoint,
            } )
            if err != nil {
                log.Error().Str("error", conversions.ToDerror(err).DebugReport()).
                    Str("serviceInstanceId", update.ServiceInstanceId).
                    Str("endpointInstanceId", endpoint.EndpointInstanceId).
                    Msg("impossible to update appEndpointInstance")
            }
        }

    }
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
    scheduledEntries := 0
    // decide the final status for this group
    for _, status := range targetGroup.Metadata.Status {
        if status == pbApplication.ServiceStatus_SERVICE_ERROR {
            groupFinalStatus = pbApplication.ServiceStatus_SERVICE_ERROR
            break
        }
        if status == pbApplication.ServiceStatus_SERVICE_TERMINATING {
            groupFinalStatus = pbApplication.ServiceStatus_SERVICE_TERMINATING
            break
        }
        if status != pbApplication.ServiceStatus_SERVICE_SCHEDULED {
            if groupFinalStatus > status {
                groupFinalStatus = status
            }
        } else {
            scheduledEntries++
        }
    }

    // if all the services are scheduled
    if scheduledEntries == len(targetGroup.Metadata.Status) {
        groupFinalStatus = pbApplication.ServiceStatus_SERVICE_SCHEDULED
    }

    if targetGroup.Status != groupFinalStatus {
        log.Debug().Str("groupInstance", targetGroup.ServiceGroupInstanceId).
            Msg(fmt.Sprintf("update service group status from %s ---> %s", targetGroup.Status, groupFinalStatus))
        targetGroup.Status = groupFinalStatus
    }

    // decide the final status for this instance
    groupsSummary := pbApplication.ServiceStatus_SERVICE_RUNNING
    scheduledEntries = 0
    for _, g := range instance.Groups {
        if g.Status == pbApplication.ServiceStatus_SERVICE_ERROR {
            groupsSummary = pbApplication.ServiceStatus_SERVICE_ERROR
            break
        }
        if g.Status == pbApplication.ServiceStatus_SERVICE_TERMINATING {
            groupsSummary = pbApplication.ServiceStatus_SERVICE_TERMINATING
            break
        }
        if g.Status != pbApplication.ServiceStatus_SERVICE_SCHEDULED {
            if g.Status < groupsSummary {
                groupsSummary = g.Status
            }
        } else {
            scheduledEntries++
        }
    }

    if scheduledEntries == len(instance.Groups) {
        groupsSummary = pbApplication.ServiceStatus_SERVICE_SCHEDULED
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
    case pbApplication.ServiceStatus_SERVICE_TERMINATING:
        finalAppStatus = pbApplication.ApplicationStatus_TERMINATING
        break
    default:
        log.Error().Interface("finalAppStatus",groupsSummary).Msg("unkown service status to be set")
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
            toReturn.Info = "impossible to queue application after failed deployment."
            if request.Info != "" {
                toReturn.Info = toReturn.Info + " ["+ request.Info + "]"
            }

        } else {
            toReturn.Status = pbApplication.ApplicationStatus_QUEUED
            toReturn.Info = "app queued after failed deployment"
            if request.Info != "" {
                toReturn.Info = toReturn.Info + " ["+ request.Info + "]"
            }
        }
    } else {
        // no more retries for this request
        log.Info().Interface("fragmentUpdate",request).Int32("numRetries",plan.DeploymentRequest.NumRetries).
            Msg("exceeded number of retries")
        toReturn.Status = pbApplication.ApplicationStatus_ERROR
        toReturn.Info = "exceeded number of retries"
        if request.Info != "" {
            toReturn.Info = toReturn.Info + " ["+ request.Info + "]"
        }
    }

    // Undeploy the application
    err := m.manager.SoftUndeploy(request.OrganizationId, request.AppInstanceId)
    if err != nil {
        log.Error().Err(err).Interface("fragmentUpdate",request).Msg("error undeploying failed application")
    }

    return toReturn
}



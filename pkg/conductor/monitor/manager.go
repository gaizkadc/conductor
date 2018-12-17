/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

// Business logic for the conductor monitor service.

package monitor

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
    "github.com/rs/zerolog/log"
    "errors"
    "fmt"
    "github.com/nalej/conductor/internal/entities"
    "context"
    "github.com/nalej/conductor/pkg/utils"
)

type Manager struct {
    pendingPlans *PendingPlans
    ConnHelper *utils.ConnectionsHelper
    AppClient pbApplication.ApplicationsClient
}

func NewManager(connHelper *utils.ConnectionsHelper) *Manager {
    // initialize clients
    pool := connHelper.GetSystemModelClients()
    if pool!=nil && len(pool.GetConnections())==0{
        log.Panic().Msg("system model clients were not started")
        return nil
    }
    conn := pool.GetConnections()[0]
    appClient := pbApplication.NewApplicationsClient(conn)
    return &Manager{ConnHelper: connHelper, AppClient: appClient,pendingPlans: NewPendingPlans()}
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
        log.Info().Msgf("deployment fragment %s was done",request.FragmentId)
        m.pendingPlans.RemoveFragment(request.FragmentId)
    }

    // If no more fragments are pending... we stop monitoring the deployment plan
    if !m.pendingPlans.PlanHasPendingFragments(request.DeploymentId) {
        log.Info().Str("deploymentId", request.DeploymentId).Msg("deployment plan was done")
        // time to delete this plan
        m.pendingPlans.RemovePendingPlan(request.DeploymentId)
        // update the application status in the system model
        req := pbApplication.UpdateAppStatusRequest{
            OrganizationId: request.OrganizationId,
            AppInstanceId: request.AppInstanceId,
            Status: pbApplication.ApplicationStatus_RUNNING,
        }
        log.Info().Str("instanceId", request.AppInstanceId).Msg("set instance to running")
        _, err := m.AppClient.UpdateAppStatus(context.Background(), &req)
        if err != nil {
            log.Error().Err(err).Interface("request", req).Msg("impossible to update app status")
            return err
        }
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



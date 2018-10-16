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
    "github.com/nalej/conductor/pkg/conductor"
    "errors"
    "fmt"
    "github.com/nalej/conductor/internal/entities"
    "context"
)

type Manager struct {
    pendingPlans *PendingPlans
    AppClient pbApplication.ApplicationsClient
}

func NewManager() *Manager {
    // initialize clients
    pool := conductor.GetSystemModelClients()
    if pool!=nil && len(pool.GetConnections())==0{
        log.Panic().Msg("system model clients were not started")
        return nil
    }
    conn := pool.GetConnections()[0]
    appClient := pbApplication.NewApplicationsClient(conn)
    return &Manager{AppClient: appClient,pendingPlans: NewPendingPlans()}
}

// Add a plan to be monitored.
func (m *Manager) AddPlanToMonitor(plan *pbConductor.DeploymentPlan) {
    m.pendingPlans.AddPendingPlan(plan)
}

func(m *Manager) UpdateFragmentStatus(request *pbConductor.DeploymentFragmentUpdateRequest) error {
    log.Debug().Msgf("monitor received fragment update %v", request)

    // Check if we are monitoring the fragment
    found := m.pendingPlans.MonitoredFragment(request.FragmentId)
    if !found {
        err := errors.New(fmt.Sprintf("fragment %s is not monitored", request.FragmentId))
        return err
    }

    if entities.DeploymentStatusToGRPC[request.Status] == entities.FRAGMENT_DONE {
        log.Info().Msgf("Deployment fragment %s was done",request.FragmentId)
        m.pendingPlans.RemoveFragment(request.FragmentId)
    }

    // If no more fragments are pending... we stop monitoring the deployment plan
    if !m.pendingPlans.PlanHasPendingFragments(request.DeploymentId) {
        log.Info().Msgf("deployment plan %s was done", request.DeploymentId)
        // time to delete this plan
        m.pendingPlans.RemovePendingPlan(request.DeploymentId)
        // update the application status in the system model
        req := pbApplication.UpdateAppStatusRequest{
            OrganizationId: request.OrganizationId,
            AppInstanceId: request.AppInstanceId,
            Status: pbApplication.ApplicationStatus_RUNNING,
        }
        m.AppClient.UpdateAppStatus(context.Background(), &req)
    }


    /*
    // Update services status
    for _, status := range request.ServicesStatus {
        log.Debug().Msgf("service %s is known to be in %s", status.InstanceId, status.Status)
    }
    */
    return nil
}

func(m *Manager) UpdateServicesStatus(request *pbConductor.DeploymentServiceUpdateRequest) error {
    log.Debug().Msgf("monitor received update service status %v", request)
        for _, update := range request.List {
        updateService := pbApplication.UpdateServiceStatusRequest{
            OrganizationId: update.OrganizationId,
            ServiceId: update.ServiceInstanceId,
            AppInstanceId: update.ApplicationInstanceId,
            Status: update.Status,
        }
        _, err := m.AppClient.UpdateServiceStatus(context.Background(), &updateService)
        if err != nil {
            log.Error().Err(err).Msgf("impossible to update service status [%v]", updateService)
            return err
        }
    }

    return nil
}



/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

// Business logic for the conductor monitor service.

package monitor

import (
    //pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
    "github.com/rs/zerolog/log"
    "github.com/nalej/conductor/pkg/conductor"
)

type Manager struct {
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
    return &Manager{AppClient: appClient}

}
/*

func(m *Manager) UpdateFragmentStatus(request *pbConductor.DeploymentFragmentUpdateRequest) error {
    log.Debug().Msgf("received fragment update status for fragment %s",request.FragmentId)

    reqStatus := pbApplication.UpdateServiceStatusRequest{
        ServiceId: request.
    }
    m.AppClient.UpdateServiceStatus(context.Background())

    return nil
}
*/
/*
func (m *Manager) AddPendingPlan(plan *pbConductor.DeploymentPlan) {
    log.Debug().Msgf("add plan of deployment %s to pending checks",plan.DeploymentId)
}

func (m *Manager) AddPendingFragment(fragmentID string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.pendingFragments[fragmentID] = entities.FRAGMENT_WAITING
}

func (m *Manager) IsPendingFragment(fragmentID string) bool{
    m.mu.Lock()
    defer m.mu.Unlock()
    _, exists := m.pendingFragments[fragmentID]
    return exists
}

func (m *Manager) SetFragmentStatus(fragmentID string, status entities.DeploymentFragmentStatus) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.pendingFragments[fragmentID] = status
}
*/
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
    "sync"
    "fmt"
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
    // Check if we are monitoring the fragment
    found := m.pendingPlans.MonitoredFragment(request.FragmentId)
    if !found {
        err := errors.New(fmt.Sprintf("fragment %s is not monitored", request.FragmentId))
        return err
    }

    // Update services status
    for _, status := range request.ServicesStatus {
        log.Debug().Msgf("service %s is known to be in %s", status.InstanceId, status.Status)
    }

    return nil
}


// Struct to control pending deployment plans
type PendingPlans struct {
    // plan_id -> deployment plan
    pending map[string]*pbConductor.DeploymentPlan
    // fragment_id -> deployment_plan_id. Index just to improve searching.
    pendingFragment map[string]string
    // service_id -> fragment_id
    pendingService map[string] string
    // mutex
    mu sync.Mutex
}

func NewPendingPlans () *PendingPlans {
    return &PendingPlans{
        pendingService: make(map[string]string,0),
        pendingFragment: make(map[string]string,0),
        pending: make(map[string]*pbConductor.DeploymentPlan),
    }
}


func (p *PendingPlans) AddPendingPlan(plan *pbConductor.DeploymentPlan) {
    log.Debug().Msgf("add plan of deployment %s to pending checks",plan.DeploymentId)
    p.mu.Lock()
    defer p.mu.Unlock()
    p.pending[plan.DeploymentId] = plan
    for _, frag := range plan.Fragments {
        p.pendingFragment[frag.FragmentId] = plan.DeploymentId
        for _, stage := range frag.Stages {
            for _, serv := range stage.Services {
                p.pendingService[serv.ServiceId] = frag.FragmentId
            }
        }
    }

}

func (p *PendingPlans) RemovePendingPlan(plan *pbConductor.DeploymentPlan) {
    log.Debug().Msgf("remove plan of deployment %s from pending checks",plan.DeploymentId)
    p.mu.Lock()
    defer p.mu.Unlock()
    // remove fragments
    for _, frag := range plan.Fragments {
        delete(p.pendingFragment, frag.FragmentId)
        for _, stage := range frag.Stages {
            for _, serv := range stage.Services {
                delete(p.pendingService, serv.ServiceId)
            }
        }
    }
    delete(p.pending, plan.DeploymentId)
}

func (p *PendingPlans) MonitoredFragment(fragmentID string) bool {
    p.mu.Lock()
    defer p.mu.Unlock()
    _, exists := p.pendingFragment[fragmentID]
    return exists
}
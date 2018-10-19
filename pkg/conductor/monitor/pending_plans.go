/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package monitor

import (
    "github.com/nalej/conductor/internal/entities"
    "sync"
    "github.com/rs/zerolog/log"
)

// Struct to control pending deployment plans
type PendingPlans struct {
    // plan_id -> deployment plan
    pending map[string]*entities.DeploymentPlan
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
        pending: make(map[string]*entities.DeploymentPlan),
    }
}


func (p *PendingPlans) AddPendingPlan(plan *entities.DeploymentPlan) {
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
    p.printStatus()
}



func (p *PendingPlans) RemovePendingPlan(deploymentId string) {
    log.Debug().Msgf("remove plan of deployment %s from pending checks",deploymentId)
    p.mu.Lock()
    defer p.mu.Unlock()
    for _, f := range p.pending[deploymentId].Fragments {
        // remove the services we find across stages
        for _, stage := range f.Stages {
            for _, serv := range stage.Services {
                delete(p.pendingService, serv.ServiceId)
            }
        }
        // remove fragments
        delete(p.pendingFragment, f.FragmentId)
    }
    // delete the plan
    delete(p.pending, deploymentId)
    p.printStatus()
}

// Check if this plan has ny pending fragment.
func (p *PendingPlans) PlanHasPendingFragments(deploymentId string) bool{
    p.mu.Lock()
    defer p.mu.Unlock()
    plan := p.pending[deploymentId]
    for _, fragment := range plan.Fragments {
        _, isPendingFragment := p.pendingFragment[fragment.FragmentId]
        if isPendingFragment {
            // we found a pending fragment
            return true
        }
    }
    // we iterated through the fragments and they are not pending
    return false
}

func (p *PendingPlans) RemoveFragment(fragmentId string){
    log.Debug().Msgf("remove fragment %s from pending", fragmentId)
    p.mu.Lock()
    p.mu.Unlock()
    // get services Id by checking the corresponding plan
    planId,_ := p.pendingFragment[fragmentId]
    plan,_ := p.pending[planId]
    for _, currentFragment := range plan.Fragments {
        if currentFragment.FragmentId == fragmentId {
            for _, stage := range currentFragment.Stages {
                for _, service := range stage.Services{
                    // delete associated services
                    delete(p.pendingService, service.ServiceId)
                }
            }
            break
        }
    }
    // delete the fragment
    delete(p.pendingFragment,fragmentId)
    p.printStatus()
}

func (p *PendingPlans) MonitoredFragment(fragmentID string) bool {
    p.mu.Lock()
    defer p.mu.Unlock()
    _, exists := p.pendingFragment[fragmentID]
    return exists
}

func (p *PendingPlans) printStatus() {
    log.Info().Msgf("%d pending plans, %d pending fragments, %d pending services",
        len(p.pending), len(p.pendingFragment), len(p.pendingService))
}
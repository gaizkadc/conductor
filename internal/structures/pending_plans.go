/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package structures

import (
    "github.com/nalej/conductor/internal/entities"
    "sync"
    "github.com/rs/zerolog/log"
)

// Struct to control Pending deployment plans
type PendingPlans struct {
    // plan_id -> deployment plan
    Pending map[string]*entities.DeploymentPlan
    // fragment_id -> deployment_plan_id. Index just to improve searching.
    PendingFragment map[string]string
    // service_id -> fragment_id
    PendingService map[string] string
    // mutex
    mu sync.Mutex
}

func NewPendingPlans () *PendingPlans {
    return &PendingPlans{
        PendingService:  make(map[string]string,0),
        PendingFragment: make(map[string]string,0),
        Pending:         make(map[string]*entities.DeploymentPlan),
    }
}


func (p *PendingPlans) AddPendingPlan(plan *entities.DeploymentPlan) {
    log.Debug().Msgf("add plan of deployment %s to Pending checks",plan.DeploymentId)
    p.mu.Lock()
    defer p.mu.Unlock()
    p.Pending[plan.DeploymentId] = plan
    for _, frag := range plan.Fragments {
        p.PendingFragment[frag.FragmentId] = plan.DeploymentId
        for _, stage := range frag.Stages {
            for _, serv := range stage.Services {
                p.PendingService[serv.ServiceId] = frag.FragmentId
            }
        }
    }
    p.printStatus()
}



func (p *PendingPlans) RemovePendingPlan(deploymentId string) {
    log.Debug().Msgf("remove plan of deployment %s from Pending checks",deploymentId)
    p.mu.Lock()
    defer p.mu.Unlock()
    for _, f := range p.Pending[deploymentId].Fragments {
        // remove the services we find across stages
        for _, stage := range f.Stages {
            for _, serv := range stage.Services {
                delete(p.PendingService, serv.ServiceId)
            }
        }
        // remove fragments
        delete(p.PendingFragment, f.FragmentId)
    }
    // delete the plan
    delete(p.Pending, deploymentId)
    p.printStatus()
}

// Check if this plan has ny Pending fragment.
func (p *PendingPlans) PlanHasPendingFragments(deploymentId string) bool{
    p.mu.Lock()
    defer p.mu.Unlock()
    plan := p.Pending[deploymentId]
    for _, fragment := range plan.Fragments {
        _, isPendingFragment := p.PendingFragment[fragment.FragmentId]
        if isPendingFragment {
            // we found a Pending fragment
            return true
        }
    }
    // we iterated through the fragments and they are not Pending
    return false
}

func (p *PendingPlans) RemoveFragment(fragmentId string){
    log.Debug().Msgf("remove fragment %s from Pending", fragmentId)
    p.mu.Lock()
    p.mu.Unlock()
    // get services Id by checking the corresponding plan
    planId,_ := p.PendingFragment[fragmentId]
    plan,_ := p.Pending[planId]
    for _, currentFragment := range plan.Fragments {
        if currentFragment.FragmentId == fragmentId {
            for _, stage := range currentFragment.Stages {
                for _, service := range stage.Services{
                    // delete associated services
                    delete(p.PendingService, service.ServiceId)
                }
            }
            break
        }
    }
    // delete the fragment
    delete(p.PendingFragment,fragmentId)
    p.printStatus()
}

func (p *PendingPlans) MonitoredFragment(fragmentID string) bool {
    p.mu.Lock()
    defer p.mu.Unlock()
    _, exists := p.PendingFragment[fragmentID]
    return exists
}

func (p *PendingPlans) printStatus() {
    log.Info().Msgf("%d Pending plans, %d Pending fragments, %d Pending services",
        len(p.Pending), len(p.PendingFragment), len(p.PendingService))
}
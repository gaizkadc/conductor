/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package monitor

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    "sync"
    "github.com/rs/zerolog/log"
)

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
}
/*
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
*/

// Check if this plan has ny pending fragment.
func (p *PendingPlans) PlanHasPendingFragments(deploymentId string) bool{
    p.mu.Lock()
    defer p.mu.Unlock()
    _, isThere := p.pendingFragment[deploymentId]
    return isThere
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
}

func (p *PendingPlans) MonitoredFragment(fragmentID string) bool {
    p.mu.Lock()
    defer p.mu.Unlock()
    _, exists := p.pendingFragment[fragmentID]
    return exists
}
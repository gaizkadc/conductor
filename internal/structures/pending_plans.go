/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package structures

import (
    "errors"
    "fmt"
    "github.com/nalej/conductor/internal/entities"
    "sync"
    "github.com/rs/zerolog/log"
)

// Struct to control Pending deployment plans
type PendingPlans struct {
    // plan_id -> deployment plan
    Pending map[string]*entities.DeploymentPlan
    // fragment_id -> deployment_plan_id. Index just to improve searching.
    PendingFragments map[string]*PendingFragment
    // service_id -> fragment_id
    PendingService map[string] string
    // Instance_id -> plan_id
    Apps map[string]string
    // mutex
    mu sync.Mutex
}

type PendingFragment struct {
    // Deployment this fragment belongs to
    DeploymentPlanID string
    // True if this fragment is pending
    IsPending bool
}

func NewPendingPlans () *PendingPlans {
    return &PendingPlans{
        PendingService:   make(map[string]string,0),
        PendingFragments: make(map[string]*PendingFragment,0),
        Pending:          make(map[string]*entities.DeploymentPlan),
        Apps:             make(map[string]string,0),
    }
}


func (p *PendingPlans) AddPendingPlan(plan *entities.DeploymentPlan) {
    log.Debug().Msgf("add plan of deployment %s to Pending checks",plan.DeploymentId)
    p.mu.Lock()
    defer p.mu.Unlock()
    p.Pending[plan.DeploymentId] = plan
    for _, frag := range plan.Fragments {
        p.PendingFragments[frag.FragmentId] = &PendingFragment{plan.DeploymentId, true}
        for _, stage := range frag.Stages {
            for _, serv := range stage.Services {
                p.PendingService[serv.ServiceId] = frag.FragmentId
                p.Apps[plan.AppInstanceId] = plan.DeploymentId
            }
        }
    }
    p.printStatus()
}


// Look for the plan pointing to an application instance and delete it
func (p *PendingPlans) RemovePendingPlanByApp(appInstanceId string) error {
    p.mu.Lock()
    targetPlanId := ""
    for planId, p := range p.Pending {
        if p.AppInstanceId == appInstanceId {
            targetPlanId = planId
            break
        }
    }
    p.mu.Unlock()

    if targetPlanId == "" {
        log.Error().Str("app_instance_id",appInstanceId).Msg("deployment plan was not found for the given app instance id")
        return errors.New(fmt.Sprintf("deployment plan was not found for the given app instance id %s", appInstanceId))
    }

    p.RemovePendingPlan(targetPlanId)
    return nil
}


func (p *PendingPlans) RemovePendingPlan(deploymentId string) {
    log.Debug().Msgf("remove plan of deployment %s from Pending checks",deploymentId)
    p.mu.Lock()
    defer p.mu.Unlock()

    _, found := p.Pending[deploymentId]
    if !found {
        log.Debug().Str("deploymentId", deploymentId).Msg("the plan was already removed")
        return
    }

    for _, f := range p.Pending[deploymentId].Fragments {
        // remove the services we find across stages
        for _, stage := range f.Stages {
            for _, serv := range stage.Services {
                delete(p.PendingService, serv.ServiceId)
            }
        }
        // remove fragments
        delete(p.PendingFragments, f.FragmentId)
    }
    // delete the plan
    delete(p.Pending, deploymentId)
    // delete the app
    delete(p.Apps, deploymentId)
    p.printStatus()
}

// Check if this plan has ny Pending fragment.
func (p *PendingPlans) PlanHasPendingFragments(deploymentId string) bool{
    p.mu.Lock()
    defer p.mu.Unlock()
    plan, found := p.Pending[deploymentId]
    if !found {
        return false
    }
    for _, fragment := range plan.Fragments {
        _, isPendingFragment := p.PendingFragments[fragment.FragmentId]
        if isPendingFragment {
            // we found a Pending fragment
            return true
        }
    }
    // we iterated through the fragments and they are not Pending
    return false
}

func (p *PendingPlans) SetFragmentNoPending(fragmentId string) {
    log.Debug().Msgf("set fragment %s to non pending", fragmentId)
    p.mu.Lock()
    p.mu.Unlock()
    // get services Id by checking the corresponding plan
    p.PendingFragments[fragmentId].IsPending = false
    p.printStatus()
}

func (p *PendingPlans) SetFragmentPending(fragmentId string) {
    log.Debug().Msgf("set fragment %s to pending", fragmentId)
    p.mu.Lock()
    p.mu.Unlock()
    // get services Id by checking the corresponding plan
    p.PendingFragments[fragmentId].IsPending = true
    p.printStatus()
}


func (p *PendingPlans) RemoveFragment(fragmentId string){
    log.Debug().Msgf("remove fragment %s from Pending", fragmentId)
    p.mu.Lock()
    p.mu.Unlock()
    // get services Id by checking the corresponding plan
    pendingPlan,_ := p.PendingFragments[fragmentId]
    plan,_ := p.Pending[pendingPlan.DeploymentPlanID]
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
    delete(p.PendingFragments,fragmentId)
    p.printStatus()
}

func (p *PendingPlans) MonitoredFragment(fragmentID string) bool {
    p.mu.Lock()
    defer p.mu.Unlock()
    _, exists := p.PendingFragments[fragmentID]
    return exists
}

func (p *PendingPlans) printStatus() {
    log.Info().Msgf("%d Pending plans, %d Pending fragments, %d Pending services",
        len(p.Pending), len(p.PendingFragments), len(p.PendingService))
}
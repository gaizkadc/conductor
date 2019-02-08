/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package plandesigner

import (
    "github.com/google/uuid"
    "github.com/nalej/conductor/internal/structures"
    pbApplication "github.com/nalej/grpc-application-go"
    pbOrganization "github.com/nalej/grpc-organization-go"
    "github.com/nalej/derrors"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/utils"
    "context"
    "github.com/rs/zerolog/log"
    "fmt"
    "sort"
)

/*
 * This plan designer designs a plan for scenarios where one or several services are indicated to run using
 * replication in all the available clusters of the organization.
 */



type SimpleReplicaPlanDesigner struct {
    // Applications client
    appClient pbApplication.ApplicationsClient
    // Organizations client
    orgClient pbOrganization.OrganizationsClient
    // Connections helper
    connHelper *utils.ConnectionsHelper
}

func NewSimpleReplicaPlanDesigner (connHelper *utils.ConnectionsHelper) PlanDesigner {
    connectionsSM := connHelper.GetSystemModelClients()
    appClient := pbApplication.NewApplicationsClient(connectionsSM.GetConnections()[0])
    orgClient := pbOrganization.NewOrganizationsClient(connectionsSM.GetConnections()[0])
    return &SimpleReplicaPlanDesigner{appClient: appClient, orgClient: orgClient, connHelper: connHelper}
}



func(p *SimpleReplicaPlanDesigner) DesignPlan(app entities.AppInstance,
score entities.DeploymentScore, request entities.DeploymentRequest) (*entities.DeploymentPlan, error) {

    // Build deployment stages for the application
    retrievedDesc,err :=p.appClient.GetAppDescriptor(context.Background(),
        &pbApplication.AppDescriptorId{OrganizationId: app.OrganizationId, AppDescriptorId: app.AppDescriptorId})
    if err!=nil{
        theErr := derrors.NewGenericError("error recovering application instance", err)
        log.Error().Err(theErr).Msg("error recovering application instance")
        return nil, theErr
    }

    // get organization name
    org, err := p.orgClient.GetOrganization(context.Background(),
        &pbOrganization.OrganizationId{OrganizationId: app.OrganizationId})
    if err != nil {
        theErr := derrors.NewGenericError("error when retrieving organization data", err)
        log.Error().Err(err).Msgf("error when retrieving organization data")
        return nil, theErr
    }

    // Get a local representation of the object
    toDeploy := entities.NewAppDescriptorFromGRPC(retrievedDesc)

    // Build nalej variables
    nalejVariables := GetDeploymentNalejVariables(org.Name, app.AppInstanceId, toDeploy)

    planId := uuid.New().String()
    log.Info().Str("planId",planId).Msg("start building the plan")

    // There must be one fragment per service group
    // Each service group with a set of stages following the stages defined in the DeploymentPlan
    // Store the group name and the corresponding deployment order.
    log.Debug().Str("appDescriptor",toDeploy.AppDescriptorId).Msg("analyze group internal dependencies")
    groupsOrder := make(map[string][][]entities.Service)
    for _, g := range toDeploy.Groups {
        log.Debug().Str("appDescriptor",toDeploy.AppDescriptorId).Str("serviceGroupId",g.ServiceGroupId).
            Msg("compute dependency graph for service group")
        dependencyGraph := NewDependencyGraph(g.Services)
        order, err := dependencyGraph.GetDependencyOrderByGroups()
        if err != nil {
            return nil, err
        }
        groupsOrder[g.Name] = order
    }

    // Build deployment matrix
    deploymentMatrix := structures.NewDeploymentMatrix(score)


    // Build fragments to deploy replicated groups
    replicatedDeployment, err := p.buildFragmentsReplicatedGroups(app, toDeploy, deploymentMatrix, groupsOrder, nalejVariables, planId, org)
    if err != nil {
        log.Error().Err(err).Msg("impossible to build deployment for groups with replica set flag enabled")
        return nil, err
    }

    log.Debug().Str("appDescriptorId",app.AppDescriptorId).Int("fragments for replication",len(replicatedDeployment)).
        Msg("deployment fragments for replicated groups already processed")

    // Build fragments to deploy everything into a single cluster
    uniqueDeployment, err := p.buildFragmentsAllGroups(app,toDeploy,deploymentMatrix, groupsOrder, nalejVariables, planId, org)
    if err != nil {
        log.Error().Err(err).Msg("impossible to build deployment for the whole application")
        return nil, err
    }
    log.Debug().Str("appDescriptorId",app.AppDescriptorId).Int("fragments for single deployment",len(replicatedDeployment)).
        Msg("deployment fragments for replicated groups already processed")


    // Now that we have all the fragments, build the deployment plan
    newPlan := entities.DeploymentPlan{
        AppInstanceId: app.AppInstanceId,
        DeploymentId: planId,
        OrganizationId: app.OrganizationId,
        Fragments: append(uniqueDeployment, replicatedDeployment...),
    }

    log.Info().Str("appDescriptorId",app.AppDescriptorId).Str("planId",newPlan.DeploymentId).
        Int("number of fragments",len(newPlan.Fragments)).
        Interface("plan",newPlan).
        Msg("a plan was generated")

    return &newPlan, nil
}


// Local function to build the fragments for groups to be replicated.
func (p *SimpleReplicaPlanDesigner) buildFragmentsReplicatedGroups(
    app entities.AppInstance,
    desc entities.AppDescriptor,
    deploymentMatrix *structures.DeploymentMatrix,
    groupsOrder map[string][][]entities.Service,
    nalejVariables map[string]string,
    planId string,
    org *pbOrganization.Organization) ([]entities.DeploymentFragment,derrors.Error) {


    // Get the list of groups with replicate set flag enabled
    replicatedGroups := p.getReplicatedGroups(desc)
    // Allocate them in the deployment matrix
    toReturn := make([]entities.DeploymentFragment, 0)
    for _, g := range replicatedGroups {
        targets, err := deploymentMatrix.FindBestTargetsForReplication(g)
        if err != nil {
            return nil, err
        }
        for _, target := range targets {
            fragments,err  := p.buildFragments(app, groupsOrder, target, nalejVariables, planId, org)
            if err != nil {
                return nil, err
            }
            toReturn = append(toReturn, fragments...)
        }
    }

    return toReturn, nil
}

// Local function to build fragments for a cluster containing all the groups.
func (p *SimpleReplicaPlanDesigner) buildFragmentsAllGroups(
    app entities.AppInstance,
    desc entities.AppDescriptor,
    deploymentMatrix *structures.DeploymentMatrix,
    groupsOrder map[string][][]entities.Service,
    nalejVariables map[string]string,
    planId string,
    org *pbOrganization.Organization) ([]entities.DeploymentFragment,derrors.Error) {

    targetCluster := deploymentMatrix.FindBestTargetForGroups(desc.Groups)

    if targetCluster == "" {
        msg := fmt.Sprintf("no available target cluster was found for app %s",app.AppInstanceId)
        err := derrors.NewGenericError(msg)
        log.Error().Err(err).Msg(msg)
        return nil, err
    }

    // Create a fragment with all the services contained in this application
    fragmentsToDeploy, err := p.buildFragments(app, groupsOrder, targetCluster, nalejVariables, planId, org)
    if err!=nil{
        return nil, derrors.NewGenericError("impossible to build deployment fragment", err)
    }
    return fragmentsToDeploy, nil
}

// This local function returns a fragment for a given list of services and its dependency graph
func (p *SimpleReplicaPlanDesigner) buildFragments(
    app entities.AppInstance,
    groupsOrder map[string][][]entities.Service,
    targetCluster string,
    nalejVariables map[string]string,
    planId string,
    org *pbOrganization.Organization,
    ) ([]entities.DeploymentFragment, derrors.Error) {

    fragments := make([]entities.DeploymentFragment,0)
    // collect stages per group and generate one fragment
    for _, group := range app.Groups {
        // UUID for this fragment
        fragmentUUID := uuid.New().String()

        order := groupsOrder[group.Name]
        // create the stages corresponding to this group
        log.Debug().Str("appDescriptor", app.AppDescriptorId).Str("groupName",group.Name).
            Interface("sequences", order).Msg("create stages for deployment sequence")
        stages := make([]entities.DeploymentStage, 0)
        for _, sequence := range order {
            // this stage must deploy the services following this order
            stage := entities.DeploymentStage{
                FragmentId: fragmentUUID,
                StageId: uuid.New().String(),
                Services: sequence,
            }
            stages = append(stages, stage)
        }

        fragment := entities.DeploymentFragment{
            OrganizationId:         app.OrganizationId,
            OrganizationName:       org.Name,
            AppInstanceId:          app.AppInstanceId,
            AppName:                app.Name,
            FragmentId:             fragmentUUID,
            DeploymentId:           planId,
            Stages:                 stages,
            NalejVariables:         nalejVariables,
            GroupServiceInstanceId: group.ServiceGroupInstanceId,
            ClusterId:              targetCluster,
        }
        fragments = append(fragments, fragment)
    }


    return fragments, nil

}

// Local function to collect the list of groups with a replica set flag enabled.
// This function returns an array with the groups.
func(p *SimpleReplicaPlanDesigner) getReplicatedGroups(toDeploy entities.AppDescriptor) (
    []entities.ServiceGroup) {
    toReturn := make([]entities.ServiceGroup,0)
    for _, g := range toDeploy.Groups{
        if g.Specs.MultiClusterReplica {
            toReturn = append(toReturn, g)
        }

    }
    return toReturn
}


// This local function finds the cluster with the largest score for all the service groups.
// TODO review this method to go for a more generic approach using the deployment score matrix
func (p *SimpleReplicaPlanDesigner) findTargetCluster(serviceGroups []entities.ServiceGroup,
    scores *entities.DeploymentScore) string {

    serviceGroupIds := make([]string,len(serviceGroups))
    for i, s := range serviceGroups {
        serviceGroupIds[i] = s.ServiceGroupId
    }
    sort.Strings(serviceGroupIds)
    // concatenate
    fullGroupId := ""
    for _, s := range serviceGroupIds {
        fullGroupId = fullGroupId + s
    }

    // Find cluster with the largest score for this fullGroupId
    maxScore := float32(0)
    bestCandidate := ""
    for _, sc := range scores.DeploymentsScore {
        targetScore, found := sc.Scores[fullGroupId]
        if !found {
            log.Debug().Str("fullDescriptorId", fullGroupId).Interface("clusterScore",sc).
                Msg("full descriptor id not found")
        } else if maxScore < targetScore{
            maxScore = targetScore
            bestCandidate = sc.ClusterId
        }
    }

    log.Debug().Str("bestCandidate",bestCandidate).Float32("score",maxScore).Msg("Best cluster found with score")
    return bestCandidate
}

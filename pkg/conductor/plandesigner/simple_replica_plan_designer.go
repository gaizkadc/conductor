/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package plandesigner

import (
    "github.com/google/uuid"
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



func(p *SimpleReplicaPlanDesigner) DesignPlan(app *pbApplication.AppInstance,
score *entities.DeploymentScore, request *entities.DeploymentRequest) (*entities.DeploymentPlan, error) {

    // Build deployment stages for the application
    toDeploy ,err :=p.appClient.GetAppDescriptor(context.Background(),
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

    // Build nalej variables
    nalejVariables := GetDeploymentNalejVariables(org.Name, app.AppInstanceId, toDeploy)

    planId := uuid.New().String()
    log.Info().Str("planId",planId).Msg("start building the plan")
    planFragments := make([]entities.DeploymentFragment,0)

    // There must be one fragment per service group
    for _, g := range toDeploy.Groups {
        log.Debug().Str("serviceGroupId",g.ServiceGroupId).Msg("compute dependency graph for service group")
        dependencyGraph := NewDependencyGraph()
    }


    // Create a deployment fragment with all the services
    allServices, allIndex := p.getServicesFromDescriptor(toDeploy)


    // Check scores are available and the application fits
    targetCluster := p.findTargetCluster(toDeploy.Groups,score)
    if targetCluster == "" {
        msg := fmt.Sprintf("no available target cluster was found for app %s",app.AppInstanceId)
        err := derrors.NewGenericError(msg)
        log.Error().Err(err).Msg(msg)
        return nil, err
    }
    // Create a fragment with all the services
    completeFragment, err := p.buildFragment(allServices, allIndex, dependencyGraph,
        app, org, targetCluster, nalejVariables, planId)
    if err!=nil{
        return nil, derrors.NewGenericError("impossible to build deployment fragment", err)
    }
    planFragments = append(planFragments, *completeFragment)

    // Create a list of services to be replicated across clusters
    replicatedServices,index := p.getReplicatedServices(toDeploy)
    log.Info().Str("appInstanceId", app.AppInstanceId).Int("numReplicatedServices",len(replicatedServices)).
        Msg("number of services to be replicated across clusters")

    if len(replicatedServices) > 0 {
        // Replicated services will be deployed in the same batch, compute dependencies
        dependencyGraphReplicated := NewDependencyGraph(replicatedServices)

        // Build one fragment per cluster
        for _, scoring := range score.DeploymentsScore {
            replicaCluster := scoring.ClusterId
            // only replicate if this is not the cluster with all the services
            if replicaCluster != targetCluster {
                fragment, err :=p.buildFragment(replicatedServices,index,dependencyGraphReplicated,app,org,
                    replicaCluster, nalejVariables,planId)
                if err!=nil{
                    return nil, derrors.NewGenericError("impossible to build deployment fragment", err)
                }
                planFragments = append(planFragments,*fragment)
            }
        }
    }


    // Now that we have all the fragments, build the deployment plan
    newPlan := entities.DeploymentPlan{
        AppInstanceId: app.AppInstanceId,
        DeploymentId: planId,
        OrganizationId: app.OrganizationId,
        Fragments: planFragments,
    }

    return &newPlan, nil
}

// This local function returns a fragment for a given list of services and its dependency graph
func (p *SimpleReplicaPlanDesigner) buildFragment(
    services []entities.Service,
    index map[string]entities.Service,
    depGraph *DependencyGraph,
    app *pbApplication.AppInstance,
    org *pbOrganization.Organization,
    targetCluster string,
    nalejVariables map[string]string,
    planId string,
    serviceGroupInstanceId string,
    ) (*entities.DeploymentFragment, derrors.Error) {

    // Split it into independent components using the dependencies graph
    groups, err := depGraph.GetDependencyOrderByGroups()
    if err != nil {
        theErr := derrors.NewGenericError("impossible to define deployment stages for app instance",err)
        log.Error().Err(theErr).Msg("impossible to define deployment order when building fragment")
        return nil, theErr
    }

    fragmentUUID := uuid.New().String()

    stages := make([]entities.DeploymentStage, len(groups))
    for stageNumber, servicesPerStage := range groups {
        inThisStage := make([]entities.Service, len(servicesPerStage))
        for i, serviceId := range servicesPerStage {
            theService := index[serviceId]
            if theService.Specs.Replicas <= 0 {
                // Force it to have a single replica
                theService.Specs.Replicas = 1
            }
            inThisStage[i] = theService
        }

        stages[stageNumber] = entities.DeploymentStage{
            FragmentId: fragmentUUID,
            StageId: uuid.New().String(),
            Services: inThisStage,
        }
    }

    fragment := entities.DeploymentFragment{
        OrganizationId: app.OrganizationId,
        OrganizationName: org.Name,
        AppInstanceId: app.AppInstanceId,
        AppName: app.Name,
        FragmentId: fragmentUUID,
        DeploymentId: planId,
        Stages: stages,
        NalejVariables: nalejVariables,
        GroupServiceInstanceId: serviceGroupInstanceId,
    }



    return &fragment, nil

}

// Local function to collect the list of services with a replica set flag enabled.
// This function returns an array with the services and a map with the serviceId and the object.
// E.g.: [serviceB, servicek1, serviceAux],   {serviceAux -> , serviceB->0,servicek1->1}
func(p *SimpleReplicaPlanDesigner) getReplicatedServices(toDeploy *pbApplication.AppDescriptor) (
    []entities.Service, map[string]entities.Service) {
    toReturn := make([]entities.Service,0)
    index := make(map[string]entities.Service,0)
    for _, s := range toDeploy.Services {
        if s.Specs.MultiClusterReplicaSet {
            serv := entities.NewServiceFromGRPC(toDeploy.AppDescriptorId, s)
            toReturn = append(toReturn, *serv)
            index[s.ServiceId]= *serv
        }
    }
    return toReturn,index
}

// Local function to build an array of services and its corresponding index map.
func(p *SimpleReplicaPlanDesigner) getServicesFromDescriptor(toDeploy *pbApplication.AppDescriptor) (
    []entities.Service, map[string]entities.Service) {
    toReturn := make([]entities.Service,0)
    index := make(map[string]entities.Service,0)
    for _, g := range toDeploy.Groups {
        for _, s := range toDeploy.Services {
            serv := entities.NewServiceFromGRPC(toDeploy.AppDescriptorId, s)
            toReturn = append(toReturn, *serv)
            index[s.ServiceId]= *serv
        }
    }

    return toReturn,index
}


func(p *SimpleReplicaPlanDesigner) getServicesFromServiceGroup(serviceGroupId string) (
    []entities.Service, map[string]entities.Service) {
    toReturn := make([]entities.Service,0)
    index := make(map[string]entities.Service,0)
    x := p.appClient.GetAppDescriptor()

    return nil, nil
}


// This local function finds the cluster with the largest score for all the service groups.
// TODO review this method to go for a more generic approach using the deployment score matrix
func (p *SimpleReplicaPlanDesigner) findTargetCluster(serviceGroups []*pbApplication.ServiceGroup,
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

    log.Debug().Str("bestCandidate",bestCandidate).Str("score",maxScore).Msg("Best cluster found with score")
    return bestCandidate
}

/*
// Find the target cluster by considering the target with the largest availability.
func(p *SimpleReplicaPlanDesigner) findTargetCluster(scores *entities.DeploymentScore) string {
    var max float32
    max = 0
    targetCluster := ""
    for _, s := range scores.DeploymentsScore {
        if s.Score > max {
            targetCluster = s.ClusterId
            max = s.Score
        }
    }
    return targetCluster
}
*/


/*
func(p *SimpleReplicaPlanDesigner) DesignPlan(app *pbApplication.AppInstance,
score *entities.DeploymentScore, request *entities.DeploymentRequest) (*entities.DeploymentPlan, error) {

    // Build deployment stages for the application
    toDeploy ,err :=p.appClient.GetAppDescriptor(context.Background(),
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

    // Build nalej variables
    nalejVariables := GetDeploymentNalejVariables(org.Name, app.AppInstanceId, toDeploy)

    planId := uuid.New().String()
    log.Info().Str("planId",planId).Msg("start building the plan")
    planFragments := make([]entities.DeploymentFragment,0)

    // Create a deployment fragment with all the services
    allServices, allIndex := p.getServicesFromDescriptor(toDeploy)
    dependencyGraph := NewDependencyGraph(allServices)

    // Check scores are available and the application fits
    targetCluster := p.findTargetCluster(score)
    if targetCluster == "" {
        msg := fmt.Sprintf("no available target cluster was found for app %s",app.AppInstanceId)
        err := derrors.NewGenericError(msg)
        log.Error().Err(err).Msg(msg)
        return nil, err
    }
    // Create a fragment with all the services
    completeFragment, err := p.buildFragment(allServices, allIndex, dependencyGraph,
        app, org, targetCluster, nalejVariables, planId)
    if err!=nil{
        return nil, derrors.NewGenericError("impossible to build deployment fragment", err)
    }
    planFragments = append(planFragments, *completeFragment)

    // Create a list of services to be replicated across clusters
    replicatedServices,index := p.getReplicatedServices(toDeploy)
    log.Info().Str("appInstanceId", app.AppInstanceId).Int("numReplicatedServices",len(replicatedServices)).
        Msg("number of services to be replicated across clusters")

    if len(replicatedServices) > 0 {
        // Replicated services will be deployed in the same batch, compute dependencies
        dependencyGraphReplicated := NewDependencyGraph(replicatedServices)

        // Build one fragment per cluster
        for _, scoring := range score.DeploymentsScore {
            replicaCluster := scoring.ClusterId
            // only replicate if this is not the cluster with all the services
            if replicaCluster != targetCluster {
                fragment, err :=p.buildFragment(replicatedServices,index,dependencyGraphReplicated,app,org,
                    replicaCluster, nalejVariables,planId)
                if err!=nil{
                    return nil, derrors.NewGenericError("impossible to build deployment fragment", err)
                }
                planFragments = append(planFragments,*fragment)
            }
        }
    }


    // Now that we have all the fragments, build the deployment plan
    newPlan := entities.DeploymentPlan{
        AppInstanceId: app.AppInstanceId,
        DeploymentId: planId,
        OrganizationId: app.OrganizationId,
        Fragments: planFragments,
    }

    return &newPlan, nil
}
 */
/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package plandesigner

import (
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/nalej/conductor/internal/entities"
    pbApplication "github.com/nalej/grpc-application-go"
    "github.com/google/uuid"
    "context"
    "github.com/rs/zerolog/log"
    "fmt"
    "errors"
)


type SimplePlanDesigner struct {
    appClient pbApplication.ApplicationsClient
}

func NewSimplePlanDesigner () PlanDesigner {
    connectionsSM := conductor.GetSystemModelClients()
    appClient := pbApplication.NewApplicationsClient(connectionsSM.GetConnections()[0])
    return &SimplePlanDesigner{appClient: appClient}
}

func (p SimplePlanDesigner) DesignPlan(app *pbApplication.AppInstance,
    score *entities.ClustersScore) (*entities.DeploymentPlan, error) {

    // Check scores are available and the application fits
    targetCluster := p.findTargetCluster(score)
    if targetCluster == "" {
        msg := fmt.Sprintf("no available target cluster was found for app %s",app.AppInstanceId)
        log.Error().Msg(msg)
        return nil, errors.New(msg)
    }

    // Build deployment stages for the application

    toDeploy ,err :=p.appClient.GetAppDescriptor(context.Background(),
        &pbApplication.AppDescriptorId{OrganizationId: app.OrganizationId, AppDescriptorId: app.AppDescriptorId})
    if err!=nil{
        log.Error().Err(err).Msg("error recovering application instance")
        return nil, err
    }
    // TODO this current version is limited to deployments contained into a single cluster

    fragmentUUID := uuid.New().String()
    index := make(map[string]entities.Service,0)

    servicesToDeploy := make([]entities.Service,len(toDeploy.Services))
    for i, serv := range toDeploy.Services {
        ent := *entities.NewServiceFromGRPC(toDeploy.AppDescriptorId,serv)
        servicesToDeploy[i] = ent
        index[serv.ServiceId] = ent
    }

    // Create dependency graph
    depGraph := NewDependencyGraph(servicesToDeploy)

    // Split it into independent components
    groups, err := depGraph.GetDependencyOrderByGroups()
    if err != nil {
        log.Error().Err(err).Msgf("impossible to define deployment stages for app instance %s",app.AppInstanceId)
        return nil, err
    }

    stages := make([]entities.DeploymentStage, len(groups))
    for stageNumber, servicesPerStage := range groups {
        inThisStage := make([]entities.Service, len(servicesPerStage))
        for i, serviceId := range servicesPerStage {
            inThisStage[i] = index[serviceId]
        }

        stages[stageNumber] = entities.DeploymentStage{
            FragmentId: fragmentUUID,
            StageId: uuid.New().String(),
            Services: inThisStage,
        }
    }

    planId := uuid.New().String()

    fragment := entities.DeploymentFragment{
        OrganizationId: app.OrganizationId,
        AppInstanceId: app.AppInstanceId,
        FragmentId: fragmentUUID,
        DeploymentId: planId,
        Stages: stages,
        ClusterId: targetCluster,

    }

    // Aggregate to a new plan
    newPlan := entities.DeploymentPlan{
        AppInstanceId: app.AppInstanceId,
        DeploymentId: planId,
        OrganizationId: app.OrganizationId,
        Fragments: []entities.DeploymentFragment{fragment},
    }

    return &newPlan, nil
}

func (p SimplePlanDesigner) findTargetCluster(scores *entities.ClustersScore) string {
    var max float32
    max = 0
    targetCluster := ""
    for _, s := range scores.Scoring {
        if s.Score > max {
            targetCluster = s.ClusterId
            max = s.Score
        }
    }
    return targetCluster
}



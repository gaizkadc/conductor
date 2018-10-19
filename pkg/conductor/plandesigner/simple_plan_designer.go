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
    score *entities.ClusterScore) (*entities.DeploymentPlan, error) {
    // Build deployment stages for the application

    toDeploy ,err :=p.appClient.GetAppDescriptor(context.Background(),
        &pbApplication.AppDescriptorId{OrganizationId: app.OrganizationId, AppDescriptorId: app.AppDescriptorId})
    if err!=nil{
        log.Error().Err(err).Msg("error recovering application instance")
        return nil, err
    }
    // TODO this version assumes everything will go into a single cluster

    fragmentUUID := uuid.New().String()
    stageUUID := uuid.New().String()

    servicesToDeploy := make([]entities.Service,len(toDeploy.Services))
    for i, serv := range toDeploy.Services {
        servicesToDeploy[i] = *entities.NewServiceFromGRPC(toDeploy.AppDescriptorId,serv)
    }

    stage := entities.DeploymentStage{
        FragmentId: fragmentUUID,
        StageId: stageUUID,
        Services: servicesToDeploy,
    }

    planId := uuid.New().String()

    fragment := entities.DeploymentFragment{
        OrganizationId: app.OrganizationId,
        AppInstanceId: app.AppInstanceId,
        FragmentId: fragmentUUID,
        DeploymentId: planId,
        Stages: []entities.DeploymentStage{stage},
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

/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package plandesigner

import (
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/nalej/conductor/internal/entities"
    pbConductor "github.com/nalej/grpc-conductor-go"
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
    score *entities.ClusterScore) (*pbConductor.DeploymentPlan, error) {
    // Build deployment stages for the application

    toDeploy ,err :=p.appClient.GetAppDescriptor(context.Background(),
        &pbApplication.AppDescriptorId{OrganizationId: app.OrganizationId, AppDescriptorId: app.AppDescriptorId})
    if err!=nil{
        log.Error().Err(err).Msgf("error recovering application instance %s", )
        return nil, err
    }
    // TODO this version assumes everything will go into a single cluster

    fragmentUUID := uuid.New().String()
    stageUUID := uuid.New().String()

    stage := pbConductor.DeploymentStage{
        FragmentId: fragmentUUID,
        StageId: stageUUID,
        Services: toDeploy.Services}

    fragment := pbConductor.DeploymentFragment{
        OrganizationId: app.OrganizationId,
        InstanceId: app.AppInstanceId,
        FragmentId: fragmentUUID,
        DeploymentId: uuid.New().String(),
        Stages: []*pbConductor.DeploymentStage{&stage},
    }

    // Aggregate to a new plan
    newPlan := pbConductor.DeploymentPlan{
        InstanceId: app.AppInstanceId,
        DeploymentId: uuid.New().String(),
        OrganizationId: app.OrganizationId,
        Fragments: []*pbConductor.DeploymentFragment{&fragment},
    }

    return &newPlan, nil
}

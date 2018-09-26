/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 */

package plandesigner

import (
    "github.com/nalej/conductor/tools"
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/nalej/conductor/internal/entities"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
    "github.com/google/uuid"
)


type SimplePlanDesigner struct {
    connections *tools.ConnectionsMap
}

func NewSimplePlanDesigner () PlanDesigner {
    connections := conductor.GetDMClients()
    return &SimplePlanDesigner{connections: connections}
}

func (p SimplePlanDesigner) DesignPlan(app *pbApplication.AppDescriptor,
    score *entities.ClusterScore) (*pbConductor.DeploymentPlan, error) {
    // Build deployment stages for the application
    // TODO this version assumes everything will go into a single cluster
    testDeploymentServices := p.getTestServices()
    fragmentUUID := uuid.New().String()
    stageUUID := uuid.New().String()

    stage := pbConductor.DeploymentStage{
        FragmentId: fragmentUUID,
        StageId: stageUUID,
        Services: testDeploymentServices}

    // TODO this is hardcoded until we can access the system model
    fragment := pbConductor.DeploymentFragment{
        FragmentId: fragmentUUID,
        DeploymentId: uuid.New().String(),
        Stages: []*pbConductor.DeploymentStage{&stage},
    }

    // Aggregate to a new plan
    newPlan := pbConductor.DeploymentPlan{
        AppId: &pbApplication.AppDescriptorId{AppDescriptorId:app.AppDescriptorId, OrganizationId: app.OrganizationId},
        DeploymentId: uuid.New().String(),
        OrganizationId: app.OrganizationId,
        Fragments: []*pbConductor.DeploymentFragment{&fragment},
    }

    return &newPlan, nil
}


func (p SimplePlanDesigner) getTestServices() []*pbConductor.Service{

    port1 := pbApplication.Port{Name: "port1", ExposedPort: 3000}
    port2 := pbApplication.Port{Name: "port2", ExposedPort: 3001}

    serv1 := pbApplication.Service{
        ServiceId: "service_001",
        Name: "test-image-1",
        Image: "nginx:1.12",
        ExposedPorts: []*pbApplication.Port{&port1, &port2},
        Labels: map[string]string { "label1":"value1", "label2":"value2"},
        Specs: &pbApplication.DeploySpecs{Replicas: 1},
    }

    serv2 := pbApplication.Service{
        ServiceId: "service_002",
        Name: "test-image-2",
        Image: "nginx:1.12",
        ExposedPorts: []*pbApplication.Port{&port1, &port2},
        Labels: map[string]string { "label1":"value1"},
        Specs: &pbApplication.DeploySpecs{Replicas: 2},
    }

    toReturn :=[]*pbConductor.Service{&serv1,&serv2}

    return toReturn
}
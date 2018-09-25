/*
 * Copyright 2018 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plandesigner

import (
    "github.com/nalej/conductor/tools"
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/nalej/conductor/internal/entities"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
    pbInfrastructure "github.com/nalej/grpc-infrastructure-go"
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
    services *pbApplication.ServiceGroup,
    score *entities.ClusterScore) (*pbConductor.DeploymentPlan, error) {
    // Build deployment stages for the application
    // TODO this version assumes everything will go into a single cluster
    testDeploymentServices := p.getTestServices()
    stage := pbConductor.DeploymentStage{
        DeploymentId: uuid.New().String(),
        StageId: uuid.New().String(),
        Services: testDeploymentServices,
        ClusterId: &pbInfrastructure.ClusterId{ClusterId:"cluster_001", OrganizationId:"org_001"}}

    // TODO this is hardcoded until we can access the system model

    // Aggregate in a new plan
    newPlan := pbConductor.DeploymentPlan{
        DeploymentId: uuid.New().String(),
        Stages: []*pbConductor.DeploymentStage{&stage},
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
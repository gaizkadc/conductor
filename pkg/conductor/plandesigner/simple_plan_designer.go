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

    // Get the application

    // Aggregate in a new plan
    newPlan := pbConductor.DeploymentPlan{
        DeploymentId: uuid.New().String(),
    }
    return &newPlan, nil
}

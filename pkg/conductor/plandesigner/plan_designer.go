/*
 * Copyright (C) 2018 Nalej Group -All Rights Reserved
 */

package plandesigner

import (
    "github.com/nalej/conductor/internal/entities"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
)

// Basic interface to be follow by any plan designer.

type PlanDesigner interface {

    // For any set of requirements an a given score elaborate a deployment plan.
    //  params:
    //   app application descriptor
    //   services group of services to deploy
    //   score obtained by musicians
    //  return:
    //   A collection of deployment plans each one designed to run in a different cluster.
    DesignPlan(app *pbApplication.AppDescriptor,
        score *entities.ClusterScore) (*pbConductor.DeploymentPlan, error)
}

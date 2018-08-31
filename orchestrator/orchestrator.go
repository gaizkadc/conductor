//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// The orchestrator is in charge of controlling the execution of the applications deployed on top of the
// Daisho platform.

package orchestrator

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    log "github.com/sirupsen/logrus"
    entitiesConductor "github.com/nalej/conductor/entities"
)

var logger = log.WithField("package", "orchestrator")

type Orchestrator interface {
    // Deploy an application following some orchestration solution.
    // params:
    //   networkId Identifier of the target network
    //   descriptor of the application to be deployed
    //   appRequest The application request.
    // return:
    //   instance representing the deployed application
    //   error if any
    Deploy(networkId string, descriptor entities.AppDescriptor, appRequest entitiesConductor.DeployAppRequest) (*entities.AppInstance, derrors.DaishoError)

    // Undeploy an already instanced application
    // params:
    //  networkId Identifier of the target network
    //  appInstance instance of an already deployed application
    // error:
    //  Error if any
    Undeploy(networkId string, appInstance entities.AppInstance) derrors.DaishoError
}

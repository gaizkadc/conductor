//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Interface for the conductor service.

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    log "github.com/sirupsen/logrus"
    entitiesConductor "github.com/daishogroup/conductor/entities"
)

var logger = log.WithField("package", "entities")

type Conductor interface {
    // Deploy a new application instance.
    //   params:
    //     networkId The network identifier.
    //     request The deploy app request.
    //   returns:
    //     An application instance
    //     An error if the application instance cannot be deployed.
    Deploy(networkId string,
        request entitiesConductor.DeployAppRequest) (*entities.AppInstance, derrors.DaishoError)

    // Undeploy a new application instance.
    //   params:
    //     networkId target network id
    //     appInstanceId Id of an already deployed instance
    //   returns:
    //      error if any
    Undeploy(networkId string, appInstanceId string) derrors.DaishoError

    // Logs get a set of log entries from the selected application.
    // params:
    //   networkId identifier of the target network
    //   instanceId identifier of the target application instance
    // return:
    //   Error if any.
    //   An array of strings.
    Logs(networkId string, instanceId string) (*entitiesConductor.LogEntries, derrors.DaishoError)
}

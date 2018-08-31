//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Apps manager.

package apps

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    entitiesConductor "github.com/nalej/conductor/entities"
)

type AppManager interface {
    // Deploy an application in a given network.
    // params:
    //   networkId identifier of the target network
    //   appName name of the application to deploy
    // return:
    //   Request to add a new instance
    //   Error if any.
    Deploy(networkId string, request entitiesConductor.DeployAppRequest) (*entities.AppInstance, derrors.DaishoError)

    // Undeploy an already deployed application.
    // params:
    //   networkId identifier of the target network
    //   instanceId identifier of the target application instance
    // return:
    //   Error if any.
    Undeploy(networkId string, instanceId string) derrors.DaishoError

    // Logs get a set of log entries from the selected application.
    // params:
    //   networkId identifier of the target network
    //   instanceId identifier of the target application instance
    // return:
    //   Error if any.
    //   An array of strings.
    Logs(networkId string, instanceId string) (*entitiesConductor.LogEntries, derrors.DaishoError)
}

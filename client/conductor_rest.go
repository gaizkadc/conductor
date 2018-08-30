//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Client to connect with conductor API REST

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "fmt"
    "github.com/daishogroup/dhttp"
    entitiesConductor "github.com/daishogroup/conductor/entities"
)

// URLS to use
const DeployAppURI = "/api/v0/app/%s/deploy"
const UndeployAppURI = "/api/v0/app/%s/%s/undeploy"
const LogsAppURI = "/api/v0/app/%s/%s/logs"

type ConductorRest struct {
    client dhttp.Client
}

func NewConductorRest(host string, port int) Conductor {
    conf := dhttp.NewRestBasicConfig(host, port)
    rest := dhttp.NewClientSling(conf)
    return &ConductorRest{rest}
}

// Deploy a new application instance.
//   params:
//     networkId The network identifier.
//     request The deploy app request.
//   returns:
//     An application instance
//     An error if the application instance cannot be deployed.
func (rest *ConductorRest) Deploy(networkId string,
    request entitiesConductor.DeployAppRequest) (*entities.AppInstance, derrors.DaishoError) {
    url := fmt.Sprintf(DeployAppURI, networkId)

    response := rest.client.Post(url, request, new(entities.AppInstance))
    if response.Ok() {
        added := response.Result.(*entities.AppInstance)
        return added, nil
    }
    return nil, derrors.NewOperationError("error deploying application",response.Error).WithParams(networkId, request)
}

func (rest *ConductorRest) Undeploy(networkId string, instanceId string) derrors.DaishoError {
    url := fmt.Sprintf(UndeployAppURI, networkId, instanceId)
    // we are not waiting for any output
    response := rest.client.Get(url, new(entities.SuccessfulOperation))
    if response.Ok() {
        return nil
    }
    return derrors.NewOperationError("error undeploying application",response.Error).WithParams(networkId, instanceId)

}

// Logs get a set of log entries from the selected application.
// params:
//   networkId identifier of the target network
//   instanceId identifier of the target application instance
// return:
//   Error if any.
//   An array of strings.
func (rest *ConductorRest) Logs(networkId string,
    instanceId string) (*entitiesConductor.LogEntries, derrors.DaishoError) {
    url := fmt.Sprintf(LogsAppURI, networkId,instanceId)

    response := rest.client.Get(url, new(entitiesConductor.LogEntries))
    if response.Ok() {
        logs := response.Result.(*entitiesConductor.LogEntries)
        return logs, nil
    }
    return nil, derrors.NewOperationError("error obtaining logs", response.Error).WithParams(networkId, instanceId)
}
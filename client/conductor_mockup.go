//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Mockup client for conductor

package client

import (
    entitiesConductor "github.com/daishogroup/conductor/entities"
    "github.com/daishogroup/derrors"
    smClient "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
)

type ConductorMockup struct {
    appClient smClient.Applications
}

func NewConductorMockup() Conductor {
    return &ConductorMockup{smClient.NewApplicationsMockup()}
}

func (mockup *ConductorMockup) Deploy(networkId string,
    request entitiesConductor.DeployAppRequest) (*entities.AppInstance, derrors.DaishoError) {
    // TODO more sophisticated logic
    err := request.IsValid()
    if err == nil {
        result := entities.NewAppInstance(networkId, request.AppDescriptorId, "clusterId", request.Name,
            request.Description, request.Label, "arguments", "10GB", entities.AppStorageDefault,
                make([]entities.ApplicationPort, 0), 80, "address")
        return result, nil
    }
    return nil, err
}

func (mockup *ConductorMockup) Undeploy(networkId string, appInstanceId string) derrors.DaishoError {
    _, err := mockup.appClient.GetInstance(networkId, appInstanceId)
    if err != nil {
        logger.Errorf("Impossible to find instance %s", err)
        return err
    }

    // Assume this is undeployed
    // TODO additional logic
    return nil
}

// Logs get a set of log entries from the selected application.
// params:
//   networkId identifier of the target network
//   instanceId identifier of the target application instance
// return:
//   Error if any.
//   An array of strings.
func (mockup *ConductorMockup) Logs(networkId string,
    instanceId string) (*entitiesConductor.LogEntries, derrors.DaishoError) {
    return nil, nil
}

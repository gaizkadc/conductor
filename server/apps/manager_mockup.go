//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Manager mockup specially designed for testing.

package apps

import (
    "github.com/nalej/conductor/errors"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/client"
    "github.com/nalej/conductor/orchestrator"
    entitiesConductor "github.com/nalej/conductor/entities"
    "github.com/daishogroup/system-model/entities"
    "github.com/nalej/conductor/asm"
)

// Manager structure using a remote system client entry and one nocmgr client.
type MockupManager struct {
    appClient     client.Applications
    clusterClient client.Cluster
    nodeClient    client.Node
    orchestrator  orchestrator.Orchestrator
}

// Create a new mockup manager.
// returns:
//   instance of a mockup manager
func NewMockupManager() AppManager {
    clusterClient := client.NewClusterMockup()
    nodeClient := client.NewNodeMockup()
    appClient := client.NewApplicationsMockup()
    asmClientFactory := asm.NewMockupClientFactory()
    orchestrator := orchestrator.NewBasicOrchestrator(clusterClient, nodeClient, appClient, asmClientFactory)
    return &RestManager{orchestrator: orchestrator,
        clusterClient: clusterClient,
        nodeClient: nodeClient,
        appClient: appClient,
    }
}

func (m *MockupManager) Deploy(networkId string, appRequest entitiesConductor.DeployAppRequest) (*entities.AppInstance, derrors.DaishoError) {
    // Call the orchestrator to deploy the application
    // Check if this application is available in this network
    // TODO think about what would happen when using a large number of available applications

    list, err := m.appClient.ListDescriptors(networkId)
    logger.Debugf("Available descriptors in network %s: %s", networkId, list)

    if err != nil {
        logger.Errorf("Problem returning the list of available descriptors from %s [%s]", networkId, err)
        return nil, nil
    }

    if len(list) == 0 {
        logger.Error("No applications available in network ", networkId)
        return nil, derrors.NewOperationError(errors.NoApplicationsAvailable).WithParams(networkId)
    }

    var deploymentResult derrors.DaishoError = nil
    var deployedInstance *entities.AppInstance = nil

    // Find the descriptor
    for _, descriptor := range list {
        if descriptor.Name == appRequest.Name {
            // found it, deploy it
            logger.Debugf("Descriptor %s found!!", appRequest.Name)
            deployedInstance, deploymentResult = m.orchestrator.Deploy(networkId, descriptor, appRequest)
            return deployedInstance, deploymentResult
        }
    }

    return nil, derrors.NewOperationError(errors.CannotDeployApp).WithParams(networkId, appRequest)
}

func (m *MockupManager) Undeploy(networkId string, instanceId string) derrors.DaishoError {
    // TODO
    return nil
}

// Logs get a set of log entries from the selected application.
// params:
//   networkId identifier of the target network
//   instanceId identifier of the target application instance
// return:
//   Error if any.
//   An array of strings.
func (m *MockupManager) Logs(networkId string, instanceId string) (*entitiesConductor.LogEntries, derrors.DaishoError) {
    return nil, nil
}

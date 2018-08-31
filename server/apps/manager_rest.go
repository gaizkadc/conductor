//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Apps manager using an API client to talk with the system model.

package apps

import (
    "github.com/nalej/conductor/errors"
    "github.com/nalej/conductor/orchestrator"
    "github.com/daishogroup/derrors"
    smclient "github.com/daishogroup/system-model/client"
    entitiesConductor "github.com/nalej/conductor/entities"
    "github.com/daishogroup/system-model/entities"
    "github.com/nalej/conductor/asm"
    loggerAggregator "github.com/nalej/conductor/logger"
)

// Manager structure using a remote system client entry and one nocmgr client.
type RestManager struct {
    appClient        smclient.Applications
    clusterClient    smclient.Cluster
    nodeClient       smclient.Node
    orchestrator     orchestrator.Orchestrator
    asmClientFactory asm.ClientFactory
    loggerClient     loggerAggregator.Client
}

// Create a new rest manager using an API-REST client to connect to the system model.
// params:
//   appClient applications client
//   clusterClient clusters client
//   nodeClient nodes client
// returns:
//   instance of an apps rest manager
func NewRestManager(appClient smclient.Applications, clusterClient smclient.Cluster,
    nodeClient smclient.Node,
    asmClientFactory asm.ClientFactory, loggerClient loggerAggregator.Client) AppManager {
    or := orchestrator.NewBasicOrchestrator(clusterClient, nodeClient, appClient, asmClientFactory)
    return &RestManager{orchestrator: or,
        clusterClient: clusterClient,
        nodeClient: nodeClient,
        appClient: appClient,
        asmClientFactory: asmClientFactory,
        loggerClient: loggerClient,
    }
}

func (m *RestManager) Deploy(networkId string,
    appRequest entitiesConductor.DeployAppRequest) (*entities.AppInstance, derrors.DaishoError) {

    var deployedInstance *entities.AppInstance = nil

    // Find the descriptor
    descriptor, err := m.appClient.GetDescriptor(networkId, appRequest.AppDescriptorId)

    if err != nil {
        logger.Debugf("Error when trying to get descriptor %s - %s",networkId, appRequest.AppDescriptorId)
        //return nil, derrors.NewOperationError(errors.CannotRetrieveApp).WithParams(networkId).WithParams(appRequest).CausedBy(err)
        return nil, derrors.NewOperationError(errors.CannotRetrieveApp,err)
    }

    logger.Debugf("Found descriptor %s", descriptor.String())
    deployedInstance, err = m.orchestrator.Deploy(networkId, *descriptor, appRequest)

    if err != nil {
        logger.Debugf("There was an error during the deployment %s", err)
        return nil, err
    }

    if deployedInstance == nil {
        return nil, derrors.NewOperationError(errors.CannotDeployApp, err)
        //return nil, derrors.NewOperationError(errors.CannotDeployApp,deploymentResult).WithParams(networkId, appRequest)
    }

    return deployedInstance, err
}

func (m *RestManager) Undeploy(networkId string, instanceId string) derrors.DaishoError {
    // Check the application is already deployed
    retrieved, err := m.appClient.GetInstance(networkId, instanceId)
    if retrieved == nil || err != nil {
        logger.Errorf("Instance %s not found with error [%s]", instanceId, err)
        return err
    }

    // Call the orchestrator to undeploy the instance
    err = m.orchestrator.Undeploy(networkId, *retrieved)
    if err != nil {
        logger.Errorf("Manager rest found an error when undeploying %s",err)
        return err
    }

    return nil
}

// Logs get a set of log entries from the selected application.
// params:
//   networkId identifier of the target network
//   instanceId identifier of the target application instance
// return:
//   Error if any.
//   An array of strings.
func (m *RestManager) Logs(networkId string, instanceId string) (*entitiesConductor.LogEntries, derrors.DaishoError) {
    // Check the application is already deployed
    retrieved, err := m.appClient.GetInstance(networkId, instanceId)
    if retrieved == nil || err != nil {
        logger.Errorf("Instance %s not found with error [%s]", retrieved, err)
        return nil, err
    }
    appmgr := m.asmClientFactory.CreateClient(retrieved.ClusterAddress, asm.SlavePort)
    pods, err := appmgr.Pods(retrieved.Name)
    if retrieved == nil || err != nil {
        return nil, err
    }
    return m.loggerClient.Logs(pods.Pods)
}

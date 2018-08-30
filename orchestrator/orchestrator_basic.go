//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// The orchestrator is in charge of controlling the execution of the applications deployed on top of the
// Daisho platform. This implementation deploys applications in the cluster with the smallest number of
// running solutions.

package orchestrator

import (
    "strings"
    "time"

    "github.com/daishogroup/conductor/asm"
    entitiesConductor "github.com/daishogroup/conductor/entities"
    "github.com/daishogroup/conductor/errors"
    "github.com/daishogroup/derrors"
    smclient "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
)

const (
    // Constant value indicating undeployed applications
    UNDEPLOYED = "UNDEPLOYED"
    DEPLOYED   = "DEPLOYED"
)

type BasicOrchestrator struct {
    clusterClient    smclient.Cluster
    nodesClient      smclient.Node
    appClient        smclient.Applications
    asmClientFactory asm.ClientFactory
    // cluster index for every cluster type
    lastClusterId map[entities.ClusterType]int
}

// Create a new basic orchestrator
// params:
//   clusterClient client to access system model clusters
//   nodesClient client to access nodes
//   appClient client to access the applications
// returns:
//   instance of an orchestrator
func NewBasicOrchestrator(clusterClient smclient.Cluster, nodesClient smclient.Node, appClient smclient.Applications,
    asmClientFactory asm.ClientFactory) Orchestrator {
    return &BasicOrchestrator{clusterClient,
        nodesClient,
        appClient,
        asmClientFactory,
        map[entities.ClusterType]int{entities.CloudType: 0, entities.GatewayType: 0, entities.EdgeType: 0}}
}

// createInstanceSM creates a new Application Instance inside the system model.
// params:
//   networkId Identifier of the target network
//   descriptor of the application to be deployed
//   appRequest The application request.
// return:
//   application instance that has been created.
//   error if any
func (o * BasicOrchestrator) createInstanceSM(
    networkId string,
    descriptor entities.AppDescriptor,
    appRequest entitiesConductor.DeployAppRequest) (* entities.AppInstance, derrors.DaishoError) {
    arguments := o.removeEdgeAddresses(appRequest.Arguments)
    instance, err := o.appClient.AddApplicationInstance(networkId, *entities.NewAddAppInstanceRequest(
        descriptor.ID, appRequest.Name, appRequest.Description, appRequest.Label,
        arguments, "1Gb", appRequest.StorageType))

    if err != nil {
        logger.Errorf("There was an error storing the application information %s", err)
        return instance, derrors.NewOperationError(errors.CannotStoreApp).
            WithParams(networkId, appRequest).CausedBy(err)
    }
    return instance, nil
}

// Deploy an application following some orchestration solution.
// params:
//   networkId Identifier of the target network
//   descriptor of the application to be deployed
//   appRequest The application request.
// return:
//   instance representing the deployed application
//   error if any
func (o *BasicOrchestrator) Deploy(
    networkId string,
    descriptor entities.AppDescriptor,
    appRequest entitiesConductor.DeployAppRequest) (*entities.AppInstance, derrors.DaishoError) {
    logger.Debug("Called deploy")
    // List of available clusters
    clustersList, err := o.clusterClient.ListByNetwork(networkId)
    logger.Debugf("Available clusters %s", clustersList)
    if err != nil {
        return nil, derrors.NewOperationError(errors.OpFail).CausedBy(err)
    }

    if clustersList == nil || len(clustersList) == 0 {
        return nil, derrors.NewOperationError(errors.ClustersNotAvailable)
    }

    // Get the target cluster type
    // TODO: Deprecate label field in favor of Labels, and use other field if required for targetClusterType.
    targetClusterType := entities.ClusterType(appRequest.Label)

    // Get the list of available clusters for that type
    listTargetClusters := make([] entities.Cluster, 0)
    for _, c := range clustersList {
        if c.Type == targetClusterType {
            listTargetClusters = append(listTargetClusters, c)
        }
    }

    if listTargetClusters == nil || len(listTargetClusters) == 0 {
        return nil, derrors.NewOperationError(errors.ClustersNotAvailable).WithParams(targetClusterType)
    }

    nextCluster := o.lastClusterId[targetClusterType] % len(clustersList)
    logger.Debugf("Next cluster to deploy is", nextCluster)

    instance, err:= o.createInstanceSM(networkId, descriptor, appRequest)
    if err != nil {
        return nil, err
    }

    // cluster were we have finally deployed the app
    deployedCluster := UNDEPLOYED
    deployedAddress := UNDEPLOYED
    attempts := 0
    // store the last observed error to return it in case of not available cluster found
    var lastError derrors.DaishoError
    for attempts < len(targetClusterType) && deployedCluster == UNDEPLOYED {
        targetCluster := listTargetClusters[nextCluster]

        logger.Debugf("Check cluster %s", targetCluster.ID)

        // Check what are we supposed to do depending on the target
        switch {
        case targetCluster.Cordon:
            logger.Debugf("Cluster %s is cordoned", targetCluster.ID)
        case targetCluster.Drain:
            logger.Debugf("Cluster %s is drained", targetCluster.ID)
        case targetCluster.Status == entities.ClusterInstalled:
            if targetCluster.Type == targetClusterType {
                logger.Debugf("Cluster %s is a candidate for deployment", targetCluster.ID)
                // Check this cluster
                var targetNode *entities.Node = nil
                targetNode, appRequest, err = o.deployAppIntoNode(
                    networkId, targetCluster, descriptor, appRequest, *instance)
                if err == nil && targetNode != nil {
                    deployedCluster = targetCluster.ID
                    deployedAddress = targetNode.PublicIP
                } else {
                    lastError = err
                }
            } else {
                logger.Debugf("Cluster %s is %s and we look for %s clusters. Ignored.", targetCluster.ID,
                    targetCluster.Type, targetClusterType)
            }
        default:
            logger.Debugf("Cluster %s has been ignored as a candidate for deployment", targetCluster.ID)
        }

        // Increase the number of attempts
        attempts++
        nextCluster = (nextCluster + 1) % len(listTargetClusters)
    }

    if deployedCluster != UNDEPLOYED {
        o.lastClusterId[targetClusterType] = nextCluster
        return o.finishAndUpdate(networkId, instance.DeployedID, deployedAddress, deployedCluster)
    }

    // We exit in an error state
    logger.Errorf("Impossible to find an available cluster to deploy %s", descriptor.ServiceName)
    o.markInstanceError(*instance)
    if lastError != nil {
        return nil, derrors.NewOperationError(errors.ClustersNotAvailable).CausedBy(lastError).WithParams(descriptor)
    }

    return nil, derrors.NewOperationError(errors.ClustersNotAvailable).WithParams(descriptor)
}

func (o * BasicOrchestrator) markInstanceError(instance entities.AppInstance) {
    updateRequest := entities.NewUpdateAppInstRequest().WithStatus(entities.AppInstError)
    _, err := o.appClient.UpdateInstance(instance.NetworkID, instance.DeployedID, *updateRequest)
    if err != nil {
        logger.Warn("Unable to assign error state to instance: " + instance.String())
    }
}

// buildLabelMap constructs the map of labels passed to appmgr so that every element inside k8s of that
// application is properly labeled. This enables other components such as the app monitoring system to
// identify the target instance to be updated.
func (o * BasicOrchestrator) buildLabelMap(networkID string, clusterID string, deployedID string) map[string]string {
    labels := make(map[string]string, 0)
    labels["networkID"] = networkID
    labels["clusterID"] = clusterID
    labels["deployedID"] = deployedID
    return labels
}

// Check if there is a node in the given cluster to deploy an application.
// params:
//    networkId target network
//    cluster cluster to check
//    descriptor application descriptor
//    appRequest client.DeployAppRequest
// return:
//    node chosen to deploy the application
//    request after possible arguments modifications
//    error if any
func (o *BasicOrchestrator) deployAppIntoNode(
    networkId string, targetCluster entities.Cluster,
    descriptor entities.AppDescriptor,
    appRequest entitiesConductor.DeployAppRequest,
    appInstance entities.AppInstance) (*entities.Node, entitiesConductor.DeployAppRequest, derrors.DaishoError) {
    availableNodes, err := o.nodesClient.List(networkId, targetCluster.ID)
    logger.Debugf("List of available nodes in cluster %s is %s", targetCluster.ID, availableNodes)
    if err != nil {
        return nil, appRequest, err
    }
    if len(availableNodes) == 0 {
        logger.Errorf("No available nodes at cluster %s", targetCluster.ID)
    }
    appRequest.Labels = o.buildLabelMap(networkId, targetCluster.ID, appInstance.DeployedID)
    logger.Debugf("Incoming appRequest is: %s", appRequest)
    // Last observed error for some feedback
    var lastError derrors.DaishoError
    for _, n := range availableNodes {
        logger.Debugf("Check node %s", n)
        // connect with the asmcli in the node and start the application
        asmClient := o.asmClientFactory.CreateClient(n.PublicIP, asm.SlavePort)
        if asmClient == nil {
            logger.Errorf("Impossible to create client")
        } else {
            if appRequest.Label == string(entities.GatewayType) {
                // If we have an gateway application, we have to get the IPs for the edge clusters.
                edgeAddresses := o.getListAddresses(networkId, entities.EdgeType)
                // add the corresponding arguments
                // ....
                if edgeAddresses == nil {
                    logger.Errorf("No available edge addresses for app %s", appRequest.Name)
                }
                additionalArgs := "edgeAddress={" + strings.Join(edgeAddresses, ",") + "}"
                if len(appRequest.Arguments) != 0 {
                    appRequest.Arguments = additionalArgs + " " + appRequest.Arguments
                } else {
                    appRequest.Arguments = additionalArgs
                }
            }
            logger.Debugf("Request app to start at %s", n.PublicIP)
            err := asmClient.Start(descriptor, appRequest)

            if err != nil {
                logger.Errorf("Impossible to start app %s in %s", descriptor.Name, n.PublicIP)
                logger.Errorf("Error: [%s]",err)
                lastError = err
            } else {
                // already deployed
                logger.Debugf("Application %s deployed in cluster %s", descriptor.Name, targetCluster.ID)
                return &n, appRequest, nil
            }
        }
    }
    return nil, appRequest, derrors.NewOperationError(errors.NodeNotAvailable).CausedBy(lastError).WithParams(targetCluster)

}

func (o *BasicOrchestrator) removeEdgeAddresses(arguments string) string {
    tokens := strings.Split(arguments, " ")
    newString := make([] string, 0)
    for _, token := range tokens {
        if strings.Contains(token, "edgeAddress") {

        } else {
            newString = append(newString, token)
        }
    }
    return strings.Join(newString, " ")
}

// This function is particularly designed for demo purposes.
// It takes the addresses of all the edge clusters and returns the resulting list.
// params:
//    networkId target network
//    clusterType type of cluster we are looking for
// returns:
//    List of addresses of edge clusters
func (o *BasicOrchestrator) getListAddresses(networkId string, clusterType entities.ClusterType) []string {
    clusters, err := o.clusterClient.ListByNetwork(networkId)
    if err != nil {
        logger.Errorf("Error checking the available clusters [%s]", err)
        return nil
    }

    var addresses [] string = nil
    // Now find the nodes for edge clusters
    for _, c := range clusters {
        if c.Type == clusterType && !c.Cordon && !c.Drain {
            nodes, err := o.nodesClient.List(networkId, c.ID)
            if err != nil {
                logger.Errorf("Error checking nodes at cluster %s\n%s", c.ID, err)
                return nil
            }

            for _, n := range nodes {
                added := false
                if n.Installed {
                    addresses = append(addresses, n.PublicIP)
                    added = true
                } else {
                    logger.Debugf("Node %s is not installed yet.", n.Name)
                }
                if added {
                    break
                }
            }
        }
    }
    return addresses
}

// This function is called after an application has been successfully deployed onto a cluster. It updates the
// information contained into the system model.
// params:
//   networkId identifier of the target network were the application has been deployed
//   descriptor Application descriptor instance
//   request
//   deployedCluster id of the cluster were the application was deployed
//   deployedAddress address of the deployment
// return:
//   application instance with all the related information
//   error if any
func (o *BasicOrchestrator) finishAndUpdate(
    networkId string,
    deployedID string,
    deployedAddress string,
    clusterID string) (*entities.AppInstance, derrors.DaishoError) {

    // update the corresponding info
    logger.Debugf("Updating instance %s deployed on %s", deployedID, deployedAddress)
    update := entities.NewUpdateAppInstRequest().WithClusterAddress(deployedAddress).WithClusterID(clusterID)
    updated, err := o.appClient.UpdateInstance(networkId, deployedID, * update)

    if err != nil {
        // TODO the application was deployed but the information was not correctly updated... what shall we do?
        // but the information was not correctly stored.
        logger.Errorf("There was an error storing the cluster information %s", err)
        return nil, derrors.NewOperationError(errors.CannotUpdateApp).WithParams(networkId, update).CausedBy(err)
    }

    return updated, nil
}

func (o * BasicOrchestrator) DeleteFromSM(networkID string, deployedID string) derrors.DaishoError {
    // remove the instance from the system model as it must be an error
    errDelete := o.appClient.DeleteInstance(networkID, deployedID)
    if errDelete != nil {
        logger.Errorf("there was an error when removing instance %s", deployedID)
        return derrors.NewOperationError(errors.CannotDeleteApp).CausedBy(errDelete).WithParams(networkID, deployedID)
    }
    return nil
}


func (o *BasicOrchestrator) Undeploy(networkId string, appInstance entities.AppInstance) derrors.DaishoError {
    logger.Debugf("called undeploy for %s instance %s", networkId, appInstance)

    if appInstance.ClusterAddress == "" {
        logger.Warnf("App %s has no cluster address, deleting from SM", appInstance.DeployedID)
        // Application was never deployed
        return o.DeleteFromSM(networkId, appInstance.DeployedID)
    }

    asmClient := o.asmClientFactory.CreateClient(appInstance.ClusterAddress, asm.SlavePort)
    if asmClient == nil {
        logger.Errorf("error when connecting with %s", appInstance.ClusterAddress)
        return derrors.NewOperationError(errors.ConnectionError)
    }

    // Check if the application exists, if not just delete it from system model.
    exists, checkErr := o.checkAppRunning(asmClient, appInstance.NetworkID, appInstance.DeployedID)
    if checkErr != nil {
        return checkErr
    }
    logger.Debugf("Application %s is running: %t", appInstance.DeployedID, exists)
    if !exists {
        logger.Warnf("App %s is not deployed on cluster %s, deleting from SM", appInstance.DeployedID, appInstance.ClusterAddress)
        return o.DeleteFromSM(networkId, appInstance.DeployedID)
    }


    attempts := 0
    stopped := false
    var err derrors.DaishoError
    // wait until the application has been stopped
    for !stopped && attempts < 60 {
        // Call asm client to force undeploy
        logger.Infof("Stopped before says %b",stopped)
        stopped, err = asmClient.Stop(appInstance.Name)

        logger.Debugf("requesting application to undeploy says %t with error [%s]", stopped, err)

        if err != nil {
            logger.Debugf("there was an error when undeploying %s", appInstance.Name)
            return derrors.NewOperationError(errors.CannotStopApp, err).WithParams(networkId, appInstance)
        }
        if !stopped {
            // Sleep a fixed period (5 seconds)
            time.Sleep(time.Millisecond * 5000)
            attempts = attempts + 1
        }
    }

    if !stopped {
        logger.Errorf("We could not undeploy %s", errors.CannotStopApp)
        return derrors.NewOperationError(errors.CannotStopApp).WithParams(networkId, appInstance)
        //return derrors.NewOperationError(errors.CannotStopApp)
    }

    // If everything was OK, remove the instance from the system model
    return o.DeleteFromSM(networkId, appInstance.DeployedID)
}

func (o * BasicOrchestrator) checkAppRunning(asmClient asm.Client, networkId string, deployedId string) (bool, derrors.DaishoError) {
    appList, err := asmClient.List()
    if err != nil {
        return false, err
    }
    logger.Debugf("Running apps %d", len(appList.Applications))
    for _, app := range appList.Applications {
        if o.checkLabels(app.Labels, networkId, deployedId) {
            return true, nil
        }
    }
    return false, nil
}

func (o * BasicOrchestrator) checkLabels(labels map[string]string, networkID string, deployedID string) bool {
    nID, exist := labels["networkID"]
    if !exist {
        return false
    }
    dID, exist := labels["deployedID"]
    if !exist {
        return false
    }
    return strings.Compare(nID, networkID) == 0 && strings.Compare(dID, deployedID) == 0
}
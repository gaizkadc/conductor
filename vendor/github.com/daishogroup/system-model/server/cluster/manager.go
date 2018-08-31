//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the cluster manager in charge of the business logic behind cluster entities.

package cluster

import (

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
)

// The Manager struct provides access to cluster related methods.
type Manager struct {
    networkProvider networkstorage.Provider
    clusterProvider clusterstorage.Provider
    appInstProvider appinststorage.Provider
}

// NewManager creates a new cluster manager.
//   params:
//     networkProvider The network storage provider.
//     clusterProvider The cluster storage provider.
//   returns:
//     A manager.
func NewManager(
    networkProvider networkstorage.Provider,
    clusterProvider clusterstorage.Provider,
    appInstProvider appinststorage.Provider) Manager {
    return Manager{networkProvider,clusterProvider, appInstProvider}
}

// AddCluster adds a new cluster to an existing network.
//   params:
//     networkID The target network identifier.
//     cluster The cluster to be added.
//   returns:
//     The added cluster.
//     An error if the network cannot be added.
func (mgr * Manager) AddCluster(networkID string, request entities.AddClusterRequest) (* entities.Cluster, derrors.DaishoError) {

    cluster := entities.ToCluster(networkID, request)

    if !mgr.networkProvider.Exists(networkID){
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if mgr.clusterProvider.Exists(cluster.ID){
        return nil, derrors.NewOperationError(errors.ClusterAlreadyExists).WithParams(request)
    }

    err := mgr.clusterProvider.Add(*cluster)
    if err == nil{
        err := mgr.networkProvider.AttachCluster(networkID, cluster.ID)
        if err == nil{
            return cluster, nil
        }
        return nil, err
    }
    return nil, err
}

// ListClusters lists the clusters inside a given network.
//   params:
//     networkID The target network identifier.
//   returns:
//     An array of clusters.
//     An error if the clusters cannot be retrieved.
func (mgr * Manager) ListClusters(networkID string) ([]entities.Cluster, derrors.DaishoError){
    clusterIDs, err := mgr.networkProvider.ListClusters(networkID)
    if err == nil {
        clusters := make([]entities.Cluster, 0, len(clusterIDs))
        failed := false
        for index := 0; index < len(clusterIDs) && !failed; index++{
            toAdd, err := mgr.clusterProvider.RetrieveCluster(clusterIDs[index])
            if err == nil{
                clusters = append(clusters, *toAdd)
            }else{
                failed = true
            }
        }
        if !failed {
            return clusters, nil
        }
        return make([]entities.Cluster, 0, 0), derrors.NewOperationError(errors.OpFail)

    }
    return nil, err
}

// GetCluster retrieves a cluster.
//   params:
//     networkId The target network identifier
//     clusterID The cluster identifier.
//   returns:
//     A cluster.
//     An error if the cluster cannot be retrieved or is not associated with the network.
func (mgr * Manager) GetCluster(networkID string, clusterID string) (* entities.Cluster, derrors.DaishoError) {
    if mgr.networkProvider.ExistsCluster(networkID, clusterID) {
        return mgr.clusterProvider.RetrieveCluster(clusterID)
    }
    return nil, derrors.NewOperationError(errors.ClusterNotAttachedToNetwork).WithParams(networkID, clusterID)
}

// UpdateCluster updates an existing cluster.
//   params:
//     networkID The network identifier.
//     clusterID The cluster identifier.
//     update The update cluster request.
//   returns:
//     The updated cluster.
//     An error if the instance cannot be update.
func (mgr * Manager) UpdateCluster(networkID string, clusterID string,
    update entities.UpdateClusterRequest) (* entities.Cluster, derrors.DaishoError) {
    if !mgr.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if !mgr.clusterProvider.Exists(clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(clusterID)
    }
    if !mgr.networkProvider.ExistsCluster(networkID, clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterNotAttachedToNetwork).WithParams(networkID, clusterID)
    }

    previous, err := mgr.clusterProvider.RetrieveCluster(clusterID)
    if err != nil{
        return nil, err
    }

    updated:=previous.Merge(update)
    err = mgr.clusterProvider.Update(* updated)
    if err != nil {
        return nil, err
    }
    return updated, nil
}

// DeleteCluster deletes a cluster.
//  params:
//      networkID The network id.
//      clusterID the cluster id.
//  returns:
//      Error if any
func (mgr * Manager) DeleteCluster(networkID string, clusterID string) derrors.DaishoError {
    if !mgr.networkProvider.ExistsCluster(networkID, clusterID) {
        return derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(networkID, clusterID)
    }

    if !mgr.clusterProvider.Exists(clusterID) {
        return derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(clusterID)
    }

    // check nodes
    nodes, err := mgr.clusterProvider.ListNodes(clusterID)
    if err!=nil{
        return err
    }
    if nodes!=nil && len(nodes) !=0 {
        return derrors.NewOperationError(errors.InvalidCondition).WithParams("node")
    }


    // check apps
    associatedApps, err := mgr.clusterHasDeployedApps(networkID, clusterID)

    if err != nil{
        return err
    }
    if *associatedApps {
        return derrors.NewOperationError(errors.InvalidCondition).WithParams("instance")
    }

    err = mgr.clusterProvider.Delete(clusterID)
    if err == nil {
        return mgr.networkProvider.DeleteCluster(networkID, clusterID)
    }
    return err
}

func (mgr * Manager) clusterHasDeployedApps(networkID string, clusterID string) (*bool, derrors.DaishoError) {
    instances, err := mgr.networkProvider.ListAppInst(networkID)
    if err != nil {
        return nil, err
    }
    for _, instanceID := range instances {
        instance, err := mgr.appInstProvider.RetrieveInstance(instanceID)
        if err != nil {
            return nil, err
        }
        if instance.ClusterID == clusterID {
            b:=true
            return &b, nil
        }
    }
    r := false
    return &r, nil
}
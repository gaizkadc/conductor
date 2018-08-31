//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Cluster Client Mockup

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/server/cluster"
)

// ClusterMockup is the mockup of the cluster client.
type ClusterMockup struct {
    networkProvider *networkstorage.MockupNetworkProvider
    clusterProvider *clusterstorage.MockupClusterProvider
    appInstProvider * appinststorage.MockupAppInstProvider
    clusterMgr      cluster.Manager
}

// NewClusterMockup creates a mockup Cluster client.
//   returns:
//     The Cluster Client.
func NewClusterMockup() Cluster {
    networkProvider := networkstorage.NewMockupNetworkProvider()
    clusterProvider := clusterstorage.NewMockupClusterProvider()
    appInstProvider := appinststorage.NewMockupAppInstProvider()
    clusterMgr := cluster.NewManager(networkProvider, clusterProvider, appInstProvider)
    return &ClusterMockup{networkProvider, clusterProvider, appInstProvider, clusterMgr}
}

// ClearMockup clears the content of this mockup.
func (mockup *ClusterMockup) ClearMockup() {
    mockup.networkProvider.Clear()
    mockup.clusterProvider.Clear()
    mockup.appInstProvider.Clear()
}

// InitMockup initializes the content of the mockup.
func (mockup *ClusterMockup) InitMockup() {
    mockup.ClearMockup()
    testLocation := "Madrid, Spain"
    testAdminName := "Admin"
    testAdminEmail := "admin@admins.com"
    testAdminPhone := "1234"

    mockup.clusterProvider.Add(* entities.NewClusterWithID("1","1", "Cluster1", "Description Cluster 1",
        entities.GatewayType, testLocation,"admin@admin.com",
        entities.ClusterInstalled, false, false))
    mockup.clusterProvider.Add(* entities.NewClusterWithID("1","2", "Cluster2", "Description Cluster 2",
        entities.GatewayType, testLocation,"admin@admin.com",
        entities.ClusterInstalled, false, false))
    mockup.clusterProvider.Add(* entities.NewClusterWithID("2","3", "Cluster3", "Description Cluster 3",
        entities.GatewayType, testLocation,"admin@admin.com",
        entities.ClusterInstalled, false, false))

    mockup.networkProvider.Add(* entities.NewNetworkWithID("1", "Network1", "Description Network 1",
        testAdminName, testAdminPhone, testAdminEmail))
    mockup.networkProvider.Add(* entities.NewNetworkWithID("2", "Network2", "Description Network 2",
        testAdminName, testAdminPhone, testAdminEmail))

    mockup.networkProvider.AttachCluster("1", "1")
    mockup.networkProvider.AttachCluster("1", "2")
    mockup.networkProvider.AttachCluster("2", "3")
}

// Add a cluster to the network.
//   params:
//     networkId The network id.
//     entity The cluster entity.
//   returns:
//	   The added network.
//     Error, if there is an internal error.
func (mockup *ClusterMockup) Add(networkID string, entity entities.AddClusterRequest) (*entities.Cluster, derrors.DaishoError) {
    return mockup.clusterMgr.AddCluster(networkID, entity)
}

// ListByNetwork obtains the list of clusters by network.
//   params:
//     networkId The network id.
//   returns:
//     The list of clusters for the selected network.
//     Error, if there is an internal error.
func (mockup *ClusterMockup) ListByNetwork(networkID string) ([] entities.Cluster, derrors.DaishoError) {
    return mockup.clusterMgr.ListClusters(networkID)
}

// Get a selected cluster
//   params:
//     networkId The network id.
//     clusterId The cluster id.
//   returns:
//     The selected cluster.
//     Error, if there is an internal error.
func (mockup *ClusterMockup) Get(networkID string, clusterID string) (*entities.Cluster, derrors.DaishoError) {
    return mockup.clusterMgr.GetCluster(networkID, clusterID)
}

// Update a selected cluster
//   params:
//     networkId The network id.
//     clusterId The cluster id.
//     update The update request.
//   returns:
//     The updated cluster.
//     Error, if there is an internal error.
func (mockup *ClusterMockup) Update(networkID string, clusterID string,
    update entities.UpdateClusterRequest) (*entities.Cluster, derrors.DaishoError) {
    return mockup.clusterMgr.UpdateCluster(networkID, clusterID, update)

}

// Delete a cluster
//  params:
//      networkID The network id.
//      clusterID the cluster id.
//  returns:
//      Error if any
func (mockup *ClusterMockup) Delete(networkID string, clusterID string) derrors.DaishoError {
    return mockup.clusterMgr.DeleteCluster(networkID, clusterID)
}

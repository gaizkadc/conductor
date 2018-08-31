//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the cluster provider mockup using in-memory storage.

package clusterstorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupClusterProvider is a mockup implementation of the cluster provider.
type MockupClusterProvider struct {
    sync.Mutex
    // Clusters indexed by cluster identifier.
    clusters map[string]entities.Cluster

    // Array of cluster by network identifier.
    nodes map[string][] string
}

// NewMockupClusterProvider creates a new mockup provider.
func NewMockupClusterProvider() *MockupClusterProvider {
    return &MockupClusterProvider{
        clusters:make(map[string]entities.Cluster),
        nodes: make(map[string][] string)}
}

// Add a new cluster to the system.
//   params:
//     network The Cluster to be added
//   returns:
//     An error if the cluster cannot be added.
func (mockup *MockupClusterProvider) Add(cluster entities.Cluster) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(cluster.ID) {
        mockup.clusters[cluster.ID] = cluster
        return nil
    }
    return derrors.NewOperationError(errors.ClusterAlreadyExists).WithParams(cluster)
}

// Update a existing cluster in the provider.
//   params:
//     cluster The Cluster to be updated, the id of the cluster must be exist.
//   returns:
//     An error if the cluster cannot be edited.
func (mockup *MockupClusterProvider) Update(cluster entities.Cluster) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(cluster.ID) {
        mockup.clusters[cluster.ID] = cluster
        return nil
    }
    return derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(cluster)
}

// Exists checks if a cluster exists in the system.
//   params:
//     clusterID The cluster identifier.
//   returns:
//     Whether the cluster exists or not.
func (mockup *MockupClusterProvider) Exists(clusterID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExists(clusterID)
}

func (mockup *MockupClusterProvider) unsafeExists(clusterID string) bool {
    _, exists := mockup.clusters[clusterID]
    return exists
}

// RetrieveCluster retrieves a given cluster.
//   params:
//     clusterID The cluster identifier.
//   returns:
//     The cluster.
//     An error if the cluster cannot be retrieved.
func (mockup *MockupClusterProvider) RetrieveCluster(clusterID string) (*entities.Cluster, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    cluster, exists := mockup.clusters[clusterID]
    if exists {
        return &cluster, nil
    }
    return nil, derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(clusterID)
}

// Delete a given cluster.
//   params:
//     clusterID The cluster identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (mockup *MockupClusterProvider) Delete(clusterID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(clusterID){
        return derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(clusterID)
    }
    delete(mockup.clusters, clusterID)
    delete(mockup.nodes, clusterID)
    return nil
}

// AttachNode links a node to an existing node.
//   params:
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     An error if the node cannot be attached.
func (mockup *MockupClusterProvider) AttachNode(clusterID string, nodeID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(clusterID) {
        if !mockup.unsafeExistsNode(clusterID, nodeID) {
            nodes, _ := mockup.nodes[clusterID]
            mockup.nodes[clusterID] = append(nodes, nodeID)
            return nil
        }
        return derrors.NewOperationError(errors.NodeAlreadyAttachedToCluster).WithParams(clusterID, nodeID)
    }
    return derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(clusterID)
}

// ListNodes lists the nodes of a cluster.
//   params:
//     clusterID The cluster identifier.
//   returns:
//     An array of node identifiers.
//     An error if the nodes cannot be retrieved.
func (mockup *MockupClusterProvider) ListNodes(clusterID string) ([]string, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    nodes, ok := mockup.nodes[clusterID]
    if ok {
        return nodes, nil
    }
    return nil, nil
}

// ExistsNode checks if a node is associated with a given cluster.
//   params:
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     Whether the node is associated to the cluster.
func (mockup *MockupClusterProvider) ExistsNode(clusterID string, nodeID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExistsNode(clusterID, nodeID)
}

func (mockup *MockupClusterProvider) unsafeExistsNode(clusterID string, nodeID string) bool {
    nodes, ok := mockup.nodes[clusterID]
    if ok {
        for _, node := range nodes {
            if node == nodeID {
                return true
            }
        }
        return false
    }
    return false
}

// DeleteNode deletes a node from an existing cluster.
//   params:
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     An error if the node cannot be removed.
func (mockup *MockupClusterProvider) DeleteNode(clusterID string, nodeID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExistsNode(clusterID, nodeID) {
        previousNodes := mockup.nodes[clusterID]
        newNodes := make([] string, 0, len(previousNodes)-1)
        for _, nID := range previousNodes {
            if nID != nodeID {
                newNodes = append(newNodes, nID)
            }
        }
        mockup.nodes[clusterID] = newNodes
        return nil
    }
    return derrors.NewOperationError(errors.NodeNotAttachedToCluster).WithParams(clusterID, nodeID)
}

// Dump obtains the list of all clusters in the system.
//   returns:
//     The list of clusters.
//     An error if the list cannot be retrieved.
func (mockup * MockupClusterProvider) Dump() ([] entities.Cluster, derrors.DaishoError) {
    result := make([] entities.Cluster, 0)
    mockup.Lock()
    defer mockup.Unlock()
    for _, cluster := range mockup.clusters {
        result = append(result, cluster)
    }
    return result, nil
}

// ReducedInfoList get a list with the reduced info.
//   returns:
//     List of the reduced info.
//     An error if the info cannot be retrieved.
func (mockup *MockupClusterProvider) ReducedInfoList() ([] entities.ClusterReducedInfo, derrors.DaishoError){
    result := make([] entities.ClusterReducedInfo, 0, len(mockup.clusters))
    mockup.Lock()
    defer mockup.Unlock()
    for _, c := range mockup.clusters {
        reducedInfo := entities.NewClusterReducedInfo(c.NetworkID,c.ID,c.Name,c.Type)
        result = append(result, *reducedInfo)
    }
    return result, nil
}

// Clear is a util function to clear the contents of the mockup.
func (mockup *MockupClusterProvider) Clear() {
    mockup.Lock()
    mockup.clusters = make(map[string]entities.Cluster)
    mockup.nodes = make(map[string][] string)
    mockup.Unlock()
}

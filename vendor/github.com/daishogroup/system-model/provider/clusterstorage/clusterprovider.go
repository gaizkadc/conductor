//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the cluster provider specification.

package clusterstorage

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Provider is the interface of the Cluster provider.
type Provider interface {
    // Add a new cluster to the system.
    //   params:
    //     cluster The Cluster to be added
    //   returns:
    //     An error if the cluster cannot be added.
    Add(cluster entities.Cluster) derrors.DaishoError

    // Update a existing cluster in the provider.
    //   params:
    //     cluster The Cluster to be updated, the id of the cluster must be exist.
    //   returns:
    //     An error if the cluster cannot be edited.
    Update(cluster entities.Cluster) derrors.DaishoError

    // Check if a cluster exists in the system.
    //   params:
    //     clusterID The cluster identifier.
    //   returns:
    //     Whether the cluster exists or not.
    Exists(clusterID string) bool

    // Retrieve a given cluster.
    //   params:
    //     clusterID The cluster identifier.
    //   returns:
    //     The cluster.
    //     An error if the cluster cannot be retrieved.
    RetrieveCluster(clusterID string) (*entities.Cluster, derrors.DaishoError)

    // Delete a given cluster.
    //   params:
    //     clusterID The cluster identifier.
    //   returns:
    //     An error if the cluster cannot be removed.
    Delete(clusterID string) derrors.DaishoError

    // Attach a node to an existing node.
    //   params:
    //     clusterID    The cluster identifier.
    //     nodeID       The node identifier.
    //   returns:
    //     An error if the node cannot be attached.
    AttachNode(clusterID string, nodeID string) derrors.DaishoError

    // List the nodes of a cluster.
    //   params:
    //     clusterID The cluster identifier.
    //   returns:
    //     An array of node identifiers.
    //     An error if the nodes cannot be retrieved.
    ListNodes(clusterID string) ([]string, derrors.DaishoError)

    // Check if a node is associated with a given cluster.
    //   params:
    //     clusterID    The cluster identifier.
    //     nodeID       The node identifier.
    //   returns:
    //     Whether the node is associated to the cluster.
    ExistsNode(clusterID string, nodeID string) bool

    // Delete a node from an existing cluster.
    //   params:
    //     clusterID    The cluster identifier.
    //     nodeID       The node identifier.
    //   returns:
    //     An error if the node cannot be removed.
    DeleteNode(clusterID string, nodeID string) derrors.DaishoError

    // Dump obtains the list of all clusters in the system.
    //   returns:
    //     The list of clusters.
    //     An error if the list cannot be retrieved.
    Dump() ([] entities.Cluster, derrors.DaishoError)

    // ReducedInfoList get a list with the reduced info.
    //   returns:
    //     List of the reduced info.
    //     An error if the info cannot be retrieved.
    ReducedInfoList() ([] entities.ClusterReducedInfo, derrors.DaishoError)
}

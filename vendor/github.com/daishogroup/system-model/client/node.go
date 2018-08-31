//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Node Client interface

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Node is a interface that represent the client for the node entity.
type Node interface {
    // Add a new node to an existing cluster.
    //   params:
    //     networkID    The target network identifier.
    //     clusterID    The target cluster identifier.
    //     node         The node to be added.
    //   returns:
    //     The added node.
    //     An error if the node cannot be added.
    Add(networkID string, clusterID string, node entities.AddNodeRequest) (*entities.Node, derrors.DaishoError)

    // List the nodes inside a given cluster.
    //   params:
    //     networkID The target network identifier.
    //     clusterID The target cluster identifier.
    //   returns:
    //     An array of nodes.
    //     An error if the nodes cannot be retrieved.
    List(networkID string, clusterID string) ([] entities.Node, derrors.DaishoError)

    // Remove a node.
    //   params:
    //     networkID    The target network identifier
    //     clusterID    The cluster identifier.
    //     nodeID       The node identifier.
    //   returns:
    //     An error if the node cannot be removed.
    Remove(networkID string, clusterID string, nodeID string) derrors.DaishoError

    // Get a node.
    //   params:
    //     networkID    The target network identifier
    //     clusterID    The cluster identifier.
    //     nodeID       The node identifier.
    //   returns:
    //     A node.
    //     An error if the node cannot be retrieved or is not associated with the cluster.
    Get(networkID string, clusterID string, nodeID string) (*entities.Node, derrors.DaishoError)

    // Update an existing node.
    //   params:
    //     networkID The network identifier.
    //     clusterID The cluster identifier.
    //     nodeID The node identifier.
    //     update The update node request.
    //   returns:
    //     The updated node.
    //     An error if the instance cannot be update.
    Update(networkID string, clusterID string, nodeID string, update entities.UpdateNodeRequest) (* entities.Node, derrors.DaishoError)

    // FilterNodes filters the set of nodes in a cluster using a set of restrictions.
    //   params:
    //     networkID The target network identifier.
    //     clusterID The target cluster identifier.
    //     filter The filtering constraints.
    //   returns:
    //     An array of nodes.
    //     An error if the nodes cannot be retrieved.
    FilterNodes(networkID string, clusterID string, filter entities.FilterNodesRequest) ([] entities.Node, derrors.DaishoError)
}

// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the node provider specification.

package nodestorage

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Provider definition of the node persistence-related methods.
type Provider interface {

    // Add a new node to the system.
    //   params:
    //     node The Node to be added
    //   returns:
    //     An error if the node cannot be added.
    Add(node entities.Node) derrors.DaishoError

    // Check if a node exists in the system.
    //   params:
    //     nodeID The node identifier.
    //   returns:
    //     Whether the node exists or not.
    Exists(nodeID string) bool

    // Retrieve a given cluster.
    //   params:
    //     nodeID The node identifier.
    //   returns:
    //     The cluster.
    //     An error if the node cannot be retrieved.
    RetrieveNode(nodeID string) (* entities.Node, derrors.DaishoError)

    // Delete a given node.
    //   params:
    //     nodeID The node identifier.
    //   returns:
    //     An error if the node cannot be removed.
    Delete(nodeID string) derrors.DaishoError

    // Update a node in the system.
    //   params:
    //     node The new node information.
    //   returns:
    //     An error if the node cannot be updated.
    Update(node entities.Node) derrors.DaishoError

    // ReducedInfoList get a list with the reduced info.
    //   returns:
    //     List of the reduced info.
    //     An error if the info cannot be retrieved.
    ReducedInfoList() ([] entities.NodeReducedInfo, derrors.DaishoError)

    // Dump obtains the list of all nodes in the system.
    //   returns:
    //     The list of nodes.
    //     An error if the list cannot be retrieved.
    Dump() ([] entities.Node, derrors.DaishoError)
}
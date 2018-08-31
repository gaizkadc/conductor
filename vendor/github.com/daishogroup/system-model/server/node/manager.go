//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the node manager in charge of the business logic behind node entities.

package node

import (
    "sort"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
)

// The Manager struct provides access to cluster related methods.
type Manager struct {
    networkProvider networkstorage.Provider
    clusterProvider clusterstorage.Provider
    nodeProvider    nodestorage.Provider
}

// NewManager creates a new node manager.
//   params:
//     networkProvider  The network storage provider.
//     clusterProvider  The cluster storage provider.
//     nodeProvider     The node storage provider.
//   returns:
//     A manager.
func NewManager(networkProvider networkstorage.Provider, clusterProvider clusterstorage.Provider,
    nodeProvider nodestorage.Provider) Manager {
    return Manager{networkProvider, clusterProvider, nodeProvider}
}

// AddNode adds a new node to an existing cluster.
//   params:
//     networkID    The target network identifier.
//     clusterID    The target cluster identifier.
//     node         The node to be added.
//   returns:
//     The added node.
//     An error if the node cannot be added.
func (manager *Manager) AddNode(networkID string, clusterID string,
    request entities.AddNodeRequest) (*entities.Node, derrors.DaishoError) {

    node := entities.ToNode(networkID, clusterID, request)

    if !manager.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if !manager.clusterProvider.Exists(clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(networkID, clusterID)
    }
    if !manager.networkProvider.ExistsCluster(networkID, clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterNotAttachedToNetwork).WithParams(networkID, clusterID)
    }

    if manager.nodeProvider.Exists(node.ID) {
        return nil, derrors.NewOperationError(errors.NodeAlreadyExists).WithParams(node)
    }

    // Add edgenet network ID to node, we need this for authorization
    // later and we don't want to have to add the network client to the
    // node manager in budo.
    network, err := manager.networkProvider.RetrieveNetwork(networkID)
    if err != nil {
        return nil, err
    }
    node.EdgenetNetworkID = network.EdgenetID

    err = manager.nodeProvider.Add(*node)
    if err == nil {
        err := manager.clusterProvider.AttachNode(clusterID, node.ID)
        if err == nil {
            return node, nil
        }
        return nil, err
    }
    return nil, err
}

// ListNodes lists the nodes inside a given cluster.
//   params:
//     networkID The target network identifier.
//     clusterID The target cluster identifier.
//   returns:
//     An array of nodes.
//     An error if the nodes cannot be retrieved.
func (manager *Manager) ListNodes(networkID string, clusterID string) ([] entities.Node, derrors.DaishoError) {
    if !manager.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if !manager.clusterProvider.Exists(clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(networkID, clusterID)
    }
    if !manager.networkProvider.ExistsCluster(networkID, clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterNotAttachedToNetwork).WithParams(networkID, clusterID)
    }


    nodeIds, err := manager.clusterProvider.ListNodes(clusterID)
    if err == nil {
        nodes := make([]entities.Node, 0, len(nodeIds))
        failed := false
        for index := 0; index < len(nodeIds) && !failed; index++ {
            toAdd, err := manager.nodeProvider.RetrieveNode(nodeIds[index])
            if err == nil {
                nodes = append(nodes, *toAdd)
            } else {
                failed = true
            }
        }

        // Sort by public IP for a good chance of having one of the master
        // nodes first. This only works as intended in a limited set of case
        // and a better solution is needed. See DP-1329.
        sort.Slice(nodes, func(i, j int) bool { return nodes[i].PublicIP < nodes[j].PublicIP})

        if !failed {
            return nodes, nil
        }
        return nil, derrors.NewOperationError(errors.OpFail).WithParams(networkID, clusterID)

    }
    return [] entities.Node{}, err
}

// RemoveNode deletes a node.
//   params:
//     networkID    The target network identifier
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     An error if the node cannot be removed.
func (manager *Manager) RemoveNode(networkID string, clusterID string, nodeID string) derrors.DaishoError {
    if !manager.networkProvider.Exists(networkID) {
        return derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if !manager.clusterProvider.Exists(clusterID) {
        return derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(networkID, clusterID)
    }
    if !manager.networkProvider.ExistsCluster(networkID, clusterID) {
        return derrors.NewOperationError(errors.ClusterNotAttachedToNetwork).WithParams(networkID, clusterID)
    }

    if manager.clusterProvider.ExistsNode(clusterID, nodeID) {
        err := manager.clusterProvider.DeleteNode(clusterID, nodeID)
        if err == nil {
            return manager.nodeProvider.Delete(nodeID)
        }
        return err
    }
    return derrors.NewOperationError(errors.NodeNotAttachedToCluster).WithParams(networkID, clusterID, nodeID)
}

// GetNode retrieves a given node.
//   params:
//     networkID    The target network identifier
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     A node.
//     An error if the node cannot be retrieved or is not associated with the cluster.
func (manager *Manager) GetNode(networkID string, clusterID string, nodeID string) (*entities.Node, derrors.DaishoError) {
    if !manager.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if !manager.clusterProvider.Exists(clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(networkID, clusterID)
    }
    if !manager.networkProvider.ExistsCluster(networkID, clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterNotAttachedToNetwork).WithParams(networkID, clusterID)
    }
    if !manager.clusterProvider.ExistsNode(clusterID, nodeID) {
        return nil, derrors.NewOperationError(errors.NodeNotAttachedToCluster).WithParams(networkID, clusterID, nodeID)
    }
    return manager.nodeProvider.RetrieveNode(nodeID)
}

// UpdateNode updates an existing node.
//   params:
//     networkID The network identifier.
//     clusterID The cluster identifier.
//     nodeID The node identifier.
//     update The update node request.
//   returns:
//     The updated node.
//     An error if the instance cannot be update.
func (manager *Manager) UpdateNode(networkID string, clusterID string, nodeID string,
    update entities.UpdateNodeRequest) (*entities.Node, derrors.DaishoError) {
    if !manager.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if !manager.clusterProvider.Exists(clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(networkID, clusterID)
    }
    if !manager.nodeProvider.Exists(nodeID) {
        return nil, derrors.NewOperationError(errors.NodeDoesNotExists).WithParams(networkID, clusterID, nodeID)
    }
    if !manager.networkProvider.ExistsCluster(networkID, clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterNotAttachedToNetwork).WithParams(networkID, clusterID)
    }
    if !manager.clusterProvider.ExistsNode(clusterID, nodeID) {
        return nil, derrors.NewOperationError(errors.NodeNotAttachedToCluster).WithParams(networkID, clusterID, nodeID)
    }

    previous, err := manager.nodeProvider.RetrieveNode(nodeID)
    if err != nil {
        return nil, err
    }

    updated := previous.Merge(update)
    err = manager.nodeProvider.Update(* updated)
    if err != nil {
        return nil, err
    }
    return updated, nil
}

// FilterNodes filters the set of nodes in a cluster using a set of restrictions.
//   params:
//     networkID The target network identifier.
//     clusterID The target cluster identifier.
//     filter The filtering constraints.
//   returns:
//     An array of nodes.
//     An error if the nodes cannot be retrieved.
func (manager *Manager) FilterNodes(networkID string, clusterID string, filter entities.FilterNodesRequest) ([] entities.Node, derrors.DaishoError) {
    if !manager.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if !manager.clusterProvider.Exists(clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(networkID, clusterID)
    }
    if !manager.networkProvider.ExistsCluster(networkID, clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterNotAttachedToNetwork).WithParams(networkID, clusterID)
    }

    nodeIds, err := manager.clusterProvider.ListNodes(clusterID)
    if err == nil {
        nodes := make([]entities.Node, 0, len(nodeIds))
        failed := false
        for index := 0; index < len(nodeIds) && !failed; index++ {
            toAdd, err := manager.nodeProvider.RetrieveNode(nodeIds[index])
            if err == nil {
                nodes = append(nodes, *toAdd)
            } else {
                failed = true
            }
        }

        if failed {
            return nil, derrors.NewOperationError(errors.OpFail).WithParams(networkID, clusterID)
        }

        return manager.applyFilters(nodes, filter), nil
    }
    return [] entities.Node{}, err
}

func (manager *Manager) applyFilters(nodes []entities.Node, filter entities.FilterNodesRequest) [] entities.Node {
    result := make([]entities.Node, 0)
    for _, n := range nodes {
        if manager.matchFilter(n, filter) {
            result = append(result, n)
        }
    }
    return result
}

func (manager *Manager) matchFilter(node entities.Node, filter entities.FilterNodesRequest) bool {
    result := true
    if filter.Labels != nil && len(*(filter.Labels)) > 0 {
        labelMap := make(map[string]bool, len(node.Labels))
        for _, l := range node.Labels {
            labelMap[l] = true
        }
        for _, l := range *filter.Labels {
            _, found := labelMap[l]
            result = result && found
        }
    }
    return result
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Node Client Rest

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "fmt"
    "github.com/daishogroup/dhttp"
)

// NodeAddURI with the URI pattern to add nodes.
const NodeAddURI = "/api/v0/node/%s/%s/add"
// NodeListByNetworkURI with the URI patter to list the nodes.
const NodeListByNetworkURI = "/api/v0/node/%s/%s/list"
const NodeFilterURI = "/api/v0/node/%s/%s/filter"
// NodeGetURI with the URI pattern to retrieve the info of a node.
const NodeGetURI = "/api/v0/node/%s/%s/%s/info"
// NodeRemoveURI with the URI pattern to remove a node.
const NodeRemoveURI = "/api/v0/node/%s/%s/%s/delete"
// NodeUpdateURI with the URI pattern to update a node.
const NodeUpdateURI = "/api/v0/node/%s/%s/%s/update"

// NodeRest structure with the REST client.
type NodeRest struct {
    client dhttp.Client
}

// NewNodeRest creates a REST Node client.
// Deprecated: Use NewNodeClientRest
func NewNodeRest(basePath string) Node {
    return NewNodeClientRest(ParseHostPort(basePath))
}

func NewNodeClientRest(host string, port int) Node {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host, port))
    return &NodeRest{rest}
}

// Add a new node to an existing cluster.
//   params:
//     networkID    The target network identifier.
//     clusterID    The target cluster identifier.
//     node         The node to be added.
//   returns:
//     The added node.
//     An error if the node cannot be added.
func (rest *NodeRest) Add(networkID string, clusterID string, node entities.AddNodeRequest) (*entities.Node, derrors.DaishoError) {
    response := rest.client.Post(fmt.Sprintf(NodeAddURI, networkID, clusterID), node, new(entities.Node))
    if response.Error != nil {
        return nil, response.Error
    }
    result :=response.Result.(*entities.Node)
    return  result, nil
}

// List the nodes inside a given cluster.
//   params:
//     networkID The target network identifier.
//     clusterID The target cluster identifier.
//   returns:
//     An array of nodes.
//     An error if the nodes cannot be retrieved.
func (rest *NodeRest) List(networkID string, clusterID string) ([] entities.Node, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(NodeListByNetworkURI, networkID, clusterID), new([] entities.Node))
    if response.Error != nil {
        return nil, response.Error
    }
    ns:=response.Result.(*[] entities.Node)
    return *ns, nil
}

// Remove a node.
//   params:
//     networkID    The target network identifier
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     An error if the node cannot be removed.
func (rest *NodeRest) Remove(networkID string, clusterID string, nodeID string) derrors.DaishoError {
    response := rest.client.Delete(fmt.Sprintf(NodeRemoveURI, networkID, clusterID, nodeID),
        new(entities.SuccessfulOperation))
    return response.Error
}

// Get a node.
//   params:
//     networkID    The target network identifier
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     A node.
//     An error if the node cannot be retrieved or is not associated with the cluster.
func (rest *NodeRest) Get(networkID string, clusterID string, nodeID string) (*entities.Node, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(NodeGetURI, networkID, clusterID, nodeID), new(entities.Node))
    if response.Error != nil {
        return nil, response.Error
    }
    n:=response.Result.(* entities.Node)
    return n, nil
}

// Update an existing node.
//   params:
//     networkID The network identifier.
//     clusterID The cluster identifier.
//     nodeID The node identifier.
//     update The update node request.
//   returns:
//     The updated node.
//     An error if the instance cannot be update.
func (rest *NodeRest) Update(networkID string, clusterID string, nodeID string, update entities.UpdateNodeRequest) (*entities.Node, derrors.DaishoError) {
    response := rest.client.Post(fmt.Sprintf(NodeUpdateURI, networkID, clusterID, nodeID), update, new(entities.Node))
    if response.Ok() {
        n := response.Result.(*entities.Node)
        return n, nil
    }
    return nil, response.Error
}

// FilterNodes filters the set of nodes in a cluster using a set of restrictions.
//   params:
//     networkID The target network identifier.
//     clusterID The target cluster identifier.
//     filter The filtering constraints.
//   returns:
//     An array of nodes.
//     An error if the nodes cannot be retrieved.
func (rest *NodeRest) FilterNodes(networkID string, clusterID string, filter entities.FilterNodesRequest) ([] entities.Node, derrors.DaishoError){
    response := rest.client.GetWithBody(fmt.Sprintf(NodeFilterURI, networkID, clusterID), filter, new([] entities.Node))
    if response.Error != nil {
        return nil, response.Error
    }
    ns:=response.Result.(*[] entities.Node)
    return *ns, nil
}
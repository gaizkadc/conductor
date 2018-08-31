//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Node Client Mockup

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/server/node"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/entities"
)

// NodeMockup structure with the mockup client.
type NodeMockup struct {
    networkProvider *networkstorage.MockupNetworkProvider
    clusterProvider *clusterstorage.MockupClusterProvider
    nodeProvider    *nodestorage.MockupNodeProvider
    nodeMgr         node.Manager
}

// NewNodeMockup creates a mockup Node client.
func NewNodeMockup() Node {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var clusterProvider = clusterstorage.NewMockupClusterProvider()
    var nodeProvider = nodestorage.NewMockupNodeProvider()
    var nodeMgr = node.NewManager(networkProvider, clusterProvider, nodeProvider)

    return &NodeMockup{networkProvider, clusterProvider, nodeProvider, nodeMgr}
}

// ClearNodeMockup clears the content of the mockup.
func (mockup *NodeMockup) ClearNodeMockup() {
    mockup.networkProvider.Clear()
    mockup.clusterProvider.Clear()
    mockup.nodeProvider.Clear()
}

// InitNodeMockup initializes the mockup.
func (mockup *NodeMockup) InitNodeMockup() {
    mockup.ClearNodeMockup()

    testIP := "0.0.0.0"
    testUsername := "username"
    testPassword := "pass"
    testSSHKey := "ABCDFGHIJK"
    testLocation := "Madrid, Spain"
    testAdminName := "Admin"
    testAdminEmail := "admin@admins.com"
    testAdminPhone := "1234"

    edgenetNode := entities.NewNodeWithID("1", "1", "1",
        "Node1", "Description Node 1", make([]string, 0),
        testIP, testIP, true, testUsername, testPassword, testSSHKey, entities.NodeUnchecked)
    edgenetNode.EdgenetNetworkID = "beefdead00123456"
    mockup.nodeProvider.Add(*edgenetNode)

    mockup.nodeProvider.Add(* entities.NewNodeWithID("1", "1", "2",
        "Node2", "Description Node 2", make([]string, 0),
        testIP, testIP, true, testUsername, testPassword, testSSHKey, entities.NodeUnchecked))

    mockup.nodeProvider.Add(* entities.NewNodeWithID("1", "2", "3",
        "Node3", "Description Node 3", make([]string, 0),
        testIP, testIP, true, testUsername, testPassword, testSSHKey, entities.NodeUnchecked))

    mockup.clusterProvider.Add(* entities.NewClusterWithID("1", "1", "Cluster1",
        "Description Cluster 1",
        entities.GatewayType, testLocation, "admin@admin.com",
        entities.ClusterCreated, false, false))
    mockup.clusterProvider.Add(* entities.NewClusterWithID("1", "2", "Cluster2",
        "Description Cluster 2",
        entities.GatewayType, testLocation, "admin@admin.com",
        entities.ClusterCreated, false, false))
    mockup.clusterProvider.Add(* entities.NewClusterWithID("2", "3", "Cluster3",
        "Description Cluster 3",
        entities.GatewayType, testLocation, "admin@admin.com",
        entities.ClusterCreated, false, false))

    mockup.clusterProvider.AttachNode("1", "1")
    mockup.clusterProvider.AttachNode("1", "2")
    mockup.clusterProvider.AttachNode("2", "3")

    mockup.networkProvider.Add(* entities.NewNetworkWithID("1", "Network1",
        "Description Network 1",
        testAdminName, testAdminPhone, testAdminEmail))
    mockup.networkProvider.Add(* entities.NewNetworkWithID("2", "Network2",
        "Description Network 2",
        testAdminName, testAdminPhone, testAdminEmail))

    mockup.networkProvider.AttachCluster("1", "1")
    mockup.networkProvider.AttachCluster("1", "2")
    mockup.networkProvider.AttachCluster("2", "3")
}

// Add a new node to an existing cluster.
//   params:
//     networkID    The target network identifier.
//     clusterID    The target cluster identifier.
//     node         The node to be added.
//   returns:
//     The added node.
//     An error if the node cannot be added.
func (mockup *NodeMockup) Add(networkID string, clusterID string,
    node entities.AddNodeRequest) (*entities.Node, derrors.DaishoError) {
    return mockup.nodeMgr.AddNode(networkID, clusterID, node)
}

// List the nodes inside a given cluster.
//   params:
//     networkID The target network identifier.
//     clusterID The target cluster identifier.
//   returns:
//     An array of nodes.
//     An error if the nodes cannot be retrieved.
func (mockup *NodeMockup) List(networkID string,
    clusterID string) ([] entities.Node, derrors.DaishoError) {
    return mockup.nodeMgr.ListNodes(networkID, clusterID)
}

// Remove a node.
//   params:
//     networkID    The target network identifier
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     An error if the node cannot be removed.
func (mockup *NodeMockup) Remove(networkID string,
    clusterID string, nodeID string) derrors.DaishoError {
    return mockup.nodeMgr.RemoveNode(networkID, clusterID, nodeID)
}

// Get a node.
//   params:
//     networkID    The target network identifier
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     A node.
//     An error if the node cannot be retrieved or is not associated with the cluster.
func (mockup *NodeMockup) Get(networkID string,
    clusterID string, nodeID string) (*entities.Node, derrors.DaishoError) {
    return mockup.nodeMgr.GetNode(networkID, clusterID, nodeID)
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
func (mockup *NodeMockup) Update(networkID string, clusterID string, nodeID string,
    update entities.UpdateNodeRequest) (*entities.Node, derrors.DaishoError) {
    return mockup.nodeMgr.UpdateNode(networkID, clusterID, nodeID, update)
}

// FilterNodes filters the set of nodes in a cluster using a set of restrictions.
//   params:
//     networkID The target network identifier.
//     clusterID The target cluster identifier.
//     filter The filtering constraints.
//   returns:
//     An array of nodes.
//     An error if the nodes cannot be retrieved.
func (mockup *NodeMockup) FilterNodes(networkID string, clusterID string, filter entities.FilterNodesRequest) ([] entities.Node, derrors.DaishoError){
    return mockup.nodeMgr.FilterNodes(networkID, clusterID, filter)
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Node Client TestSuite

package client

import (
    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/entities"
)

// TestGetExistingNode checks that an existing node can be retrieved.
func TestGetExistingNode(suite *suite.Suite, node Node) {
    n, err := node.Get("1", "1", "1")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(n, "node must not nil")
}

// TestGetNotExistingNode checks that a non existing node cannot be retrieved.
func TestGetNotExistingNode(suite *suite.Suite, node Node) {
    n, err := node.Get("1", "1", "3")
    suite.NotNil(err, "error must not be Nil")
    suite.Nil(n, "node must be nil")
}

// TestGetNodeNotExistingNetwork checks the result of requesting a node from a non
// existing network.
func TestGetNodeNotExistingNetwork(suite *suite.Suite, node Node) {
    n, err := node.Get("3", "1", "1")
    suite.NotNil(err, "error must not be Nil")
    suite.Nil(n, "node must be nil")
}

// TestGetNodeNotExistingCluster checks the result of requesting a node from a non
// existing cluster.
func TestGetNodeNotExistingCluster(suite *suite.Suite, node Node) {
    n, err := node.Get("2", "2", "1")
    suite.NotNil(err, "error must not be Nil")
    suite.Nil(n, "node must be nil")
}

// TestAddNode checks that a new node can be added and retrieved.
func TestAddNode(suite *suite.Suite, node Node) {
    toAdd := entities.NewAddNodeRequest("Cluster 4", "Description Cluster 4", []string{"l1"},
        "0.0.0.0", "0.0.0.0", true, "username", "password", "")
    addedNode, err := node.Add("1", "1", * toAdd)

    suite.Nil(err, "error must be Nil")
    suite.NotNil(addedNode, "added node must not be Nil")

    n, err := node.Get("1","1", addedNode.ID)
    suite.Nil(err, "error must be Nil")
    suite.NotNil(n, "node must not nil")
}

// TestAddNodeToNotExistingNetwork checks the result of adding a node to a non
// existing network.
func TestAddNodeToNotExistingNetwork(suite *suite.Suite, node Node) {

    toAdd := entities.NewAddNodeRequest("Cluster 5", "Description Cluster 5", []string{"l1"},
        "0.0.0.0", "0.0.0.0", true, "username", "password", "")
    addedNode, err := node.Add("3", "1", * toAdd)

    suite.NotNil(err, "error must be not Nil")
    suite.Nil(addedNode, "node must be Nil")
}

// TestAddNodeToNotExistingCluster checks the result of adding a node to a non
// existing cluster.
func TestAddNodeToNotExistingCluster(suite *suite.Suite, node Node) {
    toAdd := entities.NewAddNodeRequest("Cluster 5", "Description Cluster 5", []string{"l1"},
        "0.0.0.0", "0.0.0.0", true, "username", "password", "")
    addedNode, err := node.Add("1", "3", * toAdd)

    suite.NotNil(err, "error must be not Nil")
    suite.Nil(addedNode, "node must be Nil")
}

// TestListNodesByCluster checks that nodes of a cluster can be retrieved.
func TestListNodesByCluster(suite *suite.Suite, node Node) {
    ns, err := node.List("1","1")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(ns, "node must be not nil")
    suite.Len(ns, 2, "len must be 2")
}

// TestListNodesByClusterWithoutNodes checks the result of listing the nodes
// on a cluster without nodes.
func TestListNodesByClusterWithoutNodes(suite *suite.Suite, node Node) {
    ns, err := node.List("2","3")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(ns, "node must be not nil")
    suite.Len(ns, 0, "len must be 0")
}

// TestListNodesByClusterAfterAddAndDeleteNode checks the result of listing the nodes
// after adding and deleting nodes.
func TestListNodesByClusterAfterAddAndDeleteNode(suite *suite.Suite, node Node) {
    toAdd := entities.NewAddNodeRequest("Cluster 7", "Description Cluster 7", []string{"l1"},
        "0.0.0.0", "0.0.0.0", true, "username", "password", "")
    addedNode, err := node.Add("2", "3", * toAdd)

    ns, err := node.List("2","3")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(ns, "node must be not nil")
    suite.Len(ns, 1, "len must be 1")

    errRemove := node.Remove("2","3", addedNode.ID)
    suite.Nil(errRemove, "error must be Nil")

    fns, err := node.List("2","3")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(fns, "node must be not nil")
    suite.Len(fns, 0, "len must be 0")
}

// TestListNodesByClusterAfterAddMultipleNodeAndDeleteNode performs the above test
// several times.
func TestListNodesByClusterAfterAddMultipleNodeAndDeleteNode(suite *suite.Suite, node Node) {

    toAdd1 := entities.NewAddNodeRequest("Cluster 7", "Description Cluster 7", []string{"l1"},
        "0.0.0.0", "0.0.0.0", true, "username", "password", "")
    addedNode1, err := node.Add("2", "3", * toAdd1)

    suite.Nil(err, "error must be Nil")
    suite.NotNil(addedNode1, "added node must not be Nil")

    toAdd2 := entities.NewAddNodeRequest("Cluster 8", "Description Cluster 8", []string{"l1"},
        "0.0.0.0", "0.0.0.0", true, "username", "password", "")
    addedNode2, err := node.Add("2", "3", * toAdd2)

    suite.Nil(err, "error must be Nil")
    suite.NotNil(addedNode2, "added node must not be Nil")

    toAdd3 := entities.NewAddNodeRequest("Cluster 9", "Description Cluster 9", []string{"l1"},
        "0.0.0.0", "0.0.0.0", true, "username", "password", "")
    addedNode3, err := node.Add("2", "3", * toAdd3)

    suite.Nil(err, "error must be Nil")
    suite.NotNil(addedNode3, "added node must not be Nil")

    ns, err := node.List("2","3")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(ns, "node must be not nil")
    suite.Len(ns, 3, "len must be 3")

    errRemove := node.Remove("2","3", addedNode1.ID)
    suite.Nil(errRemove, "error must be Nil")

    fns, err := node.List("2","3")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(fns, "node must be not nil")
    suite.Len(fns, 2, "len must be 2")
}

// TestListNodesByNotExistingNetwork checks the result of listing the nodes
// of a non-existing network.
func TestListNodesByNotExistingNetwork(suite *suite.Suite, node Node) {
    cs, err := node.List("3","1")
    suite.NotNil(err, "error must not be Nil")
    suite.Nil(cs, "cluster must be nil")
}

// TestListNodesByNotExistingCluster checks the result of listing the nodes
// of a non-existing cluster.
func TestListNodesByNotExistingCluster(suite *suite.Suite, node Node) {
    cs, err := node.List("1","3")
    suite.NotNil(err, "error must not be Nil")
    suite.Nil(cs, "cluster must be nil")
}

// TestDeleteNode checks that a node can be deleted.
func TestDeleteNode(suite *suite.Suite, node Node) {
    err := node.Remove("1","1","1")
    suite.Nil(err, "error must be Nil")
    n, err := node.Get("1", "1", "1")
    suite.NotNil(err, "error must not be Nil")
    suite.Nil(n, "node must be nil")
}

// TestDeleteNodeNotExisting checks the result of deleting a non existing node.
func TestDeleteNodeNotExisting(suite *suite.Suite, node Node) {
    err := node.Remove("1","1","3")
    suite.NotNil(err, "error must not be Nil")
}

// TestDeleteNodeNotExistingNetwork checks the result of deleting a node from a non
// existing network.
func TestDeleteNodeNotExistingNetwork(suite *suite.Suite, node Node) {
    err := node.Remove("3","1","1")
    suite.NotNil(err, "error must not be Nil")
}

// TestDeleteNodesNotExistingCluster checks the result of deleting a node from a non
// existing cluster.
func TestDeleteNodesNotExistingCluster(suite *suite.Suite, node Node) {
    err := node.Remove("1","3","1")
    suite.NotNil(err, "error must not be Nil")
}

// TestUpdate checks tha a node can be updated.
func TestUpdate(suite * suite.Suite, node Node) {
    // Add a new node
    toAdd := entities.NewAddNodeRequest("New node", "Description", []string{"l1"},
        "0.0.0.0", "0.0.0.0", true, "username", "password", "")
    addedNode, err := node.Add("1", "1", * toAdd)

    suite.Nil(err, "error must be Nil")
    suite.NotNil(addedNode, "added node must not be Nil")
    // Update the information
    update := entities.NewUpdateNodeRequest().WithName("New node name").WithStatus(entities.NodeReadyToInstall)
    updated, err := node.Update("1", "1", addedNode.ID, * update)
    suite.Nil(err, "update must not fail")
    suite.NotNil(updated, "expecting updated node")
    suite.Equal(addedNode.ID, updated.ID, "Expecting same node id")
    suite.Equal("New node name", updated.Name)
    suite.Equal(entities.NodeReadyToInstall, updated.Status)
}


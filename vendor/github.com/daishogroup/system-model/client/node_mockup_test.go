//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Node Client Mockup tests

package client

import (
    "github.com/stretchr/testify/suite"
    "testing"
)

type NodeMockupTestSuite struct {
    suite.Suite
    node Node
}

func (suite *NodeMockupTestSuite) SetupSuite() {
    suite.node = NewNodeMockup()
}

func (suite *NodeMockupTestSuite) SetupTest() {
    suite.node.(*NodeMockup).InitNodeMockup()
}

func (suite *NodeMockupTestSuite) TestGetExistingNode() {
    TestGetExistingNode(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestGetNotExistingNode() {
    TestGetNotExistingNode(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestGetNodeNotExistingNetwork() {
    TestGetNodeNotExistingNetwork(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestGetNodeNotExistingCluster() {
    TestGetNodeNotExistingCluster(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestAddNode() {
    TestAddNode(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestAddNodeToNotExistingNetwork() {
    TestAddNodeToNotExistingNetwork(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestAddNodeToNotExistingCluster() {
    TestAddNodeToNotExistingCluster(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestListNodesByClusterWithoutNodes() {
    TestListNodesByClusterWithoutNodes(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestListNodesByClusterAfterAddAndDeleteNode() {
    TestListNodesByClusterAfterAddAndDeleteNode(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestListNodesByClusterAfterAddMultipleNodeAndDeleteNode(){
    TestListNodesByClusterAfterAddMultipleNodeAndDeleteNode(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestListNodesByCluster() {
    TestListNodesByCluster(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestListNodesByNotExistingNetwork() {
    TestListNodesByNotExistingNetwork(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestListNodesByNotExistingCluster() {
    TestListNodesByNotExistingCluster(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestDeleteNode() {
    TestDeleteNode(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestDeleteNodeNotExisting() {
    TestDeleteNodeNotExisting(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestDeleteNodeNotExistingNetwork() {
    TestDeleteNodeNotExistingNetwork(&suite.Suite, suite.node)
}

func (suite *NodeMockupTestSuite) TestDeleteNodesNotExistingCluster() {
    TestDeleteNodesNotExistingCluster(&suite.Suite, suite.node)
}

func (suite * NodeMockupTestSuite) TestUpdate() {
    TestUpdate(&suite.Suite, suite.node)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNodeMockup(t *testing.T) {
    suite.Run(t, new(NodeMockupTestSuite))
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Node Client integration tests

package client_external

import (
    "testing"

    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/dhttp"
)

type NodeExternalTestSuite struct {
    suite.Suite
    helper EndpointHelper
    node client.Node
}

func (suite *NodeExternalTestSuite) SetupSuite() {
    suite.helper = NewEndpointHelper()
    suite.node = client.NewNodeClientRest(suite.helper.GetListeningAddress())
    suite.helper.Start()
    dhttp.WaitURLAvailable(BaseAddress,suite.helper.port,5,"/", 1)
}

func (suite *NodeExternalTestSuite) SetupTest() {
    suite.helper.ResetProvider()
}

func (suite *NodeExternalTestSuite) TearDownSuite() {
    suite.helper.Shutdown()
}

func (suite *NodeExternalTestSuite) TestGetExistingNode() {
    client.TestGetExistingNode(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestGetNotExistingNode() {
    client.TestGetNotExistingNode(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestGetNodeNotExistingNetwork() {
    client.TestGetNodeNotExistingNetwork(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestGetNodeNotExistingCluster() {
    client.TestGetNodeNotExistingCluster(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestAddNode() {
    client.TestAddNode(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestAddNodeToNotExistingNetwork() {
    client.TestAddNodeToNotExistingNetwork(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestAddNodeToNotExistingCluster() {
    client.TestAddNodeToNotExistingCluster(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestListNodesByCluster() {
    client.TestListNodesByCluster(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestListNodesByClusterWithoutNodes() {
    client.TestListNodesByClusterWithoutNodes(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestListNodesByClusterAfterAddAndDeleteNode() {
    client.TestListNodesByClusterAfterAddAndDeleteNode(&suite.Suite, suite.node)
}

func (suite *NodeExternalTestSuite) TestListNodesByClusterAfterAddMultipleNodeAndDeleteNode(){
    client.TestListNodesByClusterAfterAddMultipleNodeAndDeleteNode(&suite.Suite, suite.node)
}

func (suite *NodeExternalTestSuite) TestListNodesByNotExistingNetwork() {
    client.TestListNodesByNotExistingNetwork(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestListNodesByNotExistingCluster() {
    client.TestListNodesByNotExistingCluster(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestDeleteNode() {
    client.TestDeleteNode(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestDeleteNodeNotExisting() {
    client.TestDeleteNodeNotExisting(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestDeleteNodeNotExistingNetwork() {
    client.TestDeleteNodeNotExistingNetwork(&suite.Suite,suite.node)
}

func (suite *NodeExternalTestSuite) TestDeleteNodesNotExistingCluster() {
    client.TestDeleteNodesNotExistingCluster(&suite.Suite,suite.node)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNodeExternal(t *testing.T) {
    suite.Run(t, new(NodeExternalTestSuite))
}
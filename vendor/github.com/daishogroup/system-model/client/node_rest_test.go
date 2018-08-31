//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Node Client Rest Test

package client

import (
    "testing"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/dhttp"
)

type NodeRestTestSuite struct {
    suite.Suite
    node Node
    rest *dhttp.ClientMockup
}

func (suite *NodeRestTestSuite) SetupSuite() {
    suite.rest = dhttp.NewClientMockup()
    suite.node = &NodeRest{suite.rest}
}

func (suite *NodeRestTestSuite) SetupTest() {
    suite.rest.Reset()
}

func (suite *NodeRestTestSuite) TestGetExistingNode() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := entities.NewNodeWithID("1", "1", "1",
            "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeInstalled)
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestGetExistingNode(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestGetNotExistingNode() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.NodeDoesNotExists))
    })
    TestGetNotExistingNode(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestGetNodeNotExistingNetwork() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.NetworkDoesNotExists))
    })
    TestGetNodeNotExistingNetwork(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestGetNodeNotExistingCluster() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.ClusterDoesNotExists))
    })
    TestGetNodeNotExistingCluster(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestAddNode() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewNodeWithID("1", "1", "3",
            "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked)
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    suite.rest.AddGet(func(path string) dhttp.Response {
        result := entities.NewNodeWithID("1", "1", "3",
            "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked)
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestAddNode(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestAddNodeToNotExistingNetwork() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.NetworkDoesNotExists))
    })
    TestAddNodeToNotExistingNetwork(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestAddNodeToNotExistingCluster() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.ClusterDoesNotExists))
    })
    TestAddNodeToNotExistingCluster(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestListNodesByCluster() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := [] entities.Node{
            *entities.NewNodeWithID("1", "1", "2",
                "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked),
            *entities.NewNodeWithID("1", "1", "3",
                "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked),
        }
        statusCode := 200
        return dhttp.NewResponse(&result, &statusCode, nil)
    })
    TestListNodesByCluster(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestListNodesByClusterWithoutNodes() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := []entities.Node{}
        statusCode := 200
        return dhttp.NewResponse(&result, &statusCode, nil)
    })
    TestListNodesByClusterWithoutNodes(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestListNodesByClusterAfterAddAndDeleteNode() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewNodeWithID("1", "1", "3",
            "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked)
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := []entities.Node{
            *entities.NewNodeWithID("1", "1", "3",
                "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked),
        }
        statusCode := 200
        return dhttp.NewResponse(&result, &statusCode, nil)
    })
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 200
        result := entities.NewSuccessfulOperation("DeleteNode")
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := []entities.Node{}
        statusCode := 200
        return dhttp.NewResponse(&result, &statusCode, nil)
    })
    TestListNodesByClusterAfterAddAndDeleteNode(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestListNodesByClusterAfterAddMultipleNodeAndDeleteNode() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewNodeWithID("1", "1", "3",
            "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked)
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewNodeWithID("1", "1", "4",
            "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked)
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewNodeWithID("1", "1", "5",
            "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked)
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := []entities.Node{
            *entities.NewNodeWithID("1", "1", "3",
                "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked),
            *entities.NewNodeWithID("1", "1", "4",
                "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked),
            *entities.NewNodeWithID("1", "1", "5",
                "Node1", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked),
        }
        statusCode := 200
        return dhttp.NewResponse(&result, &statusCode, nil)
    })
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 200
        result := entities.NewSuccessfulOperation("DeleteNode")
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := []entities.Node{
            *entities.NewNodeWithID("1", "1", "4",
                "Node4", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked),
            *entities.NewNodeWithID("1", "1", "5",
                "Node5", "Description1", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true,
            "user1", "pass1", "ssh1", entities.NodeUnchecked),
        }
        statusCode := 200
        return dhttp.NewResponse(&result, &statusCode, nil)
    })
    TestListNodesByClusterAfterAddMultipleNodeAndDeleteNode(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestListNodesByNotExistingNetwork() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.NetworkDoesNotExists))
    })
    TestListNodesByNotExistingNetwork(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestListNodesByNotExistingCluster() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.ClusterDoesNotExists))
    })
    TestListNodesByNotExistingCluster(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestDeleteNode() {
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 200
        result := entities.NewSuccessfulOperation("DeleteNode")
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.NodeDoesNotExists))
    })
    TestDeleteNode(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestDeleteNodeNotExisting() {
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.NodeDoesNotExists))
    })
    TestDeleteNodeNotExisting(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestDeleteNodeNotExistingNetwork() {
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.NetworkDoesNotExists))
    })
    TestDeleteNodeNotExistingNetwork(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestDeleteNodesNotExistingCluster() {
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.ClusterDoesNotExists))
    })
    TestDeleteNodesNotExistingCluster(&suite.Suite, suite.node)
}

func (suite *NodeRestTestSuite) TestUpdate() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewNodeWithID("1","1","id",
            "New node", "Description", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true, "username",
                "password", "", entities.NodeUnchecked)
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewNodeWithID("1","2","id",
            "New node name", "Description", make([]string, 0),
            "0.0.0.0", "0.0.0.0", true, "username",
                "password", "", entities.NodeReadyToInstall)
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestUpdate(&suite.Suite, suite.node)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNodeRest(t *testing.T) {
    suite.Run(t, new(NodeRestTestSuite))
}

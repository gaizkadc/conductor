//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Node manager tests.

package node

import (
    "testing"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/stretchr/testify/suite"
)

const (
    testNetworkName = "networkName"
    testClusterName = "clusterName"
    testNodeName    = "nodeName"
    testIP          = "0.0.0.0"
    testUsername    = "username"
    testPassword    = "password"
    testSSHKey      = "SSHKey"
    testDescription = "description"
    testAdminName   = "adminName"
    testAdminPhone  = "adminPhone"
    testAdminEmail  = "adminEmail"
    testLocation    = "Madrid, Spain"
)

type ManagerHelper struct {
    networkProvider *networkstorage.MockupNetworkProvider
    clusterProvider *clusterstorage.MockupClusterProvider
    nodeProvider    *nodestorage.MockupNodeProvider
    nodeMgr         Manager
}

func NewManagerHelper() ManagerHelper {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var clusterProvider = clusterstorage.NewMockupClusterProvider()
    var nodeProvider = nodestorage.NewMockupNodeProvider()
    var nodeMgr = NewManager(networkProvider, clusterProvider, nodeProvider)
    return ManagerHelper{networkProvider, clusterProvider, nodeProvider,
        nodeMgr}
}

type ManagerTestSuite struct {
    suite.Suite
    manager ManagerHelper
}

func (helper *ManagerTestSuite) SetupSuite() {
    managerHelper := NewManagerHelper()
    helper.manager = managerHelper
}

// The SetupTest method is called before every test on the suite.
func (helper *ManagerTestSuite) SetupTest() {
    helper.manager.networkProvider.Clear()
    helper.manager.clusterProvider.Clear()
    helper.manager.nodeProvider.Clear()
}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(ManagerTestSuite))
}

func (helper *ManagerTestSuite) addTestingCluster(networkID string, clusterID string) {
    network := entities.NewNetworkWithID(
        networkID, testNetworkName, testDescription, testAdminName, testAdminPhone, testAdminEmail)
    cluster := entities.NewClusterWithID(networkID, clusterID, testClusterName, testDescription, entities.CloudType,
        testLocation, testAdminEmail, entities.ClusterCreated, false, false)
    helper.manager.networkProvider.Add(*network)
    helper.manager.networkProvider.AttachCluster(networkID, clusterID)
    helper.manager.clusterProvider.Add(*cluster)

}

func (helper *ManagerTestSuite) TestAddNode() {
    networkID := "testAddNetwork"
    clusterID := "testAddCluster"
    helper.addTestingCluster(networkID, clusterID)
    newNode := entities.NewAddNodeRequest(testNodeName, testDescription, []string{"l1"},
        testIP, testIP, true, testUsername, testPassword, testSSHKey)
    added, err := helper.manager.nodeMgr.AddNode(networkID, clusterID, *newNode)
    helper.Nil(err, "node should be added")
    helper.NotNil(added, "node should be returned")
}

func (helper *ManagerTestSuite) TestRetrieveNode() {
    networkID := "RetrieveNetwork"
    clusterID := "RetrieveCluster"
    helper.addTestingCluster(networkID, clusterID)
    newNode := entities.NewAddNodeRequest(testNodeName, testDescription, []string{"l1"},
        testIP, testIP, true, testUsername, testPassword, testSSHKey)
    added, err := helper.manager.nodeMgr.AddNode(networkID, clusterID, *newNode)
    helper.Nil(err, "node should be added")
    helper.NotNil(added, "node must not be nil")
    retrieved, err := helper.manager.nodeMgr.GetNode(networkID,clusterID , added.ID)
    helper.Nil(err, "node should be retrieved")
    helper.EqualValues(added, retrieved, "nodes should match")
}

func (helper *ManagerTestSuite) TestListClusters() {
    networkID := "TestListNetwork"
    clusterID := "TestListClusters"
    helper.addTestingCluster(networkID, clusterID)
    numberNodes := 5
    for i := 0; i < numberNodes; i++ {
        newNode := entities.NewAddNodeRequest(testNodeName, testDescription,  []string{"l1"},
            testIP, testIP, true, testUsername, testPassword, testSSHKey)
        _, err := helper.manager.nodeMgr.AddNode(networkID, clusterID, *newNode)
        helper.Nil(err, "node should be added")
    }
    nodes, err := helper.manager.nodeMgr.ListNodes(networkID, clusterID)
    helper.Nil(err, "list should be retrieved")
    helper.NotNil(nodes, "expecting a list of nodes")
    helper.Len(nodes, numberNodes, "expecting 5 nodes")
}

func (helper *ManagerTestSuite) TestUpdateNode() {
    networkID := "testUpdateNetwork"
    clusterID := "testUpdateCluster"
    helper.addTestingCluster(networkID, clusterID)
    newNode := entities.NewAddNodeRequest(testNodeName, testDescription, []string{"l1"},
        testIP, testIP, true, testUsername, testPassword, testSSHKey)
    added, err := helper.manager.nodeMgr.AddNode(networkID, clusterID, *newNode)
    helper.Nil(err, "node should be added")
    helper.NotNil(added, "node should be returned")
    update := entities.NewUpdateNodeRequest().WithName("New node name")
    updated, err := helper.manager.nodeMgr.UpdateNode(networkID, clusterID, added.ID, * update)
    helper.Nil(err, "node should be updated")
    helper.Equal(added.ID, updated.ID)
    helper.Equal("New node name", updated.Name)
}

func (helper *ManagerTestSuite) TestUpdateToEmptyString() {
    networkID := "TestUpdateToEmptyString"
    clusterID := "TestUpdateToEmptyString"
    helper.addTestingCluster(networkID, clusterID)
    newNode := entities.NewAddNodeRequest(testNodeName, testDescription, []string{"l1"},
        testIP, testIP, true, testUsername, testPassword, testSSHKey)
    added, err := helper.manager.nodeMgr.AddNode(networkID, clusterID, *newNode)
    helper.Nil(err, "node should be added")
    helper.NotNil(added, "node should be returned")
    update := entities.NewUpdateNodeRequest().WithName("")
    updated, err := helper.manager.nodeMgr.UpdateNode(networkID, clusterID, added.ID, * update)
    helper.Nil(err, "node should be updated")
    helper.Equal(added.ID, updated.ID)
    helper.Equal("", updated.Name)
}

func (helper * ManagerTestSuite) TestDeleteNode() {
    networkID := "TestDeleteNode"
    clusterID := "TestDeleteNode"
    helper.addTestingCluster(networkID, clusterID)
    newNode := entities.NewAddNodeRequest(testNodeName, testDescription, []string{"l1"},
        testIP, testIP, true, testUsername, testPassword, testSSHKey)
    added, err := helper.manager.nodeMgr.AddNode(networkID, clusterID, *newNode)
    helper.Nil(err, "node should be added")
    helper.NotNil(added, "node should be returned")
    err = helper.manager.nodeMgr.RemoveNode(networkID, clusterID, added.ID)
    helper.Nil(err, "node should be removed")
    list, err := helper.manager.nodeMgr.ListNodes(networkID, clusterID)
    helper.Nil(err, "list should not fail")
    helper.Equal(0, len(list))
}

func (helper * ManagerTestSuite) TestFilterByLabels() {
    networkID := "TestFilterByLabels"
    clusterID := "TestFilterByLabels"
    helper.addTestingCluster(networkID, clusterID)
    newNode1 := entities.NewAddNodeRequest(testNodeName, testDescription, []string{"l1"},
        testIP, testIP, true, testUsername, testPassword, testSSHKey)
    added1, err := helper.manager.nodeMgr.AddNode(networkID, clusterID, *newNode1)
    helper.Nil(err, "node should be added")
    helper.NotNil(added1, "node should be returned")
    newNode2 := entities.NewAddNodeRequest(testNodeName, testDescription, []string{"l1", "l2"},
        testIP, testIP, true, testUsername, testPassword, testSSHKey)
    added2, err := helper.manager.nodeMgr.AddNode(networkID, clusterID, *newNode2)
    helper.Nil(err, "node should be added")
    helper.NotNil(added2, "node should be returned")
    newNode3 := entities.NewAddNodeRequest(testNodeName, testDescription, []string{"l3"},
        testIP, testIP, true, testUsername, testPassword, testSSHKey)
    added3, err := helper.manager.nodeMgr.AddNode(networkID, clusterID, *newNode3)
    helper.Nil(err, "node should be added")
    helper.NotNil(added3, "node should be returned")

    filterByL1 := entities.NewFilterNodesRequest().ByLabel("l1")
    filterByL1L2 := entities.NewFilterNodesRequest().ByLabels([]string{"l1","l2"})
    filterByL3 := entities.NewFilterNodesRequest().ByLabel("l3")
    filterNone := entities.NewFilterNodesRequest()

    filtered1, err := helper.manager.nodeMgr.FilterNodes(networkID, clusterID, *filterByL1)
    helper.Nil(err, "Filter should be applied")
    helper.Equal(2, len(filtered1), "nodes do not match")

    filtered2, err := helper.manager.nodeMgr.FilterNodes(networkID, clusterID, *filterByL1L2)
    helper.Nil(err, "Filter should be applied")
    helper.Equal(1, len(filtered2), "nodes do not match")
    helper.Equal(added2.ID, filtered2[0].ID , "nodes do not match")

    filtered3, err := helper.manager.nodeMgr.FilterNodes(networkID, clusterID, *filterByL3)
    helper.Nil(err, "Filter should be applied")
    helper.Equal(1, len(filtered3), "nodes do not match")
    helper.Equal(added3.ID, filtered3[0].ID , "nodes do not match")

    filteredNone, err := helper.manager.nodeMgr.FilterNodes(networkID, clusterID, *filterNone)
    helper.Nil(err, "Filter should be applied")
    helper.Equal(3, len(filteredNone), "nodes do not match")
}
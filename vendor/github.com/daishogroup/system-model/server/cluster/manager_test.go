//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Cluster manager tests.

package cluster

import (
    "testing"

    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"

    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
    "github.com/daishogroup/system-model/provider/appinststorage"
)

const (
    testNetworkName = "networkName"
    testClusterName = "clusterName"
    testDescription = "description"
    testAdminName   = "adminName"
    testAdminPhone  = "adminPhone"
    testAdminEmail  = "adminEmail"
    testLocation    = "Madrid, Spain"
)

type ManagerHelper struct {
    networkProvider *networkstorage.MockupNetworkProvider
    clusterProvider *clusterstorage.MockupClusterProvider
    appInstProvider *appinststorage.MockupAppInstProvider
    clusterMgr      Manager
}

func NewManagerHelper() ManagerHelper {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var clusterProvider = clusterstorage.NewMockupClusterProvider()
    var appInstProvider = appinststorage.NewMockupAppInstProvider()
    var clusterMgr = NewManager(networkProvider, clusterProvider, appInstProvider)
    return ManagerHelper{networkProvider, clusterProvider, appInstProvider, clusterMgr}
}

type TestHelper struct {
    suite.Suite
    manager ManagerHelper
}

func (helper *TestHelper) SetupSuite() {
    managerHelper := NewManagerHelper()
    helper.manager = managerHelper
}

// The SetupTest method is called before every test on the suite.
func (helper *TestHelper) SetupTest() {
    helper.manager.networkProvider.Clear()
    helper.manager.clusterProvider.Clear()
    helper.manager.appInstProvider.Clear()
}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(TestHelper))
}

func (helper *TestHelper) addTestingNetwork(id string) {
    var toAdd = entities.NewNetworkWithID(
        id, testNetworkName, testDescription, testAdminName, testAdminPhone, testAdminEmail)
    helper.manager.networkProvider.Add(*toAdd)
}

func (helper *TestHelper) TestAddCluster() {
    networkID := "testAddCluster"
    helper.addTestingNetwork(networkID)
    newCluster := entities.NewAddClusterRequest(testClusterName, testDescription, entities.CloudType,
        testLocation, testAdminEmail)
    added, err := helper.manager.clusterMgr.AddCluster(networkID, *newCluster)
    helper.Nil(err, "cluster should be added")
    helper.NotNil(added, "cluster should be returned")
}

func (helper *TestHelper) TestRetrieveCluster() {
    networkID := "RetrieveCluster"
    helper.addTestingNetwork(networkID)
    newCluster := entities.NewAddClusterRequest(testClusterName, testDescription,
        entities.CloudType, testLocation, testAdminEmail)

    added, err := helper.manager.clusterMgr.AddCluster(networkID, *newCluster)
    helper.Nil(err, "cluster should be added")
    helper.NotNil(added, "cluster should be returned")

    retrieved, err := helper.manager.clusterMgr.GetCluster(networkID, added.ID)
    helper.Nil(err, "cluster should be retrieved")
    helper.EqualValues(added, retrieved, "structs should match")
}

func (helper * TestHelper) TestUpdateCluster() {
    networkID := "TestUpdateCluster"
    helper.addTestingNetwork(networkID)
    newCluster := entities.NewAddClusterRequest(testClusterName, testDescription,
        entities.CloudType, testLocation, testAdminEmail)
    added, err := helper.manager.clusterMgr.AddCluster(networkID, *newCluster)
    helper.Nil(err, "cluster should be added")
    helper.NotNil(added, "cluster should be returned")

    update := entities.NewUpdateClusterRequest().WithClusterStatus(entities.ClusterInstalled)
    updated, err := helper.manager.clusterMgr.UpdateCluster(networkID, added.ID, * update)
    helper.Nil(err, "update should be processed")
    helper.NotNil(updated, "expecting updated cluster")
    helper.Equal(networkID, updated.NetworkID)
    helper.Equal(added.ID, updated.ID)
    helper.Equal(testClusterName, updated.Name)
    helper.Equal(testDescription, updated.Description)
    helper.Equal(entities.CloudType, updated.Type)
    helper.Equal(testLocation, updated.Location)
    helper.Equal(testAdminEmail, updated.Email)
    helper.Equal(entities.ClusterInstalled, updated.Status)
    helper.Equal(added.Drain, updated.Drain)
    helper.Equal(added.Cordon, updated.Cordon)

}

func (helper *TestHelper) TestDeleteCluster() {
    networkID := "TestDeleteCluster"
    helper.addTestingNetwork(networkID)
    newCluster := entities.NewAddClusterRequest(
        testClusterName, testDescription, entities.CloudType, testLocation,
        testAdminEmail)
    c, err := helper.manager.clusterMgr.AddCluster(networkID, *newCluster)
    assert.Nil(helper.T(), err, "cluster should be added")
    // update it to be deployed
    helper.manager.clusterMgr.UpdateCluster(networkID,c.ID,
        *entities.NewUpdateClusterRequest().WithClusterStatus(entities.ClusterInstalled))
    // delete it
    err = helper.manager.clusterMgr.DeleteCluster(networkID, c.ID)
    helper.Nil(err, "the cluster was not correctly removed")
    c, err = helper.manager.clusterMgr.GetCluster(networkID, c.ID)
    helper.NotNil(err, "cluster should not be retrieved")
    helper.Nil(c, "cluster should be nil")
    list, err := helper.manager.clusterMgr.ListClusters(networkID)
    helper.Nil(err, "list should be retrieved")
    helper.Equal(0, len(list))
}


func (helper *TestHelper) TestListClusters() {
    networkID := "TestListClusters"
    helper.addTestingNetwork(networkID)
    numberClusters := 5
    for i := 0; i < numberClusters; i++ {
        newCluster := entities.NewAddClusterRequest(
            testClusterName, testDescription, entities.CloudType, testLocation,
            testAdminEmail)
        _, err := helper.manager.clusterMgr.AddCluster(networkID, *newCluster)
        assert.Nil(helper.T(), err, "cluster should be added")
    }
    clusters, err := helper.manager.clusterMgr.ListClusters(networkID)
    helper.Nil(err, "list should be retrieved")
    helper.Equal(numberClusters, len(clusters), "expecting 5 clusters")
}

func (helper *TestHelper) TestUpdateNode() {
    networkID := "testUpdateNetwork"
    helper.addTestingNetwork(networkID)
    newCluster := entities.NewAddClusterRequest(testClusterName, testDescription,
        entities.CloudType, testLocation, testAdminEmail)
    added, err := helper.manager.clusterMgr.AddCluster(networkID, *newCluster)
    helper.Nil(err, "cluster should be added")
    helper.NotNil(added, "cluster should be returned")

    update := entities.NewUpdateClusterRequest().WithName("New cluster name")
    updated, err := helper.manager.clusterMgr.UpdateCluster(networkID, added.ID, * update)
    helper.Nil(err, "cluster should be updated")
    helper.Equal(added.ID, updated.ID)
    helper.Equal("New cluster name", updated.Name)
}

func (helper *TestHelper) TestUpdateToEmptyString() {
    networkID := "TestUpdateToEmptyString"

    helper.addTestingNetwork(networkID)
    newCluster := entities.NewAddClusterRequest(testClusterName, testDescription,
        entities.CloudType, testLocation, testAdminEmail)
    added, err := helper.manager.clusterMgr.AddCluster(networkID, *newCluster)
    helper.Nil(err, "cluster should be added")
    helper.NotNil(added, "cluster should be returned")

    update := entities.NewUpdateClusterRequest().WithName("")
    updated, err := helper.manager.clusterMgr.UpdateCluster(networkID, added.ID, * update)
    helper.Nil(err, "cluster should be updated")
    helper.Equal(added.ID, updated.ID)
    helper.Equal("", updated.Name)
}

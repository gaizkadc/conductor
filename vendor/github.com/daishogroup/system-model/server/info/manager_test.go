//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Dump manager tests

package info

import (
    "testing"

    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/provider/userstorage"
)

type ManagerHelper struct {
    networkProvider *networkstorage.MockupNetworkProvider
    clusterProvider *clusterstorage.MockupClusterProvider
    nodeProvider    *nodestorage.MockupNodeProvider
    appDescProvider *appdescstorage.MockupAppDescProvider
    appInstProvider *appinststorage.MockupAppInstProvider
    manager         Manager
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
    helper.manager.appDescProvider.Clear()
    helper.manager.appInstProvider.Clear()
}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(ManagerTestSuite))
}

func NewManagerHelper() ManagerHelper {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var clusterProvider = clusterstorage.NewMockupClusterProvider()
    var nodeProvider = nodestorage.NewMockupNodeProvider()
    var appDescProvider = appdescstorage.NewMockupAppDescProvider()
    var appInstProvider = appinststorage.NewMockupAppInstProvider()
    var userProvider = userstorage.NewMockupUserProvider()
    var dumpMgr = NewManager(networkProvider, clusterProvider, nodeProvider,
        appDescProvider, appInstProvider, userProvider)
    return ManagerHelper{networkProvider, clusterProvider, nodeProvider,
        appDescProvider, appInstProvider, dumpMgr}
}

func (helper *ManagerTestSuite) TestReducedInfo() {
    network := entities.NewNetwork("network1", "", "", "", "")
    cluster := entities.NewCluster(network.ID, "cluster1", "", entities.EdgeType,
        "", "", entities.ClusterInstalled, false, false)
    node := entities.NewNode(network.ID, cluster.ID,
        "node1", "", make([]string, 0),
        "ip0", "", false,
        "", "", "")
    descriptor := entities.NewAppDescriptor(network.ID, "descriptor1", "",
        "", "", "", 0, []string{"repo1:tag1"})
    instance := entities.NewAppInstance(network.ID, descriptor.ID, cluster.ID, "instance1", "", "",
        "", "", entities.AppStorageDefault, make([]entities.ApplicationPort, 0), 8888, "127.0.0.1")
    helper.manager.networkProvider.Add(* network)
    helper.manager.clusterProvider.Add(* cluster)
    helper.manager.nodeProvider.Add(* node)
    helper.manager.appDescProvider.Add(* descriptor)
    helper.manager.appInstProvider.Add(* instance)

    reduced, err := helper.manager.manager.ReducedInfo()
    helper.Nil(err, "expecting dump")
    helper.NotNil(reduced, "expecting dump")
    helper.Equal(1, len(reduced.Networks))
    helper.Equal(1, len(reduced.Clusters))
    helper.Equal(1, len(reduced.Nodes))
    helper.Equal(1, len(reduced.Descriptors))
    helper.Equal(1, len(reduced.Instances))
    helper.Equal("127.0.0.1",reduced.Instances[0].ClusterAddress)
    helper.Equal(8888,reduced.Instances[0].Port)
    helper.Equal("ip0",reduced.Nodes[0].PublicIP)

}

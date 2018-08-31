//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Dump manager tests

package dump

import (
    "testing"
    "time"

    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/accessstorage"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/provider/userstorage"
)

type ManagerHelper struct {
    networkProvider * networkstorage.MockupNetworkProvider
    clusterProvider * clusterstorage.MockupClusterProvider
    nodeProvider    * nodestorage.MockupNodeProvider
    appDescProvider * appdescstorage.MockupAppDescProvider
    appInstProvider * appinststorage.MockupAppInstProvider
    userProvider    * userstorage.MockupUserProvider
    accessProvider  * accessstorage.MockupUserAccessProvider
    dumpManager         Manager
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
    helper.manager.userProvider.Clear()
    helper.manager.accessProvider.Clear()
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
    var accessProvider = accessstorage.NewMockupUserAccessProvider()

    var dumpMgr = NewManager(networkProvider, clusterProvider, nodeProvider, appDescProvider, appInstProvider,
        userProvider, accessProvider)
    return ManagerHelper{networkProvider, clusterProvider, nodeProvider,
        appDescProvider, appInstProvider, userProvider,
        accessProvider,dumpMgr}
}

func (helper * ManagerTestSuite) TestDump() {
    network := entities.NewNetwork("network1", "", "", "", "")
    cluster := entities.NewCluster(network.ID, "cluster1", "", entities.EdgeType,
        "", "", entities.ClusterInstalled, false, false)
    node := entities.NewNode(network.ID, cluster.ID,
        "node1", "", make([]string, 0),
        "", "", false,
        "", "", "")
    descriptor := entities.NewAppDescriptor(network.ID, "descriptor1", "",
        "", "", "", 0, []string{"repo1:tag1"})
    instance := entities.NewAppInstance(network.ID, descriptor.ID, cluster.ID, "instance1", "", "",
        "", "", entities.AppStorageDefault, make([]entities.ApplicationPort,0), 0,"localhost")
    user := entities.NewUser("user", "9999", "email@email.com", time.Now(), time.Now().Add(time.Hour))

    access := entities.NewUserAccess("user", []entities.RoleType{entities.GlobalAdmin})
    helper.manager.networkProvider.Add(* network)
    helper.manager.clusterProvider.Add(* cluster)
    helper.manager.nodeProvider.Add(* node)
    helper.manager.appDescProvider.Add(* descriptor)
    helper.manager.appInstProvider.Add(* instance)
    helper.manager.userProvider.Add(* user)
    helper.manager.accessProvider.Add(* access)

    dump, err := helper.manager.dumpManager.Export()
    helper.Nil(err, "expecting dump")
    helper.NotNil(dump, "expecting dump")
    helper.Equal(1, len(dump.Networks))
    helper.Equal(1, len(dump.Clusters))
    helper.Equal(1, len(dump.Nodes))
    helper.Equal(1, len(dump.Descriptors))
    helper.Equal(1, len(dump.Instances))
    helper.Equal(1, len(dump.Users))
    helper.Equal(1, len(dump.UserAccesses))

}
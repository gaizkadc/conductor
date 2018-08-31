package backup
//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Backup Restore manager tests



import (
    "testing"
    "time"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/accessstorage"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/daishogroup/system-model/provider/passwordstorage"
)

type ManagerHelper struct {
    networkProvider * networkstorage.MockupNetworkProvider
    clusterProvider * clusterstorage.MockupClusterProvider
    nodeProvider    * nodestorage.MockupNodeProvider
    appDescProvider * appdescstorage.MockupAppDescProvider
    userProvider    * userstorage.MockupUserProvider
    accessProvider  * accessstorage.MockupUserAccessProvider
    passwordProvider * passwordstorage.MockupPasswordProvider
    brManager         Manager
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
    helper.manager.userProvider.Clear()
    helper.manager.accessProvider.Clear()
    helper.manager.passwordProvider.Clear()

}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(ManagerTestSuite))
}

func NewManagerHelper() ManagerHelper {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var clusterProvider = clusterstorage.NewMockupClusterProvider()
    var nodeProvider = nodestorage.NewMockupNodeProvider()
    var appDescProvider = appdescstorage.NewMockupAppDescProvider()
    var userProvider = userstorage.NewMockupUserProvider()
    var accessProvider = accessstorage.NewMockupUserAccessProvider()
    var passwordProvider = passwordstorage.NewMockupPasswordProvider()

    var brMgr = NewManager(networkProvider, clusterProvider, nodeProvider, appDescProvider,
        userProvider, accessProvider,passwordProvider)
    return ManagerHelper{networkProvider, clusterProvider, nodeProvider,
        appDescProvider,userProvider,
        accessProvider,passwordProvider,brMgr}
}


func (helper * ManagerTestSuite) TestBackup() {

    network := entities.NewNetwork("network1", "", "", "", "")
    cluster := entities.NewCluster(network.ID, "cluster1", "", entities.EdgeType,
        "", "", entities.ClusterInstalled, false, false)
    node := entities.NewNode(network.ID, cluster.ID,
        "node1", "", make([]string, 0),
        "", "", false,
        "", "", "")
    descriptor := entities.NewAppDescriptor(network.ID, "descriptor1", "",
        "", "", "", 0, []string{"repo1:tag1"})
    user := entities.NewUser("user", "9999", "email@email.com", time.Now(), time.Now().Add(time.Hour))

    access := entities.NewUserAccess(user.ID, []entities.RoleType{entities.GlobalAdmin})
    userpassword := "daisho"
    password,err := entities.NewPassword(user.ID,&userpassword)

    helper.Nil(err, "expecting password")


    helper.manager.networkProvider.Add(* network)
    helper.manager.clusterProvider.Add(* cluster)
    helper.manager.nodeProvider.Add(* node)
    helper.manager.appDescProvider.Add(* descriptor)
    helper.manager.userProvider.Add(* user)
    helper.manager.accessProvider.Add(* access)
    helper.manager.passwordProvider.Add(* password)

    backup, err := helper.manager.brManager.Export("all")
    helper.Nil(err, "expecting dump")
    helper.NotNil(backup, "expecting dump")
    helper.Equal(1, len(backup.Networks))
    helper.Equal(1, len(backup.Clusters))
    helper.Equal(1, len(backup.Nodes))
    helper.Equal(1, len(backup.AppDescriptors))
    helper.Equal(1, len(backup.Users))
    //helper.Equal(1, len(backup.UserAccesses))
   // helper.Equal(1, len(backup.Passwords))

}


func (helper * ManagerTestSuite) TestBackupUser() {

    network := entities.NewNetwork("network1", "", "", "", "")
    cluster := entities.NewCluster(network.ID, "cluster1", "", entities.EdgeType,
        "", "", entities.ClusterInstalled, false, false)
    node := entities.NewNode(network.ID, cluster.ID,
        "node1", "", make([]string, 0),
        "", "", false,
        "", "", "")
    descriptor := entities.NewAppDescriptor(network.ID, "descriptor1", "",
        "", "", "", 0, []string{"repo1:tag1"})
    user := entities.NewUser("user", "9999", "email@email.com", time.Now(), time.Now().Add(time.Hour))
    userpassword := "daisho"
    password,err := entities.NewPassword(user.ID,&userpassword)

    helper.Nil(err, "expecting password")

    access := entities.NewUserAccess(user.ID, []entities.RoleType{entities.GlobalAdmin})
    helper.manager.networkProvider.Add(* network)
    helper.manager.clusterProvider.Add(* cluster)
    helper.manager.nodeProvider.Add(* node)
    helper.manager.appDescProvider.Add(* descriptor)
    helper.manager.userProvider.Add(* user)
    helper.manager.accessProvider.Add(* access)
    helper.manager.passwordProvider.Add(* password)

    backup, err := helper.manager.brManager.Export("users")

    helper.Nil(err, "expecting Cluster backup")
    helper.NotNil(backup, "expecting cluster backup")
    helper.Equal(0, len(backup.Networks))
    helper.Equal(0, len(backup.Clusters))
    helper.Equal(0, len(backup.Nodes))
    helper.Equal(0, len(backup.AppDescriptors))
    helper.Equal(1, len(backup.Users))

}


func (helper * ManagerTestSuite) TestRestoreCluster() {


    network := entities.NewNetwork("network1", "", "", "", "")
    cluster := entities.NewCluster(network.ID, "cluster1", "", entities.EdgeType,
        "", "", entities.ClusterInstalled, false, false)
    node := entities.NewNode(network.ID, cluster.ID,
        "node1", "", make([]string, 0),
        "", "", false,
        "", "", "")
    descriptor := entities.NewAppDescriptor(network.ID, "descriptor1", "",
        "", "", "", 0, []string{"repo1:tag1"})
    user := entities.NewUser("user", "9999", "email@email.com", time.Now(), time.Now().Add(time.Hour))


    access := entities.NewUserAccess(user.ID, []entities.RoleType{entities.GlobalAdmin})
    userpassword := "daisho"
    password,err := entities.NewPassword(user.ID,&userpassword)

    helper.Nil(err, "expecting password")

    helper.manager.networkProvider.Add(* network)
    helper.manager.clusterProvider.Add(* cluster)
    helper.manager.nodeProvider.Add(* node)
    helper.manager.appDescProvider.Add(* descriptor)
    helper.manager.userProvider.Add(* user)
    helper.manager.accessProvider.Add(* access)
    helper.manager.passwordProvider.Add(* password)


    // create backup data
    backup, err := helper.manager.brManager.Export("all")

    // restore data
    err = helper.manager.brManager.Import("all", backup)
    helper.Nil(err, "expecting Cluster Import no err")

    // DP-1703: Check if application is registered to the network
    exists := helper.manager.networkProvider.ExistsAppDesc(network.ID, descriptor.ID)
    helper.True(exists)

    // get data again for verification
    backup, err = helper.manager.brManager.Export("all")

    helper.Nil(err, "expecting Cluster restore")
    helper.NotNil(backup, "expecting cluster restore data")
    helper.Equal(1, len(backup.Networks))
    helper.Equal(1, len(backup.Clusters))
    helper.Equal(1, len(backup.Nodes))
    helper.Equal(1, len(backup.AppDescriptors))
    helper.Equal(1, len(backup.Users))


}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Mockup implementation of the dump client.

package client

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/daishogroup/system-model/server/info"
)

// InfoMockup is the mockup of Info.
type InfoMockup struct {
    networkProvider *networkstorage.MockupNetworkProvider
    clusterProvider *clusterstorage.MockupClusterProvider
    nodeProvider    *nodestorage.MockupNodeProvider
    appDescProvider *appdescstorage.MockupAppDescProvider
    appInstProvider *appinststorage.MockupAppInstProvider
    userProvider    *userstorage.MockupUserProvider
    info.Manager
}

// NewInfoMockup is the basic constructor.
func NewInfoMockup() Info {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var clusterProvider = clusterstorage.NewMockupClusterProvider()
    var nodeProvider = nodestorage.NewMockupNodeProvider()
    var appDescProvider = appdescstorage.NewMockupAppDescProvider()
    var appInstProvider = appinststorage.NewMockupAppInstProvider()
    var userProvider = userstorage.NewMockupUserProvider()

    var infoMgr = info.NewManager(networkProvider, clusterProvider, nodeProvider,
        appDescProvider, appInstProvider, userProvider)
    return &InfoMockup{
        networkProvider, clusterProvider, nodeProvider,
        appDescProvider, appInstProvider, userProvider,
        infoMgr}
}

// ClearMockup remove all the elements.
func (mockup *InfoMockup) ClearMockup() {
    mockup.networkProvider.Clear()
    mockup.clusterProvider.Clear()
    mockup.nodeProvider.Clear()
    mockup.appDescProvider.Clear()
    mockup.appInstProvider.Clear()
}

// InitMockup loads a set of data.
func (mockup *InfoMockup) InitMockup() {
    mockup.ClearMockup()

    testIP := "0.0.0.0"
    testUsername := "username"
    testPassword := "pass"
    testSSHKey := "ABCDFGHIJK"
    testLocation := "Madrid, Spain"
    testAdminName := "Admin"
    testAdminEmail := "admin@admins.com"
    testAdminPhone := "1234"

    mockup.nodeProvider.Add(* entities.NewNodeWithID("1", "1", "1",
        "Node1", "Description Node 1",make([]string, 0),
        testIP, testIP, true, testUsername, testPassword, testSSHKey, entities.NodeUnchecked))

    mockup.nodeProvider.Add(* entities.NewNodeWithID("1", "1", "2",
        "Node2", "Description Node 2",make([]string, 0),
        testIP, testIP, true, testUsername, testPassword, testSSHKey, entities.NodeUnchecked))

    mockup.nodeProvider.Add(* entities.NewNodeWithID("1", "2", "3",
        "Node3", "Description Node 3",make([]string, 0),
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

    descriptor := entities.NewAppDescriptorWithID("1", "app1", "descriptor1", "",
        "", "", "", 0, []string{"repo1:tag1"})
    instance := entities.NewAppInstanceWithID("1", "inst1", descriptor.ID, "1",
        "instance1", "", "",
        "", entities.AppInstInit, "1Gb", entities.AppStorageDefault,
        make([]entities.ApplicationPort, 0), 0, "")
    mockup.appDescProvider.Add(* descriptor)
    mockup.networkProvider.RegisterAppDesc(descriptor.NetworkID, descriptor.ID)
    mockup.appInstProvider.Add(* instance)
    mockup.networkProvider.RegisterAppInst(instance.NetworkID, instance.DeployedID)
}

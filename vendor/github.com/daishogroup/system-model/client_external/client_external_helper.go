//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the helper for the client integration test.

package client_external

import (
    "context"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/server/cluster"
    "github.com/daishogroup/system-model/server/network"
    "github.com/daishogroup/system-model/server/node"
    "github.com/daishogroup/dhttp"
)

//Internal test Address.
const BaseAddress = "localhost"

// The EndpointHelper structure contains all the elements in the endpoint processing path: handler, manager and
// providers so they can be used if needed during the tests.
type EndpointHelper struct {
    networkProvider *networkstorage.MockupNetworkProvider
    clusterProvider *clusterstorage.MockupClusterProvider
    nodeProvider    *nodestorage.MockupNodeProvider
    appInstProvider *appinststorage.MockupAppInstProvider
    clusterMgr      cluster.Manager
    clusterHandler  cluster.Handler
    networkMgr      network.Manager
    networkHandler  network.Handler
    nodeMgr         node.Manager
    nodeHandler     node.Handler
    port            int
    handler         http.Handler
    srv             *http.Server
}

// Reset all the providers.
func (helper *EndpointHelper) ResetProvider() {
    helper.networkProvider.Clear()
    helper.clusterProvider.Clear()
    helper.nodeProvider.Clear()

    helper.nodeProvider.Add(* entities.NewNodeWithID("1", "1", "1",
        "Node1", "Description Node1", make([]string, 0),
        "0.0.0.0", "0.0.0.0", true,
        "username", "", "", entities.NodeUnchecked))
    helper.nodeProvider.Add(* entities.NewNodeWithID("1", "1", "2",
        "Node2", "Description Node2", make([]string, 0),
        "0.0.0.0", "0.0.0.0", true,
        "username", "", "", entities.NodeUnchecked))
    helper.nodeProvider.Add(* entities.NewNodeWithID("1", "2", "3",
        "Node3", "Description Node3", make([]string, 0),
        "0.0.0.0", "0.0.0.0", true,
        "username", "", "", entities.NodeUnchecked))

    helper.clusterProvider.Add(* entities.NewClusterWithID("1","1",
        "Cluster1", "Description Cluster 1",
        entities.GatewayType, "Madrid","admin@admin.com",
        entities.ClusterCreated, false, false))
    helper.clusterProvider.Add(* entities.NewClusterWithID("1","2",
        "Cluster2", "Description Cluster 2",
        entities.GatewayType, "Madrid","admin@admin.com",
        entities.ClusterCreated, false, false))
    helper.clusterProvider.Add(* entities.NewClusterWithID("2","3",
        "Cluster3", "Description Cluster 3",
        entities.GatewayType, "Madrid","admin@admin.com",
        entities.ClusterCreated, false, false))

    helper.clusterProvider.AttachNode("1", "1")
    helper.clusterProvider.AttachNode("1", "2")
    helper.clusterProvider.AttachNode("2", "3")

    helper.networkProvider.Add(* entities.NewNetworkWithID("1", "Network1", "Description Network1",
        "admin", "1234", "admin@admins.com"))

    helper.networkProvider.Add(* entities.NewNetworkWithID("2", "Network2", "Description Network2",
        "admin", "1234", "admin@admins.com"))

    helper.networkProvider.AttachCluster("1", "1")
    helper.networkProvider.AttachCluster("1", "2")
    helper.networkProvider.AttachCluster("2", "3")
}

func (helper *EndpointHelper) LaunchServer() {
    err := helper.srv.ListenAndServe()
    if err != nil {
        println(err.Error())
    }
}

// Start the HttpServer.
func (helper *EndpointHelper) Start() {
    go helper.LaunchServer()
}

// Shutdown the HTTPServer.
func (helper *EndpointHelper) Shutdown() {
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    helper.srv.Shutdown(ctx)
    helper.srv.Close()
}

func (helper *EndpointHelper) GetListeningAddress() (string, int) {
    return BaseAddress, helper.port
}

// Create a new EndpointHelper.
//   returns:
//     A new endpoint helper with a mockup provider.
func NewEndpointHelper() EndpointHelper {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var clusterProvider = clusterstorage.NewMockupClusterProvider()
    var nodeProvider = nodestorage.NewMockupNodeProvider()
    var appInstProvider = appinststorage.NewMockupAppInstProvider()
    var clusterMgr = cluster.NewManager(networkProvider, clusterProvider, appInstProvider)
    var clusterHandler = cluster.NewHandler(clusterMgr)
    var networkMgr = network.NewManager(networkProvider)
    var networkHandler = network.NewHandler(networkMgr)
    var nodeMgr = node.NewManager(networkProvider, clusterProvider, nodeProvider)
    var nodeHandler = node.NewHandler(nodeMgr)
    port, _ := dhttp.GetAvailablePort()
    var handler = mux.NewRouter()
    networkHandler.SetRoutes(handler)
    clusterHandler.SetRoutes(handler)
    nodeHandler.SetRoutes(handler)
    var srv = &http.Server{
        Handler: handler,
        Addr:    BaseAddress + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    return EndpointHelper{networkProvider, clusterProvider,
        nodeProvider, appInstProvider, clusterMgr, clusterHandler,
        networkMgr, networkHandler,
        nodeMgr, nodeHandler, port, handler, srv}
}

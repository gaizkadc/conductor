//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Dump handler tests.

package info

import (
    "context"
    "fmt"
    "net/http"
    "strconv"
    "testing"
    "time"

    "github.com/bitly/go-simplejson"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/suite"
    "github.com/verdverm/frisby"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/daishogroup/dhttp"
)

const (
    // Local address for the server exposing the API.
    BaseAddress = "localhost"
)

// The EndpointHelper structure contains all the elements in the endpoint processing path: handler, manager and
// providers so they can be used if needed during the tests.
type EndpointHelper struct {
    networkProvider *networkstorage.MockupNetworkProvider
    clusterProvider *clusterstorage.MockupClusterProvider
    nodeProvider    *nodestorage.MockupNodeProvider
    appDescProvider *appdescstorage.MockupAppDescProvider
    appInstProvider *appinststorage.MockupAppInstProvider
    dumpMgr         Manager
    dumpHandler     Handler
    port            int
    handler         http.Handler
    srv             *http.Server
}

// Create a new EndpointHelper.
//   returns:
//     A new endpoint helper with a mockup provider.
func NewEndpointHelper() EndpointHelper {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var clusterProvider = clusterstorage.NewMockupClusterProvider()
    var nodeProvider = nodestorage.NewMockupNodeProvider()
    var appDescProvider = appdescstorage.NewMockupAppDescProvider()
    var appInstProvider = appinststorage.NewMockupAppInstProvider()
    var userProvider = userstorage.NewMockupUserProvider()

    var dumpMgr = NewManager(networkProvider, clusterProvider, nodeProvider,
        appDescProvider, appInstProvider, userProvider)
    var dumpHandler = NewHandler(dumpMgr)
    port, _ := dhttp.GetAvailablePort()

    handler := mux.NewRouter()
    dumpHandler.SetRoutes(handler)

    var srv = &http.Server{
        Handler: handler,
        Addr:    BaseAddress + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    return EndpointHelper{networkProvider, clusterProvider, nodeProvider,
        appDescProvider, appInstProvider,
        dumpMgr, dumpHandler, port, handler, srv}
}

func (helper *EndpointHelper) LaunchServer() {
    log.Info("Launching helper on : " + BaseAddress + ":" + strconv.Itoa(helper.port))
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

// The test helper structure is used to "inherit" Suite functions and contains the REST handler and the
// endpoint structure.
type ReducedInfoHandlerTestHelper struct {
    suite.Suite
    endpoints EndpointHelper
}

// The SetupSuite is the first method called and defines the testing suite. It creates the endpoints and launches
// a server listening for client requests.
func (helper *ReducedInfoHandlerTestHelper) SetupSuite() {
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.endpoints.Start()
    dhttp.WaitURLAvailable(BaseAddress,helper.endpoints.port,5,"/", 1)
}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *ReducedInfoHandlerTestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    helper.Equal(0, frisby.Global.NumErrored, "expecting 0 failures")
}

// The SetupTest method is called before every test on the suite.
func (helper *ReducedInfoHandlerTestHelper) SetupTest() {
    helper.endpoints.networkProvider.Clear()
    helper.endpoints.clusterProvider.Clear()
    helper.endpoints.nodeProvider.Clear()
    helper.endpoints.appDescProvider.Clear()
    helper.endpoints.appInstProvider.Clear()
}

func TestReducedInfoHandlerSuite(t *testing.T) {
    suite.Run(t, new(ReducedInfoHandlerTestHelper))
}

// Obtain a target URL given an endpoint.
//   params:
//     endpoint The URL endpoint.
func (helper *ReducedInfoHandlerTestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

func (helper *ReducedInfoHandlerTestHelper) LoadTestData() {
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
        "", "", entities.AppStorageDefault, make([]entities.ApplicationPort, 0), 0, "")
    helper.endpoints.networkProvider.Add(* network)
    helper.endpoints.clusterProvider.Add(* cluster)
    helper.endpoints.nodeProvider.Add(* node)
    helper.endpoints.appDescProvider.Add(* descriptor)
    helper.endpoints.appInstProvider.Add(* instance)
}

func (helper *ReducedInfoHandlerTestHelper) TestReducedInfo() {
    helper.LoadTestData()
    f := frisby.Create("Export").Get(helper.getURL("info/reduced")).Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        networks, err := json.Get("networks").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(networks), "length does not match")
        clusters, err := json.Get("clusters").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(clusters), "length does not match")
        nodes, err := json.Get("nodes").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(nodes), "length does not match")
        descriptors, err := json.Get("descriptors").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(descriptors), "length does not match")
        instances, err := json.Get("instances").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(instances), "length does not match")
    })
}

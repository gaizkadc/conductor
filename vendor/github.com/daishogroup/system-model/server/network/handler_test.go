//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Network handler integration tests.

package network

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
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/stretchr/testify/assert"
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
    networkMgr      Manager
    networkHandler  Handler
    port            int
    handler         http.Handler
    srv             *http.Server
}

// Create a new EndpointHelper.
//   returns:
//     A new endpoint helper with a mockup provider.
func NewEndpointHelper() EndpointHelper {
    var networkMockupProvider = networkstorage.NewMockupNetworkProvider()
    var networkMgr = NewManager(networkMockupProvider)
    var networkHandler = NewHandler(networkMgr)
    port, _ := dhttp.GetAvailablePort()

    handler := mux.NewRouter()
    networkHandler.SetRoutes(handler)

    var srv = &http.Server{
        Handler: handler,
        Addr:    BaseAddress + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    return EndpointHelper{networkMockupProvider, networkMgr,
        networkHandler, port, handler, srv}
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

// The test helper structure is used to "inherit" Suite functions and contains the REST handler and the
// endpoint structure.
type TestHelper struct {
    suite.Suite
    endpoints EndpointHelper
}

// The SetupSuite is the first method called and defines the testing suite. It creates the endpoints and launches
// a server listening for client requests.
func (helper *TestHelper) SetupSuite() {
    log.Info("Creating test suite")
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.endpoints.Start()
    dhttp.WaitURLAvailable(BaseAddress,helper.endpoints.port,5,"/", 1)
}

// Obtain a target URL given an endpoint.
//   params:
//     endpoint The URL endpoint.
func (helper *TestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *TestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    helper.Equal(0, frisby.Global.NumErrored, "expecting 0 failures")
}

// The SetupTest method is called before every test on the suite.
func (helper *TestHelper) SetupTest() {
    helper.endpoints.networkProvider.Clear()
}

func (helper *TestHelper) TestAddNetwork() {
    toAdd := entities.NewAddNetworkRequest(testNetworkName, testDescription, testAdminName, testAdminPhone, testAdminEmail)
    f := frisby.Create("Add a new network").
        Post(helper.getURL("network/add")).
        SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)

    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        networkID, err := json.Get("id").String()
        helper.Nil(err, "expecting id")
        helper.NotEmpty(networkID, "a network id must be returned")
    })

    f.ExpectJson("name", testNetworkName)
    f.ExpectJson("description", testDescription)
    f.ExpectJson("adminName", testAdminName)
    f.ExpectJson("adminPhone", testAdminPhone)
    f.ExpectJson("adminEmail", testAdminEmail)

}

func (helper *TestHelper) TestAddInvalidNetwork() {
    toAdd := entities.AddNetworkRequest{}
    f := frisby.Create("Add a new network").
        Post(helper.getURL("network/add")).
        SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusBadRequest)
}

func (helper *TestHelper) TestAddInvalidEntityAsNetwork() {
    toAdd := entities.Network{}
    f := frisby.Create("Add a new network").
        Post(helper.getURL("network/add")).
        SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusBadRequest)
}

func (helper *TestHelper) TestAddNilNetwork() {
    f := frisby.Create("Add a new network").
        Post(helper.getURL("network/add")).Send()
    f.ExpectStatus(http.StatusBadRequest)
}

func (helper *TestHelper) TestListEmptyNetworks() {
    targetURL := helper.getURL("network/list")
    f := frisby.Create("List empty networks").Get(targetURL).Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        entries, err := json.Array()
        helper.Nil(err, "expecting array")
        helper.Equal(0, len(entries), "length does not match")
    })
}

func (helper *TestHelper) TestListNetworks() {
    toAdd := entities.NewNetworkWithID("n1", testNetworkName, testDescription,
        testAdminName, testAdminPhone, testAdminEmail)
    helper.endpoints.networkProvider.Add(*toAdd)
    targetURL := helper.getURL("network/list")
    f := frisby.Create("List single networks").Get(targetURL).Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        entries, err := json.Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(entries), "length does not match")
    })
}

func (helper *TestHelper) TestGetNetwork() {
    toAdd := entities.NewNetworkWithID("n2", testNetworkName, testDescription,
        testAdminName, testAdminPhone, testAdminEmail)
    helper.endpoints.networkProvider.Add(*toAdd)
    targetURL := helper.getURL("network/n2/info")
    f := frisby.Create("Get network").Get(targetURL).Send()
    f.ExpectStatus(http.StatusOK)
    f.ExpectJson("id", "n2")
    f.ExpectJson("name", testNetworkName)
    f.ExpectJson("description", testDescription)
    f.ExpectJson("adminName", testAdminName)
    f.ExpectJson("adminPhone", testAdminPhone)
    f.ExpectJson("adminEmail", testAdminEmail)
}

func (helper *TestHelper) TestDeleteNetwork() {
    // list networks
    toAdd := entities.NewNetworkWithID("n1", testNetworkName, testDescription,
        testAdminName, testAdminPhone, testAdminEmail)
    helper.endpoints.networkProvider.Add(*toAdd)

    // delete it
    targetURL := helper.getURL("network/n1/delete")
    // try to get it and check the fail
    f := frisby.Create("Delete single network").Delete(targetURL).Send()
    f.ExpectStatus(http.StatusOK)

    _, err := helper.endpoints.networkProvider.RetrieveNetwork("n1")
    assert.NotNil(helper.T(), err, "not correctly removed instance")
}

func (helper *TestHelper) TestDeleteNetworkWithClusters() {
    // list networks
    toAdd := entities.NewNetworkWithID("n1", testNetworkName, testDescription,
        testAdminName, testAdminPhone, testAdminEmail)
    helper.endpoints.networkProvider.Add(*toAdd)

    helper.endpoints.networkProvider.AttachCluster(toAdd.ID, "c1")
    // delete it
    targetURL := helper.getURL("network/n1/delete")
    // try to get it and check the fail
    f := frisby.Create("Delete network with clusters").Delete(targetURL).Send()
    f.ExpectStatus(http.StatusInternalServerError)
}

// This function launches the testing suite.
func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(TestHelper))
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package config

import (
    "context"
    "fmt"
    "net/http"
    "strconv"
    "testing"
    "time"

    "github.com/daishogroup/system-model/entities"
    log "github.com/sirupsen/logrus"

    "github.com/daishogroup/system-model/provider/configstorage"

    "github.com/gorilla/mux"
    "github.com/stretchr/testify/suite"
    "github.com/verdverm/frisby"
    "github.com/daishogroup/dhttp"
)

const (
    // Local address for the server exposing the API.
    BaseAddress = "localhost"
)

// The EndpointHelper structure contains all the elements in the endpoint processing path: handler, manager and
// providers so they can be used if needed during the tests.
type EndpointHelper struct {
    configProvider *configstorage.MockupConfigProvider
    configMgr      Manager
    configHandler  Handler
    port           int
    handler        http.Handler
    srv            *http.Server
}

// Create a new EndpointHelper.
//   returns:
//     A new endpoint helper with a mockup provider.
func NewEndpointHelper() EndpointHelper {
    var configProvider = configstorage.NewMockupConfigProvider()
    var configMgr = NewManager(configProvider)
    var configHandler = NewHandler(configMgr)
    port, _ := dhttp.GetAvailablePort()
    handler := mux.NewRouter()
    configHandler.SetRoutes(handler)

    var srv = &http.Server{
        Handler: handler,
        Addr:    BaseAddress + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    return EndpointHelper{configProvider,
        configMgr, configHandler,
        port, handler, srv}
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
type ConfigTestHelper struct {
    suite.Suite
    endpoints EndpointHelper
}

// The SetupSuite is the first method called and defines the testing suite. It creates the endpoints and launches
// a server listening for client requests.
func (helper *ConfigTestHelper) SetupSuite() {
    log.Info("Config handler test suite")
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.endpoints.Start()
    dhttp.WaitURLAvailable(BaseAddress,helper.endpoints.port,5,"/", 1)
}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *ConfigTestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    helper.Equal(0, frisby.Global.NumErrored, "expecting 0 failures")
}

// The SetupTest method is called before every test on the suite.
func (helper *ConfigTestHelper) SetupTest() {
    helper.endpoints.configProvider.Clear()
}

// This function launches the testing suite.
func TestConfigHandlerSuite(t *testing.T) {
    suite.Run(t, new(ConfigTestHelper))
}

// Obtain a target URL given an endpoint.
//   params:
//     endpoint The URL endpoint.
func (helper *ConfigTestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

func (helper *ConfigTestHelper) TestSetGetConfig() {
    toAdd := entities.NewConfig("1h")
    f := frisby.Create("Set the configuration").
        Post(helper.getURL("config/set")).
        SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)

    fr := frisby.Create("Get config").Get(helper.getURL("config/get")).Send()
    fr.ExpectStatus(http.StatusOK)
    fr.ExpectJson("logRetention", toAdd.LogRetention)
}

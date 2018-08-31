//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Credentials handler testing.

package credentials

import (
    "context"
    "fmt"
    "net/http"
    "strconv"
    "testing"
    "time"
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/suite"
    "github.com/verdverm/frisby"
    "github.com/bitly/go-simplejson"
    "github.com/daishogroup/system-model/provider/credentialsstorage"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/dhttp"
)

const (
    BaseAddress     = "localhost"
    Testuuid        = "uuid"
    TestPublicKey   = "publickey"
    TestPrivateKey  = "privatekey"
    TestDescription = "description"
    TestTypeKey     = "typekey"
)

type AppEndpointHelper struct {
    credentialsProvider *credentialsstorage.MockupCredentialsProvider
    appMgr              Manager
    appHandler          Handler
    port                int
    handler             http.Handler
    srv                 *http.Server
}

// Create a new EndpointHelper.
//   returns:
//     A new endpoint helper with a mockup provider.
func NewAppEndpointHelper() AppEndpointHelper {
    var credentialsProvider = credentialsstorage.NewMockupCredentialsProvider()
    var appManager = NewManager(credentialsProvider)
    var appHandler = NewHandler(appManager)
    port, _ := dhttp.GetAvailablePort()
    var handler = mux.NewRouter()
    appHandler.SetRoutes(handler)

    var srv = &http.Server{
        Handler: handler,
        Addr:    BaseAddress + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    return AppEndpointHelper{
        credentialsProvider,
        appManager,
        appHandler,
        port,
        handler,
        srv}
}

func (helper *AppEndpointHelper) LaunchServer() {
    err := helper.srv.ListenAndServe()
    if err != nil {
        println(err.Error())
    }
}

// Start the HttpServer.
func (helper *AppEndpointHelper) Start() {
    go helper.LaunchServer()
}

// Shutdown the HTTPServer.
func (helper *AppEndpointHelper) Shutdown() {
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    helper.srv.Shutdown(ctx)
    helper.srv.Close()
}

type AppTestHelper struct {
    suite.Suite
    endpoints AppEndpointHelper
}

func (helper *AppTestHelper) SetupSuite() {
    helper.endpoints = NewAppEndpointHelper()
    helper.endpoints.Start()
    dhttp.WaitURLAvailable(BaseAddress,helper.endpoints.port,5,"/", 1)
}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *AppTestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    helper.Equal(0, frisby.Global.NumErrored, "expecting 0 failures")
}

func (helper *AppTestHelper) SetupTest() {
    helper.endpoints.credentialsProvider.Clear()
}

// Obtain a target URL given an endpoint.
// params:
// endpoint The URL endpoint.
func (helper *AppTestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

func TestAppHandlerSuite(t *testing.T) {
    suite.Run(t, new(AppTestHelper))
}

func (helper *AppTestHelper) TestAddCredentials() {

    toAdd := entities.NewAddCredentialsRequest(Testuuid, TestPublicKey, TestPrivateKey, TestDescription, TestTypeKey)

    f := frisby.Create("Add credentials").
        Post(helper.getURL("credentials/add")).SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        helper.Nil(err, "expecting id")
    })

    // Check if it was added
    f2 := frisby.Create("Get credentials").
        Get(helper.getURL(fmt.Sprintf("credentials/%s/get", Testuuid))).Send()
    f2.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        helper.Nil(err, "expecting id")
    })
    f2.ExpectStatus(http.StatusOK)
    f2.ExpectJson("uuid", Testuuid)

    // Delete it
    f3 := frisby.Create("Delete credentials").
        Delete(helper.getURL(fmt.Sprintf("credentials/%s/delete", Testuuid))).Send()
    f3.ExpectStatus(http.StatusOK)
    // Try to retrieve it again
    f4 := frisby.Create("Get credentials").
        Get(helper.getURL(fmt.Sprintf("credentials/%s/get", Testuuid))).Send()
    f4.ExpectStatus(http.StatusInternalServerError)

}

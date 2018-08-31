//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// OAuth handler integration tests.

package oauth

import (
    "fmt"
    "net/http"
    "testing"
    "time"

    "github.com/bitly/go-simplejson"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/verdverm/frisby"

    "context"
    "strconv"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/oauthstorage"
    "github.com/daishogroup/dhttp"
)

const (
    // Local address for the server exposing the API.
    BaseAddress = "localhost"
)

// The EndpointHelper structure contains all the elements in the endpoint processing path: handler, manager and
// providers so they can be used if needed during the tests.
type EndpointHelper struct {
    oauthProvider *oauthstorage.MockupOAuthProvider
    oauthMgr      Manager
    oauthHandler  Handler
    port          int
    handler       http.Handler
    srv           *http.Server
}

// Create a new EndpointHelper.
//   returns:
//     A new endpoint helper with a mockup provider.
func NewEndpointHelper() EndpointHelper {
    var oauthProvider = oauthstorage.NewMockupOAuthProvider()
    var oauthMgr = NewManager(oauthProvider)
    var userHandler = NewHandler(oauthMgr)
    port, _ := dhttp.GetAvailablePort()

    handler := mux.NewRouter()
    userHandler.SetRoutes(handler)

    var srv = &http.Server{
        Handler: handler,
        Addr:    BaseAddress + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    return EndpointHelper{oauthProvider, oauthMgr, userHandler,
        port, handler, srv}
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
type OAuthHandlerTestHelper struct {
    suite.Suite
    endpoints EndpointHelper
}

// The SetupSuite is the first method called and defines the testing suite. It creates the endpoints and launches
// a server listening for client requests.
func (helper *OAuthHandlerTestHelper) SetupSuite() {
    log.Info("OAuth handler test suite")
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.endpoints.Start()
    dhttp.WaitURLAvailable(BaseAddress,helper.endpoints.port,5,"/", 1)
}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *OAuthHandlerTestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    assert.Equal(helper.T(), 0, frisby.Global.NumErrored, "expecting 0 failures")
}

// The SetupTest method is called before every test on the suite.
func (helper *OAuthHandlerTestHelper) SetupTest() {
    helper.endpoints.oauthProvider.Clear()
}

// This function launches the testing suite.
func TestClusterHandlerSuite(t *testing.T) {
    suite.Run(t, new(OAuthHandlerTestHelper))
}

// Obtain a target URL given an endpoint.
//   params:
//     endpoint The URL endpoint.
func (helper *OAuthHandlerTestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

func (helper *OAuthHandlerTestHelper) TestSetSecret() {

    // Add an initial empty entry
    err := helper.endpoints.oauthProvider.Add(entities.NewOAuthSecrets(TestUserID))
    helper.Nil(err, "unexpected error creating initial oauth secret list")

    request := entities.NewOAuthAddEntryRequest("app1", "clientID1", "secret1")
    //err = helper.endpoints.oauthMgr.SetSecret(TestUserID,request)
    //helper.Nil(err, "unexpected error creating password")

    f := frisby.Create("Set secret").
        Post(helper.getURL(fmt.Sprintf("oauth/%s", TestUserID))).
        SetJson(request).Send()
    f.ExpectStatus(http.StatusOK)

    // Get the secret
    f = frisby.Create("Get secret").
        Get(helper.getURL(fmt.Sprintf("oauth/%s", TestUserID))).Send()
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        _, err = json.Get("userID").String()
        assert.Nil(helper.T(), err, "expecting id")
        m, err := json.Get("entries").Map()
        assert.Nil(helper.T(), err, "expecting entries")
        _, isThere := m["app1"]
        assert.True(helper.T(), isThere)
    })
    f.ExpectJson("userID", TestUserID)

}

func (helper *OAuthHandlerTestHelper) TestDeleteSecret() {

    // Add an initial empty entry
    err := helper.endpoints.oauthProvider.Add(entities.NewOAuthSecrets(TestUserID))
    helper.Nil(err, "unexpected error creating initial oauth secret list")

    // Get
    f := frisby.Create("Get secret").
        Get(helper.getURL(fmt.Sprintf("oauth/%s", TestUserID))).Send()
    f.ExpectStatus(http.StatusOK)
    // Remove it
    f = frisby.Create("Delete secret").
        Delete(helper.getURL(fmt.Sprintf("oauth/%s", TestUserID))).Send()
    f.ExpectStatus(http.StatusOK)
    // Get must fail
    f = frisby.Create("Get secret").
        Get(helper.getURL(fmt.Sprintf("oauth/%s", TestUserID))).Send()
    f.ExpectStatus(http.StatusInternalServerError)
}

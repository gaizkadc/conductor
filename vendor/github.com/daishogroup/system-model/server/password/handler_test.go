//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Password handler integration tests.

package password

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
    "github.com/daishogroup/system-model/provider/passwordstorage"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/dhttp"
)

const (
    // Local address for the server exposing the API.
    BaseAddress = "localhost"
)

// The EndpointHelper structure contains all the elements in the endpoint processing path: handler, manager and
// providers so they can be used if needed during the tests.
type EndpointHelper struct {
    passwordProvider *passwordstorage.MockupPasswordProvider
    userMgr          Manager
    userHandler      Handler
    port             int
    handler          http.Handler
    srv              *http.Server
}

// Create a new EndpointHelper.
//   returns:
//     A new endpoint helper with a mockup provider.
func NewEndpointHelper() EndpointHelper {
    var passwordProvider = passwordstorage.NewMockupPasswordProvider()
    var userMgr = NewManager(passwordProvider)
    var userHandler = NewHandler(userMgr)
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
    return EndpointHelper{passwordProvider, userMgr, userHandler,
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
type UserHandlerTestHelper struct {
    suite.Suite
    endpoints EndpointHelper
}

// The SetupSuite is the first method called and defines the testing suite. It creates the endpoints and launches
// a server listening for client requests.
func (helper *UserHandlerTestHelper) SetupSuite() {
    log.Info("Node handler test suite")
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.endpoints.Start()
    dhttp.WaitURLAvailable(BaseAddress,helper.endpoints.port,5,"/", 1)
}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *UserHandlerTestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    assert.Equal(helper.T(), 0, frisby.Global.NumErrored, "expecting 0 failures")
}

// The SetupTest method is called before every test on the suite.
func (helper *UserHandlerTestHelper) SetupTest() {
    helper.endpoints.passwordProvider.Clear()
}

// This function launches the testing suite.
func TestClusterHandlerSuite(t *testing.T) {
    suite.Run(t, new(UserHandlerTestHelper))
}

// Obtain a target URL given an endpoint.
//   params:
//     endpoint The URL endpoint.
func (helper *UserHandlerTestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

func (helper *UserHandlerTestHelper) TestSetPassword() {

    toAdd, err := entities.NewPassword(TestUserID, nil)
    helper.Nil(err, "unexpected error creating password")
    helper.endpoints.passwordProvider.Add(*toAdd)

    thePassword := TestPassword
    req, err := entities.NewPassword(TestUserID, &thePassword)

    f := frisby.Create("Set password").
        Post(helper.getURL("password")).
        SetJson(*req).Send()

    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        _, err = json.Get("operation").String()
        assert.Nil(helper.T(), err, "expecting id")
    })
    f.ExpectJson("operation", errors.PasswordSet)
}

func (helper *UserHandlerTestHelper) TestDeletePassword() {
    thePassword := TestPassword
    toAdd, err := entities.NewPassword(TestUserID, &thePassword)
    helper.Nil(err, "unexpected error creating password")
    helper.endpoints.passwordProvider.Add(*toAdd)

    f := frisby.Create("Get password").
        Get(helper.getURL(fmt.Sprintf("password/%s", TestUserID))).Send()

    f.ExpectStatus(http.StatusOK)
    // Remove it
    f = frisby.Create("Delete password").
        Delete(helper.getURL(fmt.Sprintf("password/%s", TestUserID))).Send()
    f.ExpectStatus(http.StatusOK)

    // Try to recover again
    f = frisby.Create("Get password").
        Get(helper.getURL(fmt.Sprintf("password/%s", TestUserID))).Send()
    f.ExpectStatus(http.StatusInternalServerError)
}

func (helper *UserHandlerTestHelper) TestGetPassword() {
    thePassword := TestPassword
    toAdd, err := entities.NewPassword(TestUserID, &thePassword)
    helper.Nil(err, "unexpected error creating password")
    helper.endpoints.passwordProvider.Add(*toAdd)

    f := frisby.Create("Get password").
        Get(helper.getURL(fmt.Sprintf("password/%s", TestUserID))).Send()

    f.ExpectStatus(http.StatusOK)

    // Try to recover again
    f = frisby.Create("Get password").
        Get(helper.getURL(fmt.Sprintf("password/%s", "Not here!!!"))).Send()
    f.ExpectStatus(http.StatusInternalServerError)
}

//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// User handler integration tests.

package session

import (
    "fmt"
    "net/http"
    "testing"
    "time"

    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/verdverm/frisby"

    "context"
    "strconv"
    "github.com/daishogroup/system-model/entities"
    "github.com/bitly/go-simplejson"
    "github.com/daishogroup/dhttp"
    "github.com/daishogroup/system-model/provider/sessionstorage"
)

const (
    // Local address for the server exposing the API.
    BaseAddress = "localhost"
)

// The EndpointHelper structure contains all the elements in the endpoint processing path: handler, manager and
// providers so they can be used if needed during the tests.
type EndpointHelper struct {
    sessionProvider     *sessionstorage.MockupSessionProvider
    sessionMgr          Manager
    sessionHandler      Handler
    port             int
    handler          http.Handler
    srv              *http.Server
}

// Create a new EndpointHelper.
//   returns:
//     A new endpoint helper with a mockup provider.
func NewEndpointHelper() EndpointHelper {
    var sessionProvider = sessionstorage.NewMockupSessionProvider()
    var sessionMgr = NewManager(sessionProvider)
    var sessionHandler = NewHandler(sessionMgr)
    port, _ := dhttp.GetAvailablePort()

    handler := mux.NewRouter()
    sessionHandler.SetRoutes(handler)

    handler.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        path, err := route.GetPathTemplate()
        if err != nil {
            return err
        }
        methods, err := route.GetMethods()
        if err != nil {
            return err
        }
        methodStr := ""
        for _, m := range methods {
            methodStr = methodStr + ", " + m
        }
        log.Info(methodStr + " : " + path)
        return nil
    })

    var srv = &http.Server{
        Handler: handler,
        Addr:    BaseAddress + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    return EndpointHelper{sessionProvider,
        sessionMgr, sessionHandler, port, handler, srv}
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
type SessionHandlerTestHelper struct {
    suite.Suite
    endpoints EndpointHelper
}

// The SetupSuite is the first method called and defines the testing suite. It creates the endpoints and launches
// a server listening for client requests.
func (helper *SessionHandlerTestHelper) SetupSuite() {
    log.Info("User handler test suite")
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.endpoints.Start()
    dhttp.WaitURLAvailable(BaseAddress,helper.endpoints.port,5,"/", 1)
}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *SessionHandlerTestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    assert.Equal(helper.T(), 0, frisby.Global.NumErrored, "expecting 0 failures")
}

// The SetupTest method is called before every test on the suite.
func (helper *SessionHandlerTestHelper) SetupTest() {
    helper.endpoints.sessionProvider.Clear()
}

// This function launches the testing suite.
func TestSessionHandlerSuite(t *testing.T) {
    suite.Run(t, new(SessionHandlerTestHelper))
}

// Obtain a target URL given an endpoint.
//   params:
//     endpoint The URL endpoint.
func (helper *SessionHandlerTestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

func (helper *SessionHandlerTestHelper) TestAddSession() {

    testSession := entities.NewSession(testUserId, testExpirationTime)
    testCookie := http.Cookie{Domain: testDomain}
    testSession.AddCookie(testCookieName, testCookie)

    toAdd := entities.NewAddSessionRequest(*testSession)

    f := frisby.Create("Add a new session").
        Post(helper.getURL("session/add")).
        SetJson(toAdd).Send()

    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        sessionID, err := json.Get("id").String()
        assert.Nil(helper.T(), err, "expecting id")
        assert.NotEmpty(helper.T(), sessionID, "a user id must be returned")
    })
    f.ExpectJson("userId", testUserId)

}

func (helper *SessionHandlerTestHelper) TestGetSession() {
    // Add it directly using the provider
    testSession := entities.NewSession(testUserId, testExpirationTime)
    testCookie := http.Cookie{Domain: testDomain}
    testSession.AddCookie(testCookieName, testCookie)

    helper.endpoints.sessionProvider.Add(*testSession)

    f := frisby.Create("Get a session").
        Get(helper.getURL(fmt.Sprintf("session/%s/get", testSession.ID))).
        Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        userID, err := json.Get("userId").String()
        assert.Nil(helper.T(), err, "expecting id")
        assert.NotEmpty(helper.T(), userID, "a user id must be returned")
        assert.Equal(helper.T(), userID, testSession.UserID)
    })
    f.ExpectJson("userId", testUserId)
}


func (helper *SessionHandlerTestHelper) TestDeleteSession() {

    // Add it directly using the provider
    testSession := entities.NewSession(testUserId, testExpirationTime)
    testCookie := http.Cookie{Domain: testDomain}
    testSession.AddCookie(testCookieName, testCookie)

    helper.endpoints.sessionProvider.Add(*testSession)

    f := frisby.Create("Delete a session").
        Delete(helper.getURL(fmt.Sprintf("session/%s/delete", testSession.ID))).
        Send()
    f.ExpectStatus(http.StatusOK)
}


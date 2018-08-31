//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// User handler integration tests.

package user

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
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/daishogroup/system-model/provider/accessstorage"
    "github.com/daishogroup/system-model/provider/passwordstorage"
    "github.com/daishogroup/system-model/provider/oauthstorage"
    "github.com/daishogroup/system-model/entities"
    "github.com/bitly/go-simplejson"
    "github.com/daishogroup/dhttp"
)

const (
    // Local address for the server exposing the API.
    BaseAddress = "localhost"
)

// The EndpointHelper structure contains all the elements in the endpoint processing path: handler, manager and
// providers so they can be used if needed during the tests.
type EndpointHelper struct {
    userProvider     *userstorage.MockupUserProvider
    accessProvider   *accessstorage.MockupUserAccessProvider
    passwordProvider *passwordstorage.MockupPasswordProvider
    oauthProvider    *oauthstorage.MockupOAuthProvider
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
    var userProvider = userstorage.NewMockupUserProvider()
    var accessProvider = accessstorage.NewMockupUserAccessProvider()
    var passwordProvider = passwordstorage.NewMockupPasswordProvider()
    var oauthProvider = oauthstorage.NewMockupOAuthProvider()
    var userMgr = NewManager(userProvider, accessProvider, passwordProvider, oauthProvider)
    var userHandler = NewHandler(userMgr)
    port, _ := dhttp.GetAvailablePort()

    handler := mux.NewRouter()
    userHandler.SetRoutes(handler)

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
    return EndpointHelper{userProvider, accessProvider,
        passwordProvider, oauthProvider,
        userMgr, userHandler, port, handler, srv}
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
    dhttp.WaitURLAvailable(BaseAddress,helper.port,5,"/", 1)
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
    log.Info("User handler test suite")
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.endpoints.Start()
}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *UserHandlerTestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    assert.Equal(helper.T(), 0, frisby.Global.NumErrored, "expecting 0 failures")
}

// The SetupTest method is called before every test on the suite.
func (helper *UserHandlerTestHelper) SetupTest() {
    helper.endpoints.userProvider.Clear()
    helper.endpoints.accessProvider.Clear()
    helper.endpoints.passwordProvider.Clear()
    helper.endpoints.oauthProvider.Clear()
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

func (helper *UserHandlerTestHelper) TestAddUser() {

    toAdd := entities.NewAddUserRequest(testUserId, testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime)
    f := frisby.Create("Add a new user").
        Post(helper.getURL("user/add")).
        SetJson(toAdd).Send()

    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        userID, err := json.Get("userId").String()
        assert.Nil(helper.T(), err, "expecting id")
        assert.NotEmpty(helper.T(), userID, "a user id must be returned")
        creationTime, err := json.Get("creationTime").String()
        assert.Nil(helper.T(), err, "expecting creation time")
        assert.Equal(helper.T(), testCreationTime.Format(time.RFC3339), creationTime)
        expirationTime, err := json.Get("expirationTime").String()
        assert.Nil(helper.T(), err, "expecting expiration time")
        assert.Equal(helper.T(), testExpirationTime.Format(time.RFC3339), expirationTime)

    })
    f.ExpectJson("userId", testUserId)
    f.ExpectJson("name", testUserName)
    f.ExpectJson("email", testUserEmail)
    f.ExpectJson("phone", testUserPhone)
    f.ExpectJson("creationTime", testCreationTime.Format(time.RFC3339))
    f.ExpectJson("expirationTime", testExpirationTime.Format(time.RFC3339))

}

func (helper *UserHandlerTestHelper) TestGetUser() {

    // add directly a user to to the provider
    testUser := entities.NewUserWithID(testUserId, testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime)
    helper.endpoints.userProvider.Add(*testUser)

    f := frisby.Create("Get a user").
        Get(helper.getURL(fmt.Sprintf("user/%s/get", testUserId))).
        Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        userID, err := json.Get("userId").String()
        assert.Nil(helper.T(), err, "expecting id")
        assert.NotEmpty(helper.T(), userID, "a user id must be returned")
    })
    f.ExpectJson("userId", testUserId)
    f.ExpectJson("name", testUserName)
    f.ExpectJson("email", testUserEmail)
    f.ExpectJson("phone", testUserPhone)
}

func (helper *UserHandlerTestHelper) TestGetNonExistingUser() {
    f := frisby.Create("Get a user").
        Get(helper.getURL("user/unknown/get")).
        Send()
    f.ExpectStatus(http.StatusInternalServerError)
}

func (helper *UserHandlerTestHelper) TestDeleteUser() {
    p := "password"
    // add directly a user to to the provider
    testUser := entities.NewUserWithID(testUserId, testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime)
    helper.endpoints.userProvider.Add(*testUser)
    testRole := entities.NewUserAccess(testUserId, []entities.RoleType{entities.DeveloperType})
    helper.endpoints.accessProvider.Add(*testRole)
    testOAuth := entities.NewOAuthSecrets(testUserId)
    helper.endpoints.oauthProvider.Add(testOAuth)
    testPassword, _ := entities.NewPassword(testUserId, &p)
    helper.endpoints.passwordProvider.Add(*testPassword)

    f := frisby.Create("Get a user").
        Delete(helper.getURL(fmt.Sprintf("user/%s/delete", testUserId))).
        Send()
    f.ExpectStatus(http.StatusOK)
}

func (helper *UserHandlerTestHelper) TestDeleteNonExistingUser() {
    // add directly a user to to the provider

    f := frisby.Create("Delete a user").
        Delete(helper.getURL("user/nonexistinguser/delete")).
        Send()
    f.ExpectStatus(http.StatusInternalServerError)
}

func (helper *UserHandlerTestHelper) TestUpdateUser() {

    // add directly a user to to the provider
    testUser := entities.NewUserWithID(testUserId, testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime)
    helper.endpoints.userProvider.Add(*testUser)
    testRole := entities.NewUserAccess(testUserId, []entities.RoleType{entities.DeveloperType})
    helper.endpoints.accessProvider.Add(*testRole)
    testOAuth := entities.NewOAuthSecrets(testUserId)
    helper.endpoints.oauthProvider.Add(testOAuth)

    toAdd := entities.NewUpdateUserRequest().WithName("1").
        WithPhone("2").WithEmail("3")
    f := frisby.Create("Update an existing user").
        Post(helper.getURL(fmt.Sprintf("user/%s/update", testUserId))).SetJson(toAdd).Send()

    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        userID, err := json.Get("userId").String()
        assert.Nil(helper.T(), err, "expecting id")
        assert.NotEmpty(helper.T(), userID, "a user id must be returned")
    })
    f.ExpectJson("userId", testUserId)
    f.ExpectJson("name", "1")
    f.ExpectJson("phone", "2")
    f.ExpectJson("email", "3")
}

func (helper *UserHandlerTestHelper) TestListUsers() {
    // add directly a user to to the provider
    testUser := entities.NewUserWithID(testUserId, testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime)
    helper.endpoints.userProvider.Add(*testUser)
    testUser2 := entities.NewUserWithID("a", testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime)
    helper.endpoints.userProvider.Add(*testUser2)
    testOAuth := entities.NewOAuthSecrets(testUserId)
    helper.endpoints.oauthProvider.Add(testOAuth)

    // add roles
    helper.endpoints.accessProvider.Add(*entities.NewUserAccess(testUserId, []entities.RoleType{entities.GlobalAdmin}))
    helper.endpoints.accessProvider.Add(*entities.NewUserAccess("a", []entities.RoleType{entities.GlobalAdmin}))

    f := frisby.Create("List users").Get(helper.getURL("user/list")).Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

        userA, err := json.GetIndex(0).Get("userId").String()
        assert.Nil(helper.T(), err, "unexpected error during unmarshalling")
        userB, err := json.GetIndex(1).Get("userId").String()
        assert.Nil(helper.T(), err, "unexpected error during unmarshalling")

        assert.Equal(helper.T(), "a", userA, "unexpected username")
        assert.Equal(helper.T(), testUserId, userB, "unexpected username")
    })
}

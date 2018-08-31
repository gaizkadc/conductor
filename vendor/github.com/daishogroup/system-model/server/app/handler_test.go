//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package app

import (
    "context"
    "fmt"
    "net/http"
    "strconv"
    "testing"
    "time"

    "github.com/bitly/go-simplejson"
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/suite"
    "github.com/verdverm/frisby"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/dhttp"
)

const (
    BaseAddress           = "localhost"
    testAppDescriptorID   = "appDescId"
    testAppPersistentSize = "1Gb"
    testAppStorage        = entities.AppStorageDefault
)

type AppEndpointHelper struct {
    networkProvider *networkstorage.MockupNetworkProvider
    appDescProvider *appdescstorage.MockupAppDescProvider
    appInstProvider *appinststorage.MockupAppInstProvider
    appMgr          Manager
    appHandler      Handler
    port            int
    handler         http.Handler
    srv             *http.Server
}

// Create a new EndpointHelper.
//   returns:
//     A new endpoint helper with a mockup provider.
func NewAppEndpointHelper() AppEndpointHelper {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var appDescProvider = appdescstorage.NewMockupAppDescProvider()
    var appInstProvider = appinststorage.NewMockupAppInstProvider()
    var appManager = NewManager(networkProvider, appDescProvider, appInstProvider)
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
        networkProvider,
        appDescProvider,
        appInstProvider,
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
    helper.endpoints.networkProvider.Clear()
    helper.endpoints.appDescProvider.Clear()
    helper.endpoints.appInstProvider.Clear()
}

// Obtain a target URL given an endpoint.
// params:
// endpoint The URL endpoint.
func (helper *AppTestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

func (helper *AppTestHelper) addTestingNetwork(id string) {
    var toAdd = entities.NewNetworkWithID(
        id, testNetworkName, testDescription, testAdminName, testAdminPhone, testAdminEmail)
    helper.endpoints.networkProvider.Add(*toAdd)
}

func (helper *AppTestHelper) getTestAddDescriptorRequest() *entities.AddAppDescriptorRequest {
    return entities.NewAddAppDescriptorRequest(
        testAppDescName, testDescription, testServiceName, testServiceVersion, testLabel, testPort, []string{testImage})
}

func (helper *AppTestHelper) getTestAddInstRequest() *entities.AddAppInstanceRequest {
    return entities.NewAddAppInstanceRequest(
        testAppDescriptorID, testAppDescName, testDescription, testLabel, testArguments,
        testAppPersistentSize, testAppStorage)
}

func TestAppHandlerSuite(t *testing.T) {
    suite.Run(t, new(AppTestHelper))
}

func (helper *AppTestHelper) TestAddDescriptor() {
    networkID := "TestAddDescriptor"
    helper.addTestingNetwork(networkID)

    toAdd := helper.getTestAddDescriptorRequest()
    f := frisby.Create("Add an app descriptor").
        Post(helper.getURL("app/" + networkID + "/descriptor/add")).SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        id, err := json.Get("id").String()
        helper.Nil(err, "expecting id")
        helper.NotEmpty(id, "an id must be returned")
    })

    f.ExpectJson("name", testAppDescName)
    f.ExpectJson("description", testDescription)
    f.ExpectJson("serviceName", testServiceName)
    f.ExpectJson("serviceVersion", testServiceVersion)
}

func (helper *AppTestHelper) TestListDescriptors() {
    networkID := "TestListDescriptors"
    helper.addTestingNetwork(networkID)

    toAdd := helper.getTestAddDescriptorRequest()
    f := frisby.Create("Add an app descriptor").
        Post(helper.getURL("app/" + networkID + "/descriptor/add")).SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)

    fl := frisby.Create("List descriptors").
        Get(helper.getURL("app/" + networkID + "/descriptor/list")).Send()

    fl.ExpectStatus(http.StatusOK)
    fl.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        entries, err := json.Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(entries), "length does not match")
    })

}

func (helper *AppTestHelper) TestGetDescriptor() {
    networkID := "TestGetDescriptor"
    helper.addTestingNetwork(networkID)

    toAdd := helper.getTestAddDescriptorRequest()
    f := frisby.Create("Add an app descriptor").
        Post(helper.getURL("app/" + networkID + "/descriptor/add")).SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)

    result, _ := f.Resp.Json()
    descriptorID, _ := result.Get("id").String()

    fl := frisby.Create("Get descriptor").
        Get(helper.getURL("app/" + networkID + "/descriptor/" + descriptorID + "/info")).Send()

    fl.ExpectStatus(http.StatusOK)
    fl.ExpectJson("id", descriptorID)
}

func (helper *AppTestHelper) TestDeleteDescriptor() {
    networkID := "TestDeleteDescriptor"
    helper.addTestingNetwork(networkID)

    toAdd := helper.getTestAddDescriptorRequest()
    f := frisby.Create("Add an app descriptor").
        Post(helper.getURL("app/" + networkID + "/descriptor/add")).SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)

    result, _ := f.Resp.Json()
    descriptorID, _ := result.Get("id").String()

    fl := frisby.Create("Delete descriptor").
        Delete(helper.getURL("app/" + networkID + "/descriptor/" + descriptorID + "/delete")).Send()

    fl.ExpectStatus(http.StatusOK)
}

func (helper *AppTestHelper) TestAddInstance() {
    networkID := "TestAddInstance"
    helper.addTestingNetwork(networkID)

    toAdd := helper.getTestAddInstRequest()
    f := frisby.Create("Add an app instance").
        Post(helper.getURL("app/" + networkID + "/instance/add")).SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)

    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        id, err := json.Get("deployedId").String()
        helper.Nil(err, "expecting id")
        helper.NotEmpty(id, "an id must be returned")
    })

    f.ExpectJson("appDescriptorId", testAppDescriptorID)
    f.ExpectJson("name", testAppDescName)
    f.ExpectJson("description", testDescription)
    f.ExpectJson("label", testLabel)
    f.ExpectJson("arguments", testArguments)

}

func (helper *AppTestHelper) TestListInstances() {
    networkID := "TestListInstances"
    helper.addTestingNetwork(networkID)

    toAdd := helper.getTestAddInstRequest()
    f := frisby.Create("Add an app instance").
        Post(helper.getURL("app/" + networkID + "/instance/add")).SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)

    fl := frisby.Create("List instances").Get(helper.getURL("app/" + networkID + "/instance/list")).Send()

    fl.ExpectStatus(http.StatusOK)
    fl.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        entries, err := json.Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(entries), "length does not match")
    })
}

func (helper *AppTestHelper) TestGetInstance() {
    networkID := "TestGetInstance"
    helper.addTestingNetwork(networkID)

    toAdd := helper.getTestAddInstRequest()
    f := frisby.Create("Add an app instance").
        Post(helper.getURL("app/" + networkID + "/instance/add")).SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)

    result, _ := f.Resp.Json()
    deployedID, _ := result.Get("deployedId").String()

    fl := frisby.Create("Get instance").
        Get(helper.getURL("app/" + networkID + "/instance/" + deployedID + "/info")).Send()

    fl.ExpectStatus(http.StatusOK)
    fl.ExpectJson("deployedId", deployedID)
}

func (helper *AppTestHelper) TestUpdateInstance() {
    networkID := "TestUpdateInstance"
    helper.addTestingNetwork(networkID)

    toAdd := helper.getTestAddInstRequest()
    f := frisby.Create("Add an app instance").
        Post(helper.getURL("app/" + networkID + "/instance/add")).SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)

    result, _ := f.Resp.Json()
    deployedID, _ := result.Get("deployedId").String()

    update := entities.NewUpdateAppInstRequest().WithClusterID("newCluster")

    fu := frisby.Create("Update instance").
        Post(helper.getURL("app/" + networkID + "/instance/" + deployedID + "/update")).SetJson(update).Send()
    fu.ExpectStatus(http.StatusOK)
    fu.ExpectJson("deployedId", deployedID)
    fu.ExpectJson("clusterId", "newCluster")

}

func (helper *AppTestHelper) TestDeleteInstance() {
    networkID := "TestDeleteInstance"
    helper.addTestingNetwork(networkID)

    toAdd := helper.getTestAddInstRequest()
    f := frisby.Create("Add an app instance").
        Post(helper.getURL("app/" + networkID + "/instance/add")).SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)

    result, _ := f.Resp.Json()
    deployedID, _ := result.Get("deployedId").String()

    fu := frisby.Create("Remove instance").
        Delete(helper.getURL("app/" + networkID + "/instance/" + deployedID + "/delete")).Send()
    fu.ExpectStatus(http.StatusOK)

}

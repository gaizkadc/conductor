//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Cluster handler integration tests.

package cluster

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
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
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
    appInstProvider *appinststorage.MockupAppInstProvider
    clusterMgr      Manager
    clusterHandler  Handler
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
    var appInstProvider = appinststorage.NewMockupAppInstProvider()
    var clusterMgr = NewManager(networkProvider, clusterProvider, appInstProvider)
    var clusterHandler = NewHandler(clusterMgr)
    port, _ := dhttp.GetAvailablePort()
    handler := mux.NewRouter()
    clusterHandler.SetRoutes(handler)

    var srv = &http.Server{
        Handler: handler,
        Addr:    BaseAddress + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    return EndpointHelper{networkProvider, clusterProvider, appInstProvider, clusterMgr,
        clusterHandler, port, handler, srv}
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
type ClusterTestHelper struct {
    suite.Suite
    endpoints EndpointHelper
}

// The SetupSuite is the first method called and defines the testing suite. It creates the endpoints and launches
// a server listening for client requests.
func (helper *ClusterTestHelper) SetupSuite() {
    log.Info("Cluster handler test suite")
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.endpoints.Start()
    dhttp.WaitURLAvailable(BaseAddress,helper.endpoints.port,5,"/", 1)
}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *ClusterTestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    helper.Equal(0, frisby.Global.NumErrored, "expecting 0 failures")
}

// The SetupTest method is called before every test on the suite.
func (helper *ClusterTestHelper) SetupTest() {
    helper.endpoints.networkProvider.Clear()
    helper.endpoints.clusterProvider.Clear()
    helper.endpoints.appInstProvider.Clear()
}

// This function launches the testing suite.
func TestClusterHandlerSuite(t *testing.T) {
    suite.Run(t, new(ClusterTestHelper))
}

// Obtain a target URL given an endpoint.
//   params:
//     endpoint The URL endpoint.
func (helper *ClusterTestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

func (helper *ClusterTestHelper) addTestNetwork(id string) {
    var toAdd = entities.NewNetworkWithID(
        id, testNetworkName, testDescription, testAdminName, testAdminPhone, testAdminEmail)
    helper.endpoints.networkProvider.Add(*toAdd)
}

func (helper *ClusterTestHelper) addTestCluster(networkID string, newCluster entities.Cluster) {
    var toAdd = entities.NewNetworkWithID(
        networkID, testNetworkName, testDescription, testAdminName, testAdminPhone, testAdminEmail)
    helper.endpoints.networkProvider.Add(*toAdd)
    helper.endpoints.clusterProvider.Add(newCluster)
    helper.endpoints.networkProvider.AttachCluster(networkID, newCluster.ID)
}

func (helper *ClusterTestHelper) TestAddCluster() {
    networkID := "TestAddCluster"
    helper.addTestNetwork(networkID)
    toAdd := entities.NewAddClusterRequest(testClusterName, testDescription, entities.CloudType,
        testLocation, testAdminEmail)
    f := frisby.Create("Add a new cluster").
        Post(helper.getURL("cluster/" + networkID + "/add")).
        SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        networkID, err := json.Get("id").String()
        helper.Nil(err, "expecting id")
        helper.NotEmpty(networkID, "a network id must be returned")
    })

    f.ExpectJson("name", testClusterName)
    f.ExpectJson("description", testDescription)
    f.ExpectJson("type", string(entities.CloudType))
    f.ExpectJson("location", testLocation)

}

func (helper *ClusterTestHelper) TestAddClusterInvalidType() {
    networkID := "TestAddClusterInvalidType"
    helper.addTestNetwork(networkID)
    toAdd := entities.NewAddClusterRequest(testClusterName, testDescription, "",
        testLocation, testAdminEmail)
    f := frisby.Create("Add a new cluster").
        Post(helper.getURL("cluster/" + networkID + "/add")).
        SetJson(toAdd).Send()
    f.ExpectStatus(http.StatusBadRequest)
}

func (helper *ClusterTestHelper) TestListCluster() {
    helper.TestAddCluster()
    f := frisby.Create("List clusters").Get(helper.getURL("cluster/TestAddCluster/list")).Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        entries, err := json.Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(entries), "length does not match")
    })
}

func (helper *ClusterTestHelper) TestGetCluster() {
    // Add a new network
    networkID := "testGetCluster"
    helper.addTestNetwork(networkID)
    // Add a new cluster
    clusterID := "testClusterId"

    newCluster := entities.NewClusterWithID(networkID, clusterID, testClusterName,
        testDescription, entities.CloudType, testLocation,
        testAdminEmail, entities.ClusterCreated,
        false, false)
    helper.addTestCluster(networkID, *newCluster)

    // Retrieve the cluster
    targetEndpoint := "cluster/" + networkID + "/" + clusterID + "/info"
    f := frisby.Create("Get cluster").Get(helper.getURL(targetEndpoint)).Send()
    helper.Equal(http.StatusOK, f.Resp.StatusCode)
    f.ExpectJson("id", clusterID)

}

func (helper *ClusterTestHelper) TestDeleteCluster() {
    // Add a new network
    networkID := "testDeleteCluster"
    helper.addTestNetwork(networkID)
    // Add a new cluster
    clusterID := "testClusterId"
    newCluster := entities.NewClusterWithID(networkID, clusterID, testClusterName, testDescription,
        entities.CloudType, testLocation, testAdminEmail,
        entities.ClusterInstalled, false, false)
    helper.addTestCluster(networkID, *newCluster)

    // Remove the cluster
    targetEndpoint := "cluster/" + networkID + "/" + clusterID + "/delete"
    f := frisby.Create("Delete cluster").Delete(helper.getURL(targetEndpoint)).Send()
    f.ExpectStatus(http.StatusOK)
}

func (helper *ClusterTestHelper) TestDeleteClusterWithNodes() {
    // Add a new network
    networkID := "testDeleteClusterWithNodes"
    helper.addTestNetwork(networkID)
    // Add a new cluster
    clusterID := "testClusterId"
    newCluster := entities.NewClusterWithID(networkID, clusterID, testClusterName, testDescription,
        entities.CloudType, testLocation, testAdminEmail,
        entities.ClusterInstalled, false, false)

    helper.addTestCluster(networkID, *newCluster)
    // Add some nodes
    helper.endpoints.clusterMgr.clusterProvider.AttachNode(clusterID, "node1")

    // Remove the cluster
    targetEndpoint := "cluster/" + networkID + "/" + clusterID + "/delete"
    f := frisby.Create("Delete cluster with nodes").Delete(helper.getURL(targetEndpoint)).Send()
    f.ExpectStatus(http.StatusInternalServerError)
}

func (helper *ClusterTestHelper) TestUpdateCluster() {
    // Add a new network
    networkID := "testUpdateCluster"
    helper.addTestNetwork(networkID)
    // Add a new cluster
    clusterID := "testClusterId"
    newCluster := entities.NewClusterWithID(networkID, clusterID, testClusterName,
        testDescription, entities.CloudType, testLocation, testAdminEmail,
        entities.ClusterCreated, false, false)
    helper.addTestCluster(networkID, * newCluster)
    // Remove the cluster
    targetEndpoint := "cluster/" + networkID + "/" + clusterID + "/update"

    updateRequest := entities.NewUpdateClusterRequest().WithClusterStatus(entities.ClusterInstalled)
    f := frisby.Create("Update cluster").Post(helper.getURL(targetEndpoint)).SetJson(updateRequest).Send()
    f.ExpectStatus(http.StatusOK)

    updated, err := helper.endpoints.clusterProvider.RetrieveCluster(clusterID)
    helper.Nil(err, "error should be nil")

    helper.Equal(networkID, updated.NetworkID)
    helper.Equal(clusterID, updated.ID)
    helper.Equal(testClusterName, updated.Name)
    helper.Equal(testDescription, updated.Description)
    helper.Equal(entities.CloudType, updated.Type)
    helper.Equal(testLocation, updated.Location)
    helper.Equal(testAdminEmail, updated.Email)
    helper.Equal(entities.ClusterInstalled, updated.Status)
    helper.Equal(newCluster.Drain, updated.Drain)
    helper.Equal(newCluster.Cordon, updated.Cordon)

}

func (helper *ClusterTestHelper) TestUpdateClusterInvalidStatus() {
    // Add a new network
    networkID := "TestUpdateClusterInvalidStatus"
    helper.addTestNetwork(networkID)
    // Add a new cluster
    clusterID := "testClusterId"
    newCluster := entities.NewClusterWithID(networkID, clusterID, testClusterName,
        testDescription, entities.CloudType, testLocation, testAdminEmail,
        entities.ClusterCreated, false, false)
    helper.addTestCluster(networkID, * newCluster)
    // Remove the cluster
    targetEndpoint := "cluster/" + networkID + "/" + clusterID + "/update"

    updateRequest := entities.NewUpdateClusterRequest().WithClusterStatus("Invalid")
    f := frisby.Create("Update cluster").Post(helper.getURL(targetEndpoint)).SetJson(updateRequest).Send()
    f.ExpectStatus(http.StatusBadRequest)

}

func (helper *ClusterTestHelper) TestUpdateClusterInvalidType() {
    // Add a new network
    networkID := "TestUpdateClusterInvalidType"
    helper.addTestNetwork(networkID)
    // Add a new cluster
    clusterID := "testClusterId"
    newCluster := entities.NewClusterWithID(networkID, clusterID, testClusterName,
        testDescription, entities.CloudType, testLocation, testAdminEmail,
        entities.ClusterCreated, false, false)
    helper.addTestCluster(networkID, * newCluster)
    // Remove the cluster
    targetEndpoint := "cluster/" + networkID + "/" + clusterID + "/update"

    updateRequest := entities.NewUpdateClusterRequest().WithType("Invalid")
    f := frisby.Create("Update cluster").Post(helper.getURL(targetEndpoint)).SetJson(updateRequest).Send()
    f.ExpectStatus(http.StatusBadRequest)

}

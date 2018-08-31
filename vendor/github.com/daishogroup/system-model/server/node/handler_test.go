//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Node handler integration tests.

package node

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

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "context"
    "strconv"
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
    nodeMgr         Manager
    nodeHandler     Handler
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
    var nodeMgr = NewManager(networkProvider, clusterProvider, nodeProvider)
    var nodeHandler = NewHandler(nodeMgr)
    port, _ := dhttp.GetAvailablePort()

    handler := mux.NewRouter()
    nodeHandler.SetRoutes(handler)

    var srv = &http.Server{
        Handler: handler,
        Addr:    BaseAddress + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    return EndpointHelper{networkProvider, clusterProvider, nodeProvider,
        nodeMgr, nodeHandler, port, handler, srv}
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
type NodeHandlerTestHelper struct {
    suite.Suite
    endpoints EndpointHelper
}

// The SetupSuite is the first method called and defines the testing suite. It creates the endpoints and launches
// a server listening for client requests.
func (helper *NodeHandlerTestHelper) SetupSuite() {
    log.Info("Node handler test suite")
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.endpoints.Start()
    dhttp.WaitURLAvailable(BaseAddress,helper.endpoints.port,5,"/", 1)
}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *NodeHandlerTestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    assert.Equal(helper.T(), 0, frisby.Global.NumErrored, "expecting 0 failures")
}

// The SetupTest method is called before every test on the suite.
func (helper *NodeHandlerTestHelper) SetupTest() {
    helper.endpoints.networkProvider.Clear()
    helper.endpoints.clusterProvider.Clear()
}

// This function launches the testing suite.
func TestClusterHandlerSuite(t *testing.T) {
    suite.Run(t, new(NodeHandlerTestHelper))
}

// Obtain a target URL given an endpoint.
//   params:
//     endpoint The URL endpoint.
func (helper *NodeHandlerTestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

func (helper *NodeHandlerTestHelper) addTestingCluster(networkID string, clusterID string) {
    network := entities.NewNetworkWithID(
        networkID, testNetworkName, testDescription, testAdminName, testAdminPhone, testAdminEmail)
    cluster := entities.NewClusterWithID(networkID, clusterID, testClusterName, testDescription, entities.CloudType,
        testLocation, testAdminEmail,entities.ClusterCreated,false,false)
    helper.endpoints.networkProvider.Add(*network)
    helper.endpoints.networkProvider.AttachCluster(networkID, clusterID)
    helper.endpoints.clusterProvider.Add(*cluster)
}

func (helper *NodeHandlerTestHelper) addTestingNode(networkID string, clusterID string, nodeID string) {
    network := entities.NewNetworkWithID(
        networkID, testNetworkName, testDescription, testAdminName, testAdminPhone, testAdminEmail)
    cluster := entities.NewClusterWithID(networkID, clusterID, testClusterName, testDescription, entities.CloudType,
        testLocation, testAdminEmail,entities.ClusterCreated,false,false)
    node := entities.NewNodeWithID(networkID, clusterID, nodeID, testNodeName, testDescription, make([]string, 0),
        testIP, testIP, true, testUsername, "", "", entities.NodeInstalled)
    helper.endpoints.networkProvider.Add(*network)
    helper.endpoints.networkProvider.AttachCluster(networkID, clusterID)
    helper.endpoints.clusterProvider.Add(*cluster)
    helper.endpoints.nodeProvider.Add(*node)
    helper.endpoints.clusterProvider.AttachNode(clusterID, nodeID)
}

func (helper *NodeHandlerTestHelper) TestAddNode() {
    networkID := "TestAddNetwork"
    clusterID := "TestAddCluster"
    helper.addTestingCluster(networkID, clusterID)
    toAdd := entities.NewAddNodeRequest(
        testNodeName, testDescription,  []string{"l1"}, testIP, testIP, true, testUsername, testPassword, testSSHKey)
    f := frisby.Create("Add a new node").
        Post(helper.getURL("node/" + networkID + "/" + clusterID + "/add")).
        SetJson(toAdd).Send()

    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        networkID, err := json.Get("id").String()
        assert.Nil(helper.T(), err, "expecting id")
        assert.NotEmpty(helper.T(), networkID, "a network id must be returned")
    })

    f.ExpectJson("name", testNodeName)
    f.ExpectJson("description", testDescription)

}

func (helper *NodeHandlerTestHelper) TestListNode() {
    networkID := "TestAddNetwork"
    clusterID := "TestAddCluster"
    helper.TestAddNode()
    f := frisby.Create("List nodes").Get(
        helper.getURL("node/" + networkID + "/" + clusterID + "/list")).Send()
    assert.Equal(helper.T(), http.StatusOK, f.Resp.StatusCode)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        entries, err := json.Array()
        assert.Nil(helper.T(), err, "expecting array")
        assert.Equal(helper.T(), 1, len(entries), "length does not match")
    })
    f.ExpectStatus(http.StatusOK)
}

func (helper *NodeHandlerTestHelper) TestGetNode() {
    // Add a new network
    networkID := "testGetNetwork"
    clusterID := "testGetCluster"
    helper.addTestingCluster(networkID, clusterID)

    newNode := entities.NewAddNodeRequest(testNodeName, testDescription,  []string{"l1"},
        testIP, testIP, true,
        testUsername, testPassword, testSSHKey)
    addNode, err := helper.endpoints.nodeMgr.AddNode(networkID, clusterID, *newNode)

    helper.Nil(err, "error must be nil")
    helper.NotNil(addNode, "node must not be nil")

    // Retrieve the cluster
    targetEndpoint := "node/" + networkID + "/" + clusterID + "/" + addNode.ID + "/info"
    f := frisby.Create("Get node").Get(helper.getURL(targetEndpoint)).Send()
    f.ExpectStatus(http.StatusOK)
    f.ExpectJson("id", addNode.ID)
}

func (helper *NodeHandlerTestHelper) TestDeleteNode() {
    // Add a new network
    networkID := "testGetNetwork"
    clusterID := "testGetCluster"
    helper.addTestingCluster(networkID, clusterID)

    // Add a new cluster
    newNode := entities.NewAddNodeRequest(testNodeName, testDescription, []string{"l1"},
        testIP, testIP, true,
        testUsername, testPassword, testSSHKey)
    addNode, err := helper.endpoints.nodeMgr.AddNode(networkID, clusterID, *newNode)
    helper.Nil(err, "error must be nil")
    helper.NotNil(addNode, "node must not be nil")

    // Remove the cluster
    targetEndpoint := "node/" + networkID + "/" + clusterID + "/" + addNode.ID + "/delete"
    f := frisby.Create("Delete node").Delete(helper.getURL(targetEndpoint)).Send()

    f.ExpectStatus(http.StatusOK)

    _, err = helper.endpoints.nodeProvider.RetrieveNode(addNode.ID)
    assert.NotNil(helper.T(), err, "cluster should have been removed")
}

func (helper *NodeHandlerTestHelper) TestInvalidUpdateStatus() {
    networkID := "testUpdateNetwork"
    clusterID := "testUpdateCluster"
    nodeID := "testInvalidNodeStatus"
    helper.addTestingNode(networkID, clusterID, nodeID)
    updateRequest := entities.NewUpdateNodeRequest().WithStatus("Invalid")
    targetEndpoint := "node/" + networkID + "/" + clusterID + "/" + nodeID + "/update"
    f := frisby.Create("Update node").Post(helper.getURL(targetEndpoint)).SetJson(updateRequest).Send()
    f.ExpectStatus(http.StatusBadRequest)
}

func (helper *NodeHandlerTestHelper) TestFilterNode() {
    networkID := "TestAddNetwork"
    clusterID := "TestAddCluster"
    helper.addTestingCluster(networkID, clusterID)
    toAdd := entities.NewAddNodeRequest(
        testNodeName, testDescription,  []string{"l1"}, testIP, testIP, true, testUsername, testPassword, testSSHKey)
    f := frisby.Create("Add a new node").
        Post(helper.getURL("node/" + networkID + "/" + clusterID + "/add")).
        SetJson(toAdd).Send()

    f.ExpectStatus(http.StatusOK)
    nodeID := ""
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        nID, err := json.Get("id").String()
        assert.Nil(helper.T(), err, "expecting id")
        assert.NotEmpty(helper.T(), nID, "a node id must be returned")
        nodeID = nID
    })

    filter := entities.NewFilterNodesRequest().ByLabel("l1")
    filterRequest := frisby.Create("filter nodes").
        Get(helper.getURL("node/" + networkID + "/" + clusterID + "/filter")).
        SetJson(filter).Send()
    filterRequest.ExpectStatus(http.StatusOK)
    filterRequest.ExpectJsonLength("", 1)
}

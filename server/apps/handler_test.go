//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Network handler integration tests.

package apps

import (
    "fmt"
    "net/http"
    "testing"

    "github.com/bitly/go-simplejson"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/suite"
    "github.com/verdverm/frisby"

    "github.com/daishogroup/conductor/asm"
    entitiesConductor "github.com/daishogroup/conductor/entities"
    logger2 "github.com/daishogroup/conductor/logger"
    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/server"

    "github.com/daishogroup/dhttp"
)

const (
    // Local address for the server exposing API
    TestAddressHost = "localhost"

    // Simply to indicate whether the testing app should run
    TestPort = 5555
)

type EndpointHelper struct {
    systemModelURI     string
    conductorURI       string
    appsHandler        Handler
    networkClient      client.Network
    applicationsClient client.Applications
    clusterClient      client.Cluster
    nodeClient         client.Node
    systemModelService server.Service
}

// Struct containing a set of elements to be defined in order to run the testing
// This is mainly a set of ids
type TestingElements struct {
    network       entities.Network
    cluster1      entities.Cluster
    cluster2      entities.Cluster
    node11        entities.Node
    node12        entities.Node
    node21        entities.Node
    appDescriptor entities.AppDescriptor
}

var testingElements TestingElements

func NewEndpointHelper() EndpointHelper {

    // Get system model URI
    systemModelPort, _ := dhttp.GetAvailablePort()
    systemModelURI := fmt.Sprintf("http://%s:%d", TestAddressHost, systemModelPort)
    logger.Debugf("System model generated API %s", systemModelURI)

    // Get conductor URI
    conductorPort, _ := dhttp.GetAvailablePort()
    conductorURI := fmt.Sprintf("%s:%d", TestAddressHost, conductorPort)
    logger.Debugf("Conductor generated API %s", conductorURI)

    //Get logger URI
    loggerURI := fmt.Sprintf("http://%s:%d", TestAddressHost, systemModelPort)

    // clients
    networkClient := client.NewNetworkRest(systemModelURI)
    applicationsClient := client.NewApplicationRest(systemModelURI)
    clusterClient := client.NewClusterRest(systemModelURI)
    nodeClient := client.NewNodeRest(systemModelURI)
    loggerClient := logger2.NewRestClient(loggerURI)

    var appsMgr = NewRestManager(applicationsClient, clusterClient, nodeClient,
        asm.NewMockupClientFactory(), loggerClient)
    var appsHandler = NewHandler(appsMgr)

    service := server.Service{Configuration: server.Config{Port: uint16(systemModelPort),
        UseInMemoryProviders: true}}

    // Create the service for the system model
    return EndpointHelper{systemModelURI, conductorURI, appsHandler,
        networkClient, applicationsClient, clusterClient, nodeClient,
        service}
}

type AppsTestHelper struct {
    suite.Suite
    endpoints EndpointHelper
    handler   *mux.Router
}

func (helper *AppsTestHelper) SetupSuite() {
    logger.Debug("Creating test suite")
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.handler = mux.NewRouter()
    helper.endpoints.appsHandler.SetRoutes(helper.handler)
    // Run system model api
    helper.endpoints.systemModelService.Run()
    go http.ListenAndServe(helper.endpoints.conductorURI, helper.handler)
}

// Method invoked after testing completion.
func (helper *AppsTestHelper) TearDownSuite() {
    // stop system model
    helper.endpoints.systemModelService.Finalize(true)
    frisby.Global.PrintReport()
    helper.Equal(0, frisby.Global.NumErrored, "expecting 0 failures")
}

// To be run before every test.
func (helper *AppsTestHelper) SetupTest() {
    logger.Debug("Stop system model service")
    helper.endpoints.systemModelService.Finalize(true)
    logger.Debug("Start system model service")
    helper.endpoints.systemModelService.Run()
    logger.Debug("System model started")

    // Add network
    n, err := helper.endpoints.networkClient.Add(*entities.NewAddNetworkRequest("1", "desc1",
        "adminName", "adminPhone", "adminEmail"))
    helper.Nil(err, "Impossible to initialize network")

    // Add cluster
    c1, err := helper.endpoints.clusterClient.Add(n.ID, *entities.NewAddClusterRequest(
        "1", "desc", entities.EdgeType, "Madrid", "email1"))
    helper.Nil(err, "Impossible to add cluster 1")
    // Change to deployed
    update := *entities.NewUpdateClusterRequest().WithClusterStatus(entities.ClusterInstalled)
    helper.endpoints.clusterClient.Update(n.ID, c1.ID, update)

    c2, err := helper.endpoints.clusterClient.Add(n.ID, *entities.NewAddClusterRequest(
        "2", "desc", entities.EdgeType, "Boston", "email2"))
    helper.Nil(err, "Impossible to add cluster 2")
    // change to deployed
    helper.endpoints.clusterClient.Update(n.ID, c2.ID, update)

    // Add descriptor
    appDescriptor, err := helper.endpoints.applicationsClient.AddApplicationDescriptor(n.ID,
        *entities.NewAddAppDescriptorRequest(
            "app1", "desc", "service1", "1.0", "edge", TestPort, []string {"image:version"}))
    helper.Nil(err, "Impossible to initialize descriptor")

    // Add nodes
    labels := []string{"master"}
    n11, err := helper.endpoints.nodeClient.Add(n.ID, c1.ID,
        *entities.NewAddNodeRequest("node1-1", "desc", labels, "0.0.0.0", "0.0.0.0", true,
            "user", "pass", "sshkey"))
    helper.Nil(err, "Impossible to initialize node11")
    n12, err := helper.endpoints.nodeClient.Add(n.ID, c1.ID,
        *entities.NewAddNodeRequest("node1-2", "desc", labels, "0.0.0.1", "0.0.0.0", true,
            "user", "pass", "sshkey"))
    helper.Nil(err, "Impossible to initialize node12")
    n21, err := helper.endpoints.nodeClient.Add(n.ID, c2.ID,
        *entities.NewAddNodeRequest("node2-1", "desc", labels, "0.0.0.2", "0.0.0.0", true,
            "user", "pass", "sshkey"))
    helper.Nil(err, "Impossible to initialize node21")

    testingElements = TestingElements{*n, *c1, *c2, *n11, *n12, *n21,
        *appDescriptor}
}

// This launches the testing suite
func TestAppsHandlerSuite(t *testing.T) {
    logrus.SetLevel(logrus.DebugLevel)
    suite.Run(t, new(AppsTestHelper))
}

// Obtain a target URL given an endpoint.
//   params:
//     uri The URI of the endpoint
//     endpoint The URL endpoint.
func getUrl(uri string, endpoint string) string {
    return fmt.Sprintf("http://%s/api/v0/%s", uri, endpoint)
}

// Test a basic deployment of an application on an empty platform.
func (helper *AppsTestHelper) TestDeploy() {

    logger.Debug("Start testing...")

    //instanceRequest := entities.NewAddAppInstanceRequest(testingElements.appDescriptor.Id,
    //	testingElements.appDescriptor.Name, testingElements.appDescriptor.Description, "label", "arguments")

    instanceRequest := entitiesConductor.NewDeployAppRequest(
        testingElements.appDescriptor.Name,
        testingElements.appDescriptor.ID,
        testingElements.appDescriptor.Description,
        "edge", make(map[string]string, 0), "arguments", "1Gb", entities.AppStorageDefault)

    logger.Debug("testingElements.appDescriptor: " + testingElements.appDescriptor.String())
    logger.Debug("instanceRequest: " + instanceRequest.String())
    logger.Debug("testingElements.network.ID: " + testingElements.network.ID)
    targetUrl := getUrl(helper.endpoints.conductorURI, fmt.Sprintf("app/%s/deploy", testingElements.network.ID))
    logger.Debug("targetURL: " + targetUrl)
    f := frisby.Create("Deploy application on top of network").Post(targetUrl).SetJson(instanceRequest).Send()
    if len(f.Errors()) > 0 {
        logger.Warn("Errors found, ", len(f.Errors()), f.Error())
    }
    helper.Equal(0, len(f.Errors()), "Post should not fail")

    f.ExpectStatus(http.StatusOK)
    f.PrintBody()
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        deployedIP, err := json.Get("clusterAddress").String()
        helper.Nil(err, "Error parsing deployed cluster id")
        helper.Equal(testingElements.node11.PublicIP, deployedIP, "Application not deployed on expected cluster 1")
    })

    // Check application instance is there
    instances, err := helper.endpoints.applicationsClient.ListInstances(testingElements.network.ID)
    helper.Nil(err, "There was an error requesting the available application descriptors")
    helper.Equal(1, len(instances), "Unexpected number of deployed applications")

    // Launch a second application and check we deploy in cluster 2
    f = frisby.Create("Deploy application on top of network").Post(targetUrl).SetJson(instanceRequest).Send()
    helper.Equal(http.StatusOK, f.Resp.StatusCode)

    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        deployedIP, err := json.Get("clusterAddress").String()
        helper.Nil(err, "Error parsing deployed cluster id")
        helper.Equal(testingElements.node21.PublicIP, deployedIP, "Application not deployed on expected cluster 2")
    })
    // Check application instance is there
    instances, err = helper.endpoints.applicationsClient.ListInstances(testingElements.network.ID)
    helper.Nil(err, "There was an error requesting the available application descriptors")
    helper.Equal(2, len(instances), "Unexpected number of deployed applications")

    // Launch a third application and check we deploy in cluster 1 again
    f = frisby.Create("Deploy application on top of network").Post(targetUrl).SetJson(instanceRequest).Send()
    helper.Equal(http.StatusOK, f.Resp.StatusCode)

    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        deployedIP, err := json.Get("clusterAddress").String()
        helper.Nil(err, "Error parsing deployed cluster id")
        helper.Equal(testingElements.node11.PublicIP, deployedIP, "Application not deployed on expected cluster 1 after third application")
    })
    // Check application instance is there
    instances, err = helper.endpoints.applicationsClient.ListInstances(testingElements.network.ID)
    helper.Nil(err, "There was an error requesting the available application descriptors")
    helper.Equal(3, len(instances), "Unexpected number of deployed applications")

}

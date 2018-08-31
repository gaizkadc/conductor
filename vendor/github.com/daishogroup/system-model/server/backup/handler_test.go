package backup

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Backuprestore handler tests.

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
    "github.com/daishogroup/system-model/provider/accessstorage"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/daishogroup/system-model/provider/passwordstorage"
    "github.com/daishogroup/dhttp"
)

const (
    // Local address for the server exposing the API.
    BaseAddress = "localhost"
)

// The EndpointHelper structure contains all the elements in the endpoint processing path: handler, manager and
// providers so they can be used if needed during the tests.
type EndpointHelper struct {
    networkProvider  *networkstorage.MockupNetworkProvider
    clusterProvider  *clusterstorage.MockupClusterProvider
    nodeProvider     *nodestorage.MockupNodeProvider
    appDescProvider  *appdescstorage.MockupAppDescProvider
    userProvider     *userstorage.MockupUserProvider
    accessProvider   *accessstorage.MockupUserAccessProvider
    passwordProvider *passwordstorage.MockupPasswordProvider
    brMgr            Manager
    brHandler        Handler
    port             int
    handler          http.Handler
    srv              *http.Server
}

// Create a new EndpointHelper.
//   returns:
//     A new endpoint helper with a mockup provider.
func NewEndpointHelper() EndpointHelper {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var clusterProvider = clusterstorage.NewMockupClusterProvider()
    var nodeProvider = nodestorage.NewMockupNodeProvider()
    var appDescProvider = appdescstorage.NewMockupAppDescProvider()
    var userProvider = userstorage.NewMockupUserProvider()
    var accessProvider = accessstorage.NewMockupUserAccessProvider()
    var passwordProvider = passwordstorage.NewMockupPasswordProvider()

    var brMgr = NewManager(networkProvider, clusterProvider, nodeProvider, appDescProvider,
        userProvider, accessProvider, passwordProvider)
    var brHandler = NewHandler(brMgr)
    port, _ := dhttp.GetAvailablePort()

    handler := mux.NewRouter()
    brHandler.SetRoutes(handler)

    var srv = &http.Server{
        Handler: handler,
        Addr:    BaseAddress + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    return EndpointHelper{networkProvider, clusterProvider, nodeProvider,
        appDescProvider, userProvider,
        accessProvider, passwordProvider, brMgr, brHandler, port, handler, srv}
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
type BRHandlerTestHelper struct {
    suite.Suite
    endpoints EndpointHelper
}

// The SetupSuite is the first method called and defines the testing suite. It creates the endpoints and launches
// a server listening for client requests.
func (helper *BRHandlerTestHelper) SetupSuite() {
    endpointHelper := NewEndpointHelper()
    helper.endpoints = endpointHelper
    helper.endpoints.Start()
    dhttp.WaitURLAvailable(BaseAddress,helper.endpoints.port,5,"/", 1)

}

// Last method called on the Suite. Use this to shutdown services if required.
func (helper *BRHandlerTestHelper) TearDownSuite() {
    helper.endpoints.Shutdown()
    frisby.Global.PrintReport()
    helper.Equal(0, frisby.Global.NumErrored, "expecting 0 failures")
}

// The SetupTest method is called before every test on the suite.
func (helper *BRHandlerTestHelper) SetupTest() {
    helper.endpoints.networkProvider.Clear()
    helper.endpoints.clusterProvider.Clear()
    helper.endpoints.nodeProvider.Clear()
    helper.endpoints.appDescProvider.Clear()
    helper.endpoints.userProvider.Clear()
    helper.endpoints.accessProvider.Clear()
    helper.endpoints.passwordProvider.Clear()
}

func TestDumpHandlerSuite(t *testing.T) {
    suite.Run(t, new(BRHandlerTestHelper))
}

// Obtain a target URL given an endpoint.
//   params:
//     endpoint The URL endpoint.
func (helper *BRHandlerTestHelper) getURL(endpoint string) string {
    return fmt.Sprintf("http://%s:%d/api/v0/%s", BaseAddress, helper.endpoints.port, endpoint)
}

func (helper *BRHandlerTestHelper) LoadTestData() {
    network := entities.NewNetwork("network1", "", "", "", "")
    cluster := entities.NewCluster(network.ID, "cluster1", "", entities.EdgeType,
        "", "", entities.ClusterInstalled, false, false)
    node := entities.NewNode(network.ID, cluster.ID,
        "node1", "", make([]string, 0),
        "", "", false,
        "", "", "")
    descriptor := entities.NewAppDescriptor(network.ID, "descriptor1", "",
        "", "", "", 0, []string{"repo1:tag1"})
    user := entities.NewUser("user", "9999", "email@email.com", time.Now(), time.Now().Add(time.Hour))

    access := entities.NewUserAccess(user.ID, []entities.RoleType{entities.GlobalAdmin})
    userpassword := "daisho"
    password, err := entities.NewPassword(user.ID, &userpassword)

    helper.Nil(err, "expecting password")

    helper.endpoints.networkProvider.Add(* network)
    helper.endpoints.clusterProvider.Add(* cluster)
    helper.endpoints.nodeProvider.Add(* node)
    helper.endpoints.appDescProvider.Add(* descriptor)
    helper.endpoints.userProvider.Add(* user)
    helper.endpoints.accessProvider.Add(* access)
    helper.endpoints.passwordProvider.Add(* password)

}

func (helper *BRHandlerTestHelper) TestExport() {
    helper.LoadTestData()
    f := frisby.Create("Export").Get(helper.getURL("backup/all/create")).Send()
    f.ExpectStatus(http.StatusOK)
    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        networks, err := json.Get("networks").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(networks), "length does not match")
        clusters, err := json.Get("clusters").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(clusters), "length does not match")
        nodes, err := json.Get("nodes").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(nodes), "length does not match")
        descriptors, err := json.Get("appdesc").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(descriptors), "length does not match")
        //   users, err := json.Get("users").Array()
        //   helper.Equal(1, len(users), "expecting users")
        //   access, err := json.Get("access").Array()
        //   helper.Equal(1, len(access), "expecting access")
        //   password, err := json.Get("password").Array()
        //   helper.Equal(1, len(password), "expecting access")
    })

}

func (helper *BRHandlerTestHelper) LoadImportData() (*entities.BackupRestore) {
    network := entities.NewNetwork("network1", "", "", "", "")
    cluster := entities.NewCluster(network.ID, "cluster1", "", entities.EdgeType,
        "", "", entities.ClusterInstalled, false, false)
    node := entities.NewNode(network.ID, cluster.ID,
        "node1", "", make([]string, 0), "", "", false,
        "", "", "")
    descriptor := entities.NewAppDescriptor(network.ID, "descriptor1", "",
        "", "", "", 0, []string{"repo1:tag1"})

    user := entities.NewUser("user", "9999", "email@email.com", time.Now(), time.Now().Add(time.Hour))

    access := entities.NewUserAccess(user.ID, []entities.RoleType{entities.GlobalAdmin})
    userpassword := "daisho"
    password, err := entities.NewPassword(user.ID, &userpassword)

    helper.Nil(err, "expecting password")

    helper.endpoints.networkProvider.Add(* network)
    helper.endpoints.clusterProvider.Add(* cluster)
    helper.endpoints.nodeProvider.Add(* node)
    helper.endpoints.appDescProvider.Add(* descriptor)

    helper.endpoints.userProvider.Add(* user)
    helper.endpoints.accessProvider.Add(* access)
    helper.endpoints.passwordProvider.Add(* password)

    // create backup entity
    networks, _ := helper.endpoints.networkProvider.ListNetworks()
    backup := entities.NewBackup()
    backup.AddNetworks(networks)

    clusters, _ := helper.endpoints.clusterProvider.Dump()
    backup.AddClusters(clusters)

    nodes, _ := helper.endpoints.nodeProvider.Dump()
    backup.AddNodes(nodes)

    descriptors, _ := helper.endpoints.appDescProvider.Dump()
    backup.AddAppDescriptors(descriptors)

    users, _ := helper.endpoints.userProvider.Dump()
    accesses, _ := helper.endpoints.accessProvider.Dump()
    passwords, _ := helper.endpoints.passwordProvider.Dump()

    for _, user := range users {
        password := entities.Password{}
        access := entities.UserAccess{}

        for _, access = range accesses {
            if user.ID == access.UserID {
                break
            }
        }
        for _, password = range passwords {
            if user.ID == password.UserID {
                break

            }
        }
        backup.AddUsers(entities.BackupUser{User: user, Access: access, Password: password})
    }

    // Build backup user
    helper.SetupTest() // cleanup for restoring

    return backup
}

func (helper *BRHandlerTestHelper) TestRestore() {
    backup := helper.LoadImportData()

    f := frisby.Create("Import").Post(helper.getURL("backup/all/restore")).SetJson(backup).Send()
    f.ExpectStatus(http.StatusOK)

    f = frisby.Create("Export").Get(helper.getURL("backup/all/create")).Send()
    f.ExpectStatus(http.StatusOK)

    f.AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
        networks, err := json.Get("networks").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(networks), "length does not match")
        clusters, err := json.Get("clusters").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(clusters), "length does not match")
        nodes, err := json.Get("nodes").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(nodes), "length does not match")
        descriptors, err := json.Get("appdesc").Array()
        helper.Nil(err, "expecting array")
        helper.Equal(1, len(descriptors), "length does not match")
        users, err := json.Get("users").Array()
        helper.Equal(1, len(users), "expecting users")
    })

}

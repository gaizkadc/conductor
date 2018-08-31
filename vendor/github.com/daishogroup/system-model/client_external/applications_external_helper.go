//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the helper for the applications client integration test.

package client_external

import (
    "context"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"

    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/server/app"
    "github.com/daishogroup/dhttp"
)

type ApplicationsEndpointHelper struct {
    networkProvider * networkstorage.MockupNetworkProvider
    appDescProvider * appdescstorage.MockupAppDescProvider
    appInstProvider * appinststorage.MockupAppInstProvider

    appManager app.Manager
    appHandler app.Handler

    port int
    handler         http.Handler
    srv             * http.Server
}

func (helper * ApplicationsEndpointHelper) LaunchServer(){
    err := helper.srv.ListenAndServe()
    if err != nil {
        println(err.Error())
    }
}

// Start the HttpServer.
func (helper * ApplicationsEndpointHelper) Start() {
    go helper.LaunchServer()
}

// Shutdown the HTTPServer.
func (helper * ApplicationsEndpointHelper) Shutdown() {
    ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
    helper.srv.Shutdown(ctx)
    helper.srv.Close()
}

// Clear mockup providers
func (helper * ApplicationsEndpointHelper) ResetProviders() {
    helper.networkProvider.Clear()
    helper.appDescProvider.Clear()
    helper.appInstProvider.Clear()
}

func (helper * ApplicationsEndpointHelper) GetListeningAddress() (string, int) {
    return BaseAddress ,helper.port
}

// Add a test network.
//   params:
//     networkId The network identifier.
func (helper * ApplicationsEndpointHelper) AddNetwork(networkId string) {
    toAdd := entities.NewNetworkWithID(
        networkId, client.TestName, client.TestDescription, "", "", "")
    helper.networkProvider.Add(*toAdd)
}

// Add a test application descriptor.
//   params:
//     networkId The network identifier.
//     descriptorId The application descriptor identifier.
func (helper * ApplicationsEndpointHelper) AddDescriptor(networkId string, descriptorId string){
    toAdd := entities.NewAppDescriptorWithID(networkId,
        descriptorId, client.TestName, client.TestDescription,
        client.TestServiceName, client.TestServiceVersion, client.TestLabel, client.TestPort, []string {client.TestImage})
    helper.appDescProvider.Add(*toAdd)
    helper.networkProvider.RegisterAppDesc(networkId, descriptorId)
}

// Add a test application instance.
//   params:
//     networkId The network identifier.
//     deployedId The deployed application instance identifier.
func (helper * ApplicationsEndpointHelper) AddInstance(networkId string, deployedId string){
    toAdd := entities.NewAppInstanceWithID(networkId,
        deployedId, client.TestDescriptorID, "",
        client.TestName, client.TestDescription, client.TestLabel, client.TestArguments,
        entities.AppInstReady, client.TestPersistenceSize, client.TestStorageType,
        make([]entities.ApplicationPort, 0), client.TestPort, "")
    helper.appInstProvider.Add(*toAdd)
    helper.networkProvider.RegisterAppInst(networkId, deployedId)
}


// Create a new ApplicationsEndpointHelper.
//   returns:
//     A new endpoint helper with mockup providers.
func NewApplicationsEndpointHelper() ApplicationsEndpointHelper {
    var networkProvider = networkstorage.NewMockupNetworkProvider()
    var appDescProvider = appdescstorage.NewMockupAppDescProvider()
    var appInstProvider = appinststorage.NewMockupAppInstProvider()

    var appManager = app.NewManager(networkProvider, appDescProvider, appInstProvider)
    var appHandler = app.NewHandler(appManager)

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

    return ApplicationsEndpointHelper{
        networkProvider, appDescProvider, appInstProvider,
        appManager, appHandler,
        port, handler, srv}
}
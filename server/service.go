//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Main entry to launch the conductor service.

package server

import (
    "fmt"
    "net/http"

    "github.com/daishogroup/conductor/server/apps"
    smclient "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "github.com/daishogroup/conductor/asm"
    "github.com/daishogroup/conductor/logger"
    "github.com/daishogroup/dhttp"
)

// Logger for the Service struct.
var loggerService = log.WithField(
    "package", "api",
).WithField(
    "struct", "Service",
).WithField(
    "file", "service.go",
)

// Api service object.
type Service struct {
    Configuration Config
}

// Name of the service.
func (app *Service) Name() string {
    return "Conductor API Service."
}

// Service description.
func (app *Service) Description() string {
    return "API service of the Conductor project."
}

// Run the service, launch the REST service handler.
func (app *Service) Run() error {
    loggerService.Info("Starting Conductor Service.")
    app.Configuration.Print()

    router := mux.NewRouter()
    // Create the BUDO network handler
    // The network manager points to the system model API
    appClient := smclient.NewApplicationRest(app.Configuration.SystemModelAddress)
    clusterClient := smclient.NewClusterRest(app.Configuration.SystemModelAddress)
    nodeClient := smclient.NewNodeRest(app.Configuration.SystemModelAddress)
    asmClientFactory := asm.NewRestClientFactory()
    loggerClient := logger.NewRestClient(app.Configuration.LoggerAddress)
    conductorManager := apps.NewRestManager(appClient, clusterClient, nodeClient, asmClientFactory, loggerClient)
    // Instanciate the api handler
    conductorHandler := apps.NewHandler(conductorManager)
    conductorHandler.SetRoutes(router)

    router.HandleFunc("/ping", app.ping).Methods("GET")
    loggerService.Info("Ready to serve HTTP")
    app.listRoutes(router)

    // Workaround for CORS headers
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"POST", "GET"})

    // Run the HTTP server
    //go http.ListenAndServe(fmt.Sprintf(":%d", app.Configuration.Port),router)
    go http.ListenAndServe(fmt.Sprintf(":%d", app.Configuration.Port), handlers.CORS(originsOk, headersOk, methodsOk)(router))
    return nil
}

// The ping endpoint returns a 200 Ok response for external services to check that the service is up.
// TODO If used for liveness probes, enhance the check with provider status.
func (app *Service) ping(w http.ResponseWriter, r *http.Request) {
    dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation("ping"))
}

func (app *Service) listRoutes(handler *mux.Router) {
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
}

// Finalize the service.
func (app *Service) Finalize(killSignal bool) {
    loggerService.Info("Finalize Conductor Api Service.")
}

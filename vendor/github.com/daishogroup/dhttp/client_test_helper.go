//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package dhttp

import (
    "context"
    "net/http"
    "os"
    "time"
    "strconv"
    "github.com/gorilla/mux"
    "github.com/daishogroup/derrors"
    log "github.com/sirupsen/logrus"
)

var loggerTestHelper = log.WithField("package", "client")

const maxTries = 10
const sleepSeconds = 1

// ClientTestHelper is the struct that contains all the components needed to run the tests below.
type ClientTestHelper struct {
    host   string
    port   int
    router *mux.Router
    srv    *http.Server
}

// LaunchServer is the method to run the test server.
func (helper *ClientTestHelper) launchServer() {
    err := helper.srv.ListenAndServe()
    if err != http.ErrServerClosed {
        loggerTestHelper.Error(err.Error())
        os.Exit(1)
    }
}

// Shutdown the HTTPServer.
func (helper *ClientTestHelper) Shutdown() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    helper.srv.Shutdown(ctx)
    helper.srv.Close()
}

// Start the HttpServer.
func (helper *ClientTestHelper) Start() {
    go helper.launchServer()
    helper.waitTestPing(maxTries, sleepSeconds)
}

func (helper *ClientTestHelper) waitTestPing(tries int, seconds int) {
    var i = 0
    var exit = false
    for !exit && i < tries {
        conf:= NewRestBasicConfig(helper.host, helper.port)
        c := NewClientSling(conf)
        url := GetURL(helper.host, helper.port, "/test/ping")
        result := &SuccessfulOperation{}
        response := c.Get(url, result)
        if response.Error == nil {
            exit = true
        } else {
            time.Sleep(time.Duration(seconds) * time.Second)
        }
        i++
    }
    if !exit {
        panic(derrors.NewConnectionError("Server is not launched"))
    }
}

func addTestRoutes(router *mux.Router) {
    loggerTestHandler.Debug("Setting cluster routes")
    router.HandleFunc("/test/ping", testPing).Methods("GET")
}

// test ping add test ping endpoint.
//   params:
//     w The HTTP writer structure.
//     r The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func testPing(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    loggerTestHandler.Debug("testPing")
    RespondWithJSON(w, http.StatusOK, NewSuccessfulOperation("test-ping"))
}

// NewClientTestHelper builds a new ClientTestHelper.
//   returns:
//	   New endpoint helper
func NewClientTestHelper(handler Handler, host string) *ClientTestHelper {
    port, _ := GetAvailablePort()
    router := mux.NewRouter()
    addTestRoutes(router)

    handler.SetRoutes(router)

    var srv = &http.Server{
        Handler: router,
        Addr:    host + ":" + strconv.Itoa(port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    return &ClientTestHelper{host, port, router, srv}
}

// Handler is the interface of a basic handler.
type Handler interface {
    SetRoutes(router *mux.Router)
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package config

import (

    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/errors"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/dhttp"
)

var logger = log.WithField("package", "server.config")

// Handler structure that contains the link with the underlying manager.
type Handler struct {
    manager Manager
}

// NewHandler creates a new handler for cluster endpoints.
//   params:
//     manager The cluster manager.
//   returns:
//     A new handler.
func NewHandler(manager Manager) Handler {
    return Handler{manager}
}

// SetRoutes registers the endpoints of this handler.
//   params:
//     The REST handler.
func (handler * Handler) SetRoutes(router * mux.Router) {
    logger.Info("Setting config routes")
    router.HandleFunc("/api/v0/config/set", handler.setConfig).Methods("POST")
    router.HandleFunc("/api/v0/config/get", handler.getConfig).Methods("GET")
}

// Set the configuration.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler * Handler) setConfig(w http.ResponseWriter, r * http.Request){
    toAdd := &entities.Config{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil{
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    }else{
        err = json.Unmarshal(b, &toAdd)
        if err == nil && toAdd.Valid(){
            logger.Debug("Setting new config: " + toAdd.String())
            err := handler.manager.SetConfig(*toAdd)
            if err != nil {
                dhttp.RespondWithError(w, http.StatusInternalServerError, err)
            }else{
                dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation("setConfig"))
            }
        }else{
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(toAdd, errors.InvalidEntity))
        }
    }
}

// Get the configuration.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler * Handler) getConfig(w http.ResponseWriter, r * http.Request){
    config, err := handler.manager.GetConfig()
    if err != nil {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }else{
        dhttp.RespondWithJSON(w, http.StatusOK, config)
    }
}


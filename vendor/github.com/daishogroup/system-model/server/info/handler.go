//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Dump handler.

package info

import (
    "net/http"

    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/dhttp"
)

var logger = log.WithField("package", "server.info")

// Handler structure that contains the link with the underlying manager.
type Handler struct {
    manager Manager
}

// NewHandler obtains a new handler.
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
func (handler *Handler) SetRoutes(router *mux.Router) {
    logger.Info("Setting cluster routes")
    router.HandleFunc("/api/v0/info/reduced", handler.reducedInfo).Methods("GET")
    router.HandleFunc("/api/v0/info/summary", handler.summaryInfo).Methods("GET")
    router.HandleFunc("/api/v0/info/{networkID}/reduced", handler.reducedInfoByNetwork).Methods("GET")
}

// ReducedInfo of the system model.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) reducedInfo(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("reducedInfo")
    dump, err := handler.manager.ReducedInfo()
    if err != nil {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
    dhttp.RespondWithJSON(w, http.StatusOK, dump)
}

// SummaryInfo of the system model.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) summaryInfo(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("summaryInfo")
    dump, err := handler.manager.SummaryInfo()
    if err != nil {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
    dhttp.RespondWithJSON(w, http.StatusOK, dump)
}

// Get a specific network.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK with the network if the operation is successful or the error.
func (handler *Handler) reducedInfoByNetwork(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("reducedInfoByNetwork")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID != "" {
        logger.Debug("GetNetwork: " + networkID)
        network, err := handler.manager.ReducedInfoByNetwork(networkID)
        if err == nil {
            dhttp.RespondWithJSON(w, http.StatusOK, network)
        } else {
            dhttp.RespondWithError(w, http.StatusInternalServerError, err)
        }
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError,
            derrors.NewOperationError(errors.MissingRESTParameter).
                WithParams("networkID"))
    }
}

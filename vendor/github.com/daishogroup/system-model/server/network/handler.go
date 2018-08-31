//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the network handler specification in charge of validating client requests for network
// operations.

package network

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strconv"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/errors"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/dhttp"
)

var logger = log.WithField("package", "server.network")

// Handler structure that contains the link with the underlying manager.
type Handler struct {
    manager Manager
}

// NewHandler creates a new handler.
//   params:
//     manager The network manager.
//   returns:
//     A new manager.
func NewHandler(manager Manager) Handler {
    return Handler{manager}
}

// SetRoutes registers the endpoints of this handler.
//   params:
//     The REST handler.
func (handler *Handler) SetRoutes(router *mux.Router) {
    logger.Info("Setting network routes")
    router.HandleFunc("/api/v0/network/add", handler.addNetworkEndpoint).Methods("POST")
    router.HandleFunc("/api/v0/network/list", handler.listNetworks).Methods("GET")
    router.HandleFunc("/api/v0/network/{networkID}/info", handler.getNetwork).Methods("GET")
    router.HandleFunc("/api/v0/network/{networkID}/delete", handler.deleteNetwork).Methods("DELETE")
}

// The Echo method returns the parameter q as response.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) addNetworkEndpoint(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("addNetworkEndpoint")
    addNetworkRequest := &entities.AddNetworkRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &addNetworkRequest)
        if err == nil && addNetworkRequest.IsValid() {
            logger.Debug("Adding new network: " + addNetworkRequest.String())
            added, err2 := handler.manager.AddNetwork(*addNetworkRequest)
            if added != nil {
                logger.Debug("added: ", added)
            }
            if err2 != nil {
                logger.Debugf(">>>>>> err: %#v", err2)
            }
            if err2 != nil {
                dhttp.RespondWithError(w, http.StatusInternalServerError, err2)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, added)
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(addNetworkRequest, errors.InvalidEntity))
        }
    }

}

// List the networks in the system model.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) listNetworks(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("listNetworks")
    networks, err := handler.manager.ListNetworks()
    if err == nil {
        logger.Debug("List result: " + strconv.Itoa(len(networks)))
        dhttp.RespondWithJSON(w, http.StatusOK, networks)
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Get a specific network.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK with the network if the operation is successful or the error.
func (handler *Handler) getNetwork(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("getNetworks")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID != "" {
        logger.Debug("GetNetwork: " + networkID)
        network, err := handler.manager.GetNetwork(networkID)
        if err == nil {
            dhttp.RespondWithJSON(w, http.StatusOK, network)
        } else {
            dhttp.RespondWithError(w, http.StatusInternalServerError, err)
        }
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
    }
}

// Delete a specific network.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK with the network if the operation is successful or the error.
func (handler *Handler) deleteNetwork(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("deleteNetwork")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID != "" {
        logger.Debug("DeleteNetwork: " + networkID)
        err := handler.manager.DeleteNetwork(networkID)
        if err == nil {
            dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation("DeleteNetwork"))
        } else {
            dhttp.RespondWithError(w, http.StatusInternalServerError, err)
        }
    } else {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
    }
}

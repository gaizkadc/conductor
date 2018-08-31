//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the node handler specification in charge of validating client requests for node
// operations.

package node

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "github.com/daishogroup/dhttp"
)

var logger = log.WithField("package", "server.node")

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

    router.HandleFunc("/api/v0/node/{networkID}/{clusterID}/add", handler.addNode).
        Methods("POST")

    router.HandleFunc("/api/v0/node/{networkID}/{clusterID}/list", handler.listNodes).
        Methods("GET")

    router.HandleFunc("/api/v0/node/{networkID}/{clusterID}/filter", handler.filterNodes).
        Methods("GET")

    router.HandleFunc("/api/v0/node/{networkID}/{clusterID}/{nodeID}/delete", handler.deleteNode).
        Methods("DELETE")

    router.HandleFunc("/api/v0/node/{networkID}/{clusterID}/{nodeID}/info", handler.getNode).
        Methods("GET")

    router.HandleFunc("/api/v0/node/{networkID}/{clusterID}/{nodeID}/update", handler.updateNode).
        Methods("POST")
}

// Add a new node to an existing cluster.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) addNode(w http.ResponseWriter, r *http.Request) {
    logger.Debug("addNode")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    clusterID := vars["clusterID"]
    if clusterID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("clusterID"))
        return
    }

    addNodeRequest := &entities.AddNodeRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &addNodeRequest)
        if err == nil && addNodeRequest.IsValid() {
            logger.Debug("Adding new node: " + addNodeRequest.String())
            added, err := handler.manager.AddNode(networkID, clusterID, *addNodeRequest)
            if err != nil {
                dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, added)
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(addNodeRequest, errors.InvalidEntity))
        }
    }
}

// List the nodes of an existing cluster.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) listNodes(w http.ResponseWriter, r *http.Request) {
    logger.Debug("listNodes")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    clusterID := vars["clusterID"]
    if clusterID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("clusterID"))
        return
    }
    nodes, err := handler.manager.ListNodes(networkID, clusterID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, nodes)
    } else {
        dhttp.RespondWithError(w, http.StatusBadRequest, err)
    }
}

// Delete an existing node from an existing cluster.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) deleteNode(w http.ResponseWriter, r *http.Request) {
    logger.Debug("deleteNode")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter))
        return
    }
    clusterID := vars["clusterID"]
    if clusterID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter))
        return
    }
    nodeID := vars["nodeID"]
    if nodeID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter))
        return
    }
    err := handler.manager.RemoveNode(networkID, clusterID, nodeID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation("DeleteNode"))
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Get the cluster information.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) getNode(w http.ResponseWriter, r *http.Request) {
    logger.Debug("getNode")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    clusterID := vars["clusterID"]
    if clusterID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("clusterID"))
        return
    }
    nodeID := vars["nodeID"]
    if nodeID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("nodeID"))
        return
    }
    cluster, err := handler.manager.GetNode(networkID, clusterID, nodeID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, cluster)
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Update an existing node.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) updateNode(w http.ResponseWriter, r *http.Request) {
    logger.Debug("updateNode")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    clusterID := vars["clusterID"]
    if clusterID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("clusterID"))
        return
    }
    nodeID := vars["nodeID"]
    if nodeID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("nodeID"))
        return
    }

    updateRequest := &entities.UpdateNodeRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &updateRequest)
        if err == nil && updateRequest.IsValid() {
            logger.Debug("Updating node: " + nodeID + " with " + updateRequest.String())
            updated, err := handler.manager.UpdateNode(networkID, clusterID, nodeID, * updateRequest)
            if err != nil {
                dhttp.RespondWithError(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, updated)
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(updateRequest, errors.InvalidEntity))
        }
    }
}

// Filter the nodes using a set of criterias.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) filterNodes(w http.ResponseWriter, r *http.Request) {
    logger.Debug("filterNodes")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    clusterID := vars["clusterID"]
    if clusterID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("clusterID"))
        return
    }

    filterRequest := &entities.FilterNodesRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &filterRequest)
        if err == nil {
            nodes, err := handler.manager.FilterNodes(networkID, clusterID, *filterRequest)
            if err != nil {
                dhttp.RespondWithError(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, nodes)
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(filterRequest, errors.InvalidEntity))
        }
    }
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the cluster handler specification in charge of validating client requests for cluster
// operations.

package cluster

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

var logger = log.WithField("package", "server.cluster")

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
func (handler *Handler) SetRoutes(router *mux.Router) {
    logger.Info("Setting cluster routes")
    router.HandleFunc("/api/v0/cluster/{networkID}/add", handler.addCluster).Methods("POST")
    router.HandleFunc("/api/v0/cluster/{networkID}/list", handler.listClusters).Methods("GET")
    router.HandleFunc("/api/v0/cluster/{networkID}/{clusterID}/delete", handler.deleteCluster).Methods("DELETE")
    router.HandleFunc("/api/v0/cluster/{networkID}/{clusterID}/info", handler.getCluster).Methods("GET")
    router.HandleFunc("/api/v0/cluster/{networkID}/{clusterID}/update", handler.updateCluster).Methods("POST")
}

// Add a new cluster to an existing network.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) addCluster(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("addCluster")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }

    addClusterRequest := &entities.AddClusterRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &addClusterRequest)
        if err == nil && addClusterRequest.IsValid() {
            logger.Debug("Adding new cluster: " + addClusterRequest.String())
            added, err := handler.manager.AddCluster(networkID, *addClusterRequest)
            if err != nil {
                dhttp.RespondWithError(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, added)
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(addClusterRequest, errors.InvalidEntity))
        }
    }
}

// Add a new cluster to an existing network.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) listClusters(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("listClusters")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    logger.Debug("ListClusters: " + networkID)
    clusters, err := handler.manager.ListClusters(networkID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, clusters)
    } else {
        dhttp.RespondWithError(w, http.StatusBadRequest, err)
    }
}

// Delete an existing cluster from an existing network.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) deleteCluster(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("deleteCluster")
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
    logger.Debug("DeleteCluster: " + networkID + " => " + clusterID)
    err := handler.manager.DeleteCluster(networkID, clusterID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation("DeleteCluster"))
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Get the cluster information.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) getCluster(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("getCluster")
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
    logger.Debug("GetCluster: " + networkID + " => " + clusterID)
    cluster, err := handler.manager.GetCluster(networkID, clusterID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, cluster)
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Update an existing cluster.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) updateCluster(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("updateCluster")
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

    updateRequest := &entities.UpdateClusterRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &updateRequest)
        if err == nil && updateRequest.IsValid() {
            logger.Debug("Updating cluster: " + clusterID + " with " + updateRequest.PrettyString())
            updated, err := handler.manager.UpdateCluster(networkID, clusterID, * updateRequest)
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

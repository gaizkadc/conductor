//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Handler for the apps entrypoint.

package apps

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/daishogroup/conductor/errors"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    entitiesConductor "github.com/daishogroup/conductor/entities"
    "github.com/daishogroup/dhttp"

)

var logger = log.WithField("package", "apps.handler")

type Handler struct {
    manager AppManager
}

func NewHandler(manager AppManager) Handler {
    return Handler{manager}
}

// Register the endpoints of this handler.
//   params:
//     The REST handler.
func (h *Handler) SetRoutes(router *mux.Router) {
    logger.Info("Setting conductor application routes")
    router.HandleFunc("/api/v0/app/{networkId}/deploy", h.deployApp).Methods("POST")
    router.HandleFunc("/api/v0/app/{networkId}/{instanceId}/undeploy", h.undeployApp).Methods("GET")
    router.HandleFunc("/api/v0/app/{networkId}/{instanceId}/logs", h.logs).Methods("GET")
}

// Deploy an application in a given network.
// params:
//   w http writer
//   r htt request
// returns:
//   A 200 OK if the operation was successful. 500 if error occurred.
func (h *Handler) deployApp(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("Called deployApps")
    vars := mux.Vars(r)
    networkID := vars["networkId"]

    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter))
        return
    }

    deployAppRequest := &entitiesConductor.DeployAppRequest{}
    b, err := ioutil.ReadAll(r.Body)

    if err != nil {
        logger.Errorf("Impossible to read JSON body")
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.UnmarshalError, err))
    } else {
        json.Unmarshal(b, &deployAppRequest)
        err := deployAppRequest.IsValid()
        if err == nil {
            logger.Debug("Deploy update: " + deployAppRequest.String())
            deployed, err := h.manager.Deploy(networkID, *deployAppRequest)
            if err != nil {
                logger.Debugf("Error found during deployment process")
                dhttp.RespondWithError(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, deployed)
            }
        } else {
            logger.Error("The instance application request was not valid")
            dhttp.RespondWithError(w, http.StatusBadRequest, err)
        }
    }
}

// Undeploy an already deployed application.
// params:
//   w http writer
//   r htt request
// returns:
//   A 200 OK if the operation was successful. 500 if error occurred.
func (h *Handler) undeployApp(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("Called undeployApp")
    vars := mux.Vars(r)
    networkID := vars["networkId"]
    instanceID := vars["instanceId"]

    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter))
        return
    }

    if instanceID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter))
        return
    }

    err := h.manager.Undeploy(networkID, instanceID)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    } else {
        dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation("undeploy"))
    }

}

// Log recover logs entry from a specific app.
// params:
//   w http writer
//   r htt request
// returns:
//   A 200 OK if the operation was successful. 500 if error occurred.
func (h *Handler) logs(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("Called getAppLogs")
    vars := mux.Vars(r)
    networkID := vars["networkId"]
    instanceID := vars["instanceId"]

    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter))
        return
    }

    if instanceID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter))
        return
    }

    logs, err := h.manager.Logs(networkID, instanceID)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    } else {
        dhttp.RespondWithJSON(w, http.StatusOK, logs)
    }

}

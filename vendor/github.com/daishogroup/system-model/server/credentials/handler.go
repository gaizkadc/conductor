//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Credentials requests handler.

package credentials

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "fmt"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "github.com/daishogroup/dhttp"
)

var logger = log.WithField("package", "server.credentials")

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
    logger.Info("Setting credentials routes")
    router.HandleFunc("/api/v0/credentials/add", handler.addCredentials).Methods("POST")
    router.HandleFunc("/api/v0/credentials/{uuid}/get", handler.getCredentials).Methods("GET")
    router.HandleFunc("/api/v0/credentials/{uuid}/delete", handler.deleteCredentials).Methods("DELETE")
}

// Add a new user.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) addCredentials(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("addCredentials")

    addCredentialsRequest := &entities.AddCredentialsRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &addCredentialsRequest)
        if err == nil && addCredentialsRequest.IsValid() {
            logger.Debugf("Adding new credentials: #%v", addCredentialsRequest)
            err := handler.manager.AddCredentials(*addCredentialsRequest)
            if err != nil {
                dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation(errors.CredentialsAdded))
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(addCredentialsRequest, errors.InvalidEntity))
        }
    }
}

// Get credentials.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) getCredentials(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("getCredentials")

    vars := mux.Vars(r)
    uuid := vars["uuid"]

    if uuid == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("uuid"))
        return
    }

    logger.Debug(fmt.Sprintf("Get user credentials: %s", uuid))
    added, err := handler.manager.GetCredentials(uuid)
    if err != nil {
        dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
    } else {
        dhttp.RespondWithJSON(w, http.StatusOK, added)
    }
}

// Delete credentials
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) deleteCredentials(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("deleteCredentials")

    vars := mux.Vars(r)
    uuid := vars["uuid"]

    if uuid == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("uuid"))
        return
    }

    logger.Debug(fmt.Sprintf("Delete credentials: %s", uuid))
    err := handler.manager.DeleteCredentials(uuid)
    if err != nil {
        dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
    } else {
        dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation(errors.CredentialsRemoved))
    }
}

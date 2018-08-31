//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the handler for the oauth secrets manager.

package oauth

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

var logger = log.WithField("package", "server.oauth")

// Handler structure that contains the link with the underlying manager.
type Handler struct {
    manager Manager
}

// NewHandler creates a new handler.
//   params:
//     manager The OAuth manager.
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
    router.HandleFunc("/api/v0/oauth/{userID}", handler.setSecret).
        Methods("POST")
    router.HandleFunc("/api/v0/oauth/{userID}", handler.deleteUser).
        Methods("DELETE")
    router.HandleFunc("/api/v0/oauth/{userID}", handler.getSecrets).
        Methods("GET")
}

// Set a secret for an existing user.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) setSecret(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("setSecret")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    setEntryRequest := &entities.OAuthAddEntryRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &setEntryRequest)
        if err == nil {
            logger.Debugf("Setting a new app secret: %s", setEntryRequest)
            err := handler.manager.SetSecret(userID, *setEntryRequest)
            if err != nil {
                dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation(errors.OAuthEntrySet))
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(setEntryRequest, errors.InvalidEntity))
        }
    }
}

// Delete a user entry with all his secrets.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("deleteUser")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    err := handler.manager.DeleteSecrets(userID)
    if err != nil {
        dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
    } else {
        dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation(errors.OAuthUserDeleted))
    }

}

// Get secrets from an existing user.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) getSecrets(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("getPassword")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    toReturn, err := handler.manager.GetSecrets(userID)
    if err != nil {
        dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
    } else {
        dhttp.RespondWithJSON(w, http.StatusOK, toReturn)
    }

}

// TODO Add indivual application control such as delete application information, etc.

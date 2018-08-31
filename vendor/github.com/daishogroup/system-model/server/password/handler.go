//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the handler for the password manager.

package password

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

var logger = log.WithField("package", "server.password")

// Handler structure that contains the link with the underlying manager.
type Handler struct {
    manager Manager
}

// NewHandler creates a new handler.
//   params:
//     manager The password manager.
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
    router.HandleFunc("/api/v0/password", handler.setPassword).
        Methods("POST")
    router.HandleFunc("/api/v0/password/{userID}", handler.deletePassword).
        Methods("DELETE")
    router.HandleFunc("/api/v0/password/{userID}", handler.getPassword).
        Methods("GET")
}

// Set a password.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) setPassword(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("setPassword")

    setPasswordRequest := &entities.Password{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &setPasswordRequest)
        if err == nil {
            logger.Debugf("Setting a password: %s", setPasswordRequest)
            err := handler.manager.SetPassword(*setPasswordRequest)
            if err != nil {
                dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation(errors.PasswordSet))
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(setPasswordRequest, errors.InvalidEntity))
        }
    }
}

// Delete a password.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) deletePassword(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("deletePassword")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    err := handler.manager.DeletePassword(userID)
    if err != nil {
        dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
    } else {
        dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation(errors.PasswordDeleted))
    }

}

// Get a password.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) getPassword(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("getPassword")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    toReturn, err := handler.manager.GetPassword(userID)
    if err != nil {
        dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
    } else {
        dhttp.RespondWithJSON(w, http.StatusOK, toReturn)
    }

}

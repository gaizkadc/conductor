//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the access manager in charge of granting users privileges.

package access

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/gorilla/mux"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    log "github.com/sirupsen/logrus"
    "github.com/daishogroup/dhttp"
)

var logger = log.WithField("package", "server.access")

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

    router.HandleFunc("/api/v0/access/{userID}/add", handler.addAccess).Methods("POST")
    router.HandleFunc("/api/v0/access/{userID}/set", handler.setAccess).Methods("POST")
    router.HandleFunc("/api/v0/access/{userID}/get", handler.getAccess).Methods("GET")
    router.HandleFunc("/api/v0/access/{userID}/delete", handler.deleteAccess).Methods("DELETE")
    router.HandleFunc("/api/v0/access/list", handler.listAccess).Methods("GET")
}

// Add new privilege for a user.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) addAccess(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("addAccess")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    addUserRequest := &entities.AddUserAccessRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &addUserRequest)
        if err == nil {
            logger.Debugf("Adding new access: %s", addUserRequest)
            added, err := handler.manager.AddAccess(userID, *addUserRequest)
            if err != nil {
                dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, added)
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(addUserRequest, errors.InvalidEntity))
        }
    }
}

// Set user privileges.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) setAccess(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("setAccess")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    addUserRequest := &entities.AddUserAccessRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &addUserRequest)
        if err == nil {
            logger.Debugf("Setting access: %s", addUserRequest)
            added, err := handler.manager.SetAccess(userID, *addUserRequest)
            if err != nil {
                dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, added)
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(addUserRequest, errors.InvalidEntity))
        }
    }
}

// Get existing user privileges.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) getAccess(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("getAccess")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    user, err := handler.manager.GetAccess(userID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, user)
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Delete an existing user.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) deleteAccess(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("deleteAccess")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    err := handler.manager.DeleteAccess(userID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation("DeleteAccess"))
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Get the list of existing users with their roles
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) listAccess(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("listAccess")

    entries, err := handler.manager.ListAccess()
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, entries)
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }

}

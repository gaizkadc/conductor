//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the user handler specification in charge of creating and modifying system users.

package user

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

var logger = log.WithField("package", "server.user")

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
    router.HandleFunc("/api/v0/user/add", handler.addUser).Methods("POST")
    router.HandleFunc("/api/v0/user/{userID}/get", handler.getUser).Methods("GET")
    router.HandleFunc("/api/v0/user/{userID}/delete", handler.deleteUser).Methods("DELETE")
    router.HandleFunc("/api/v0/user/{userID}/update", handler.updateUser).Methods("POST")
    router.HandleFunc("/api/v0/user/list", handler.getList).Methods("GET")
}

// Add a new user.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) addUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("addUser")

    addUserRequest := &entities.AddUserRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &addUserRequest)
        if err == nil && addUserRequest.IsValid() {
            logger.Debug("Adding new user: " + addUserRequest.String())
            added, err := handler.manager.AddUser(*addUserRequest)
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

// Get an existing user.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) getUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("getUser")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    user, err := handler.manager.GetUser(userID)
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
func (handler *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("deleteUser")

    vars := mux.Vars(r)
    userID := vars["userID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("userID"))
        return
    }

    err := handler.manager.DeleteUser(userID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation("DeleteSecrets"))
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Update a user.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("updateUser")

    vars := mux.Vars(r)
    userID := vars["userID"]

    updateUserRequest := &entities.UpdateUserRequest{}

    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &updateUserRequest)
        if err == nil && updateUserRequest.IsValid() {
            logger.Debug(fmt.Sprintf("Update user: %s", userID))
            added, err := handler.manager.UpdateUser(userID, *updateUserRequest)
            if err != nil {
                dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, added)
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(updateUserRequest, errors.InvalidEntity))
        }
    }
}

// Get the list of users.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) getList(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("getList")
    users, err := handler.manager.ListUsers()
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, users)
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

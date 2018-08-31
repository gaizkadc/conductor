//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the user handler specification in charge of creating and modifying system users.

package session

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

var logger = log.WithField("package", "server.session")

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
    router.HandleFunc("/api/v0/session/add", handler.addSession).Methods("POST")
    router.HandleFunc("/api/v0/session/{sessionID}/get", handler.getSession).Methods("GET")
    router.HandleFunc("/api/v0/session/{sessionID}/delete", handler.deleteSession).Methods("DELETE")
}

// Add a new user.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) addSession(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("addSession")

    addSessionRequest := &entities.AddSessionRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    } else {
        err = json.Unmarshal(b, &addSessionRequest)
        logger.Infof("%v", addSessionRequest)
        if err == nil && addSessionRequest.IsValid() {
            logger.Debugf("Adding new session: %v ", addSessionRequest)
            added, err := handler.manager.AddSession(*addSessionRequest)
            if err != nil {
                dhttp.RespondWithJSON(w, http.StatusInternalServerError, err)
            } else {
                dhttp.RespondWithJSON(w, http.StatusOK, added)
            }
        } else {
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(addSessionRequest, errors.InvalidEntity))
        }
    }
}

// Get an existing user.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) getSession(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("getSession")

    vars := mux.Vars(r)
    sessionID := vars["sessionID"]

    if sessionID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("sessionID"))
        return
    }

    session, err := handler.manager.GetSession(sessionID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, session)
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Delete a session.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) deleteSession(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    logger.Debug("deleteSession")

    vars := mux.Vars(r)
    userID := vars["sessionID"]

    if userID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("sessionID"))
        return
    }

    err := handler.manager.DeleteSession(userID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation("DeleteSession"))
    } else {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}
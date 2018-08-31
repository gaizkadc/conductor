package backup

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Backup Restore handler.


import (
    "net/http"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "encoding/json"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/dhttp"
)

var logger = log.WithField("package", "server.dump")

// Handler structure that contains the link with the underlying manager.
type Handler struct {
    manager Manager
}

// NewHandler creates a new handler.
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
func (handler * Handler) SetRoutes(router * mux.Router) {
    logger.Info("Setting cluster routes")
    router.HandleFunc("/api/v0/backup/{component}/create", handler.Export).Methods("GET")
    router.HandleFunc("/api/v0/backup/{component}/restore", handler.Import).Methods("POST")
}

// Backup the system model.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler * Handler) Export(w http.ResponseWriter, r * http.Request) {
    defer r.Body.Close()
    vars := mux.Vars(r)

    component, ok := vars["component"]
    if !ok {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewGenericError("Component not present"))
        return
    }

    backup, derr := handler.manager.Export(component)
    if derr != nil {
        dhttp.RespondWithError(w, http.StatusInternalServerError, derr)
        return
    }

    dhttp.RespondWithJSON(w, http.StatusOK, backup)
}


// Restore the system model.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler * Handler) Import(w http.ResponseWriter, r * http.Request) {
    defer r.Body.Close()
    vars := mux.Vars(r)

    component, ok := vars["component"]
    if !ok {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewGenericError("Component not present"))
        return
    }

    entity := entities.BackupRestore{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusInternalServerError, derrors.NewGenericError("HTTP body error", err))
        return
    }

    jErr := json.Unmarshal(b, &entity)
    if jErr != nil {
        dhttp.RespondWithError(w, http.StatusInternalServerError, derrors.NewGenericError("Error unmarshalling request", jErr))
        return
    }

    derr := handler.manager.Import(component, &entity)
    if derr != nil {
        dhttp.RespondWithError(w, http.StatusInternalServerError, derr)
        return
    }

    dhttp.RespondWithJSON(w, http.StatusOK, nil)
}

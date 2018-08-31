//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Dump handler.

package dump

import (
    "net/http"

    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"

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
func (handler *Handler) SetRoutes(router *mux.Router) {
    logger.Info("Setting cluster routes")
    router.HandleFunc("/api/v0/dump/export", handler.export).Methods("GET")
}

// Dump the system model.
//   params:
//     request The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (handler *Handler) export(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    dump, err := handler.manager.Export()
    if err != nil {
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
    dhttp.RespondWithJSON(w, http.StatusOK, dump)
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package app

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

var logger = log.WithField("package", "server.app")

// Handler structure that contains the link with the underlying manager.
type Handler struct {
    manager Manager
}

// NewHandler obtains a new handler.
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
func (h* Handler) SetRoutes(router * mux.Router) {
    logger.Info("Setting application routes")
    router.HandleFunc("/api/v0/app/{networkID}/descriptor/add", h.addDescriptor).Methods("POST")
    router.HandleFunc("/api/v0/app/{networkID}/descriptor/list", h.listDescriptor).Methods("GET")
    router.HandleFunc("/api/v0/app/{networkID}/descriptor/{appDescriptorID}/info", h.getDescriptor).Methods("GET")
    router.HandleFunc("/api/v0/app/{networkID}/descriptor/{appDescriptorID}/delete", h.deleteDescriptor).Methods("DELETE")
    router.HandleFunc("/api/v0/app/{networkID}/instance/add", h.addInstance).Methods("POST")
    router.HandleFunc("/api/v0/app/{networkID}/instance/list", h.listInstances).Methods("GET")
    router.HandleFunc("/api/v0/app/{networkID}/instance/{appInstanceID}/info", h.getInstance).Methods("GET")
    router.HandleFunc("/api/v0/app/{networkID}/instance/{appInstanceID}/update", h.updateInstance).Methods("POST")
    router.HandleFunc("/api/v0/app/{networkID}/instance/{appInstanceID}/delete", h.deleteInstance).Methods("DELETE")
}

// Add a new descriptor to an existing network.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h * Handler) addDescriptor(w http.ResponseWriter, r * http.Request){
    defer r.Body.Close()
    logger.Debug("addDescriptor")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    addDescriptorRequest := &entities.AddAppDescriptorRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    }else{
        err = json.Unmarshal(b, &addDescriptorRequest)
        if err == nil && addDescriptorRequest.IsValid() {
            logger.Debug("Adding new descriptor: " + addDescriptorRequest.String())
            added, err := h.manager.AddApplicationDescriptor(networkID, * addDescriptorRequest)
            if err != nil {
                dhttp.RespondWithError(w, http.StatusInternalServerError, err)
            }else{
                dhttp.RespondWithJSON(w, http.StatusOK, added)
            }
        }else{
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(addDescriptorRequest, errors.InvalidEntity))
        }
    }
}

// List all the descriptors in an existing network.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h * Handler) listDescriptor(w http.ResponseWriter, r * http.Request){
    defer r.Body.Close()
    logger.Debug("listDescriptor")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    descriptors, err := h.manager.ListDescriptors(networkID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, descriptors)
    }else{
        dhttp.RespondWithError(w, http.StatusBadRequest, err)
    }
}

// Get the information of an application descriptor.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h * Handler) getDescriptor(w http.ResponseWriter, r * http.Request){
    defer r.Body.Close()
    logger.Debug("getDescriptor")

    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    appDescriptorID := vars["appDescriptorID"]
    if appDescriptorID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("appDescriptorID"))
        return
    }

    descriptor, err := h.manager.GetDescriptor(networkID, appDescriptorID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, descriptor)
    }else{
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Delete an application descriptor.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h * Handler) deleteDescriptor(w http.ResponseWriter, r * http.Request){
    defer r.Body.Close()
    logger.Debug("deleteDescriptor")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    appDescriptorID := vars["appDescriptorID"]
    if appDescriptorID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("appDescriptorID"))
        return
    }

    err := h.manager.DeleteDescriptor(networkID, appDescriptorID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK,entities.NewSuccessfulOperation("DeleteDescriptor"))
    }else{
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}


// Add a new instance to an existing network.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h * Handler) addInstance(w http.ResponseWriter, r * http.Request){
    defer r.Body.Close()
    logger.Debug("addInstance")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }

    addInstanceRequest := &entities.AddAppInstanceRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    }else{
        err = json.Unmarshal(b, &addInstanceRequest)
        if err == nil && addInstanceRequest.IsValid() {
            logger.Debug("Adding new instance: " + addInstanceRequest.String())
            added, err := h.manager.AddApplicationInstance(networkID,* addInstanceRequest)
            if err != nil {
                dhttp.RespondWithError(w, http.StatusInternalServerError, err)
            }else{
                dhttp.RespondWithJSON(w, http.StatusOK, added)
            }
        }else{
            msg := "invalid instance request"
            logger.Error(msg + " => " + addInstanceRequest.String())
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(addInstanceRequest, errors.InvalidEntity))
        }
    }
}

// List all the instances in an existing network.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h * Handler) listInstances(w http.ResponseWriter, r * http.Request){
    defer r.Body.Close()
    logger.Debug("listInstances")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    instances, err := h.manager.ListInstances(networkID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, instances)
    }else{
        dhttp.RespondWithError(w, http.StatusBadRequest, err)
    }
}

// Get the instance information.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h * Handler) getInstance(w http.ResponseWriter, r * http.Request){
    defer r.Body.Close()
    logger.Debug("getInstance")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    appInstanceID := vars["appInstanceID"]
    if appInstanceID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("appInstanceID"))
        return
    }

    instance, err := h.manager.GetInstance(networkID, appInstanceID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK, instance)
    }else{
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

// Update the information of an existing instance.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h * Handler) updateInstance(w http.ResponseWriter, r * http.Request){
    defer r.Body.Close()
    logger.Debug("updateInstance")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    appInstanceID := vars["appInstanceID"]
    if appInstanceID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("appInstanceID"))
        return
    }

    updateInstanceRequest := &entities.UpdateAppInstanceRequest{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.IOError))
    }else{
        err = json.Unmarshal(b, &updateInstanceRequest)
        if err == nil && updateInstanceRequest.IsValid() {
            logger.Debug("Updating instance: " + appInstanceID + " with " + updateInstanceRequest.String())
            updated, err := h.manager.UpdateInstance(networkID, appInstanceID, * updateInstanceRequest)
            if err != nil {
                dhttp.RespondWithError(w, http.StatusInternalServerError, err)
            }else{
                dhttp.RespondWithJSON(w, http.StatusOK, updated)
            }
        }else{
            dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewEntityError(updateInstanceRequest, errors.InvalidEntity))
        }
    }

}

// Delete an existing instance.
//   params:
//     request The HTTP request structure
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h * Handler) deleteInstance(w http.ResponseWriter, r * http.Request){
    defer r.Body.Close()
    logger.Debug("deleteInstance")
    vars := mux.Vars(r)
    networkID := vars["networkID"]
    if networkID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("networkID"))
        return
    }
    appInstanceID := vars["appInstanceID"]
    if appInstanceID == "" {
        dhttp.RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(errors.MissingRESTParameter).WithParams("appInstanceID"))
        return
    }

    err := h.manager.DeleteInstance(networkID, appInstanceID)
    if err == nil {
        dhttp.RespondWithJSON(w, http.StatusOK,entities.NewSuccessfulOperation("DeleteInstance"))
    }else{
        dhttp.RespondWithError(w, http.StatusInternalServerError, err)
    }
}

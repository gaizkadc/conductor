//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Applications client rest implementation

package client

import (

    "fmt"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"

    "github.com/daishogroup/dhttp"
)

// AddDescriptorURI with the URI pattern to add application descriptors.
const AddDescriptorURI = "/api/v0/app/%s/descriptor/add"
// ListDescriptorsURI with the URI pattern to list application descriptors.
const ListDescriptorsURI = "/api/v0/app/%s/descriptor/list"
// GetDescriptorURI with the URI pattern to get an application descriptor.
const GetDescriptorURI = "/api/v0/app/%s/descriptor/%s/info"
// DeleteDescriptorURI with the URI pattern to delete an application descriptor.
const DeleteDescriptorURI = "/api/v0/app/%s/descriptor/%s/delete"
// AddInstanceURI with the URI pattern to add new application instances.
const AddInstanceURI = "/api/v0/app/%s/instance/add"
// ListInstancesURI with the URI pattern to list instances.
const ListInstancesURI = "/api/v0/app/%s/instance/list"
// GetInstanceURI with the URI pattern to get an application instance.
const GetInstanceURI = "/api/v0/app/%s/instance/%s/info"
// UpdateInstanceURI with the URI pattern to update an existing application instance.
const UpdateInstanceURI = "/api/v0/app/%s/instance/%s/update"
// DeleteInstanceURI with the URI pattern to delete an application instance.
const DeleteInstanceURI = "/api/v0/app/%s/instance/%s/delete"

// ApplicationsRest structure with the rest client.
type ApplicationsRest struct {
    client dhttp.Client
}

// NewApplicationRest creates a REST applications client.
// Deprecated: Use NewApplicationClientRest
func NewApplicationRest(basePath string) Applications {
    return NewApplicationClientRest(ParseHostPort(basePath))
}

func NewApplicationClientRest(host string, port int) Applications {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host,port))
    return &ApplicationsRest{rest}
}

// AddApplicationDescriptor adds a new application descriptor.
//   params:
//     networkId The network identifier.
//     descriptor The application descriptor.
//   returns:
//     The added application descriptor.
//     An error if the application cannot be added.
func (rest * ApplicationsRest) AddApplicationDescriptor(networkID string, descriptor entities.AddAppDescriptorRequest) (*entities.AppDescriptor, derrors.DaishoError){
    url := fmt.Sprintf(AddDescriptorURI, networkID)
    response := rest.client.Post(url, descriptor, new(entities.AppDescriptor))
    if response.Ok() {
        added := response.Result.(*entities.AppDescriptor)
        return added, nil
    }
    return nil, response.Error
}

// ListDescriptors obtains a list of all the application descriptors available for a given network.
//   params:
//     networkId The network identifier.
//   returns:
//     An array of application descriptors.
//     An error in case the list cannot be retrieved.
func (rest * ApplicationsRest) ListDescriptors(networkID string) ([] entities.AppDescriptor, derrors.DaishoError){
    url := fmt.Sprintf(ListDescriptorsURI, networkID)
    response := rest.client.Get(url, new([] entities.AppDescriptor))
    if response.Ok() {
        descriptors := response.Result.(* [] entities.AppDescriptor)
        return *descriptors, nil
    }
    return nil, response.Error
}

// GetDescriptor obtains an application descriptor
//   params:
//     networkId The network identifier.
//     appDescriptorId The application descriptor identifier.
//   returns:
//     An application descriptor
//     An error if the descriptor cannot be retrieved or is not associated with the network.
func (rest * ApplicationsRest) GetDescriptor(networkID string, appDescriptorID string) (* entities.AppDescriptor, derrors.DaishoError){
    url := fmt.Sprintf(GetDescriptorURI, networkID, appDescriptorID)
    response := rest.client.Get(url, new(entities.AppDescriptor))
    if response.Ok() {
        descriptor := response.Result.(* entities.AppDescriptor)
        return descriptor, nil
    }
    return nil, response.Error
}

// DeleteDescriptor deletes an application descriptor
//   params:
//     networkID The network identifier.
//     appDescriptorID The application descriptor identifier.
//   returns:
//     An error if the descriptor cannot be retrieved or is not associated with the network.
func (rest * ApplicationsRest) DeleteDescriptor(networkID string, appDescriptorID string) derrors.DaishoError {
    url := fmt.Sprintf(DeleteDescriptorURI, networkID, appDescriptorID)
    response := rest.client.Delete(url, new(entities.SuccessfulOperation))
    return response.Error
}

// AddApplicationInstance adds a new application instance.
//   params:
//     networkID The network identifier.
//     instance The application instance.
//   returns:
//     The added instance.
//     An error if the application instance cannot be added.
func (rest * ApplicationsRest) AddApplicationInstance(networkID string, instance entities.AddAppInstanceRequest) (* entities.AppInstance, derrors.DaishoError){
    url := fmt.Sprintf(AddInstanceURI, networkID)
    response := rest.client.Post(url, instance, new(entities.AppInstance))
    if response.Ok() {
        added := response.Result.(* entities.AppInstance)
        return added, nil
    }
    return nil, response.Error
}

// ListInstances obtains a list of all the application instances inside a network.
//   params:
//     networkId The network identifier.
//   returns:
//     An array of application instances.
//     An error if the list cannot be retrieved.
func (rest * ApplicationsRest) ListInstances(networkID string) ([] entities.AppInstance, derrors.DaishoError){
    url := fmt.Sprintf(ListInstancesURI, networkID)
    response := rest.client.Get(url, new([] entities.AppInstance))
    if response.Ok() {
        instances := response.Result.(* [] entities.AppInstance)
        return *instances, nil
    }
    return nil, response.Error
}

// GetInstance retrieves an application instance.
//   params:
//     networkId The network identifier.
//     appInstanceId The application instance identifier.
//   returns:
//     An application instance.
//     An error if the instance cannot be retrieved.
func (rest * ApplicationsRest) GetInstance(networkID string, appInstanceID string) (* entities.AppInstance, derrors.DaishoError){
    url := fmt.Sprintf(GetInstanceURI, networkID, appInstanceID)
    response := rest.client.Get(url, new(entities.AppInstance))
    if response.Ok() {
        instance := response.Result.(* entities.AppInstance)
        return instance, nil
    }
    return nil, response.Error
}

// UpdateInstance updates an application instance.
//   params:
//     networkId The network identifier.
//     appInstanceId The application instance identifier.
//     request The application update request.
//   returns:
//     The updated application instance.
//     An error if the instance cannot be retrieved.
func (rest * ApplicationsRest) UpdateInstance(
    networkID string, appInstanceID string, request entities.UpdateAppInstanceRequest) (* entities.AppInstance, derrors.DaishoError){
    url := fmt.Sprintf(UpdateInstanceURI, networkID, appInstanceID)
    response := rest.client.Post(url, request, new(entities.AppInstance))
    if response.Ok() {
        instance := response.Result.(* entities.AppInstance)
        return instance, nil
    }
    return nil, response.Error
}

// DeleteInstance deletes an application instance.
//   params:
//     networkId The network identifier.
//     appInstanceId The application instance identifier.
//   returns:
//     An error if the instance cannot be removed.
func (rest * ApplicationsRest) DeleteInstance(networkID string, appInstanceID string) derrors.DaishoError {
    url := fmt.Sprintf(DeleteInstanceURI, networkID, appInstanceID)
    response := rest.client.Delete(url, new(entities.SuccessfulOperation))
    return response.Error
}

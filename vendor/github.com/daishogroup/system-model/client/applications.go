//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Applications Client interface
package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Applications is an interface that represents the client for application related operations.
type Applications interface {

    // AddApplicationDescriptor adds a new application descriptor.
    //   params:
    //     networkID The network identifier.
    //     descriptor The application descriptor.
    //   returns:
    //     The added application descriptor.
    //     An error if the application cannot be added.
    AddApplicationDescriptor(networkID string, descriptor entities.AddAppDescriptorRequest) (*entities.AppDescriptor, derrors.DaishoError)

    // ListDescriptors lists all the application descriptors available for a given network.
    //   params:
    //     networkID The network identifier.
    //   returns:
    //     An array of application descriptors.
    //     An error in case the list cannot be retrieved.
    ListDescriptors(networkID string) ([] entities.AppDescriptor, derrors.DaishoError)

    // GetDescriptor gets an application descriptor
    //   params:
    //     networkID The network identifier.
    //     appDescriptorID The application descriptor identifier.
    //   returns:
    //     An application descriptor
    //     An error if the descriptor cannot be retrieved or is not associated with the network.
    GetDescriptor(networkID string, appDescriptorID string) (* entities.AppDescriptor, derrors.DaishoError)

    // DeleteDescriptor deletes an application descriptor
    //   params:
    //     networkID The network identifier.
    //     appDescriptorID The application descriptor identifier.
    //   returns:
    //     An error if the descriptor cannot be retrieved or is not associated with the network.
    DeleteDescriptor(networkID string, appDescriptorID string) derrors.DaishoError

    // AddApplicationInstance adds a new application instance.
    //   params:
    //     networkID The network identifier.
    //     instance The application instance.
    //   returns:
    //     The added instance.
    //     An error if the application instance cannot be added.
    AddApplicationInstance(networkID string, instance entities.AddAppInstanceRequest) (* entities.AppInstance, derrors.DaishoError)

    // ListInstances obtains the list of all the application instances inside a network.
    //   params:
    //     networkID The network identifier.
    //   returns:
    //     An array of application instances.
    //     An error if the list cannot be retrieved.
    ListInstances(networkID string) ([] entities.AppInstance, derrors.DaishoError)

    // GetInstance retrieves an application instance.
    //   params:
    //     networkID The network identifier.
    //     appInstanceID The application instance identifier.
    //   returns:
    //     An application instance.
    //     An error if the instance cannot be retrieved.
    GetInstance(networkID string, appInstanceID string) (* entities.AppInstance, derrors.DaishoError)

    // UpdateInstance updates an application instance.
    //   params:
    //     networkID The network identifier.
    //     appInstanceID The application instance identifier.
    //     request The application update request.
    //   returns:
    //     The updated application instance.
    //     An error if the instance cannot be retrieved.
    UpdateInstance(networkID string, appInstanceID string, request entities.UpdateAppInstanceRequest) (* entities.AppInstance, derrors.DaishoError)

    // DeleteInstance deletes an application instance.
    //   params:
    //     networkID The network identifier.
    //     appInstanceID The application instance identifier.
    //   returns:
    //     An error if the instance cannot be removed.
    DeleteInstance(networkID string, appInstanceID string) derrors.DaishoError

}

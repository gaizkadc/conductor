//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Definition of the operations of an application descriptor provider.

package appdescstorage

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

//Provider is the interface od the AppDescriptorProvider.
type Provider interface {
    // Add a new application descriptor to the system.
    //   params:
    //     descriptor The application descriptor to be added
    //   returns:
    //     An error if the descriptor cannot be added.
    Add(descriptor entities.AppDescriptor) derrors.DaishoError

    // Check if an application descriptor exists in the system.
    //   params:
    //     descriptorID The application descriptor identifier.
    //   returns:
    //     Whether the descriptor exists or not.
    Exists(descriptorID string) bool

    // Retrieve a given application descriptor.
    //   params:
    //     descriptorID The application descriptor identifier.
    //   returns:
    //     The application descriptor.
    //     An error if the descriptor cannot be retrieved.
    RetrieveDescriptor(descriptorID string) (* entities.AppDescriptor, derrors.DaishoError)

    // Delete a given application descriptor.
    //   params:
    //     descriptorID The application descriptor identifier.
    //   returns:
    //     An error if the application descriptor cannot be removed.
    Delete(descriptorID string) derrors.DaishoError

    // Dump obtains the list of all app descriptors in the system.
    //   returns:
    //     The list of AppDescriptors.
    //     An error if the list cannot be retrieved.
    Dump() ([] entities.AppDescriptor, derrors.DaishoError)

    // ReducedInfoList get a list with the reduced info.
    //   returns:
    //     List of the reduced app info.
    //     An error if the descriptor cannot be retrieved.
    ReducedInfoList() ([] entities.AppDescriptorReducedInfo, derrors.DaishoError)
}

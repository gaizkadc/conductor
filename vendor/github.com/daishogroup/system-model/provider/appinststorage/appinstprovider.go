//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Definition of the operations of an application instance provider.

package appinststorage

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

//Provider is the interface of the AppInstance providers.
type Provider interface {

    // Add a new application instance to the system.
    //   params:
    //     instance The application instance to be added
    //   returns:
    //     An error if the instance cannot be added.
    Add(instance entities.AppInstance) derrors.DaishoError

    // Update an instance in the system.
    //   params:
    //     instance The new instance information. The instance identifier will be used and cannot be modified.
    //   returns:
    //     An error if the instance cannot be updated.
    Update(instance entities.AppInstance) derrors.DaishoError

    // Check if an application instance exists in the system.
    //   params:
    //     instanceID The application instance identifier.
    //   returns:
    //     Whether the instance exists or not.
    Exists(instanceID string) bool

    // Retrieve a given application instance.
    //   params:
    //     instanceID The application instance identifier.
    //   returns:
    //     The application instance.
    //     An error if the instance cannot be retrieved.
    RetrieveInstance(instanceID string) (* entities.AppInstance, derrors.DaishoError)

    // Delete a given instance.
    //   params:
    //     instanceID The application instance identifier.
    //   returns:
    //     An error if the instance cannot be removed.
    Delete(instanceID string) derrors.DaishoError

    // Dump obtains the list of all application instances in the system.
    //   returns:
    //     The list of AppInstance.
    //     An error if the list cannot be retrieved.
    Dump() ([] entities.AppInstance, derrors.DaishoError)

    // ReducedInfoList get a list with the reduced info.
    //   returns:
    //     List of the reduced info.
    //     An error if the info cannot be retrieved.
    ReducedInfoList() ([] entities.AppInstanceReducedInfo, derrors.DaishoError)
}
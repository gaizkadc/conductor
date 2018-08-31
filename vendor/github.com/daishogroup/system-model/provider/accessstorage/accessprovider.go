
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the access provider specification.

package accessstorage

import (
"github.com/daishogroup/derrors"
"github.com/daishogroup/system-model/entities"
)

// Provider definition of the node persistence-related methods.
type Provider interface {

    // Add new access entry to the system.
    //   params:
    //     user The user to be added
    //   returns:
    //     An error if the user cannot be added.
    Add(user entities.UserAccess) derrors.DaishoError

    // Check if a user exists in the system.
    //   params:
    //     userID The user identifier.
    //   returns:
    //     Whether the user exists or not.
    Exists(userID string) bool

    // Retrieve a given user access.
    //   params:
    //     userID The user identifier.
    //   returns:
    //     The user.
    //     An error if the user cannot be retrieved.
    RetrieveAccess(userID string) (* entities.UserAccess, derrors.DaishoError)

    // Delete a given user access.
    //   params:
    //     userID The user identifier.
    //   returns:
    //     An error if the user cannot be removed.
    Delete(userID string) derrors.DaishoError

    // Update a user in the system.
    //   params:
    //     node The new user information.
    //   returns:
    //     An error if the user cannot be updated.
    Update(user entities.UserAccess) derrors.DaishoError

    // ReducedInfoList get a list with the reduced info.
    //   returns:
    //     List of the reduced info.
    //     An error if the info cannot be retrieved.
    ReducedInfoList() ([] entities.UserAccessReducedInfo, derrors.DaishoError)

    // Dump obtains the list of all access in the system.
    //   returns:
    //     The list of user.
    //     An error if the list cannot be retrieved.
    Dump() ([] entities.UserAccess, derrors.DaishoError)
}

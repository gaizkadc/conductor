
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the user provider specification.

package userstorage

import (
"github.com/daishogroup/derrors"
"github.com/daishogroup/system-model/entities"
)

// Provider definition of the node persistence-related methods.
type Provider interface {

    // Add a new user to the system.
    //   params:
    //     user The user to be added
    //   returns:
    //     An error if the user cannot be added.
    Add(user entities.User) derrors.DaishoError

    // Check if a user exists in the system.
    //   params:
    //     userID The user identifier.
    //   returns:
    //     Whether the user exists or not.
    Exists(userID string) bool

    // Retrieve a given user.
    //   params:
    //     userID The user identifier.
    //   returns:
    //     The user.
    //     An error if the user cannot be retrieved.
    RetrieveUser(userID string) (* entities.User, derrors.DaishoError)

    // Delete a given user.
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
    Update(user entities.User) derrors.DaishoError

    // ReducedInfoList get a list with the reduced info.
    //   returns:
    //     List of the reduced info.
    //     An error if the info cannot be retrieved.
    ReducedInfoList() ([] entities.UserReducedInfo, derrors.DaishoError)

    // Dump obtains the list of all user in the system.
    //   returns:
    //     The list of user.
    //     An error if the list cannot be retrieved.
    Dump() ([] entities.User, derrors.DaishoError)
}


// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the user provider specification.

package passwordstorage

import (
"github.com/daishogroup/derrors"
"github.com/daishogroup/system-model/entities"
)

// Provider definition of the node persistence-related methods.
type Provider interface {

    // Add a new user to the system.
    //   params:
    //     password The user to be added
    //   returns:
    //     An error if the user cannot be added.
    Add(password entities.Password) derrors.DaishoError

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
    //     The password.
    //     An error if the user cannot be retrieved.
    RetrievePassword(userID string) (* entities.Password, derrors.DaishoError)

    // Delete an existing password.
    //   params:
    //     userID The user identifier.
    //   returns:
    //     An error if the user cannot be removed.
    Delete(userID string) derrors.DaishoError

    // Update a password in the system.
    //   params:
    //     node The new user information.
    //   returns:
    //     An error if the user cannot be updated.
    Update(password entities.Password) derrors.DaishoError


    // Dump obtains the list of all passwords in the system.
    //   returns:
    //     The list of user.
    //     An error if the list cannot be retrieved.
    Dump() ([] entities.Password, derrors.DaishoError)
}

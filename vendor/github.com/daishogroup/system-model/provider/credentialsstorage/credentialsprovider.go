
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the user credentials provider specification.

package credentialsstorage

import (
"github.com/daishogroup/derrors"
"github.com/daishogroup/system-model/entities"
)

// Provider definition of the node persistence-related methods.
type Provider interface {

    // Add a new user to the system.
    //   params:
    //     credentials User credentials.
    //   returns:
    //     An error if any.
    Add(credentials entities.Credentials) derrors.DaishoError

    // Check if a user already has credentials.
    //   params:
    //     uuID The user identifier.
    //   returns:
    //     Whether the user exists or not.
    Exists(uuID string) bool

    // Retrieve the credentials.
    //   params:
    //     uuID The user identifier.
    //   returns:
    //     The user credentials.
    //     An error if the user cannot be retrieved.
    Retrieve(uuID string) (* entities.Credentials, derrors.DaishoError)

    // Delete credentials.
    //   params:
    //     uuID The user identifier.
    //   returns:
    //     An error if the user cannot be removed.
    Delete(uuID string) derrors.DaishoError

    // Update a user in the system.
    //   params:
    //     node The new credential information.
    //   returns:
    //     An error if the user cannot be updated.
    Update(credentials entities.Credentials) derrors.DaishoError

}

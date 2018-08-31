
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the user provider specification.

package oauthstorage

import (
"github.com/daishogroup/derrors"
"github.com/daishogroup/system-model/entities"
)

// Provider definition of the OAuth secrets persistence-related methods.
type Provider interface {

    // Add new secrets to the system.
    //   params:
    //     secrets The secrets to be added.
    //   returns:
    //     An error if the user cannot be added.
    Add(secrets entities.OAuthSecrets) derrors.DaishoError

    // Check if a user exists in the system.
    //   params:
    //     userID The user identifier.
    //   returns:
    //     Whether the user exists or not.
    Exists(userID string) bool

    // Retrieve a given user secrets.
    //   params:
    //     userID The user identifier.
    //   returns:
    //     The secrets.
    //     An error if the user cannot be retrieved.
    Retrieve(userID string) (* entities.OAuthSecrets, derrors.DaishoError)

    // Delete an existing secrets collection.
    //   params:
    //     userID The user identifier.
    //   returns:
    //     An error if the user cannot be removed.
    Delete(userID string) derrors.DaishoError

    // Update a secrets collection in the system.
    //   params:
    //     secrets The collection of secrets
    //   returns:
    //     An error if the user cannot be updated.
    Update(secrets entities.OAuthSecrets) derrors.DaishoError


    // Dump obtains the list of all secrets in the system.
    //   returns:
    //     The list of user.
    //     An error if the list cannot be retrieved.
    Dump() ([] entities.OAuthSecrets, derrors.DaishoError)
}

//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the oauth Client interface

package client

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/derrors"
)

type OAuth interface {
    // Set the OAuth information entry for a certain app.
    //  params:
    //   userID user identifier.
    //   setEntryRequest request with all the information required.
    //  return:
    //   Error if any.
    SetSecret(userID string, setEntryRequest entities.OAuthAddEntryRequest) derrors.DaishoError

    // Delete the set of secrets of an existing user.
    //  params:
    //   userID user identifier.
    //  return:
    //   Error if any.
    DeleteSecrets(userID string) derrors.DaishoError

    // Get the set of secrets of an existing user.
    //  params:
    //   userID The user identifier.
    //  return:
    //   Set of oauth secrets.
    //   Error if any.
    GetSecrets(userID string) (*entities.OAuthSecrets, derrors.DaishoError)
}
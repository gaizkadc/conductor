//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Interface for the credentials client.

package client

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/derrors"
)

type Credentials interface {


    // Add existing credential,
    //  params:
    //   request: Credentials to be added.
    //  return:
    //   Error if any.
    Add(request entities.AddCredentialsRequest) derrors.DaishoError

    // Get existing credential.
    //  params:
    //   uuid: Identifier of the target credential.
    //  return:
    //   The credentials object or error if any.
    Get(uuid string) (*entities.Credentials, derrors.DaishoError)

    // Delete an existing credential.
    //  params:
    //   uuid: Identifier of the target credential.
    //  return:
    //   Error if any.
    Delete(uuid string) derrors.DaishoError

}
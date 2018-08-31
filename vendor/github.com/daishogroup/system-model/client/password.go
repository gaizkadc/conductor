//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the oauth Client interface

package client

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/derrors"
)

type Password interface {
    // Get a user oauth
    // params:
    //  userID user identifier.
    // return:
    //  Password entity.
    //  Error if any.
    GetPassword(userID string) (*entities.Password, derrors.DaishoError)

    // Set a oauth entry with a new value.
    //  params:
    //   oauth New oauth entity.
    //  return:
    //   Error if any.
    SetPassword(password entities.Password) derrors.DaishoError

    // Delete an existing oauth.
    //  params:
    //   userID the user identifier
    //  return:
    //   Error if any.
    DeletePassword(userID string) derrors.DaishoError
}
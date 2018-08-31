//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the Access Client interface

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Access elements interface
type Access interface {

    // AddAccess adds a new user access
    //   params:
    //     userID  The user identifier.
    //     request The add user request.
    //   returns:
    //     The user access entity.
    //     An error if the application cannot be added.
    AddAccess(userID string, request entities.AddUserAccessRequest) (*entities.UserAccess, derrors.DaishoError)

    // SetAccess sets the access roles for an existing user.
    //   params:
    //     userID  The user identifier.
    //     request The add user request.
    //   returns:
    //     The user access entity.
    //     An error if the application cannot be added.
    SetAccess(userID string, request entities.AddUserAccessRequest) (*entities.UserAccess, derrors.DaishoError)

    // GetAccess gets user access privilege entry for an existing user.
    //   params:
    //     userID The user identifier.
    //   returns:
    //     UserAccess values
    //     An error in case the list cannot be retrieved.
    GetAccess(userID string) (* entities.UserAccess, derrors.DaishoError)

    // DeleteAccess deletes user access privilege.
    //   params:
    //     userID The user id.
    //   returns:
    //     Error if any.
    DeleteAccess(networkID string) derrors.DaishoError

    // ListAccess gets a list of user privileges.
    //   returns:
    //     Complete list of users with their access roles.
    //     An error if the user does not exist.
    ListAccess() ([]entities.UserAccessReducedInfo, derrors.DaishoError)

}


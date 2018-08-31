//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the User Client interface

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// User interface that defines the operations available on Users.
type User interface {
    // Add a new user
    // params:
    //    user User request.
    // return:
    //    Generated user.
    //    Error if any.
    Add(user entities.AddUserRequest) (*entities.User, derrors.DaishoError)

    // Get an existing user data.
    // params:
    //    userID User identifier.
    // return:
    //    Found user data.
    //    Error if any.
    Get(userID string) (*entities.User, derrors.DaishoError)

    // Delete an existing user entry.
    // params:
    //    userID User identifier.
    // return:
    //    Error if any.
    Delete(userID string) derrors.DaishoError

    // Update an existing user entry.
    // params:
    //   userID User identifier.
    //   updateRequest Update request.
    // return:
    //   Error if any.
    Update(userID string, update entities.UpdateUserRequest) (*entities.User, derrors.DaishoError)

    // ListUsers returns the list of users with their corresponding role/privileges.
    //  return:
    //   List of users with their roles.
    ListUsers() ([]entities.UserExtended, derrors.DaishoError)

}
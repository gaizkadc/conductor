//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the User Rest Client

package client

import (
    "fmt"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/dhttp"
)

const (
    UserAddURI    = "/api/v0/user/add"
    UserGetURI    = "/api/v0/user/%s/get"
    UserDeleteURI = "/api/v0/user/%s/delete"
    UserUpdateURI = "/api/v0/user/%s/update"
    UserListURI   = "/api/v0/user/list"
)

// Client Rest for user resources.
type UserRest struct {
    client dhttp.Client
}

// Deprecated: Use NewUserClientRest
func NewUserRest(basePath string) User {
    return NewUserClientRest(ParseHostPort(basePath))
}

func NewUserClientRest(host string, port int) User {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host, port))
    return &UserRest{rest}
}

// Add a new user
// params:
//    user User request.
// return:
//    Generated user.
//    Error if any.
func (rest *UserRest) Add(user entities.AddUserRequest) (*entities.User, derrors.DaishoError) {
    response := rest.client.Post(UserAddURI, user, new(entities.User))
    if response.Error != nil {
        return nil, response.Error
    } else {
        result := response.Result.(*entities.User)
        return result, nil
    }
}

// Get an existing user data.
// params:
//    userId User identifier.
// return:
//    Found user data.
//    Error if any.
func (rest *UserRest) Get(userId string) (*entities.User, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(UserGetURI, userId), new(entities.User))
    if response.Error != nil {
        return nil, response.Error
    } else {
        result := response.Result.(*entities.User)
        return result, nil
    }
}

// Delete an existing user entry.
// params:
//    userId User identifier.
// return:
//    Error if any.
func (rest *UserRest) Delete(userId string) derrors.DaishoError {
    response := rest.client.Delete(fmt.Sprintf(UserDeleteURI, userId), new(entities.SuccessfulOperation))
    if response.Error != nil {
        return response.Error
    } else {
        return nil
    }
}

// Update an existing user entry.
// params:
//   userId User identifier.
//   updateRequest Update request.
// return:
//   Error if any.
func (rest *UserRest) Update(userId string, update entities.UpdateUserRequest) (*entities.User, derrors.DaishoError) {
    response := rest.client.Post(fmt.Sprintf(UserUpdateURI, userId), update, new(entities.User))
    if response.Error != nil {
        return nil, response.Error
    } else {
        result := response.Result.(*entities.User)
        return result, nil
    }
}

// Return the list of users with their corresponding role/privileges.
//  return:
//   List of users with their roles.
func (rest *UserRest) ListUsers() ([]entities.UserExtended, derrors.DaishoError) {
    response := rest.client.Get(UserListURI, new([]entities.UserExtended))
    if response.Error != nil {
        return nil, response.Error
    } else {
        result := response.Result.(*[]entities.UserExtended)
        return *result, nil
    }
}

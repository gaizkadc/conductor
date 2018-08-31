//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// API REST client for access entities

package client

import (
    "fmt"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/dhttp"
)

// AddAccessURI with the URI pattern to add user access.
const AddAccessURI = "/api/v0/access/%s/add"
// SetAccessURI with the URI pattern to set user access.
const SetAccessURI = "/api/v0/access/%s/set"
// GetAccessURI with the URI pattern to get user access.
const GetAccessURI = "/api/v0/access/%s/get"
// DeleteAccessURI with the URI pattern to delete user access.
const DeleteAccessURI = "/api/v0/access/%s/delete"
// ListAccessURI with the URI pattern to list user access.
const ListAccessURI = "/api/v0/access/list"

// AccessRest structure with the client to make the requests.
type AccessRest struct {
    client dhttp.Client
}

// NewAccessRest creates an Access client using REST requests.
// Deprecated: Use NewAccessClientRest
func NewAccessRest(basePath string) Access {
    return NewAccessClientRest(ParseHostPort(basePath))
}

func NewAccessClientRest(host string, port int) Access {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host, port))
    return &AccessRest{rest}
}

// AddAccess adds a new user access.
//   params:
//     userID  The user identifier.
//     request The add user request.
//   returns:
//     The user access entity.
//     An error if the application cannot be added.
func (rest *AccessRest) AddAccess(userID string, request entities.AddUserAccessRequest) (*entities.UserAccess, derrors.DaishoError) {
    url := fmt.Sprintf(AddAccessURI, userID)
    response := rest.client.Post(url, request, new(entities.UserAccess))
    if response.Ok() {
        added := response.Result.(*entities.UserAccess)
        return added, nil
    }
    return nil, response.Error
}

// SetAccess sets the access roles for an existing user.
//   params:
//     userID  The user identifier.
//     request The add user request.
//   returns:
//     The user access entity.
//     An error if the application cannot be added.
func (rest *AccessRest) SetAccess(userID string, request entities.AddUserAccessRequest) (*entities.UserAccess, derrors.DaishoError) {
    url := fmt.Sprintf(SetAccessURI, userID)
    response := rest.client.Post(url, request, new(entities.UserAccess))
    if response.Ok() {
        added := response.Result.(*entities.UserAccess)
        return added, nil
    }
    return nil, response.Error
}

// GetAccess gets user access privilege entry for an existing user.
//   params:
//     userID The user identifier.
//   returns:
//     UserAccess values
//     An error in case the list cannot be retrieved.
func (rest *AccessRest) GetAccess(userID string) (*entities.UserAccess, derrors.DaishoError) {
    url := fmt.Sprintf(GetAccessURI, userID)
    response := rest.client.Get(url, new(entities.UserAccess))
    if response.Ok() {
        descriptor := response.Result.(*entities.UserAccess)
        return descriptor, nil
    }
    return nil, response.Error
}

// DeleteAccess deletes user access privilege.
//   params:
//     userID The user id.
//   returns:
//     Error if any.
func (rest *AccessRest) DeleteAccess(userID string) derrors.DaishoError {
    url := fmt.Sprintf(DeleteAccessURI, userID)
    response := rest.client.Delete(url, new(entities.SuccessfulOperation))
    return response.Error
}

// ListAccess gets a list of user privileges.
//   returns:
//     Complete list of users with their access roles.
//     An error if the user does not exist.
func (rest *AccessRest) ListAccess() ([]entities.UserAccessReducedInfo, derrors.DaishoError) {
    response := rest.client.Get(ListAccessURI, new([] entities.UserAccessReducedInfo))
    if response.Ok() {
        list := response.Result.(*[] entities.UserAccessReducedInfo)
        return *list, nil
    }
    return nil, response.Error
}

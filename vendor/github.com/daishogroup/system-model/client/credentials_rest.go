//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// REST implementation of the credentials client.

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "fmt"
    "github.com/daishogroup/dhttp"
)

const CredentialsAddURI = "/api/v0/credentials/add"
const CredentialsGetURI = "/api/v0/credentials/%s/get"
const CredentialsDeleteURI = "/api/v0/credentials/%s/delete"

type CredentialsRest struct {
    client dhttp.Client
}

// Deprecated: Use NewCredentialsClientRest
func NewCredentialsRest(basePath string) Credentials {
    return NewCredentialsClientRest(ParseHostPort(basePath))
}

func NewCredentialsClientRest(host string, port int) Credentials {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host, port))
    return &CredentialsRest{rest}
}

// Add existing credential.
//  params:
//   request: Credentials to be added.
//  return:
//   Error if any.
func (c *CredentialsRest) Add(request entities.AddCredentialsRequest) derrors.DaishoError {
    resp := c.client.Post(CredentialsAddURI, request, new(entities.SuccessfulOperation))
    if resp.Error != nil {
        return resp.Error
    } else {
        return nil
    }
}

// Get existing credential.
//  params:
//   uuid: Identifier of the target credential.
//  return:
//   The credentials object or error if any.
func (c *CredentialsRest) Get(uuid string) (*entities.Credentials, derrors.DaishoError) {
    resp := c.client.Get(fmt.Sprintf(CredentialsGetURI, uuid), new(entities.Credentials))
    if resp.Error != nil {
        return nil, resp.Error
    } else {
        return resp.Result.(*entities.Credentials), nil
    }
}

// Delete an existing credential.
//  params:
//   uuid: Identifier of the target credential.
//  return:
//   Error if any.
func (c *CredentialsRest) Delete(uuid string) derrors.DaishoError {
    resp := c.client.Delete(fmt.Sprintf(CredentialsDeleteURI, uuid), new(entities.SuccessfulOperation))
    if resp.Error != nil {
        return resp.Error
    }
    return nil
}

//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the oauth Client interface

package client

import (
    "fmt"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/dhttp"
)

const (
    PasswordGetURI    = "/api/v0/password/%s"
    PasswordSetURI    = "/api/v0/password"
    PasswordDeleteURI = "/api/v0/password/%s"
)

type PasswordRest struct {
    client dhttp.Client
}

// Deprecated: Use NewPasswordClientRest
func NewPasswordRest(basePath string) Password {
    return NewPasswordClientRest(ParseHostPort(basePath))
}

func NewPasswordClientRest(host string, port int) Password {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host, port))
    return &PasswordRest{rest}
}

// Get a user oauth
// params:
//  userID user identifier.
// return:
//  Password entity.
//  Error if any.
func (rest *PasswordRest) GetPassword(userID string) (*entities.Password, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(PasswordGetURI, userID), new(entities.Password))
    if response.Error != nil {
        return nil, response.Error
    } else {
        result := response.Result.(*entities.Password)
        return result, nil
    }

}

// Set a oauth entry with a new value.
//  params:
//   oauth New oauth entity.
//  return:
//   Error if any.
func (rest *PasswordRest) SetPassword(password entities.Password) derrors.DaishoError {
    response := rest.client.Post(PasswordSetURI, password, new(entities.SuccessfulOperation))
    if response.Error != nil {
        return response.Error
    }
    return nil
}

// Delete an existing oauth.
//  params:
//   userID the user identifier
//  return:
//   Error if any.
func (rest *PasswordRest) DeletePassword(userID string) derrors.DaishoError {
    response := rest.client.Delete(fmt.Sprintf(PasswordDeleteURI, userID), new(entities.SuccessfulOperation))
    if response.Error != nil {
        return response.Error
    } else {
        return nil
    }
}

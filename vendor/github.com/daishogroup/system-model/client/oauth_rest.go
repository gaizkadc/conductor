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
    OAuthSetSecrets = "/api/v0/oauth/%s"
    OAuthGetSecrets = "/api/v0/oauth/%s"
    OAuthDeleteSecrets = "/api/v0/oauth/%s"
)


type OAuthRest struct {
    client dhttp.Client
}
// Deprecated: Use NewOAuthClientRest
func NewOAuthRest(basePath string) OAuth{
    return NewOAuthClientRest(ParseHostPort(basePath))
}

func NewOAuthClientRest(host string, port int) OAuth{
    rest:=dhttp.NewClientSling(dhttp.NewRestBasicConfig(host,port))
    return &OAuthRest{rest}
}


// Set the OAuth information entry for a certain app.
//  params:
//   userID user identifier.
//   setEntryRequest request with all the information required.
//  return:
//   Error if any.
func(rest *OAuthRest) SetSecret(userID string, setEntryRequest entities.OAuthAddEntryRequest) derrors.DaishoError {
    response := rest.client.Post(fmt.Sprintf(OAuthSetSecrets, userID), setEntryRequest, new(entities.SuccessfulOperation))
    if response.Error != nil {
        return response.Error
    } else {
        return nil
    }
}

// Delete the set of secrets of an existing user.
//  params:
//   userID user identifier.
//  return:
//   Error if any.
func(rest *OAuthRest) DeleteSecrets(userID string) derrors.DaishoError {
    response := rest.client.Delete(fmt.Sprintf(OAuthDeleteSecrets, userID), new(entities.SuccessfulOperation))
    if response.Error != nil {
        return response.Error
    } else {
        return nil
    }
}

// Get the set of secrets of an existing user.
//  params:
//   userID The user identifier.
//  return:
//   Set of oauth secrets.
//   Error if any.
func(rest *OAuthRest) GetSecrets(userID string) (*entities.OAuthSecrets, derrors.DaishoError){
    response := rest.client.Get(fmt.Sprintf(OAuthGetSecrets, userID), new(entities.OAuthSecrets))
    if response.Error != nil {
        return nil, response.Error
    } else {
        return response.Result.(*entities.OAuthSecrets),nil
    }
}

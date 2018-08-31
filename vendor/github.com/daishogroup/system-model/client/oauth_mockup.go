//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the oauth Client interface

package client

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/provider/oauthstorage"
    "github.com/daishogroup/system-model/server/oauth"
)

type OAuthMockup struct {
    oauthProvider *oauthstorage.MockupOAuthProvider
    oauthMgr      oauth.Manager
}

func NewOAuthMockup() OAuth{
    oauthProvider := oauthstorage.NewMockupOAuthProvider()
    oauthMgr := oauth.NewManager(oauthProvider)
    return &OAuthMockup{oauthProvider: oauthProvider, oauthMgr: oauthMgr}
}

func(m *OAuthMockup) ClearMockup() {
    m.oauthProvider.Clear()
}

func (m *OAuthMockup) InitMockup() {
    m.oauthProvider.Clear()
    m.oauthProvider.Add(entities.NewOAuthSecrets("userDefault"))
}

// Set the OAuth information entry for a certain app.
//  params:
//   userID user identifier.
//   setEntryRequest request with all the information required.
//  return:
//   Error if any.
func(m *OAuthMockup) SetSecret(userID string, setEntryRequest entities.OAuthAddEntryRequest) derrors.DaishoError {
    return m.oauthMgr.SetSecret(userID, setEntryRequest)
}

// Delete the set of secrets of an existing user.
//  params:
//   userID user identifier.
//  return:
//   Error if any.
func(m *OAuthMockup) DeleteSecrets(userID string) derrors.DaishoError {
    return m.oauthMgr.DeleteSecrets(userID)
}

// Get the set of secrets of an existing user.
//  params:
//   userID The user identifier.
//  return:
//   Set of oauth secrets.
//   Error if any.
func(m *OAuthMockup) GetSecrets(userID string) (*entities.OAuthSecrets, derrors.DaishoError){
    return m.oauthMgr.GetSecrets(userID)
}

//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the credentials Client Mockup

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/credentialsstorage"
    "github.com/daishogroup/system-model/server/credentials"
)

type CredentialMockup struct {
    credentialProvider *credentialsstorage.MockupCredentialsProvider
    credentialsMgr      credentials.Manager
}

func NewCredentialsMockup() Credentials {
    var credentialsProvider = credentialsstorage.NewMockupCredentialsProvider()
    var credentialsMgr = credentials.NewManager(credentialsProvider)

    return &CredentialMockup{credentialsProvider, credentialsMgr}
}

func (mockup *CredentialMockup) ClearMockup() {
    mockup.credentialProvider.Clear()
}


// Add existing credential,
//  params:
//   request: Credentials to be added.
//  return:
//   Error if any.
func (mockup *CredentialMockup) Add(request entities.AddCredentialsRequest) derrors.DaishoError {
    return mockup.credentialsMgr.AddCredentials(request)
}

// Get existing credential.
//  params:
//   uuid: Identifier of the target credential.
//  return:
//   The credentials object or error if any.
func (mockup *CredentialMockup) Get(uuid string) (*entities.Credentials, derrors.DaishoError) {
    return mockup.credentialsMgr.GetCredentials(uuid)
}

// Delete an existing credential.
//  params:
//   uuid: Identifier of the target credential.
//  return:
//   Error if any.
func (mockup *CredentialMockup) Delete(uuid string) derrors.DaishoError {
    return mockup.credentialsMgr.DeleteCredentials(uuid)
}
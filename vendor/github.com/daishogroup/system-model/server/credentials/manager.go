//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Credentials manager.

package credentials

import (
    "github.com/daishogroup/system-model/provider/credentialsstorage"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/derrors"
)

// The Manager struct provides access to user related methods.
type Manager struct {
    credentialsProvider credentialsstorage.Provider
}

// NewManager creates a new user manager.
//   params:
//     credentialsProvider: credentials provider
//   returns:
//     A manager.
func NewManager(userCredentialsProvider credentialsstorage.Provider ) Manager {
    return Manager{userCredentialsProvider}
}

// Add existing credential,
//  params:
//   request: Credentials to be added.
//  return:
//   Error if any.
func(m *Manager) AddCredentials(request entities.AddCredentialsRequest) derrors.DaishoError {
    newEntry := entities.NewCredentials(request.UUID, request.PublicKey, request.PrivateKey,
        request.Description, request.TypeKey)
    return m.credentialsProvider.Add(*newEntry)
}

// Get existing credential.
//  params:
//   uuid: Identifier of the target credential.
//  return:
//   The credentials object or error if any.
func(m *Manager) GetCredentials(uuid string) (*entities.Credentials, derrors.DaishoError) {
    return m.credentialsProvider.Retrieve(uuid)
}

// Delete an existing credential.
//  params:
//   uuid: Identifier of the target credential.
//  return:
//   Error if any.
func(m *Manager) DeleteCredentials(uuid string) derrors.DaishoError {
    return m.credentialsProvider.Delete(uuid)
}
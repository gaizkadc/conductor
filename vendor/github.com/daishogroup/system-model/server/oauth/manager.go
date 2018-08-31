//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the manager implementation for the oauth handler.

package oauth

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/provider/oauthstorage"
    "github.com/daishogroup/system-model/errors"
)

type Manager struct {
    provider oauthstorage.Provider
}

// Build a new manager with the corresponding provider.
//  params:
//   provider OAuth provider.
//  return:
//   Instantiated manager.
func NewManager(provider oauthstorage.Provider) Manager {
    return Manager{provider}
}

// Set the OAuth information entry for a certain app.
//  params:
//   userID user identifier.
//   setEntryRequest request with all the information required.
//  return:
//   Error if any.
func(m *Manager) SetSecret(userID string, setEntryRequest entities.OAuthAddEntryRequest) derrors.DaishoError {
    if !m.provider.Exists(userID) {
        return derrors.NewOperationError(errors.UserDoesNotExist)
    }

    instance,err := m.provider.Retrieve(userID)
    if err != nil {
        return derrors.NewOperationError(errors.UserDoesNotExist, err)
    }
    err = instance.AddEntry(setEntryRequest.AppName,entities.NewOAuthEntry(setEntryRequest.ClientID, setEntryRequest.Secret))
    if err != nil {
        return err
    }
    err = m.provider.Update(*instance)
    if err != nil {
        return err
    }

    return nil
}

// Delete the set of secrets of an existing user.
//  params:
//   userID user identifier.
//  return:
//   Error if any.
func(m *Manager) DeleteSecrets(userID string) derrors.DaishoError {
    if !m.provider.Exists(userID) {
        return derrors.NewOperationError(errors.UserDoesNotExist)
    }
    return m.provider.Delete(userID)
}

// Get the set of secrets of an existing user.
//  params:
//   userID The user identifier.
//  return:
//   Set of oauth secrets.
//   Error if any.
func(m *Manager) GetSecrets(userID string) (*entities.OAuthSecrets,derrors.DaishoError) {
    if !m.provider.Exists(userID) {
        return nil, derrors.NewOperationError(errors.UserDoesNotExist)
    }
    return m.provider.Retrieve(userID)
}
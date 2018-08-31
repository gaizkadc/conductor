//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the manager implementation for the password handler.

package password

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/provider/passwordstorage"
)

type Manager struct {
    provider passwordstorage.Provider
}

// Build a new manager with the corresponding provider.
//  params:
//   provider Password provider.
//  return:
//   Instantiated manager.
func NewManager(provider passwordstorage.Provider) Manager {
    return Manager{provider}
}


// Get a user password
// params:
//  userID user identifier.
// return:
//  Password entity.
//  Error if any.
func(m *Manager) GetPassword(userID string) (*entities.Password, derrors.DaishoError){
   return m.provider.RetrievePassword(userID)
}

// Set a password entry with a new value.
//  params:
//   password New password entity.
//  return:
//   Error if any.
func(m *Manager) SetPassword(password entities.Password) derrors.DaishoError{
    return m.provider.Update(password)
}

// Delete an existing password.
//  params:
//   userID the user identifier
//  return:
//   Error if any.
func(m *Manager) DeletePassword(userID string) derrors.DaishoError{
    return m.provider.Delete(userID)
}


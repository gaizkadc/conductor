//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the session manager in charge of the business logic behind session entities.

package session

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/system-model/provider/sessionstorage"
)

// The Manager struct provides access to session related methods.
type Manager struct {
    sessionProvider     sessionstorage.Provider
}

// Create a new session manager.
//   params:
//     sessionProvider     The session storage provider.
//   returns:
//     A manager.
func NewManager(sessionProvider sessionstorage.Provider) Manager {
    return Manager{sessionProvider}
}

// AddNode adds a new node to an existing cluster.
//   params:
//     user         The user to be added.
//   returns:
//     The added node.
//     An error if the node cannot be added.
func (manager *Manager) AddSession(request entities.AddSessionRequest) (*entities.Session, derrors.DaishoError) {
    if manager.sessionProvider.Exists(request.NewSession.ID) {
        logger.Errorf("Session already exists")
        return nil, derrors.NewOperationError(errors.SessionAlreadyExists).WithParams(request.NewSession.ID)
    }
    err := manager.sessionProvider.Add(request.NewSession)
    if err != nil {
        logger.Errorf("Error adding session fot %s", request.NewSession.ID)
        return nil, err
    }
    return &request.NewSession, nil
}


// Get an existing user.
//   params:
//     sessionID   The session id.
//   returns:
//     The session entity.
//     An error if the user does not exist.
func (manager *Manager) GetSession(sessionID string) (*entities.Session, derrors.DaishoError) {
    return manager.sessionProvider.Retrieve(sessionID)
}

// Delete an existing user.
//   params:
//     userID   The user id.
//   returns:
//     An error if the user does not exist.
func (manager *Manager) DeleteSession(sessionID string) derrors.DaishoError {
    // delete user
    err := manager.sessionProvider.Delete(sessionID)
    if err != nil {
        logger.Errorf("Error deleting session entry for %s", sessionID)
        return err
    }

    return err
}

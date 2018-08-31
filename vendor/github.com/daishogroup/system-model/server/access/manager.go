//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the user manager in charge of the business logic behind user access.

package access

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/system-model/provider/accessstorage"
)


// The Manager struct provides access to access related methods.
type Manager struct {
    accessProvider accessstorage.Provider
}

// NewManager creates a new access manager.
//   params:
//     accessProvider     The node storage provider.
//   returns:
//     A manager.
func NewManager(accessProvider accessstorage.Provider) Manager {
    return Manager{accessProvider}
}

// AddAccess adds new roles to an existing user.
//   params:
//     userID         The user to be added.
//     request        The access request.
//   returns:
//     The added access.
//     An error if the node cannot be added.
func (manager *Manager) AddAccess(userID string, request entities.AddUserAccessRequest) (*entities.UserAccess, derrors.DaishoError) {
    // Check if the users exists or not
    if !manager.accessProvider.Exists(userID) {
        return nil, derrors.NewOperationError(errors.UserDoesNotExist)
    }

    currentEntry, err := manager.accessProvider.RetrieveAccess(userID)
    if err != nil {
        // entry dissappeared between the previous call and this one
        return nil, derrors.NewOperationError(errors.UserDoesNotHaveRoles)
    }
    currentRoles := make(map[entities.RoleType]bool)
    for _, r := range currentEntry.Roles {
        currentRoles[r]=true
    }
    newRoles := currentEntry.Roles
    // fill it
    for _, incomingRole := range request.Roles {
        _, exists := currentRoles[incomingRole]
        if !exists {
            newRoles = append(newRoles, incomingRole)
        }
    }
    // update
    u := entities.NewUserAccess(userID, newRoles)
    err = manager.accessProvider.Update(*u)

    if err != nil {
        return nil, err
    }
    return u, nil
}

// SetAccess set the roles for an existing user removing any previous configuration.
//   params:
//     userID         The user to be added.
//     request        The access request.
//   returns:
//     The added access.
//     An error if the node cannot be added.
func (manager *Manager) SetAccess(userID string, request entities.AddUserAccessRequest) (*entities.UserAccess, derrors.DaishoError) {
    // Check if the users exists or not

    u := entities.NewUserAccess(userID, request.Roles)

    var err derrors.DaishoError
    if !manager.accessProvider.Exists(userID) {
        err = manager.accessProvider.Add(*u)
    } else {
        // update
        err = manager.accessProvider.Update(*u)
    }

    if err != nil {
        return nil, err
    }
    return u, nil
}



// Get an existing access.
//   params:
//     userId   The user id.
//   returns:
//     The access entity.
//     An error if the user does not exist.
func (manager *Manager) GetAccess(userId string) (*entities.UserAccess, derrors.DaishoError) {
    return manager.accessProvider.RetrieveAccess(userId)
}

// Delete an existing user access.
//   params:
//     userId   The user id.
//   returns:
//     An error if the user does not exist.
func (manager *Manager) DeleteAccess(userId string) derrors.DaishoError {
    return manager.accessProvider.Delete(userId)
}

// Get the list of accesses.
//   returns:
//     Complete list of users with their access roles.
//     An error if the user does not exist.
func (manager *Manager) ListAccess() ([]entities.UserAccessReducedInfo, derrors.DaishoError) {
    return manager.accessProvider.ReducedInfoList()
}
//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the user manager in charge of the business logic behind user entities.

package user

import (
    "sort"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/daishogroup/system-model/provider/accessstorage"
    "github.com/daishogroup/system-model/provider/passwordstorage"
    "github.com/daishogroup/system-model/provider/oauthstorage"
)

// The Manager struct provides access to user related methods.
type Manager struct {
    userProvider        userstorage.Provider
    accessProvider      accessstorage.Provider
    passwordProvider    passwordstorage.Provider
    oauthProvider       oauthstorage.Provider
}

// NewManager creates a new user manager.
//   params:
//     userProvider     The node storage provider.
//     accessProvider   The access storage provider.
//     passwordProvider The password storage provider.
//     oauthProvider    The oauth storage provider.
//   returns:
//     A manager.
func NewManager(userProvider userstorage.Provider, accessProvider accessstorage.Provider,
    passwordProvider passwordstorage.Provider, oauthProvider oauthstorage.Provider) Manager {
    return Manager{userProvider, accessProvider, passwordProvider,
    oauthProvider}
}

// AddNode adds a new node to an existing cluster.
//   params:
//     user         The user to be added.
//   returns:
//     The added node.
//     An error if the node cannot be added.
func (manager *Manager) AddUser(request entities.AddUserRequest) (*entities.User, derrors.DaishoError) {
    if manager.userProvider.Exists(request.ID) {
        logger.Errorf("Error user already exists.")
        return nil, derrors.NewOperationError(errors.UserAlreadyExists).WithParams(request.ID)
    }
    u := entities.NewUserWithID(request.ID, request.Name, request.Phone, request.Email,
        request.CreationTime.Time, request.ExpirationTime.Time)
    err := manager.userProvider.Add(*u)
    if err != nil {
        logger.Errorf("Error adding new user %s",request.ID)
        return nil, err
    }

    // Add the OAuth entry
    err = manager.oauthProvider.Add(entities.NewOAuthSecrets(request.ID))
    if err != nil {
        logger.Errorf("Error adding new oauth secrets for %s",request.ID)
        return nil, err
    }

    // Add the corresponding list of accesses
    err = manager.accessProvider.Add(*entities.NewUserAccess(request.ID, make([]entities.RoleType,0)))
    if err != nil{
        logger.Errorf("Error adding new user access for %s",request.ID)
        return nil, err
    }
    // Add the password entry
    pass, err := entities.NewPassword(request.ID, nil)
    if err != nil {
        return nil, derrors.NewOperationError(errors.UserPasswordError, err)
    }
    err = manager.passwordProvider.Add(*pass)
    if err!= nil {
        logger.Errorf("Error adding new password for %s",request.ID)
        return nil, derrors.NewOperationError(errors.UserPasswordError).CausedBy(err)
    }
    return u, nil
}


// Get an existing user.
//   params:
//     userID   The user id.
//   returns:
//     The user entity.
//     An error if the user does not exist.
func (manager *Manager) GetUser(userID string) (*entities.User, derrors.DaishoError) {
    return manager.userProvider.RetrieveUser(userID)
}

// Delete an existing user.
//   params:
//     userID   The user id.
//   returns:
//     An error if the user does not exist.
func (manager *Manager) DeleteUser(userID string) derrors.DaishoError {
    // delete user
    err := manager.userProvider.Delete(userID)
    if err != nil {
        logger.Errorf("Error deleting user entry for %s",userID)
        return err
    }
    // delete his oauth entries
    err = manager.oauthProvider.Delete(userID)
    if err != nil {
        logger.Errorf("Error deleting oauth entry for %s",userID)
        return err
   }

    // delete his roles
    err = manager.accessProvider.Delete(userID)
    if err != nil {
        logger.Errorf("Error deleting access entry for %s",userID)
        return err
    }
    // delete his password entry
    err = manager.passwordProvider.Delete(userID)
    if err!= nil {
        logger.Errorf("Error deleting password entry for %s",userID)
    }
    return err
}

// Update an existing user.
//  params:
//    userID The user id.
//    request Update request.
//  return:
//    New user entity after update.
//    Errors if any.
func (manager *Manager) UpdateUser(userID string, request entities.UpdateUserRequest) (*entities.User, derrors.DaishoError) {

    current, err := manager.userProvider.RetrieveUser(userID)
    if err != nil {
        return nil, err
    }
    updatedUser := current.Merge(request)
    err = manager.userProvider.Update(*updatedUser)
    if err != nil {
        return nil, err
    }

    return updatedUser, err
}


// Return the list of users with their corresponding role/privileges.
//  return:
//   List of users with their roles.
func (manager *Manager) ListUsers() ([]entities.UserExtended, derrors.DaishoError) {
    toReturn := make([] entities.UserExtended,0)

    users, err := manager.userProvider.Dump()
    if err != nil {
        return nil, err
    }
    if len(users)==0{
        return toReturn, nil
    }
    access, err := manager.accessProvider.Dump()

    if len(users)!= len(access){
        return nil, derrors.NewOperationError("Number of users and roles mismatch")
    }

    // sort users
    sort.Slice(users, func(i, j int) bool {
        return users[i].ID < users[j].ID
    })
    sort.Slice(access, func(i, j int) bool {
        return access[i].UserID < access[j].UserID
    })

    // Join both structures to create something useful to be returned
    for index, a := range(users){
        if !access[index].Roles[0].IsInternalUser() {
            toReturn = append(toReturn, entities.NewUserExtended(a.ID, a.Name, a.Phone, a.Email,
                a.CreationTime.Time, a.ExpirationTime.Time, access[index].Roles))
        }
    }

    return toReturn, nil
}

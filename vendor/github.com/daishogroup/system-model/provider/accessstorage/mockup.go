//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Mockup access provider.

package accessstorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupUserProvider is a mockup in-memory implementation of an access provider.
type MockupUserAccessProvider struct {
    sync.Mutex
    // Users indexed by ID
    users map[string]entities.UserAccess
}

// It creates a mockup provider for node operations.
func NewMockupUserAccessProvider() *MockupUserAccessProvider {
    return &MockupUserAccessProvider{users:make(map[string]entities.UserAccess)}
}

// Add a new user to the system.
//   params:
//     user The User to be added
//   returns:
//     An error if the node cannot be added.
func (mockup *MockupUserAccessProvider) Add(user entities.UserAccess) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(user.UserID){
        mockup.users[user.UserID] = user
        return nil
    }
    return derrors.NewOperationError(errors.UserAlreadyExists).WithParams(user)
}

// Exists checks if a node exists in the system.
//   params:
//     userID The Node identifier.
//   returns:
//     Whether the user exists or not.
func (mockup *MockupUserAccessProvider) Exists(userID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExists(userID)
}

func (mockup *MockupUserAccessProvider) unsafeExists(userID string) bool {
    _, exists := mockup.users[userID]
    return exists
}

// RetrieveNode retrieves a given user.
//   params:
//     userID The user identifier.
//   returns:
//     The cluster.
//     An error if the cluster cannot be retrieved.
func (mockup *MockupUserAccessProvider) RetrieveAccess(userID string) (*entities.UserAccess, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    node, exists := mockup.users[userID]
    if exists {
        return &node, nil
    }
    return nil, derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}

// Delete access entry for a given user.
//   params:
//     userID The user identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (mockup *MockupUserAccessProvider) Delete(userID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    _, exists := mockup.users[userID]
    if exists {
        delete(mockup.users, userID)
        return nil
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}

// Update access in the system.
//   params:
//     node The new node information.
//   returns:
//     An error if the instance cannot be updated.
func (mockup *MockupUserAccessProvider) Update(user entities.UserAccess) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(user.UserID) {
        mockup.users[user.UserID] = user
        return nil
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(user)
}

// ReducedInfoList get a list with the reduced info.
//   returns:
//     List of the reduced info.
//     An error if the info cannot be retrieved.
func (mockup *MockupUserAccessProvider) ReducedInfoList() ([] entities.UserAccessReducedInfo, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    result := make([] entities.UserAccessReducedInfo, 0, len(mockup.users))
    for _, n := range mockup.users {
        reducedInfo := entities.NewUserAccessReducedInfo(n.UserID, n.Roles)
        result = append(result, *reducedInfo)
    }
    return result, nil
}

// Dump obtains the list of all users in the system.
//   returns:
//     The list of users.
//     An error if the list cannot be retrieved.
func (mockup * MockupUserAccessProvider) Dump() ([] entities.UserAccess, derrors.DaishoError) {
    result := make([] entities.UserAccess, 0)
    mockup.Lock()
    defer mockup.Unlock()
    for _, node := range mockup.users {
        result = append(result, node)
    }
    return result, nil
}

// Clear is a util function to clear the contents of the mockup.
func (mockup *MockupUserAccessProvider) Clear() {
    mockup.Lock()
    mockup.users = make(map[string]entities.UserAccess)
    mockup.Unlock()
}

//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Mockup users provider.

package userstorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupUserProvider is a mockup in-memory implementation of a user provider.
type MockupUserProvider struct {
    sync.Mutex
    // Users indexed by ID
    users map[string]entities.User
}

// NewMockupUserProvider creates a mockup provider for node operations.
func NewMockupUserProvider() *MockupUserProvider {
    return &MockupUserProvider{users:make(map[string]entities.User)}
}

// Add a new user to the system.
//   params:
//     user The User to be added
//   returns:
//     An error if the node cannot be added.
func (mockup *MockupUserProvider) Add(user entities.User) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(user.ID){
        mockup.users[user.ID] = user
        return nil
    }
    return derrors.NewOperationError(errors.UserAlreadyExists).WithParams(user)
}

// Exists checks if a node exists in the system.
//   params:
//     userID The Node identifier.
//   returns:
//     Whether the user exists or not.
func (mockup *MockupUserProvider) Exists(userID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExists(userID)
}

func (mockup *MockupUserProvider) unsafeExists(userID string) bool {
    _, exists := mockup.users[userID]
    return exists
}

// RetrieveNode retrieves a given user.
//   params:
//     userID The user identifier.
//   returns:
//     The cluster.
//     An error if the cluster cannot be retrieved.
func (mockup *MockupUserProvider) RetrieveUser(userID string) (*entities.User, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    node, exists := mockup.users[userID]
    if exists {
        return &node, nil
    }
    return nil, derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}

// Delete a given user.
//   params:
//     userID The user identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (mockup *MockupUserProvider) Delete(userID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    _, exists := mockup.users[userID]
    if exists {
        delete(mockup.users, userID)
        return nil
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}

// Update a node in the system.
//   params:
//     node The new node information.
//   returns:
//     An error if the instance cannot be updated.
func (mockup *MockupUserProvider) Update(user entities.User) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(user.ID) {
        mockup.users[user.ID] = user
        return nil
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(user)
}

// ReducedInfoList get a list with the reduced info.
//   returns:
//     List of the reduced info.
//     An error if the info cannot be retrieved.
func (mockup *MockupUserProvider) ReducedInfoList() ([] entities.UserReducedInfo, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    result := make([] entities.UserReducedInfo, 0, len(mockup.users))
    for _, n := range mockup.users {
        reducedInfo := entities.NewUserReducedInfo(n.ID, n.Email)
        result = append(result, *reducedInfo)
    }
    return result, nil
}

// Dump obtains the list of all users in the system.
//   returns:
//     The list of users.
//     An error if the list cannot be retrieved.
func (mockup * MockupUserProvider) Dump() ([] entities.User, derrors.DaishoError) {
    result := make([] entities.User, 0)
    mockup.Lock()
    defer mockup.Unlock()
    for _, node := range mockup.users {
        result = append(result, node)
    }
    return result, nil
}

// Clear is a util function to clear the contents of the mockup.
func (mockup *MockupUserProvider) Clear() {
    mockup.Lock()
    mockup.users = make(map[string]entities.User)
    mockup.Unlock()
}

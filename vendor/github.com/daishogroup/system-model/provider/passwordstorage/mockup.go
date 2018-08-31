//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Mockup password provider.

package passwordstorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupUserProvider is a mockup in-memory implementation of a user provider.
type MockupPasswordProvider struct {
    sync.Mutex
    // Passwords indexed by ID
    passwords map[string]entities.Password
}

// NewMockupUserProvider creates a mockup provider for node operations.
func NewMockupPasswordProvider() *MockupPasswordProvider {
    return &MockupPasswordProvider{passwords: make(map[string]entities.Password)}
}

// Add a new user to the system.
//   params:
//     user The User to be added
//   returns:
//     An error if the node cannot be added.
func (mockup *MockupPasswordProvider) Add(password entities.Password) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(password.UserID){
        mockup.passwords[password.UserID] = password
        return nil
    }
    return derrors.NewOperationError(errors.UserAlreadyExists).WithParams(password)
}

// Exists checks if a node exists in the system.
//   params:
//     userID The Node identifier.
//   returns:
//     Whether the user exists or not.
func (mockup *MockupPasswordProvider) Exists(userID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExists(userID)
}

func (mockup *MockupPasswordProvider) unsafeExists(userID string) bool {
    _, exists := mockup.passwords[userID]
    return exists
}

// RetrieveNode retrieves a given user.
//   params:
//     userID The user identifier.
//   returns:
//     The cluster.
//     An error if the cluster cannot be retrieved.
func (mockup *MockupPasswordProvider) RetrievePassword(userID string) (*entities.Password, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    password, exists := mockup.passwords[userID]
    if exists {
        return &password, nil
    }
    return nil, derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}

// Delete a given user.
//   params:
//     userID The user identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (mockup *MockupPasswordProvider) Delete(userID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    _, exists := mockup.passwords[userID]
    if exists {
        delete(mockup.passwords, userID)
        return nil
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}

// Update a node in the system.
//   params:
//     node The new node information.
//   returns:
//     An error if the instance cannot be updated.
func (mockup *MockupPasswordProvider) Update(password entities.Password) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(password.UserID) {
        mockup.passwords[password.UserID] = password
        return nil
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(password)
}

// Dump obtains the list of all users in the system.
//   returns:
//     The list of users.
//     An error if the list cannot be retrieved.
func (mockup * MockupPasswordProvider) Dump() ([] entities.Password, derrors.DaishoError) {
    result := make([] entities.Password, 0)
    mockup.Lock()
    defer mockup.Unlock()
    for _, password := range mockup.passwords {
        result = append(result, password)
    }
    return result, nil
}

// Clear is a util function to clear the contents of the mockup.
func (mockup *MockupPasswordProvider) Clear() {
    mockup.Lock()
    mockup.passwords = make(map[string]entities.Password)
    mockup.Unlock()
}

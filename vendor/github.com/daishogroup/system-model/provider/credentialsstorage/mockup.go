//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Mockup credentials provider.

package credentialsstorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupCredentialsProvider is a mockup in-memory implementation of a user provider.
type MockupCredentialsProvider struct {
    sync.Mutex
    // Users indexed by ID
    credentials map[string]entities.Credentials
}

// NewMockupUserProvider creates a mockup provider for node operations.
func NewMockupCredentialsProvider() *MockupCredentialsProvider {
    return &MockupCredentialsProvider{credentials:make(map[string]entities.Credentials)}
}

// Add a new user to the system.
//   params:
//     user The User to be added
//   returns:
//     An error if the node cannot be added.
func (mockup *MockupCredentialsProvider) Add(credentials entities.Credentials) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(credentials.UUID){
        mockup.credentials[credentials.UUID] = credentials
        return nil
    }
    return derrors.NewOperationError(errors.CredentialsAlreadyExist).WithParams(credentials)
}

// Exists checks if a node exists in the system.
//   params:
//     userID The Node identifier.
//   returns:
//     Whether the user exists or not.
func (mockup *MockupCredentialsProvider) Exists(uuid string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExists(uuid)
}

func (mockup *MockupCredentialsProvider) unsafeExists(uuid string) bool {
    _, exists := mockup.credentials[uuid]
    return exists
}

// RetrieveNode retrieves a given user.
//   params:
//     userID The user identifier.
//   returns:
//     The cluster.
//     An error if the cluster cannot be retrieved.
func (mockup *MockupCredentialsProvider) Retrieve(uuid string) (*entities.Credentials, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    node, exists := mockup.credentials[uuid]
    if exists {
        return &node, nil
    }
    return nil, derrors.NewOperationError(errors.CredentialsDoNotExist).WithParams(uuid)
}

// Delete a given user.
//   params:
//     userID The user identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (mockup *MockupCredentialsProvider) Delete(userID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    _, exists := mockup.credentials[userID]
    if exists {
        delete(mockup.credentials, userID)
        return nil
    }
    return derrors.NewOperationError(errors.CredentialsDoNotExist).WithParams(userID)
}

// Update a node in the system.
//   params:
//     node The new node information.
//   returns:
//     An error if the instance cannot be updated.
func (mockup *MockupCredentialsProvider) Update(credentials entities.Credentials) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(credentials.UUID) {
        mockup.credentials[credentials.UUID] = credentials
        return nil
    }
    return derrors.NewOperationError(errors.CredentialsDoNotExist).WithParams(credentials)
}


// Clear is a util function to clear the contents of the mockup.
func (mockup *MockupCredentialsProvider) Clear() {
    mockup.Lock()
    mockup.credentials = make(map[string]entities.Credentials)
    mockup.Unlock()
}

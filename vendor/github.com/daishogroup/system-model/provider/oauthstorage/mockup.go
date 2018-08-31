//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Mockup password provider.

package oauthstorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupUserProvider is a mockup in-memory implementation of a oauth provider.
type MockupOAuthProvider struct {
    sync.Mutex
    // Secrets indexed by user ID
    secrets map[string]entities.OAuthSecrets
}

// NewMockupUserProvider creates a mockup provider for node operations.
func NewMockupOAuthProvider() *MockupOAuthProvider {
    return &MockupOAuthProvider{secrets: make(map[string]entities.OAuthSecrets)}
}

// Add a new user to the system.
//   params:
//     user The User to be added
//   returns:
//     An error if the node cannot be added.
func (mockup *MockupOAuthProvider) Add(secrets entities.OAuthSecrets) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(secrets.UserID){
        mockup.secrets[secrets.UserID] = secrets
        return nil
    }
    return derrors.NewOperationError(errors.UserAlreadyExists).WithParams(secrets)
}

// Exists checks if a node exists in the system.
//   params:
//     userID The Node identifier.
//   returns:
//     Whether the user exists or not.
func (mockup *MockupOAuthProvider) Exists(userID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExists(userID)
}

func (mockup *MockupOAuthProvider) unsafeExists(userID string) bool {
    _, exists := mockup.secrets[userID]
    return exists
}

// RetrieveNode retrieves a given user.
//   params:
//     userID The user identifier.
//   returns:
//     The cluster.
//     An error if the cluster cannot be retrieved.
func (mockup *MockupOAuthProvider) Retrieve(userID string) (*entities.OAuthSecrets, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    secrets, exists := mockup.secrets[userID]
    if exists {
        return &secrets, nil
    }
    return nil, derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}

// Delete a given user.
//   params:
//     userID The user identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (mockup *MockupOAuthProvider) Delete(userID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    _, exists := mockup.secrets[userID]
    if exists {
        delete(mockup.secrets, userID)
        return nil
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}

// Update a node in the system.
//   params:
//     node The new node information.
//   returns:
//     An error if the instance cannot be updated.
func (mockup *MockupOAuthProvider) Update(secrets entities.OAuthSecrets) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(secrets.UserID) {
        mockup.secrets[secrets.UserID] = secrets
        return nil
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(secrets)
}

// Dump obtains the list of all users in the system.
//   returns:
//     The list of users.
//     An error if the list cannot be retrieved.
func (mockup *MockupOAuthProvider) Dump() ([] entities.OAuthSecrets, derrors.DaishoError) {
    result := make([] entities.OAuthSecrets, 0)
    mockup.Lock()
    defer mockup.Unlock()
    for _, secrets := range mockup.secrets {
        result = append(result, secrets)
    }
    return result, nil
}

// Clear is a util function to clear the contents of the mockup.
func (mockup *MockupOAuthProvider) Clear() {
    mockup.Lock()
    mockup.secrets = make(map[string]entities.OAuthSecrets)
    mockup.Unlock()
}

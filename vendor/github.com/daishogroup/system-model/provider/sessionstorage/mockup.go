//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Mockup sessions provider.

package sessionstorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupSessionProvider is a mockup in-memory implementation of a session provider.
type MockupSessionProvider struct {
    sync.Mutex
    // sessions indexed by ID
    sessions map[string]entities.Session
}

// NewMockupsessionProvider creates a mockup provider for node operations.
func NewMockupSessionProvider() *MockupSessionProvider {
    return &MockupSessionProvider{sessions:make(map[string]entities.Session)}
}

// Add a new session to the system.
//   params:
//     session The session to be added
//   returns:
//     An error if the node cannot be added.
func (mockup *MockupSessionProvider) Add(session entities.Session) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(session.ID){
        mockup.sessions[session.ID] = session
        return nil
    }
    return derrors.NewOperationError(errors.SessionAlreadyExists).WithParams(session)
}

// Exists checks if a node exists in the system.
//   params:
//     sessionID The Node identifier.
//   returns:
//     Whether the session exists or not.
func (mockup *MockupSessionProvider) Exists(sessionID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExists(sessionID)
}

func (mockup *MockupSessionProvider) unsafeExists(sessionID string) bool {
    _, exists := mockup.sessions[sessionID]
    return exists
}

// RetrieveNode retrieves a given session.
//   params:
//     sessionID The session identifier.
//   returns:
//     The cluster.
//     An error if the cluster cannot be retrieved.
func (mockup *MockupSessionProvider) Retrieve(sessionID string) (*entities.Session, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    node, exists := mockup.sessions[sessionID]
    if exists {
        return &node, nil
    }
    return nil, derrors.NewOperationError(errors.SessionDoesNotExist).WithParams(sessionID)
}

// Delete a given session.
//   params:
//     sessionID The session identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (mockup *MockupSessionProvider) Delete(sessionID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    _, exists := mockup.sessions[sessionID]
    if exists {
        delete(mockup.sessions, sessionID)
        return nil
    }
    return derrors.NewOperationError(errors.SessionDoesNotExist).WithParams(sessionID)
}

// Update a node in the system.
//   params:
//     node The new node information.
//   returns:
//     An error if the instance cannot be updated.
func (mockup *MockupSessionProvider) Update(session entities.Session) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(session.ID) {
        mockup.sessions[session.ID] = session
        return nil
    }
    return derrors.NewOperationError(errors.SessionDoesNotExist).WithParams(session)
}


// Dump obtains the list of all sessions in the system.
//   returns:
//     The list of sessions.
//     An error if the list cannot be retrieved.
func (mockup * MockupSessionProvider) Dump() ([] entities.Session, derrors.DaishoError) {
    result := make([] entities.Session, 0)
    mockup.Lock()
    defer mockup.Unlock()
    for _, node := range mockup.sessions {
        result = append(result, node)
    }
    return result, nil
}

// Clear is a util function to clear the contents of the mockup.
func (mockup *MockupSessionProvider) Clear() {
    mockup.Lock()
    mockup.sessions = make(map[string]entities.Session)
    mockup.Unlock()
}

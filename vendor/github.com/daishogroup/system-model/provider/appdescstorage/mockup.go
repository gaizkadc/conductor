//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Mockup application descriptor provider.

package appdescstorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupAppDescProvider is the mockup version of the AppDescProvider.
type MockupAppDescProvider struct {
    sync.Mutex
    // Application descriptors indexed by descriptor identifier.
    applications map[string]entities.AppDescriptor
}

// NewMockupAppDescProvider is the main builder of the MockupAppDescProvider.
func NewMockupAppDescProvider() *MockupAppDescProvider {
    return &MockupAppDescProvider{applications: make(map[string]entities.AppDescriptor)}
}

// Add a new application descriptor to the system.
//   params:
//     descriptor The application descriptor to be added
//   returns:
//     An error if the descriptor cannot be added.
func (mockup *MockupAppDescProvider) Add(descriptor entities.AppDescriptor) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(descriptor.ID) {
        mockup.applications[descriptor.ID] = descriptor
        return nil
    }
    return derrors.NewOperationError(errors.AppDescAlreadyExists).WithParams(descriptor)
}

// Exists checks if an application descriptor exists in the system.
//   params:
//     descriptorID The application descriptor identifier.
//   returns:
//     Whether the descriptor exists or not.
func (mockup *MockupAppDescProvider) Exists(descriptorID string) bool {
    mockup.Lock()
    _, exists := mockup.applications[descriptorID]
    mockup.Unlock()
    return exists
}

func (mockup *MockupAppDescProvider) unsafeExists(descriptorID string) bool {
    _, exists := mockup.applications[descriptorID]
    return exists
}

// RetrieveDescriptor retrieves a given application descriptor.
//   params:
//     descriptorID The application descriptor identifier.
//   returns:
//     The application descriptor.
//     An error if the descriptor cannot be retrieved.
func (mockup *MockupAppDescProvider) RetrieveDescriptor(descriptorID string) (*entities.AppDescriptor, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    descriptor, exists := mockup.applications[descriptorID]
    if exists {
        return &descriptor, nil
    }
    return nil, derrors.NewOperationError(errors.AppDescDoesNotExists).WithParams(descriptorID)

}

// Delete a given application descriptor.
//   params:
//     instanceID The application descriptor identifier.
//   returns:
//     An error if the application descriptor cannot be removed.
func (mockup *MockupAppDescProvider) Delete(descriptorID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    _, exists := mockup.applications[descriptorID]
    if exists {
        delete(mockup.applications, descriptorID)
        return nil
    }
    return derrors.NewOperationError(errors.AppDescDoesNotExists).WithParams(descriptorID)
}

// ReducedInfoList get a list with the reduced info.
//   returns:
//     List of the reduced app info.
//     An error if the descriptor cannot be retrieved.
func (mockup *MockupAppDescProvider) ReducedInfoList() ([] entities.AppDescriptorReducedInfo, derrors.DaishoError) {
    result := make([] entities.AppDescriptorReducedInfo, 0, len(mockup.applications))
    mockup.Lock()
    for _, appDesc := range mockup.applications {
        reducedInfo := entities.NewAppDescriptorReducedInfo(appDesc.NetworkID, appDesc.ID, appDesc.Name)
        result = append(result, *reducedInfo)
    }
    mockup.Unlock()
    return result, nil
}

// Dump obtains the list of all app descriptors in the system.
//   returns:
//     The list of AppDescriptors.
//     An error if the list cannot be retrieved.
func (mockup * MockupAppDescProvider) Dump() ([] entities.AppDescriptor, derrors.DaishoError) {
    result := make([] entities.AppDescriptor, 0)
    mockup.Lock()
    for _, app := range mockup.applications {
        result = append(result, app)
    }
    mockup.Unlock()
    return result, nil
}

// Clear is a util function to clear the contents of the mockup provider.
func (mockup *MockupAppDescProvider) Clear() {
    mockup.Lock()
    mockup.applications = make(map[string]entities.AppDescriptor)
    mockup.Unlock()
}

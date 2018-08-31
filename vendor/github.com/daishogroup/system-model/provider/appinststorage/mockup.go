//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Mockup application instance provider.

package appinststorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupAppInstProvider is the mockup version of the AppInstProvider.
type MockupAppInstProvider struct {
    sync.Mutex
    // Instances indexed by application instance identifier.
    instances map[string]entities.AppInstance
}

// NewMockupAppInstProvider creates a new instance provider.
func NewMockupAppInstProvider() *MockupAppInstProvider {
    return &MockupAppInstProvider{instances: make(map[string]entities.AppInstance)}
}

// Add a new application instance to the system.
//   params:
//     instance The application instance to be added
//   returns:
//     An error if the instance cannot be added.
func (mockup *MockupAppInstProvider) Add(instance entities.AppInstance) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(instance.DeployedID) {
        mockup.instances[instance.DeployedID] = instance
        return nil
    }
    return derrors.NewOperationError(errors.AppInstAlreadyExists).WithParams(instance)
}

// Update an instance in the system.
//   params:
//     instance The new instance information. The instance identifier will be used and cannot be modified.
//   returns:
//     An error if the instance cannot be updated.
func (mockup *MockupAppInstProvider) Update(instance entities.AppInstance) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(instance.DeployedID) {
        mockup.instances[instance.DeployedID] = instance
        return nil
    }
    return derrors.NewOperationError(errors.AppInstDoesNotExists).WithParams(instance)
}

// Exists checks if an application instance exists in the system.
//   params:
//     instanceID The application instance identifier.
//   returns:
//     Whether the instance exists or not.
func (mockup *MockupAppInstProvider) Exists(instanceID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExists(instanceID)
}

func (mockup *MockupAppInstProvider) unsafeExists(instanceID string) bool {
    _, exists := mockup.instances[instanceID]
    return exists
}

// RetrieveInstance retrieves a given application instance.
//   params:
//     instanceID The application instance identifier.
//   returns:
//     The application instance.
//     An error if the instance cannot be retrieved.
func (mockup *MockupAppInstProvider) RetrieveInstance(instanceID string) (*entities.AppInstance, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    instance, exists := mockup.instances[instanceID]
    if exists {
        return &instance, nil
    }
    return nil, derrors.NewOperationError(errors.AppInstDoesNotExists).WithParams(instanceID)
}

// Delete a given instance.
//   params:
//     instanceID The application instance identifier.
//   returns:
//     An error if the instance cannot be removed.
func (mockup *MockupAppInstProvider) Delete(instanceID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    _, exists := mockup.instances[instanceID]
    if exists {
        delete(mockup.instances, instanceID)
        return nil
    }
    return derrors.NewOperationError(errors.AppInstDoesNotExists).WithParams(instanceID)
}

// ReducedInfoList get a list with the reduced info.
//   returns:
//     List of the reduced info.
//     An error if the info cannot be retrieved.
func (mockup *MockupAppInstProvider) ReducedInfoList() ([] entities.AppInstanceReducedInfo, derrors.DaishoError) {
    result := make([] entities.AppInstanceReducedInfo, 0, len(mockup.instances))
    mockup.Lock()
    defer mockup.Unlock()
    for _, appInstance := range mockup.instances {
        reducedInfo := entities.NewAppInstanceReducedInfo(appInstance.NetworkID, appInstance.ClusterID,
            appInstance.AppDescriptorID, appInstance.DeployedID, appInstance.Name, appInstance.Description,
            appInstance.Ports, appInstance.Port, appInstance.ClusterAddress)
        result = append(result, *reducedInfo)
    }
    return result, nil
}

// Dump obtains the list of all application instances in the system.
//   returns:
//     The list of AppInstance.
//     An error if the list cannot be retrieved.
func (mockup *MockupAppInstProvider) Dump() ([] entities.AppInstance, derrors.DaishoError) {
    result := make([] entities.AppInstance, 0)
    mockup.Lock()
    defer mockup.Unlock()
    for _, app := range mockup.instances {
        result = append(result, app)
    }
    return result, nil
}

// Clear is a util function to clear the contents of the mockup provider.
func (mockup *MockupAppInstProvider) Clear() {
    mockup.Lock()
    mockup.instances = make(map[string]entities.AppInstance)
    mockup.Unlock()
}

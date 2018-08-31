//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package configstorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupClusterProvider is a mockup implementation of the config provider.
type MockupConfigProvider struct {
    sync.Mutex
    // configuration
    config * entities.Config
}

// NewMockupClusterProvider creates a new mockup provider.
func NewMockupConfigProvider() *MockupConfigProvider {
    return &MockupConfigProvider{}
}

// Store the configuration.
//   params:
//     config The Config to be stored.
//   returns:
//     An error if the config cannot be added.
func (mockup *MockupConfigProvider) Store(config entities.Config) derrors.DaishoError {
    mockup.config = entities.NewConfig(config.LogRetention)
    return nil
}

// Retrieve the current configuration.
//   returns:
//     The config.
//     An error if the config cannot be retrieved.
func (mockup *MockupConfigProvider) Get() (*entities.Config, derrors.DaishoError){
    if mockup.config != nil {
        return mockup.config, nil
    }
    return nil, derrors.NewOperationError(errors.ConfigDoesNotExists)
}

// Clear is a util function to clear the contents of the mockup.
func (mockup *MockupConfigProvider) Clear() {
    mockup.Lock()
    mockup.config = nil
    mockup.Unlock()
}

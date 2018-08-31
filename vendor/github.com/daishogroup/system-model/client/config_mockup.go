//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/configstorage"
    "github.com/daishogroup/system-model/server/config"
)

type ConfigMockup struct {
    configProvider *configstorage.MockupConfigProvider
    configManager config.Manager
}

func NewConfigMockup() Config {
    var configProvider = configstorage.NewMockupConfigProvider()
    var configManager = config.NewManager(configProvider)
    return &ConfigMockup{configProvider, configManager}
}

func (mockup *ConfigMockup) ClearMockup() {
    mockup.configProvider.Clear()
}

func (mockup *ConfigMockup) InitMockup() {
    defaultConfig := entities.NewConfig("168h")
    mockup.configProvider.Store(*defaultConfig)
}

// Store the configuration.
//   params:
//     config The Config to be stored.
//   returns:
//     An error if the config cannot be added.
func (mockup *ConfigMockup) Set(config entities.Config) derrors.DaishoError{
    return mockup.configManager.SetConfig(config)
}

// Retrieve the current configuration.
//   returns:
//     The config.
//     An error if the config cannot be retrieved.
func (mockup *ConfigMockup) Get() (*entities.Config, derrors.DaishoError){
    return mockup.configManager.GetConfig()
}

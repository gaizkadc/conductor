//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package config

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/configstorage"
)

// The Manager struct provides access to cluster related methods.
type Manager struct {
    configProvider configstorage.Provider
}

func NewManager(
    configProvider configstorage.Provider) Manager {
    return Manager{configProvider}
}

func (mgr * Manager) GetConfig() (* entities.Config, derrors.DaishoError) {
    return mgr.configProvider.Get()
}

func (mgr * Manager) SetConfig(config entities.Config) derrors.DaishoError {
    return mgr.configProvider.Store(config)
}

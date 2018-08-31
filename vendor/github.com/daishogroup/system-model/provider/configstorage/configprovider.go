//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package configstorage

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Provider is the interface of the Config provider.
type Provider interface {
    // Store the configuration.
    //   params:
    //     config The Config to be stored.
    //   returns:
    //     An error if the config cannot be added.
    Store(config entities.Config) derrors.DaishoError

    // Retrieve the current configuration.
    //   returns:
    //     The config.
    //     An error if the config cannot be retrieved.
    Get() (*entities.Config, derrors.DaishoError)

}



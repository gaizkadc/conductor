//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Config client interface.
type Config interface {

    // Set the configuration.
    //   params:
    //     config The Config to be stored.
    //   returns:
    //     An error if the config cannot be added.
    Set(config entities.Config) derrors.DaishoError

    // Retrieve the current configuration.
    //   returns:
    //     The config.
    //     An error if the config cannot be retrieved.
    Get() (*entities.Config, derrors.DaishoError)

}
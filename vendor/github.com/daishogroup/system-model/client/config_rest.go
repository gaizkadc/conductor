//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package client

import (
    "fmt"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/dhttp"
)

const ConfigSetURI = "/api/v0/config/set"
const ConfigGetURI = "/api/v0/config/get"

// Client Rest for Network resources.
type ConfigRest struct {
    client dhttp.Client
}

//Deprecated: Use NewConfigClientRest
func NewConfigRest(basePath string) Config {
    return NewConfigClientRest(ParseHostPort(basePath))
}

func NewConfigClientRest(host string, port int) Config {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host, port))
    return &ConfigRest{rest}
}

// Set the configuration.
//   params:
//     config The Config to be stored.
//   returns:
//     An error if the config cannot be added.
func (rest *ConfigRest) Set(config entities.Config) derrors.DaishoError {
    response := rest.client.Post(fmt.Sprintf(ConfigSetURI), config, new(entities.SuccessfulOperation))
    if response.Error != nil {
        return response.Error
    }
    return nil
}

// Retrieve the current configuration.
//   returns:
//     The config.
//     An error if the config cannot be retrieved.
func (rest *ConfigRest) Get() (*entities.Config, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(ConfigGetURI), new(entities.Config))
    if response.Error != nil {
        return nil, response.Error
    } else {
        n := response.Result.(*entities.Config)
        return n, nil
    }
}

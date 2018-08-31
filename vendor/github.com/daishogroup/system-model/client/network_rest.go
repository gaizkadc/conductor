//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Network Client Rest

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "fmt"
    "github.com/daishogroup/dhttp"
)

const NetworkAddURI = "/api/v0/network/add"
const NetworkListURI = "/api/v0/network/list"
const NetworkGetURI = "/api/v0/network/%s/info"
const NetworkDeleteURI = "/api/v0/network/%s/delete"

// Client Rest for Network resources.
type NetworkRest struct {
    client dhttp.Client
}
// NewNetworkRest get a Network client that uses REST protocol.
//   params:
//     basePath Full base path
//   returns:
//     The Network Client.
// Deprecated: Use NewNodeClientRest
func NewNetworkRest(basePath string) Network {
    return NewNetworkClientRest(ParseHostPort(basePath))
}

// NewNetworkClientRest get a Network client that uses REST protocol.
//   params:
//     host Host name
//     port Port
//   returns:
//     The Network Client.
func NewNetworkClientRest(host string, port int) Network {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host, port))
    return &NetworkRest{rest}
}

// Add a network to the system model.
//   params:
//     entity The cluster entity.
//   returns:
//     The added network.
//     Error, if there is an internal error.
func (rest *NetworkRest) Add(entity entities.AddNetworkRequest) (*entities.Network, derrors.DaishoError) {
    response := rest.client.Post(NetworkAddURI, entity, new(entities.Network))
    if response.Error != nil {
        return nil, response.Error
    } else {
        result := response.Result.(*entities.Network)
        return result, nil
    }
}

// List of the networks in the system model.
//   returns:
//     The list of networks.
//     Error, if there is an internal error.
func (rest *NetworkRest) List() ([] entities.Network, derrors.DaishoError) {
    response := rest.client.Get(NetworkListURI, new([] entities.Network))
    if response.Error != nil {
        return nil, response.Error
    } else {
        ns := response.Result.(*[] entities.Network)
        return *ns, nil
    }
}

// Get the info of the selected network.
//   params:
//     networkId The network id.
//   returns:
//     The selected network.
//     Error, if there is an internal error.
func (rest *NetworkRest) Get(networkID string) (*entities.Network, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(NetworkGetURI, networkID), new(entities.Network))
    if response.Error != nil {
        return nil, response.Error
    } else {
        n := response.Result.(*entities.Network)
        return n, nil
    }
}

// Delete a target network
//  params:
//      networkID The network ID
//  returns:
//      Error if any
func (rest *NetworkRest) Delete(networkID string) derrors.DaishoError {
    response := rest.client.Delete(fmt.Sprintf(NetworkDeleteURI, networkID), new(entities.SuccessfulOperation))
    if response.Error != nil {
        return response.Error
    }
    return nil
}

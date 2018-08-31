//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Network Client interface

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)
// Network is a interface that represent the client for the network entity.
type Network interface {
    // Add a network to the system model.
    //   params:
    //     entity The cluster entity.
    //   returns:
    //     The added network.
    //     Error, if there is an internal error.
    Add(entity entities.AddNetworkRequest) (* entities.Network, derrors.DaishoError)

    // List of the networks in the system model.
    //   returns:
    //     The list of networks.
    //     Error, if there is an internal error.
    List() ([] entities.Network, derrors.DaishoError)

    // Get the info of the selected network.
    //   params:
    //     networkId The network id.
    //   returns:
    //     The selected network.
    //     Error, if there is an internal error.
    Get(networkID string) (* entities.Network, derrors.DaishoError)

    // Delete a target network
    //  params:
    //      networkID The network ID
    //  returns:
    //      Error if any
    Delete(networkID string) derrors.DaishoError
}

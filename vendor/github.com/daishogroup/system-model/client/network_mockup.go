//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Network Client Mockup

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

//NetworkMockup is a mockup for the Network Resource.
type NetworkMockup struct {
    //Internal data.
    data [] entities.Network
}

// Add a network to the system model.
//   params:
//     entity The cluster entity.
//   returns:
//     The added network.
//     Error, if there is an internal error.
func (mockup *NetworkMockup) Add(entity entities.AddNetworkRequest) (* entities.Network, derrors.DaishoError) {
    n := entities.ToNetwork(entity)
    mockup.data = append(mockup.data, *n)
    return n,nil
}

// List of the networks in the system model.
//   returns:
//     The list of networks.
//     Error, if there is an internal error.
func (mockup *NetworkMockup) List() ([] entities.Network, derrors.DaishoError) {
    return mockup.data, nil
}

// Get the info of the selected network.
//   params:
//     networkId The network id.
//   returns:
//     The selected network.
//     Error, if there is an internal error.
func (mockup *NetworkMockup) Get(networkID string) (*entities.Network, derrors.DaishoError) {
    for _, network := range mockup.data {
        if network.ID == networkID {
            return &network, nil
        }
    }
    return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}


// Delete a target network
//  params:
//      networkID The network ID
//  returns:
//      Error if any
func (mockup *NetworkMockup) Delete(networkID string) derrors.DaishoError {
    newList := make([] entities.Network,0)
    removed :=false
    for _, network := range mockup.data{
        if network.ID != networkID {
            newList = append(newList,network)
        }else{
            removed=true
        }
    }
    if removed {
        mockup.data = newList
        return nil
    }

    return derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}

// NewNetworkMockup get a Mockup Network client.
//   returns:
//     The Network Client.
func NewNetworkMockup() Network {
    net1 := entities.NewNetworkWithID("1","n1","d1","a1","ap1","ae1")
    net1.EdgenetID = "beefdead00123456"
    net2 := entities.NewNetworkWithID("2","n2","d2","a2","ap2","ae2")
    return &NetworkMockup{
        data: []entities.Network{
            *net1,
            *net2,
        },
    }
}

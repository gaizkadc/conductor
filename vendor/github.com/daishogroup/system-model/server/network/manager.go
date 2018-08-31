//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the network manager in charge of the business logic behind network entities.

package network

import (

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/system-model/provider/networkstorage"
)

// Manager struct provides access to network related providers.
type Manager struct {
    networkProvider networkstorage.Provider
}

// NewManager creates a new network manager.
//   params:
//     networkProvider The network storage provider.
//   returns:
//     A manager.
func NewManager(networkProvider networkstorage.Provider) Manager {
    return Manager{networkProvider}
}

// AddNetwork adds a new network to the system.
//   params:
//     network The network to be added.
//   returns:
//     The added network.
//     An error if the network cannot be added.
func (mgr *Manager) AddNetwork(request entities.AddNetworkRequest) (*entities.Network, derrors.DaishoError) {
    network := entities.ToNetwork(request)
    if mgr.networkProvider.Exists(network.ID) {
        return nil, derrors.NewOperationError(errors.NetworkAlreadyExists)
    }

    err := mgr.networkProvider.Add(*network)
    if err == nil {
        return network, nil
    }
    return nil, err

}

// ListNetworks obtains the list of networks in the system.
//   returns:
//     An array of networks.
//     An error if the networks cannot be retrieved.
func (mgr *Manager) ListNetworks() ([]entities.Network, derrors.DaishoError) {
    return mgr.networkProvider.ListNetworks()
}

// GetNetwork retrieves a selected network.
//   params:
//     networkID The network identifier.
//   returns:
//     The selected network.
//     An error if the network cannot be retrieved.
func (mgr *Manager) GetNetwork(networkID string) (*entities.Network, derrors.DaishoError) {
    return mgr.networkProvider.RetrieveNetwork(networkID)
}

// DeleteNetwork delete a specified network.
//   params:
//     networkID The network identifier.
//   returns:
//     Error if any.
func (mgr *Manager) DeleteNetwork(networkID string) derrors.DaishoError {

    // Check no clusters are attached
    clusters, err := mgr.networkProvider.ListClusters(networkID)
    if err != nil {
        return err
    }

    if clusters != nil && len(clusters) != 0 {
        return derrors.NewOperationError(errors.InvalidCondition).WithParams("cluster")
    }

    // check no attached instances
    instances, err := mgr.networkProvider.ListAppInst(networkID)
    if err != nil {
        return err
    }
    if instances != nil && len(instances) != 0 {
        return derrors.NewOperationError(errors.InvalidCondition).WithParams("instance")
    }

    // check no attached descriptors
    descriptors, err := mgr.networkProvider.ListAppDesc(networkID)
    if err != nil {
        return err
    }
    if descriptors != nil && len(descriptors) != 0 {
        return derrors.NewOperationError(errors.InvalidCondition).WithParams("descriptor")
    }

    // everything is OK, proceed
    return mgr.networkProvider.DeleteNetwork(networkID)
}

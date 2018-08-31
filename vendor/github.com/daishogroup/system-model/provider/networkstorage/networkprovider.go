//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the network provider specification.

package networkstorage

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Provider definition of the network persistence-related methods.
type Provider interface {

    // Add a new network to the system.
    //   params:
    //     network The Network to be added
    //   returns:
    //     An error if the network cannot be added.
    Add(network entities.Network) derrors.DaishoError

    // Exists checks if a network exists in the system.
    //   params:
    //     networkID The network identifier.
    //   returns:
    //     Whether the network exists or not.
    Exists(networkID string) bool

    // RetrieveNetwork retrieves a given network.
    //   params:
    //     networkID The network identifier.
    //   returns:
    //     The network.
    //     An error if the network cannot be retrieved.
    RetrieveNetwork(networkID string) (* entities.Network, derrors.DaishoError)

    // ListNetworks retrieves all the networks in the system.
    //   returns:
    //     An array of networks.
    //     An error if the networks cannot be retrieved.
    ListNetworks() ([]entities.Network, derrors.DaishoError)

    // Delete a given network
    // returns:
    //  Error if any
    DeleteNetwork(networkID string) derrors.DaishoError

    // AttachCluster attaches a cluster to an existing network.
    //   params:
    //     networkID The network identifier.
    //     clusterID The cluster identifier.
    //   returns:
    //     An error if the cluster cannot be attached.
    AttachCluster(networkID string, clusterID string) derrors.DaishoError

    // ListClusters lists the clusters of a network.
    //   params:
    //     networkID The network identifier.
    //   returns:
    //     An array of cluster identifiers.
    //     An error if the clusters cannot be retrieved.
    ListClusters(networkID string) ([]string, derrors.DaishoError)

    // ExistsCluster checks if a cluster is associated with a given network.
    //   params:
    //     networkID The network identifier.
    //     clusterID The cluster identifier.
    //   returns:
    //     Whether the cluster is associated to the network.
    ExistsCluster(networkID string, clusterID string) bool

    // DeleteCluster deletes a cluster from an existing network.
    //   params:
    //     networkID The network identifier.
    //     clusterID The cluster identifier.
    //   returns:
    //     An error if the cluster cannot be removed.
    DeleteCluster(networkID string, clusterID string) derrors.DaishoError

// TODO Add methods related to operators
    // AttachOperator(networkId, userId)
    // RetrieveOperator(networkId)

    // TODO App management methods need to be revisited once we define how the user will upload apps and how to manage repositories.

    // RegisterAppDesc registers a new application descriptor inside a given network.
    //   params:
    //     networkID The network identifier.
    //     appDescriptorID The application descriptor identifier.
    //   returns:
    //     An error if the descriptor cannot be registered.
    RegisterAppDesc(networkID string, appDescriptorID string) derrors.DaishoError

    // ListAppDesc lists all the application descriptors in a given network.
    //   params:
    //     networkID The network identifier.
    //   returns:
    //     An array of application descriptor identifiers.
    //     An error if the list cannot be retrieved.
    ListAppDesc(networkID string) ([] string, derrors.DaishoError)

    // ExistsAppDesc checks if an application descriptor exists in a network.
    //   params:
    //     networkID The network identifier.
    //     appDescriptorID The application descriptor identifier.
    //   returns:
    //     Whether the application exists in the given network.
    ExistsAppDesc(networkID string, appDescriptorID string) bool

    // DeleteAppDescriptor deletes an application descriptor from a network.
    //   params:
    //     networkID The network identifier.
    //     appDescriptorID The application descriptor identifier.
    //   returns:
    //     An error if the application descriptor cannot be removed.
    DeleteAppDescriptor(networkID string, appDescriptorID string) derrors.DaishoError

    // RegisterAppInst registers a new application instance inside a given network.
    //   params:
    //     networkID The network identifier.
    //     appInstanceID The application instance identifier.
    //   returns:
    //     An error if the descriptor cannot be registered.
    RegisterAppInst(networkID string, appInstanceID string) derrors.DaishoError

    // ListAppInst list all the application instances in a given network.
    //   params:
    //     networkID The network identifier.
    //   returns:
    //     An array of application descriptor identifiers.
    //     An error if the list cannot be retrieved.
    ListAppInst(networkID string) ([] string, derrors.DaishoError)

    // ExistsAppInst checks if an application instance exists in a network.
    //   params:
    //     networkID The network identifier.
    //     appInstanceID The application instance identifier.
    //   returns:
    //     Whether the application exists in the given network.
    ExistsAppInst(networkID string, appInstanceID string) bool

    // DeleteAppInstance deletes an application instance from a network.
    //   params:
    //     networkID The network identifier.
    //     appInstanceID The application instance identifier.
    //   returns:
    //     An error if the application cannot be removed.
    DeleteAppInstance(networkID string, appInstanceID string) derrors.DaishoError

    // ReducedInfoList get a list with the reduced info.
    //   returns:
    //     List of the reduced info.
    //     An error if the info cannot be retrieved.
    ReducedInfoList() ([] entities.NetworkReducedInfo, derrors.DaishoError)
}
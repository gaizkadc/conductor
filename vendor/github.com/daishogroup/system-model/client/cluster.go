//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Cluster Client interface


package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)


// Cluster is an interface that represent the client for the cluster entity.
type Cluster interface {
    // Add a cluster to the network.
    //   params:
    //     networkId The network id.
    //     entity The cluster entity.
    //   returns:
    //	   The added network.
    //     Error, if there is an internal error.
    Add(networkID string, entity entities.AddClusterRequest) (*entities.Cluster, derrors.DaishoError)


    // List the cluster by network.
    //   params:
    //     networkId The network id.
    //   returns:
    //     The list of clusters for the selected network.
    //     Error, if there is an internal error.
    ListByNetwork(networkID string) ([] entities.Cluster, derrors.DaishoError)

    // Get a selected cluster
    //   params:
    //     networkId The network id.
    //     clusterId The cluster id.
    //   returns:
    //     The selected cluster.
    //     Error, if there is an internal error.
    Get(networkID string, clusterID string) (*entities.Cluster, derrors.DaishoError)

    // Update a selected cluster
    //   params:
    //     networkId The network id.
    //     clusterId The cluster id.
    //     update The update request.
    //   returns:
    //     The updated cluster.
    //     Error, if there is an internal error.
    Update(networkID string, clusterID string, update entities.UpdateClusterRequest) (*entities.Cluster,derrors.DaishoError)

    // Delete a cluster
    //  params:
    //      networkID The network id.
    //      clusterID the cluster id.
    //  returns:
    //      Error if any
    Delete(networkID string, clusterID string) derrors.DaishoError
}




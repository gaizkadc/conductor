//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Cluster Client Rest

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"

    "fmt"
    "github.com/daishogroup/dhttp"
)

// ClusterAddURI with the URI pattern to add new clusters.
const ClusterAddURI = "/api/v0/cluster/%s/add"
// ClusterListByNetworkURI with the URI pattern to list existing networks.
const ClusterListByNetworkURI = "/api/v0/cluster/%s/list"
// ClusterGetURI with the URI pattern to get the cluster information.
const ClusterGetURI = "/api/v0/cluster/%s/%s/info"
// ClusterUpdateURI with the URI pattern to update the fields of a cluster.
const ClusterUpdateURI = "/api/v0/cluster/%s/%s/update"
// ClusterDeleteURI with the URI pattern to delete a cluster.
const ClusterDeleteURI = "/api/v0/cluster/%s/%s/delete"

// ClusterRest structure with the Rest client.
type ClusterRest struct {
    client dhttp.Client
}

// Add a cluster to the network.
//   params:
//     networkId The network id.
//     entity The cluster entity.
//   returns:
//	   The added network.
//     Error, if there is an internal error.
func (rest *ClusterRest) Add(networkID string, entity entities.AddClusterRequest) (*entities.Cluster, derrors.DaishoError) {
    response := rest.client.Post(fmt.Sprintf(ClusterAddURI, networkID), entity, new(entities.Cluster))
    if response.Error != nil {
        return nil, response.Error
    }
    result := response.Result.(*entities.Cluster)
    return result, nil
}

// ListByNetwork obtains the list of clusters by network.
//   params:
//     networkId The network id.
//   returns:
//     The list of clusters for the selected network.
//     Error, if there is an internal error.
func (rest *ClusterRest) ListByNetwork(networkID string) ([] entities.Cluster, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(ClusterListByNetworkURI, networkID), new([] entities.Cluster))
    if response.Error != nil {
        return nil, response.Error
    }
    cs := response.Result.(*[] entities.Cluster)
    return *cs, nil
}

// Get a selected cluster
//   params:
//     networkId The network id.
//     clusterId The cluster id.
//   returns:
//     The selected cluster.
//     Error, if there is an internal error.
func (rest *ClusterRest) Get(networkID string, clusterID string) (*entities.Cluster, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(ClusterGetURI, networkID, clusterID), new(entities.Cluster))
    if response.Error != nil {
        return nil, response.Error
    }
    c := response.Result.(*entities.Cluster)
    return c, nil
}

// Update a selected cluster
//   params:
//     networkId The network id.
//     clusterId The cluster id.
//     update The update request.
//   returns:
//     The updated cluster.
//     Error, if there is an internal error.
func (rest *ClusterRest) Update(networkID string, clusterID string,
    update entities.UpdateClusterRequest) (*entities.Cluster, derrors.DaishoError) {
    response := rest.client.Post(fmt.Sprintf(ClusterUpdateURI, networkID, clusterID), update, new(entities.Cluster))
    if response.Error != nil {
        return nil, response.Error
    }
    c := response.Result.(*entities.Cluster)
    return c, nil
}

// Delete a cluster
//  params:
//      networkID The network id.
//      clusterID the cluster id.
//  returns:
//      Error if any
func (rest *ClusterRest) Delete(networkID string, clusterID string) derrors.DaishoError{
    response := rest.client.Delete(fmt.Sprintf(ClusterDeleteURI, networkID, clusterID), new(entities.SuccessfulOperation))
    if response.Error != nil {
        return response.Error
    }
    return nil
}
// Deprecated: Use NewClusterClientRest
func NewClusterRest(basePath string) Cluster {
    return NewClusterClientRest(ParseHostPort(basePath))
}

// NewClusterRest creates a Cluster client that uses REST protocol.
//   params:
//     basePath Full bash path
//   returns:
//     The Cluster Client.
func NewClusterClientRest(host string, port int) Cluster {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host,port))
    return &ClusterRest{rest}
}

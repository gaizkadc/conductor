//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Cluster Client TestSuite

package client

import (
    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/entities"
)

// TestGetExistingCluster checks that existing cluster can be retrieved.
func TestGetExistingCluster(suite *suite.Suite, cluster Cluster) {
    c, err := cluster.Get("1", "1")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(c, "cluster must not nil")
}

// TestGetNotExistingCluster checks the result of retrieving a non existing cluster.
func TestGetNotExistingCluster(suite *suite.Suite, cluster Cluster) {
    c, err := cluster.Get("1", "3")
    suite.NotNil(err, "error must not be Nil")
    suite.Nil(c, "cluster must be nil")
}

// TestGetClusterNotExistingNetwork checks the result of retriving a cluster from a non
// existing network.
func TestGetClusterNotExistingNetwork(suite *suite.Suite, cluster Cluster) {
    c, err := cluster.Get("3", "1")
    suite.NotNil(err, "error must not be Nil")
    suite.Nil(c, "cluster must be nil")
}

// TestAddCluster checks that clusters can be added and retrieved.
func TestAddCluster(suite *suite.Suite, cluster Cluster) {
    addedCluster, err := cluster.Add("1", entities.AddClusterRequest{
        Name:        "Cluster4",
        Description: "Description Cluster 4",
        Type:        entities.GatewayType,
        Location:    "Madrid",
    })
    suite.Nil(err, "error must be Nil")
    suite.NotNil(addedCluster, "added cluster must not be Nil")

    c, err := cluster.Get("1", addedCluster.ID)
    suite.Nil(err, "error must be Nil")
    suite.NotNil(c, "cluster must not nil")
}

// TestAddClusterToNotExistingNetwork checks the result of adding a cluster to a non
// existing network.
func TestAddClusterToNotExistingNetwork(suite *suite.Suite, cluster Cluster) {
    c, err := cluster.Add("3", entities.AddClusterRequest{
        Name:        "Cluster5",
        Description: "Description Cluster 5",
        Type:        entities.GatewayType,
        Location:    "Madrid",
    })
    suite.NotNil(err, "error must be not Nil")
    suite.Nil(c, "cluster must be Nil")
}

// TestListByNetwork checks that cluster can be listed.
func TestListByNetwork(suite *suite.Suite, cluster Cluster) {
    cs, err := cluster.ListByNetwork("1")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(cs, "cluster must be not nil")
    suite.Len(cs, 2, "len must be 2")
}

// TestListByNotExistingNetwork checks the result of listing clusters from a non
// existing network.
func TestListByNotExistingNetwork(suite *suite.Suite, cluster Cluster) {
    cs, err := cluster.ListByNetwork("3")
    suite.NotNil(err, "error must not be Nil")
    suite.Nil(cs, "cluster must be nil")
}

// TestUpdateCluster checks that clusters can be updated.
func TestUpdateCluster(suite *suite.Suite, cluster Cluster) {
    addedCluster, err := cluster.Add("1", entities.AddClusterRequest{
        Name:        "Cluster7",
        Description: "Description Cluster 7",
        Type:        entities.GatewayType,
        Location:    "Madrid",
    })
    suite.Nil(err, "error must be Nil")
    suite.NotNil(addedCluster, "added cluster must not be Nil")

    updateRequest := entities.NewUpdateClusterRequest().WithName("newUpdate")
    _, errUpdate := cluster.Update("1", addedCluster.ID, *updateRequest)
    suite.Nil(errUpdate, "errUpdate must be Nil")

    c, err := cluster.Get("1", addedCluster.ID)
    suite.Nil(err, "error must be Nil")
    suite.NotNil(c, "cluster must not be nil")
    suite.Equal("newUpdate", c.Name, "name must be updated")
}

// TestDeleteCluster checks that clusters can be deleted.
func TestDeleteCluster(suite *suite.Suite, cluster Cluster) {

    initialClusters, err := cluster.ListByNetwork("1")
    suite.Nil(err, "List should be retrieved")

    addedCluster, err := cluster.Add("1", *entities.NewAddClusterRequest(
        "Cluster8", "Description cluster 8", entities.GatewayType, "madrid", "email"))

    suite.Nil(err, "error must be Nil")
    suite.NotNil(addedCluster, "added cluster must not be Nil")

    addedCluster, err = cluster.Update("1", addedCluster.ID,
        *entities.NewUpdateClusterRequest().WithClusterStatus(entities.ClusterInstalled))
    suite.Nil(err, "error updated")

    // delete it
    err = cluster.Delete("1", addedCluster.ID)
    suite.Nil(err,"the cluster was not correctly removed")

    // try to get it
    retrieved, err := cluster.Get("1",addedCluster.ID)
    suite.Nil(retrieved,"something was returned after delete")
    suite.NotNil(err, "a error should have been generated")

    clusters, err := cluster.ListByNetwork("1")
    suite.Nil(err, "List should be retrieved")
    suite.Equal(len(initialClusters), len(clusters), "Expecting 0 clusters after remove")

}
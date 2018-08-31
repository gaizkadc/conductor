//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Tests on the Dump structure.

package entities

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func GetTestCluster(networkID string, clusterID string) * Cluster {
    return NewClusterWithID(
        networkID, clusterID, "Name", "Description",
        EdgeType, "Madrid, Spain", "admin@admins.com",
        ClusterInstalled, false, false)
}

func TestDump_AddClusters(t *testing.T) {
    clusters := make([] Cluster, 0)
    clusters = append(clusters, * GetTestCluster("n1", "c1"))
    clusters = append(clusters, * GetTestCluster("n1", "c2"))
    newClusters := make([] Cluster, 0)
    newClusters = append(newClusters, * GetTestCluster("n2", "c1"))
    newClusters = append(newClusters, * GetTestCluster("n2", "c2"))

    dump := NewDump(make([] Network, 0), clusters, make([] Node, 0), make([] AppDescriptor, 0), make([] AppInstance, 0),
        make([] User, 0), make([] UserAccess, 0))
    assert.Equal(t, 2, len(dump.Clusters))
    dump.AddClusters(newClusters)
    assert.Equal(t, 4, len(dump.Clusters))
}

func TestDump_AddCluster(t *testing.T) {
    dump := NewDumpWithNetworks(make([] Network, 0))
    dump.AddCluster(* GetTestCluster("n1", "c1"))
    assert.Equal(t, 1, len(dump.Clusters))
    dump.AddCluster(* GetTestCluster("n1", "c2"))
    assert.Equal(t, 2, len(dump.Clusters))
}

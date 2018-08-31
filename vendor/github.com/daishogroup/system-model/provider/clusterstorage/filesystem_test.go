//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package clusterstorage

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
    "testing"
    "io/ioutil"
    "os"
)

func TestFileSystemProvider_Clusters(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testCluster")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    toAdd1 := entities.NewCluster("n1", "c1", "d1", entities.EdgeType, "l",
        "email1", entities.ClusterCreated, false, false)
    toAdd2 := entities.NewCluster("n2", "c1", "d1", entities.EdgeType, "l",
        "email1", entities.ClusterCreated, false, false)
    err := provider.Add(* toAdd1)
    assert.Nil(t, err, "cluster should be added")
    err = provider.Add(* toAdd2)
    assert.Nil(t, err, "cluster should be added")

    r1, err := provider.RetrieveCluster(toAdd1.ID)
    assert.Nil(t, err, "cluster should be retrieved")
    assert.EqualValues(t, toAdd1, r1, "cluster should match")
    r2, err := provider.RetrieveCluster(toAdd2.ID)
    assert.Nil(t, err, "cluster should be retrieved")
    assert.EqualValues(t, toAdd2, r2, "cluster should match")

    clusters, err := provider.Dump()
    assert.Nil(t, err, "cluster list should be retrieved")
    assert.Equal(t, 2, len(clusters))

}

func TestFileSystemProvider_Nodes(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testCluster")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    cluster := entities.NewCluster("n1", "c1", "d1", entities.EdgeType, "l",
        "email1", entities.ClusterCreated, false, false)
    err := provider.Add(* cluster)
    assert.Nil(t, err, "cluster should be added")

    nodes, err := provider.ListNodes(cluster.ID)
    assert.Nil(t, err, "nodes should be retrieved")
    assert.Equal(t, 0, len(nodes))

    n1 := entities.NewNode("n1", cluster.ID,
        "n1", "d1", make([]string, 0),
        "ip", "ip",
        false, "u1", "p1", "s1")
    n2 := entities.NewNode("n1", cluster.ID,
        "n1", "d1", make([]string, 0),
        "ip", "ip",
        false, "u1", "p1", "s1")

    err = provider.AttachNode(cluster.ID, n1.ID)
    assert.Nil(t, err, "node should be attached")
    err = provider.AttachNode(cluster.ID, n2.ID)
    assert.Nil(t, err, "nodes should be attached")
    nodes, err = provider.ListNodes(cluster.ID)
    assert.Nil(t, err, "nodes should be retrieved")
    assert.Equal(t, 2, len(nodes))
}

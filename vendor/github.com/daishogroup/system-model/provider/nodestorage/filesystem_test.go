//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package nodestorage

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
    "testing"
    "os"
    "io/ioutil"
)

func TestFileSystemProvider_Nodes(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testNodes")
    defer os.RemoveAll(dir)
    labels := make([]string, 0)
    labels = append(labels, "l1")
    provider := NewFileSystemProvider(dir)
    toAdd1 := entities.NewNode("n1", "n1", "n1", "d1", labels, "ip", "ip",
        false, "u1", "p1", "s1")
    toAdd2 := entities.NewNode("n1", "n1", "n1", "d1", labels, "ip", "ip",
        false, "u1", "p1", "s1")
    err := provider.Add(* toAdd1)
    assert.Nil(t, err, "node should be added")
    err = provider.Add(* toAdd2)
    assert.Nil(t, err, "node should be added")

    r1, err := provider.RetrieveNode(toAdd1.ID)
    assert.Nil(t, err, "node should be retrieved")
    assert.EqualValues(t, toAdd1, r1, "node should match")
    r2, err := provider.RetrieveNode(toAdd2.ID)
    assert.Nil(t, err, "node should be retrieved")
    assert.EqualValues(t, toAdd2, r2, "node should match")

    nodes, err := provider.Dump()
    assert.Nil(t, err, "node list should be retrieved")
    assert.Equal(t, 2, len(nodes))

}


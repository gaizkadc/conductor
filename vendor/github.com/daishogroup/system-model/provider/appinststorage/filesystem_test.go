//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package appinststorage

import (
    "io/ioutil"
    "os"
    "testing"

    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
)

func TestFileSystemProvider_AppInst(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testAppInst")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    toAdd1 := entities.NewAppInstance("n1", "d1", "c1", "n1", "d1", "l1",
        "arg","1Gb", entities.AppStorageDefault, make([]entities.ApplicationPort, 0), 0, "0.0.0.0")
    toAdd2 := entities.NewAppInstance("n1", "d2", "c2", "n2", "d2", "l2",
        "arg","1Gb", entities.AppStorageDefault, make([]entities.ApplicationPort, 0), 0, "0.0.0.0")

    err := provider.Add(* toAdd1)
    assert.Nil(t, err, "instance should be added")
    err = provider.Add(* toAdd2)
    assert.Nil(t, err, "instance should be added")

    r1, err := provider.RetrieveInstance(toAdd1.DeployedID)
    assert.Nil(t, err, "instance should be retrieved")
    assert.EqualValues(t, toAdd1, r1, "instance should match")
    r2, err := provider.RetrieveInstance(toAdd2.DeployedID)
    assert.Nil(t, err, "instance should be retrieved")
    assert.EqualValues(t, toAdd2, r2, "instance should match")

    instances, err := provider.Dump()
    assert.Nil(t, err, "instance list should be retrieved")
    assert.Equal(t, 2, len(instances))

}
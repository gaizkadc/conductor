//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package appdescstorage

import (
    "io/ioutil"
    "os"
    "testing"

    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
)

func TestFileSystemProvider_AppDesc(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testAppDesc")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    toAdd1 := entities.NewAppDescriptor("n1", "n1", "d1","s1","1.0",
        "l",0, []string{"repo1:tag1"})
    toAdd2 := entities.NewAppDescriptor("n1", "n2", "d2","s2","1.0",
        "l",0, []string{"repo1:tag1"})


    err := provider.Add(* toAdd1)
    assert.Nil(t, err, "descriptor should be added")
    err = provider.Add(* toAdd2)
    assert.Nil(t, err, "descriptor should be added")

    r1, err := provider.RetrieveDescriptor(toAdd1.ID)
    assert.Nil(t, err, "descriptor should be retrieved")
    assert.EqualValues(t, toAdd1, r1, "descriptor should match")
    r2, err := provider.RetrieveDescriptor(toAdd2.ID)
    assert.Nil(t, err, "descriptor should be retrieved")
    assert.EqualValues(t, toAdd2, r2, "descriptor should match")

    descriptors, err := provider.Dump()
    assert.Nil(t, err, "descriptor list should be retrieved")
    assert.Equal(t, 2, len(descriptors))

}
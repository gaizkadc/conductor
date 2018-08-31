//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Test the file system provider for users.

package accessstorage

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
    "testing"
    "os"
    "io/ioutil"
)

func TestFileSystemProvider_Users(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testUsers")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    user1 := entities.NewUserAccess("user1", []entities.RoleType{entities.GlobalAdmin})
    user2 := entities.NewUserAccess("user2", []entities.RoleType{entities.OperatorType})

    err := provider.Add(*user1)
    assert.Nil(t, err, "user access should be added")
    err = provider.Add(*user2)
    assert.Nil(t, err, "user access should be added")

    r1, err := provider.RetrieveAccess(user1.UserID)
    assert.Nil(t, err, "access should be retrieved")
    assert.EqualValues(t, user1, r1, "node should match")
    r2, err := provider.RetrieveAccess(user2.UserID)
    assert.Nil(t, err, "access should be retrieved")
    assert.EqualValues(t, user2, r2, "node should match")

    accesses, err := provider.Dump()
    assert.Nil(t, err, "user access list should be retrieved")
    assert.Equal(t, 2, len(accesses))

}


//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Test the file system provider for users.

package userstorage

import (
    "time"
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
    "testing"
    "os"
    "io/ioutil"
)

func TestFileSystemProvider_Users(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testUsers")
    defer os.RemoveAll(dir)

    creationTime := time.Now()
    expirationTime := creationTime.Add(time.Hour)

    provider := NewFileSystemProvider(dir)
    user1 := entities.NewUser("user1", "phone1", "email1", creationTime, expirationTime)
    user2 := entities.NewUser("user2", "phone2", "email2", creationTime, expirationTime)

    err := provider.Add(* user1)
    assert.Nil(t, err, "user should be added")
    err = provider.Add(* user2)
    assert.Nil(t, err, "user should be added")

    r1, err := provider.RetrieveUser(user1.ID)
    assert.Nil(t, err, "user should be retrieved")
    assert.Equal(t, user1.String(), r1.String(), "node should match")
    r2, err := provider.RetrieveUser(user2.ID)
    assert.Nil(t, err, "user should be retrieved")
    assert.Equal(t, user2.String(), r2.String(), "node should match")

    nodes, err := provider.Dump()
    assert.Nil(t, err, "user list should be retrieved")
    assert.Equal(t, 2, len(nodes))

}


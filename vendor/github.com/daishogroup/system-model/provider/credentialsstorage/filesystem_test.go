//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Test the file system provider for credentials.

package credentialsstorage

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
    "testing"
    "os"
    "io/ioutil"
)

func TestFileSystemProvider_UsersCredentials(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testUserCredentials")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)

credentials1 := entities.NewCredentials("user1","publickey1","privatekey2", "desc", "type")
    credentials2 := entities.NewCredentials("user2","publickey2","privatekey2", "desc", "type")


    err := provider.Add(*credentials1)
    assert.Nil(t, err, "credentials1 should be added")
    err = provider.Add(* credentials2)
    assert.Nil(t, err, "credentials2 should be added")

    r1, err := provider.Retrieve(credentials1.UUID)
    assert.Nil(t, err, "user should be retrieved")
    assert.EqualValues(t, credentials1, r1, "node should match")
    r2, err := provider.Retrieve(credentials2.UUID)
    assert.Nil(t, err, "user should be retrieved")
    assert.EqualValues(t, credentials2, r2, "node should match")

    // remove r1
    err = provider.Delete(credentials1.UUID)
    assert.Nil(t, err, "credentials were not removed")
    // try to retrieve it
    r3, err := provider.Retrieve(credentials1.UUID)
    assert.Nil(t, r3, "a nil credential should be returned")
    assert.NotNil(t, err, "an error must be returned")

}


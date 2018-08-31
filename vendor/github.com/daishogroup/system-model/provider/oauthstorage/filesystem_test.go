//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Test the file system provider for users.

package oauthstorage

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
    "testing"
    "os"
    "io/ioutil"
)

func TestFileSystemProvider_OAuthSecrets(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testOAuthSecrets")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    secrets := entities.NewOAuthSecrets("user1")
    secrets.AddEntry("app1", entities.OAuthEntry{"clientid","secret"})
    secrets.AddEntry("app2", entities.OAuthEntry{"clientid2", "secret2"})

    // Store it
    err := provider.Add(secrets)
    assert.Nil(t, err, "secrets should be added")
    // Try to store it again
    err = provider.Add(secrets)
    assert.NotNil(t, err, "secrets should not be added")

    // Try to retrieve it
    r1, err := provider.Retrieve(secrets.UserID)
    assert.Nil(t, err, "secrets should be retrieved")
    assert.EqualValues(t, "user1", r1.UserID, "userId should match")
    assert.Equal(t, 2, len(r1.Entries), "number of entries does not match")

    // Add a second entry
    secrets2 := entities.NewOAuthSecrets("user2")
    secrets2.AddEntry("app1", entities.OAuthEntry{"clientid12","secret12"})
    err = provider.Add(secrets2)
    assert.Nil(t, err, "secrets should be added")
    // Try to retrieve it
    r2, err := provider.Retrieve(secrets2.UserID)
    assert.Equal(t, 1, len(r2.Entries), "number of entries does not match")
    assert.Equal(t, "clientid12", r2.Entries["app1"].ClientID, "unexpected client id")
    assert.Equal(t, "secret12", r2.Entries["app1"].Secret, "unexpected secret")

    // Dump it
    everything, err := provider.Dump()
    assert.Nil(t, err, "secrets list should be retrieved")
    assert.Equal(t, 2, len(everything))

    // Remove an entry
    err = provider.Delete("user2")
    assert.Nil(t, err, "unexpected problem removing entry")
    // Try to remove again
    err = provider.Delete("user2")
    assert.NotNil(t, err, "expected problem removing entry")




}


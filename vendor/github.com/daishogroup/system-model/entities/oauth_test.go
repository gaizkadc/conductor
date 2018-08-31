//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Set of tests for OAuth.

package entities

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestAdditionOAuth(t *testing.T) {
    entrySet := NewOAuthSecrets("testUser")
    assert.Equal(t, 0, len(entrySet.Entries), "unexpected number of initial entries")
    // Add something
    entryA := OAuthEntry{ClientID:"client1",Secret:"secret1"}
    err:= entrySet.AddEntry("app1", entryA)
    assert.Nil(t, err, "unexpected error adding first entry")
    assert.Equal(t, 1, len(entrySet.Entries), "unexpected number entries")
    // try to reinsert
    err= entrySet.AddEntry("app1", entryA)
    assert.NotNil(t, err, "an error was expected")
    // Add another one
    entryB := OAuthEntry{ClientID:"client2",Secret:"secret2"}
    err= entrySet.AddEntry("app2", entryB)
    assert.Nil(t, err, "unexpected error adding second entry")
    assert.Equal(t, 2, len(entrySet.Entries), "unexpected number entries")
}

func TestDeletionOAuth(t *testing.T) {
    entrySet := NewOAuthSecrets("testUser")
    assert.Equal(t, 0, len(entrySet.Entries), "unexpected number of initial entries")
    // Add something
    entryA := OAuthEntry{ClientID:"client1",Secret:"secret1"}
    err:= entrySet.AddEntry("app1", entryA)
    assert.Nil(t, err, "unexpected error adding first entry")
    assert.Equal(t, 1, len(entrySet.Entries), "unexpected number entries")
    // Remove it
    err = entrySet.DeleteEntry("app1")
    assert.Nil(t, err, "unexpected error when deleting",err)
    // Try to remove it again
    err = entrySet.DeleteEntry("app1")
    assert.NotNil(t, err, "expected error when deleting")
}
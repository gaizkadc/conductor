//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Test the file system provider for users.

package passwordstorage

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
    "testing"
    "os"
    "io/ioutil"
)

func TestFileSystemProvider_Users(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testPasswords")
    defer os.RemoveAll(dir)

    p1 := "password1"
    p2 := "password2"

    provider := NewFileSystemProvider(dir)
    pass1,err  := entities.NewPassword("user1",&p1)
    assert.Nil(t, err, "unexpected error")
    pass2,err  := entities.NewPassword("user2",&p2)
    assert.Nil(t, err, "unexpected error")

    err = provider.Add(* pass1)
    assert.Nil(t, err, "password should be added")
    err = provider.Add(* pass2)
    assert.Nil(t, err, "password should be added")

    r1, err := provider.RetrievePassword(pass1.UserID)
    assert.Nil(t, err, "user should be retrieved")
    assert.EqualValues(t, "user1", r1.UserID, "userId should match")
    assert.True(t, r1.CompareWith(p1), "password does not match")
    r2, err := provider.RetrievePassword(pass2.UserID)
    assert.Nil(t, err, "user should be retrieved")
    assert.EqualValues(t, "user2", r2.UserID, "userId should match")
    assert.True(t, r2.CompareWith(p2), "password does not match")

    // Check failing password
    assert.False(t, r1.CompareWith("error"),"unexpected matching passwords")

    passwords, err := provider.Dump()
    assert.Nil(t, err, "passwords list should be retrieved")
    assert.Equal(t, 2, len(passwords))

}


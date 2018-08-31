//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Test the file system provider for users.

package sessionstorage

import (
    "time"
    "net/http"
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
    "testing"
    "os"
    "io/ioutil"
)

func TestFileSystemProvider_Sessions(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testSessions")
    defer os.RemoveAll(dir)

    expiration := time.Now()

    provider := NewFileSystemProvider(dir)
    session1 := entities.NewSession("user1",expiration)
    cookie1 := http.Cookie{
        Domain:"localhost",
    }
    session1.AddCookie("testCookie", cookie1)

    session2 := entities.NewSession("user2",expiration)
    cookie2 := http.Cookie{
        Domain: "localhost",
        Path: "/something/something",
    }
    session1.AddCookie("testCookie", cookie2)


    err := provider.Add(* session1)
    assert.Nil(t, err, "session should be added")
    err = provider.Add(* session2)
    assert.Nil(t, err, "session should be added")

    r1, err := provider.Retrieve(session1.ID)
    assert.Nil(t, err, "session should be retrieved")
    assert.Equal(t, r1.ExpirationDate.UTC(), expiration.UTC())

    r2, err := provider.Retrieve(session2.ID)
    assert.Nil(t, err, "user should be retrieved")
    assert.Equal(t, r2.ExpirationDate.UTC(), expiration.UTC())

    nodes, err := provider.Dump()
    assert.Nil(t, err, "user list should be retrieved")
    assert.Equal(t, 2, len(nodes))

}


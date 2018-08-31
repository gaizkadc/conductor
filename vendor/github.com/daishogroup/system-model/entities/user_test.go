//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package entities

import (
    "github.com/stretchr/testify/assert"
    "testing"
    "time"

)

func TestNewUserWithId(t *testing.T) {
    id := "u1"
    name := "The Admin"
    phone := "1234 1234 1234"
    email := "admin@admins.com"
    creation := time.Date(2010, time.January,1,1,1,1,0, time.UTC)
    expiration := creation.Add(time.Hour)
    user := NewUserWithID(id, name, phone, email, creation, expiration)

    assert.NotNil(t, user, "User should be defined")
    assert.Equal(t, user.CreationTime.Time, creation)
    //expected := "&entities.User{ID:\"u1\", Name:\"The Admin\", Phone:\"1234 1234 1234\", Email:\"admin@admins.com\"}"
    //assert.Equal(t, expected, user.String())

}

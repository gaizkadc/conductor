//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
// User password testing.
//

package entities

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestMatchingPassword(t *testing.T) {
    pass := "first"
    startingPassword,err := NewPassword("u1", &pass)
    assert.Nil(t, err, "unexpected error creating password")
    assert.True(t, startingPassword.CompareWith("first"), "unexpected difference")
    assert.False(t, startingPassword.CompareWith("non-first"), "unexpected difference")

    emptyPassword, err := NewPassword("u1",nil)
    assert.Nil(t, err, "unexpected error creating password")
    assert.False(t, emptyPassword.CompareWith("first"), "unexpected difference")

}
//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package entities

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestValidDuration(t *testing.T) {
    config := NewConfig("24h")
    assert.True(t, config.Valid())
}

func TestInvalidDuration(t *testing.T){
    config := NewConfig("1w")
    assert.False(t, config.Valid())
}
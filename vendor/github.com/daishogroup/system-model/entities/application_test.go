//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package entities

import (
    "encoding/json"
    "testing"

    "github.com/stretchr/testify/assert"
)

// DP-868 json.Unmarshal fails to recognize enums, it treats them as they inner type.
func TestInvalidUpdateAppInstUnmarshal(t *testing.T) {
    // Type is not valid
    payload1 := []byte(`{"status":"NotFound"}`)
    c1 := UpdateAppInstanceRequest{}
    err:= json.Unmarshal(payload1, &c1)
    assert.Nil(t, err, "Should return error, but we know it does not")
    assert.False(t, c1.IsValid())
}
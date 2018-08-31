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

func TestUpdateNodeRequest(t *testing.T) {
    update := NewUpdateNodeRequest()
    assert.Empty(t, update.Name, "Name must be null")
    update.WithName("newName")
    assert.NotNil(t, update.Name, "Name must be modified")
    assert.Equal(t, "newName", *update.Name)
}

func TestFluentUpdateNodeRequest(t *testing.T) {
    update := NewUpdateNodeRequest().WithName("newName").
        WithStatus(NodeReadyToInstall).WithUsername("newUser")
    assert.Equal(t, "newName", * update.Name)
    assert.Equal(t, "newUser", * update.Username)
    assert.Equal(t, NodeReadyToInstall, * update.Status)
}

func TestMerge(t *testing.T) {
    node := NewNode("n1", "c1",
        "oldName", "oldDesc", make([]string, 0),
        "0.0.0.0", "0.0.0.0", false,
        "oldUser", "oldPass", "oldKey")
    newLabels := make([]string, 0)
    newLabels = append(newLabels, "newLabel")
    update := NewUpdateNodeRequest().WithName("newName").WithDescription("newDesc").
        WithPublicIP("1.1.1.1").WithPrivateIP("1.1.1.1").WithInstalled(true).
        WithUsername("newUser").WithPassword("newPass").WithSSHKey("newKey").
        WithStatus(NodeInstalling).WithEdgenet("deadbeef00", "1.0.0.1").WithLabels(newLabels)
    updatedNode := node.Merge(* update)
    assert.Equal(t, "newName", updatedNode.Name)
    assert.Equal(t, "newDesc", updatedNode.Description)
    assert.Equal(t, "1.1.1.1", updatedNode.PublicIP)
    assert.Equal(t, "1.1.1.1", updatedNode.PrivateIP)
    assert.True(t, updatedNode.Installed)
    assert.Equal(t, "newUser", updatedNode.Username)
    assert.Equal(t, "newPass", updatedNode.Password)
    assert.Equal(t, "newKey", updatedNode.SSHKey)
    assert.Equal(t, NodeInstalling, updatedNode.Status)
    assert.Equal(t, "deadbeef00", updatedNode.EdgenetAddress)
    assert.Equal(t, "1.0.0.1", updatedNode.EdgenetIP)
    assert.Equal(t, newLabels, updatedNode.Labels)
}

// DP-868 json.Unmarshal fails to recognize enums, it treats them as they inner type.
func TestInvalidUpdateNodeUnmarshal(t *testing.T) {
    // Type is not valid
    payload1 := []byte(`{"status":"NotFound"}`)
    c1 := UpdateNodeRequest{}
    err:= json.Unmarshal(payload1, &c1)
    assert.Nil(t, err, "Should return error, but we know it does not")
    assert.False(t, c1.IsValid())
}

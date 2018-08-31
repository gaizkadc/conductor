//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package entities

import (
    "encoding/json"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestUpdateClusterRequest(t *testing.T) {
    update := NewUpdateClusterRequest()
    assert.Nil(t, update.Name, "Name must be null")
    assert.Nil(t, update.Description)
    assert.Nil(t, update.Type)
    assert.Nil(t, update.Location)
    assert.Nil(t, update.Email)
    assert.Nil(t, update.Status)
    assert.Nil(t, update.Drain)
    assert.Nil(t, update.Cordon)
    update.WithName("newName")
    assert.NotNil(t, update.Name, "Name must be modified")
    assert.Equal(t, "newName", *update.Name)
    assert.Nil(t, update.Description)
    assert.Nil(t, update.Type)
    assert.Nil(t, update.Location)
    assert.Nil(t, update.Email)
    assert.Nil(t, update.Status)
    assert.Nil(t, update.Drain)
    assert.Nil(t, update.Cordon)
}

func TestFluentUpdateClusterRequest(t *testing.T) {
    update := NewUpdateClusterRequest().WithName("newName").WithType(CloudType).WithCordon(true)
    assert.Equal(t, "newName", * update.Name)
    assert.Equal(t, CloudType, * update.Type)
    assert.Equal(t, true, * update.Cordon)
}

func TestClusterMerge(t *testing.T) {
    cluster := NewClusterWithID("1", "1", "oldName", "oldDescription",
        CloudType, "oldLocation", "oldMail", ClusterCreated,
        false, false)
    update := NewUpdateClusterRequest().WithName("newName").WithDescription("newDescription").
        WithType(GatewayType).WithLocation("newLocation").WithEmail("newEmail").
        WithClusterStatus(ClusterInstalling).WithDrain(true).WithCordon(true)
    mergedCluster := cluster.Merge(* update)
    assert.Equal(t, "newName", mergedCluster.Name)
    assert.Equal(t, "newDescription", mergedCluster.Description)
    assert.Equal(t, GatewayType, mergedCluster.Type)
    assert.Equal(t, "newLocation", mergedCluster.Location)
    assert.Equal(t, "newEmail", mergedCluster.Email)
    assert.Equal(t, ClusterInstalling, mergedCluster.Status)
    assert.True(t, mergedCluster.Drain)
    assert.True(t, mergedCluster.Cordon)

}

func TestNewClusterWithId(t *testing.T) {
    networkID := "n1"
    id := "c1"
    name := "Test cluster"
    description := "Cluster description"
    clusterType := CloudType
    location := "California"
    mail := "admin@admins.com"
    clusterStatus := ClusterCreated
    drain := false
    cordon := false
    cluster := NewClusterWithID(networkID, id, name, description, clusterType, location,
        mail, clusterStatus, drain, cordon)
    assert.NotNil(t, cluster, "Cluster should be defined")
    expected := "&entities.Cluster{NetworkID:\"n1\", ID:\"c1\", Name:\"Test cluster\", " +
        "Description:\"Cluster description\", " +
        "Type:\"cloud\", Location:\"California\", Email:\"admin@admins.com\", " +
        "Status:\"Created\", Drain:false, Cordon:false}"
    assert.Equal(t, expected, cluster.String())
}

// DP-868 json.Unmarshal fails to recognize enums, it treats them as they inner type.
func TestInvalidAddClusterRequestUnmarshal(t *testing.T) {
    // This should fail as both enums are missing an required fields are missing
    payload1 := []byte(`{"num":6.13,"strs":["a","b"]}`)
    c1 := AddClusterRequest{}
    err:= json.Unmarshal(payload1, &c1)
    assert.Nil(t, err, "Should return error, but we know it does not")
    assert.False(t, c1.IsValid())

    // Type is not valid
    payload2 := []byte(`{"name":"n1","description":"desc1","type":"NotFound","location":"Location","email":"mail@mail.com"}`)
    c2 := AddClusterRequest{}
    err = json.Unmarshal(payload2, &c2)
    assert.Nil(t, err, "Should return error, but we know it does not")
    assert.False(t, c2.IsValid())

    // Type is missing
    payload3 := []byte(`{"name":"n1","description":"desc1","location":"Location","email":"mail@mail.com"}`)
    c3 := AddClusterRequest{}
    err = json.Unmarshal(payload3, &c3)
    assert.Nil(t, err, "Should return error, but we know it does not")
    assert.False(t, c3.IsValid())

}

// DP-868 json.Unmarshal fails to recognize enums, it treats them as they inner type.
func TestInvalidUpdateClusterUnmarshal(t *testing.T) {
    // Type is not valid
    payload1 := []byte(`{"type":"NotFound"}`)
    c1 := UpdateClusterRequest{}
    err:= json.Unmarshal(payload1, &c1)
    assert.Nil(t, err, "Should return error, but we know it does not")
    assert.False(t, c1.IsValid())

    // Status is not valid
    payload2 := []byte(`{"status":"NotFound"}`)
    c2 := UpdateClusterRequest{}
    err = json.Unmarshal(payload2, &c2)
    assert.Nil(t, err, "Should return error, but we know it does not")
    assert.False(t, c2.IsValid())

    // Both are invalid
    payload3 := []byte(`{"type":"NotFound","status":"NotFound"}`)
    c3 := UpdateClusterRequest{}
    err = json.Unmarshal(payload3, &c3)
    assert.Nil(t, err, "Should return error, but we know it does not")
    assert.False(t, c3.IsValid())
}

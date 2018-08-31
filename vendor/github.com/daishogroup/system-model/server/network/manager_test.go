//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Network manager tests.

package network

import (
    "testing"
    "github.com/stretchr/testify/assert"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/networkstorage"
)

const (
    testNetworkName = "name"
    testDescription = "description"
    testAdminName   = "adminName"
    testAdminPhone  = "adminPhone"
    testAdminEmail  = "adminEmail"
)

var networkProvider = networkstorage.NewMockupNetworkProvider()
var networkMgr = NewManager(networkProvider)

func TestManager_AddNetwork(t *testing.T) {
    networkProvider.Clear()
    toAdd := entities.NewAddNetworkRequest(testNetworkName, testDescription, testAdminName,
        testAdminPhone, testAdminEmail)
    added, err := networkMgr.AddNetwork(*toAdd)
    assert.Nil(t, err, "network should be added")
    assert.NotNil(t, added, "added network should be retrieved")
    // Check default clusters
    clusters, err := networkProvider.ListClusters(added.ID)
    assert.Nil(t, err, "clusters should be retrieved")
    assert.Equal(t, 0, len(clusters), "Expecting 0 default clusters")
}

func TestManager_GetNetwork(t *testing.T) {
    networkProvider.Clear()
    toAdd := entities.NewAddNetworkRequest(testNetworkName, testDescription, testAdminName,
        testAdminPhone, testAdminEmail)
    added, err := networkMgr.AddNetwork(*toAdd)
    assert.Nil(t, err, "network should be added")
    assert.NotNil(t, added, "added network should be retrieved")

    retrieved, err := networkMgr.GetNetwork(added.ID)
    assert.Nil(t, err, "network should exist")
    assert.EqualValues(t, added.ID, retrieved.ID, "networks should match")
}

func TestManager_ListNetworks(t *testing.T) {
    networkProvider.Clear()
    numberNetworks := 5
    for index := 0; index < numberNetworks; index ++ {
        toAdd := entities.NewAddNetworkRequest(testNetworkName, testDescription, testAdminName,
            testAdminPhone, testAdminEmail)
        _, err := networkMgr.AddNetwork(*toAdd)
        assert.Nil(t, err, "network should be added")
    }
    networks, err := networkMgr.ListNetworks()
    assert.Nil(t, err, "networks should be retrieved")
    assert.Equal(t, numberNetworks, len(networks), "number of networks should match.")
}

func TestManager_DeleteNetwork(t *testing.T){
    networkProvider.Clear()
    numberNetworks := 5
    idsList := make([]string,0)
    for index := 0; index < numberNetworks; index ++ {
        toAdd := entities.NewAddNetworkRequest(testNetworkName, testDescription, testAdminName,
            testAdminPhone, testAdminEmail)
        response, err := networkMgr.AddNetwork(*toAdd)
        assert.Nil(t, err, "network should be added")
        idsList = append(idsList,response.ID)
    }

    networks, err := networkMgr.ListNetworks()
    assert.Nil(t, err, "error retrieving the list of networks")
    assert.Equal(t, numberNetworks, len(networks),"the number of networks should match")

    err = networkMgr.DeleteNetwork(idsList[0])
    assert.Nil(t, err, "error deleting network")

    networks, err = networkMgr.ListNetworks()
    assert.Equal(t, numberNetworks-1, len(networks), "unexpected number of networks after retrieval")

    retrieved, err := networkMgr.GetNetwork(idsList[0])
    assert.NotNil(t, err, "network should not exists")
    assert.Nil(t, retrieved, "no object expected")
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Network Client TestSuite

package client

import (
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/entities"
)

func TestGetExistingNetwork(suite *suite.Suite, network Network) {
    n, err := network.Get("1")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(n, "network must not be Nil")
}

func TestGetNotExistingNetwork(suite *suite.Suite, network Network) {
    n, err := network.Get("3")
    suite.NotNil(err, "error must not be Nil")
    suite.Nil(n, "network must be nil")
}

func TestGetNetworkList(suite *suite.Suite, network Network) {
    ns, err := network.List()
    suite.Nil(err, "error must be Nil")
    suite.NotNil(ns, "network must not be Nil")
    suite.Len(ns, 2, "len must be 2")
}

func TestAddNetwork(suite *suite.Suite, network Network) {
    addedNetwork, err := network.Add(entities.AddNetworkRequest{
        Name:        "Network3",
        Description: "Description of Network3",
    })
    suite.Nil(err, "error must be Nil")
    suite.NotNil(addedNetwork,"addedNetwork must not be Nil")

    n,err := network.Get(addedNetwork.ID)
    suite.Nil(err, "error must be Nil")
    suite.NotNil(n,"n must not be Nil")
}

func TestDeleteNetwork(suite *suite.Suite, network Network) {
    addedNetwork, err := network.Add(entities.AddNetworkRequest{
        Name:        "Network3",
        Description: "Description of Network3",
    })
    suite.Nil(err, "error must be Nil")
    suite.NotNil(addedNetwork,"addedNetwork must not be Nil")

    err = network.Delete(addedNetwork.ID)
    suite.Nil(err, "error must be Nil")

    retrieved, err := network.Get(addedNetwork.ID)
    suite.Nil(retrieved, "unexpectedly something was returned after delete")
    suite.NotNil(err, "an error must be specified")
}
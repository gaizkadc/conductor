//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Network Client Mockup tests

package client

import (
    "github.com/stretchr/testify/suite"
    "testing"
)

type NetworkMockupTestSuite struct {
    suite.Suite
    network Network
}

func (suite *NetworkMockupTestSuite) SetupTest() {
    suite.network = NewNetworkMockup()
}

func (suite *NetworkMockupTestSuite) TestGetExistingNetwork() {
    TestGetExistingNetwork(&suite.Suite,suite.network)
}

func (suite *NetworkMockupTestSuite) TestGetNotExistingNetwork() {
    TestGetNotExistingNetwork(&suite.Suite,suite.network)
}

func (suite *NetworkMockupTestSuite) TestGetNetworkList() {
    TestGetNetworkList(&suite.Suite,suite.network)
}

func (suite *NetworkMockupTestSuite) TestAddNetwork() {
    TestAddNetwork(&suite.Suite,suite.network)
}

func (suite *NetworkMockupTestSuite) TestDeleteNetwork() {
    TestDeleteNetwork(&suite.Suite,suite.network)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNetworkMockup(t *testing.T) {
    suite.Run(t, new(NetworkMockupTestSuite))
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Network Client integration tests

package client_external

import (
    "testing"

    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/client"
    log "github.com/sirupsen/logrus"
    "github.com/daishogroup/dhttp"
)

type NetworkExternalTestSuite struct {
    suite.Suite
    helper EndpointHelper
    network client.Network
}

func (suite *NetworkExternalTestSuite) SetupSuite() {
    suite.helper = NewEndpointHelper()
    log.SetLevel(log.DebugLevel)
    suite.network = client.NewNetworkClientRest(suite.helper.GetListeningAddress())
    suite.helper.Start()
    dhttp.WaitURLAvailable(BaseAddress,suite.helper.port,5,"/", 1)
}

func (suite *NetworkExternalTestSuite) SetupTest() {
    suite.helper.ResetProvider()
}

func (suite *NetworkExternalTestSuite) TearDownSuite() {
    suite.helper.Shutdown()
}

func (suite *NetworkExternalTestSuite) TestGetExistingNetwork() {
    client.TestGetExistingNetwork(&suite.Suite,suite.network)
}

func (suite *NetworkExternalTestSuite) TestGetNotExistingNetwork() {
    client.TestGetNotExistingNetwork(&suite.Suite,suite.network)
}

func (suite *NetworkExternalTestSuite) TestGetNetworkList() {
    client.TestGetNetworkList(&suite.Suite,suite.network)
}

func (suite *NetworkExternalTestSuite) TestAddNetwork() {
    client.TestAddNetwork(&suite.Suite,suite.network)
}


// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNetworkExternal(t *testing.T) {
    suite.Run(t, new(NetworkExternalTestSuite))
}
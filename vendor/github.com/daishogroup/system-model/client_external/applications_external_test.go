//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the applications client integration test.

package client_external

import (
    "testing"

    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/dhttp"
)

type ApplicationsExternalTestSuite struct {
    suite.Suite
    helper  ApplicationsEndpointHelper
    applications client.Applications
}

func (suite *ApplicationsExternalTestSuite) SetupSuite() {
    suite.helper = NewApplicationsEndpointHelper()
    suite.loadTestingData()
    suite.applications = client.NewApplicationClientRest(suite.helper.GetListeningAddress())
    suite.helper.Start()
    dhttp.WaitURLAvailable(BaseAddress,suite.helper.port,5,"/", 1)

}

func (suite *ApplicationsExternalTestSuite) loadTestingData() {
    suite.helper.AddNetwork(client.TestAppNetworkID)
    suite.helper.AddDescriptor(client.TestAppNetworkID, client.TestDescriptorID)
    suite.helper.AddInstance(client.TestAppNetworkID, client.TestInstanceID)
    suite.helper.AddInstance(client.TestAppNetworkID, client.TestInstanceIDUpdate)
    suite.helper.AddNetwork(client.TestNetworkIDDelete)
    suite.helper.AddInstance(client.TestNetworkIDDelete, client.TestInstanceIDDelete)
}

func (suite *ApplicationsExternalTestSuite) SetupTest() {
    suite.helper.ResetProviders()
    suite.loadTestingData()
}

func (suite *ApplicationsExternalTestSuite) TearDownSuite() {
    suite.helper.Shutdown()
}

func TestApplicationsExternal(t *testing.T) {
    suite.Run(t, new(ApplicationsExternalTestSuite))
}

func (suite * ApplicationsExternalTestSuite) TestAddDescriptor() {
    client.TestAddDescriptor(&suite.Suite, suite.applications)
}

func (suite * ApplicationsExternalTestSuite) TestListDescriptors() {
    client.TestListDescriptors(&suite.Suite, suite.applications)
}

func (suite * ApplicationsExternalTestSuite) TestGetDescriptor() {
    client.TestGetDescriptor(&suite.Suite, suite.applications)
}

func (suite * ApplicationsExternalTestSuite) TestAddInstance() {
    client.TestAddInstance(&suite.Suite, suite.applications)
}

func (suite * ApplicationsExternalTestSuite) TestListInstances() {
    client.TestListInstances(&suite.Suite, suite.applications)
}

func (suite * ApplicationsExternalTestSuite) TestGetInstance() {
    client.TestGetInstance(&suite.Suite, suite.applications)
}

func (suite * ApplicationsExternalTestSuite) TestUpdateInstance() {
    client.TestUpdateInstance(&suite.Suite, suite.applications)
}

func (suite * ApplicationsExternalTestSuite) TestDeleteInstance() {
    client.TestDeleteInstance(&suite.Suite, suite.applications)
}



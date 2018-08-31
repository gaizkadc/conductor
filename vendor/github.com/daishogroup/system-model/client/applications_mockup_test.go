//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Application mockup tests

package client

import (

    "testing"

    "github.com/stretchr/testify/suite"

)

type ApplicationMockupTestSuite struct {
    suite.Suite
    client Applications
    mockup * ApplicationsMockup
}

func (suite * ApplicationMockupTestSuite) loadTestStructs() {
    suite.mockup.AddTestNetwork(TestAppNetworkID)
    suite.mockup.AddTestDescriptor(TestAppNetworkID, TestDescriptorID)
    suite.mockup.AddTestInstance(TestAppNetworkID, TestDescriptorID, TestInstanceID)
    suite.mockup.AddTestInstance(TestAppNetworkID, TestDescriptorID, TestInstanceIDUpdate)
    suite.mockup.AddTestNetwork(TestNetworkIDDelete)
    suite.mockup.AddTestInstance(TestNetworkIDDelete, TestDescriptorID, TestInstanceIDDelete)
}

func (suite * ApplicationMockupTestSuite) SetupTest() {
    suite.client = NewApplicationsMockup()
    suite.mockup = suite.client.(*ApplicationsMockup)
    suite.loadTestStructs()
}

func (suite * ApplicationMockupTestSuite) TestAddDescriptor() {
    TestAddDescriptor(&suite.Suite, suite.client)
}

func (suite * ApplicationMockupTestSuite) TestListDescriptors() {
    TestListDescriptors(&suite.Suite, suite.client)
}

func (suite * ApplicationMockupTestSuite) TestGetDescriptor() {
    TestGetDescriptor(&suite.Suite, suite.client)
}

func (suite * ApplicationMockupTestSuite) TestDeleteDescriptor() {
    TestDeleteDescriptor(&suite.Suite, suite.client)
}

func (suite * ApplicationMockupTestSuite) TestAddInstance() {
    TestAddInstance(&suite.Suite, suite.client)
}

func (suite * ApplicationMockupTestSuite) TestListInstances() {
    TestListInstances(&suite.Suite, suite.client)
}

func (suite * ApplicationMockupTestSuite) TestGetInstance() {
    TestGetInstance(&suite.Suite, suite.client)
}

func (suite * ApplicationMockupTestSuite) TestUpdateInstance() {
    TestUpdateInstance(&suite.Suite, suite.client)
}

func (suite * ApplicationMockupTestSuite) TestDeleteInstance() {
    TestDeleteInstance(&suite.Suite, suite.client)
}

func TestApplicationMockup(t *testing.T){
    suite.Run(t, new(ApplicationMockupTestSuite))
}
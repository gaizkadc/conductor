//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Cluster Client integration tests

package client_external

import (
    "testing"

    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/dhttp"
)

type ClusterExternalTestSuite struct {
    suite.Suite
    helper  EndpointHelper
    cluster client.Cluster
}

func (suite *ClusterExternalTestSuite) SetupSuite() {
    suite.helper = NewEndpointHelper()
    suite.cluster = client.NewClusterClientRest(suite.helper.GetListeningAddress())
    suite.helper.Start()
    dhttp.WaitURLAvailable(BaseAddress,suite.helper.port,5,"/", 1)
}

func (suite *ClusterExternalTestSuite) SetupTest() {
    suite.helper.ResetProvider()
}

func (suite *ClusterExternalTestSuite) TearDownSuite() {
    suite.helper.Shutdown()
}

func (suite *ClusterExternalTestSuite) TestGetExistingCluster() {
    client.TestGetExistingCluster(&suite.Suite, suite.cluster)
}

func (suite *ClusterExternalTestSuite) TestGetNotExistingCluster() {
    client.TestGetNotExistingCluster(&suite.Suite, suite.cluster)
}

func (suite *ClusterExternalTestSuite) TestGetClusterNotExistingNetwork() {
    client.TestGetClusterNotExistingNetwork(&suite.Suite, suite.cluster)
}

func (suite *ClusterExternalTestSuite) TestAddCluster() {
    client.TestAddCluster(&suite.Suite, suite.cluster)
}

func (suite *ClusterExternalTestSuite) TestAddClusterToNotExistingNetwork() {
    client.TestAddClusterToNotExistingNetwork(&suite.Suite, suite.cluster)
}

func (suite *ClusterExternalTestSuite) TestListByNetwork() {
    client.TestListByNetwork(&suite.Suite, suite.cluster)
}

func (suite *ClusterExternalTestSuite) TestListByNotExistingNetwork() {
    client.TestListByNotExistingNetwork(&suite.Suite, suite.cluster)
}

func (suite *ClusterExternalTestSuite) TestUpdateCluster() {
    client.TestUpdateCluster(&suite.Suite,suite.cluster)
}

func TestClusterExternal(t *testing.T) {
    suite.Run(t, new(ClusterExternalTestSuite))
}

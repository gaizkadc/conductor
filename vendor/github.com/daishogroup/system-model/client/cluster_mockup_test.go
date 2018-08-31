//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Cluster Client Mockup Tests

package client



import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClusterMockupTestSuite struct {
	suite.Suite
	cluster Cluster
}

func (suite *ClusterMockupTestSuite) SetupSuite() {
    suite.cluster = NewClusterMockup()
}

func (suite *ClusterMockupTestSuite) SetupTest() {
	suite.cluster.(*ClusterMockup).ClearMockup()
	suite.cluster.(*ClusterMockup).InitMockup()
}

func (suite *ClusterMockupTestSuite) TestGetExistingCluster() {
	TestGetExistingCluster(&suite.Suite,suite.cluster)
}

func (suite *ClusterMockupTestSuite) TestGetNotExistingCluster() {
	TestGetNotExistingCluster(&suite.Suite,suite.cluster)
}

func (suite *ClusterMockupTestSuite) TestGetClusterNotExistingNetwork() {
	TestGetClusterNotExistingNetwork(&suite.Suite,suite.cluster)
}


func (suite *ClusterMockupTestSuite) TestAddCluster() {
	TestAddCluster(&suite.Suite,suite.cluster)
}


func (suite *ClusterMockupTestSuite) TestAddClusterToNotExistingNetwork()  {
	TestAddClusterToNotExistingNetwork(&suite.Suite,suite.cluster)
}

func (suite *ClusterMockupTestSuite) TestListByNetwork() {
	TestListByNetwork(&suite.Suite,suite.cluster)
}

func (suite *ClusterMockupTestSuite) TestListByNotExistingNetwork() {
	TestListByNotExistingNetwork(&suite.Suite,suite.cluster)
}

func (suite *ClusterMockupTestSuite) TestUpdateCluster() {
	TestUpdateCluster(&suite.Suite,suite.cluster)
}

func (suite *ClusterMockupTestSuite) TestDeleteCluster() {
	TestDeleteCluster(&suite.Suite,suite.cluster)
}

func TestClusterMockup(t *testing.T) {
	suite.Run(t, new(ClusterMockupTestSuite))
}



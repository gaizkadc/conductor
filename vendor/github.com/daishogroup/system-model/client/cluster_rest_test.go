//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Cluster Client Rest Test

package client

import (
    "testing"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/dhttp"
)

type ClusterRestTestSuite struct {
    suite.Suite
    cluster Cluster
    rest    *dhttp.ClientMockup
}

func (suite *ClusterRestTestSuite) SetupSuite() {
    rest := dhttp.NewClientMockup()
    suite.rest = rest
    cluster := ClusterRest{rest}
    suite.cluster = &cluster

    suite.Equal(suite.rest,cluster.client, "rest client must be equals")
}

func (suite *ClusterRestTestSuite) SetupTest() {
    suite.rest.Reset()
}

func (suite *ClusterRestTestSuite) TestGetExistingCluster() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := *entities.NewClusterWithID("1","1", "Cluster1", "Description Cluster 1",
            entities.GatewayType, "Madrid", "admin@admin.com",
            entities.ClusterCreated, false, false)
        statusCode := 200
        return dhttp.NewResponse(&result,&statusCode, nil)
    })
    TestGetExistingCluster(&suite.Suite, suite.cluster)
}

func (suite *ClusterRestTestSuite) TestGetNotExistingCluster() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil,&statusCode, derrors.NewOperationError(errors.ClusterDoesNotExists))
    })
    TestGetNotExistingCluster(&suite.Suite, suite.cluster)
}

func (suite *ClusterRestTestSuite) TestGetClusterNotExistingNetwork() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil,&statusCode, derrors.NewOperationError(errors.NetworkDoesNotExists))
    })
    TestGetClusterNotExistingNetwork(&suite.Suite, suite.cluster)
}

func (suite *ClusterRestTestSuite) TestAddCluster() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := *entities.NewClusterWithID("1","4", "Cluster4", "Description Cluster 3",
            entities.GatewayType, "Oregon", "admin@admin.com",
            entities.ClusterCreated, false, false)
        return dhttp.NewResponse(&result,&statusCode, nil)
    })
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := *entities.NewClusterWithID("1","4", "Cluster4", "Description Cluster 3",
            entities.GatewayType, "Oregon", "admin@admin.com",
            entities.ClusterCreated, false, false)
        statusCode := 200
        return dhttp.NewResponse(&result,&statusCode, nil)
    })
    TestAddCluster(&suite.Suite, suite.cluster)
}

func (suite *ClusterRestTestSuite) TestAddClusterToNotExistingNetwork() {
    suite.rest.AddPost(func(path string,body interface{}) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil,&statusCode, derrors.NewOperationError(errors.NetworkDoesNotExists))
    })
    TestAddClusterToNotExistingNetwork(&suite.Suite, suite.cluster)
}

func (suite *ClusterRestTestSuite) TestListByNetwork() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := [] entities.Cluster{
            *entities.NewClusterWithID("1","1", "Cluster1", "Description Cluster 1",
                entities.GatewayType, "Madrid", "admin@admin.com",
                entities.ClusterCreated, false, false),
            *entities.NewClusterWithID("1","2", "Cluster2", "Description Cluster 2",
                entities.GatewayType, "Boston", "admin@admin.com",
                entities.ClusterCreated, false, false),
        }
        statusCode := 200
        return dhttp.NewResponse(&result,&statusCode, nil)
    })
    TestListByNetwork(&suite.Suite, suite.cluster)
}

func (suite *ClusterRestTestSuite) TestListByNotExistingNetwork() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil,&statusCode, derrors.NewOperationError(errors.NetworkDoesNotExists))
    })
    TestListByNotExistingNetwork(&suite.Suite, suite.cluster)
}

func (suite *ClusterRestTestSuite) TestUpdateCluster() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := *entities.NewClusterWithID("1","4", "Cluster4", "Description Cluster 3",
            entities.GatewayType, "Oregon", "admin@admin.com",
            entities.ClusterInstalled, false, false)
        return dhttp.NewResponse(&result,&statusCode, nil)
    })
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := *entities.NewClusterWithID("1","4", "newUpdate", "Description Cluster 3",
            entities.GatewayType, "Oregon", "admin@admin.com",
            entities.ClusterInstalled, false, false)
        return dhttp.NewResponse(&result,&statusCode, nil)
    })
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := *entities.NewClusterWithID("1","4", "newUpdate", "Description Cluster 3",
            entities.GatewayType, "Oregon", "admin@admin.com",
            entities.ClusterInstalled, false, false)
        return dhttp.NewResponse(&result,&statusCode, nil)
    })
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 200
        result := *entities.NewClusterWithID("1","4", "newUpdate", "Description Cluster 3",
            entities.GatewayType, "Oregon", "admin@admin.com",
            entities.ClusterInstalled, false, false)
        return dhttp.NewResponse(&result,&statusCode, nil)
    })
    TestUpdateCluster(&suite.Suite,suite.cluster)
}


func (suite *ClusterRestTestSuite) TestDeleteCluster() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := [] entities.Cluster{
            *entities.NewClusterWithID("1","1", "Cluster1", "Description Cluster 1",
                entities.GatewayType, "Madrid", "admin@admin.com",
                entities.ClusterCreated, false, false),
            *entities.NewClusterWithID("1","2", "Cluster2", "Description Cluster 2",
                entities.GatewayType, "Boston", "admin@admin.com",
                entities.ClusterCreated, false, false),
        }
        statusCode := 200
        return dhttp.NewResponse(&result,&statusCode, nil)
    })
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 200
        result := entities.NewSuccessfulOperation("DeleteCluster")
        return dhttp.NewResponse(result,&statusCode,nil)
    })
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := *entities.NewClusterWithID("n1","1", "cluster 8", "Description cluster 8",
            entities.GatewayType, "madrid", "email", entities.ClusterInstalled,
                false, false)
        return dhttp.NewResponse(&result, &statusCode, nil)
    })
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 500
        return dhttp.NewResponse(nil,&statusCode, derrors.NewOperationError(errors.ClusterDoesNotExists))
    })
    suite.rest.AddGet(func(path string) dhttp.Response {
        result := [] entities.Cluster{
            *entities.NewClusterWithID("1","1", "Cluster1", "Description Cluster 1",
                entities.GatewayType, "Madrid", "admin@admin.com",
                entities.ClusterCreated, false, false),
            *entities.NewClusterWithID("1","2", "Cluster2", "Description Cluster 2",
                entities.GatewayType, "Boston", "admin@admin.com",
                entities.ClusterCreated, false, false),
        }
        statusCode := 200
        return dhttp.NewResponse(&result,&statusCode, nil)
    })

    TestDeleteCluster(&suite.Suite, suite.cluster)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestClusterRest(t *testing.T) {
    suite.Run(t, new(ClusterRestTestSuite))
}

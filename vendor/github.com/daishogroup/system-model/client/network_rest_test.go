//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Network Client Rest Test


package client

import (
    "testing"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/dhttp"
)

type NetworkRestTestSuite struct {
    suite.Suite
    network Network
    rest    *dhttp.ClientMockup
}

func (suite *NetworkRestTestSuite) SetupSuite() {
    rest := dhttp.NewClientMockup()
    suite.rest = rest
    cluster := NetworkRest{rest}
    suite.network = &cluster

    suite.Equal(suite.rest,cluster.client, "rest client must be equals")
}

func (suite *NetworkRestTestSuite) SetupTest() {
    suite.rest.Reset()
}

func (suite *NetworkRestTestSuite) TestGetExistingNetwork() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        result:=entities.NewNetworkWithID("1","n1","d1","a1","ap1","ae1")
        statusCode:=200
        return dhttp.NewResponse(result,&statusCode,nil)
    })
    TestGetExistingNetwork(&suite.Suite, suite.network)
}

func (suite *NetworkRestTestSuite) TestGetNotExistingNetwork() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode:=404
        return dhttp.NewResponse(nil,&statusCode, derrors.NewOperationError(errors.NetworkDoesNotExists))
    })
    TestGetNotExistingNetwork(&suite.Suite, suite.network)
}

func (suite *NetworkRestTestSuite) TestGetNetworkList() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        result:=[] entities.Network{
            *entities.NewNetworkWithID("1","n1","d1","a1","ap1","ae1"),
            *entities.NewNetworkWithID("2","n2","d2","a2","ap2","ae2"),
        }
        statusCode:=200
        return dhttp.NewResponse(&result,&statusCode,nil)
    })
    TestGetNetworkList(&suite.Suite, suite.network)
}

func (suite *NetworkRestTestSuite) TestAddNetwork() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode:=200
        result:=entities.NewNetworkWithID("3","n3","d3","a3","ap3","ae3")
        return dhttp.NewResponse(result,&statusCode,nil)
    })

    suite.rest.AddGet(func(path string) dhttp.Response {
        result:=entities.NewNetworkWithID("3","n3","d3","a3","ap3","ae3")
        statusCode:=200
        return dhttp.NewResponse(result,&statusCode,nil)
    })
    TestAddNetwork(&suite.Suite, suite.network)
}

func (suite *NetworkRestTestSuite) TestDeleteNetwork() {
    suite.rest.AddDelete(func(path string) dhttp.Response{
        statusCode:=200
        result := entities.NewSuccessfulOperation("DeleteNetwork")
        return dhttp.NewResponse(result,&statusCode, nil)
    })
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        result:=*entities.NewNetworkWithID("3","Network3","d3","a3","ap3","ae3")
        statusCode:=200
        return dhttp.NewResponse(&result,&statusCode,nil)
    })
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode:=500
        return dhttp.NewResponse(nil,&statusCode, derrors.NewOperationError(errors.NetworkDoesNotExists))
    })

    TestDeleteNetwork(&suite.Suite, suite.network)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNetworkRest(t *testing.T) {
    suite.Run(t, new(NetworkRestTestSuite))
}

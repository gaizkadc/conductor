//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Testing client to connect with conductor API REST

package client

import (
    "testing"

    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/dhttp"
    "github.com/daishogroup/system-model/entities"
)

type ConductorRestTestSuite struct {
    suite.Suite
    client Conductor
    rest   *dhttp.ClientMockup
}

func (suite *ConductorRestTestSuite) SetupSuite() {
    rest := dhttp.NewClientMockup()
    suite.rest = rest
    client := ConductorRest{rest}
    suite.client = &client
}

func (suite *ConductorRestTestSuite) SetupTest() {
    suite.rest.Reset()
}

func TestConductorRest(t *testing.T) {
    suite.Run(t, new(ConductorRestTestSuite))
}

func (suite *ConductorRestTestSuite) TestDeploy() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewAppInstance("1", "1", "1", "name",
            "description", "label", "arguments", "10GB", entities.AppStorageDefault,
            make([]entities.ApplicationPort, 0), 80, "address")
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestDeploy(&suite.Suite, suite.client)
}

func (suite *ConductorRestTestSuite) TestUndeploy() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 200
        result := ""
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestUndeploy(&suite.Suite, suite.client)
}

//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package dhttp

import (
    "github.com/daishogroup/derrors"
    "github.com/stretchr/testify/suite"
    "testing"
)

type ClientSlingTestSuite struct {
    ClientTestSuite
    Handler *TestHandler
    helper  *ClientTestHelper
}

func (suite *ClientSlingTestSuite) SetupSuite() {
    suite.Handler = NewTestHandler()
    suite.helper = NewClientTestHelper(suite.Handler, host)
    config := NewRestBasicConfig(suite.helper.host, suite.helper.port)
    // Add header for test
    config.Headers = map[string]string{
        "X-Daisho-Test": "test123",
    }
    suite.client = NewClientSling(config)
    suite.helper.Start()
}

func (suite *ClientSlingTestSuite) TearDownSuite() {
    suite.helper.Shutdown()
}
func (suite *ClientSlingTestSuite) SetupTest() {
    suite.Handler.Err = false
}

func (suite *ClientSlingTestSuite) TestGetNotValidError() {
    suite.Handler.Err = true
    url := suite.GetURL("/test/get")
    result := &SuccessfulOperation{}
    response := suite.client.Get(url, result)
    suite.Error(response.Error, "Must be a error")
    suite.NotEqual(response.Error.Type(), derrors.ConnectionErrorType)
}

func (suite *ClientSlingTestSuite) TestDeleteNotValidError() {
    suite.Handler.Err = true
    url := suite.GetURL("/test/delete")
    result := &SuccessfulOperation{}
    response := suite.client.Delete(url, result)
    suite.Error(response.Error, "Must be a error")
    suite.NotEqual(response.Error.Type(), derrors.ConnectionErrorType)
}


func TestClientSling(t *testing.T) {
    suite.Run(t, new(ClientSlingTestSuite))
}

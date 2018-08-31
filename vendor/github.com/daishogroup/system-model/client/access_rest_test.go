//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the Access Client TestSuite

package client

import (
    "testing"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/dhttp"
    "github.com/stretchr/testify/suite"
)

type AccessRestTestSuite struct {
    suite.Suite
    client Access
    rest * dhttp.ClientMockup
}

func (suite * AccessRestTestSuite) SetupSuite() {
    rest := dhttp.NewClientMockup()
    suite.rest = rest
    client := AccessRest{rest}
    suite.client = &client
}

func (suite * AccessRestTestSuite) SetupTest() {
    suite.rest.Reset()
}

func TestAccessRest(t *testing.T){
    suite.Run(t, new(AccessRestTestSuite))
}

func (suite * AccessRestTestSuite) TestAddUserAccess(){
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response{
        result := entities.NewUserAccess("user",[]entities.RoleType{entities.GlobalAdmin})
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestAddUserAccess(&suite.Suite, suite.client)
}

func (suite * AccessRestTestSuite) TestGetUserAccess(){

    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response{
        result := entities.NewUserAccess("user",[]entities.RoleType{entities.GlobalAdmin})
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    suite.rest.AddGet(func(path string) dhttp.Response{
        result := entities.NewUserAccess("user",[]entities.RoleType{entities.GlobalAdmin})
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestGetUserAccess(&suite.Suite, suite.client)
}

func (suite * AccessRestTestSuite) TestDeleteUserAccess(){

    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response{
        result := entities.NewUserAccess("user",[]entities.RoleType{entities.GlobalAdmin})
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    suite.rest.AddDelete(func(path string) dhttp.Response{
        result := entities.NewSuccessfulOperation(errors.AccessDeleted)
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    suite.rest.AddGet(func(path string) dhttp.Response{
        statusCode := 500
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.UserDoesNotExist))
    })
    TestDeleteUserAccess(&suite.Suite, suite.client)
}
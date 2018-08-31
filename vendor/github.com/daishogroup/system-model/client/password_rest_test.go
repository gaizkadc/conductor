//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the User Client Rest Test

package client

import (
    "testing"

    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/dhttp"
)

type PasswordRestTestSuite struct {
    suite.Suite
    password Password
    rest *dhttp.ClientMockup
}

func (suite *PasswordRestTestSuite) SetupSuite() {
    suite.rest = dhttp.NewClientMockup()
    suite.password = &PasswordRest{suite.rest}
}

func (suite *PasswordRestTestSuite) SetupTest() {
    suite.rest.Reset()

}

func (suite *PasswordRestTestSuite) TestSetPassword() {

    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        return dhttp.NewResponse(entities.NewSuccessfulOperation(errors.PasswordSet), &statusCode, nil)
    })

    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 200
        p := "anotherPassword"
        toReturn, _ := entities.NewPassword("userDefault",&p)
        return dhttp.NewResponse(*toReturn, &statusCode, nil)
    })

    TestSetPassword(&suite.Suite, suite.password)
}


func (suite *PasswordRestTestSuite) TestGetPassword() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 200
        p := "apassword"
        result, _ := entities.NewPassword("userDefault",&p)
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    TestGetPassword(&suite.Suite, suite.password)
}

func (suite *PasswordRestTestSuite) TestDeletePassword() {
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 200
        return dhttp.NewResponse(entities.NewSuccessfulOperation(errors.PasswordDeleted), &statusCode, nil)
    })

    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 500
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.UserDoesNotExist))
    })

    TestDeletePassword(&suite.Suite, suite.password)
}


// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPasswordRest(t *testing.T) {
    suite.Run(t, new(UserRestTestSuite))
}
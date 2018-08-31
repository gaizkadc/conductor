//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the User Client Rest Test

package client

import (
    "testing"

    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/dhttp"
)

type UserRestTestSuite struct {
    suite.Suite
    user User
    rest *dhttp.ClientMockup
}

func (suite *UserRestTestSuite) SetupSuite() {
    suite.rest = dhttp.NewClientMockup()
    suite.user = &UserRest{suite.rest}
}

func (suite *UserRestTestSuite) SetupTest() {
    suite.rest.Reset()

}

func (suite *UserRestTestSuite) TestAddUser() {

    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewUserWithID("userTestId", "userTest", "9999 6666",
            "userTest@email.com", testUserCreationTime, testUserExpirationTime)
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    TestAddUser(&suite.Suite, suite.user)
}


func (suite *UserRestTestSuite) TestAddAlreadyExistingUser() {

    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 404
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.UserAlreadyExists))
    })

    TestAddAlreadyExistingUser(&suite.Suite, suite.user)
}

func (suite *UserRestTestSuite) TestGetUser() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewUserWithID("userTestId", "userTest", "9999 6666",
            "userTest@email.com", testUserCreationTime, testUserExpirationTime)
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 200
        result := entities.NewUserWithID("userTestId", "userTest", "9999 6666",
            "userTest@email.com", testUserCreationTime, testUserExpirationTime)
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    TestGetUser(&suite.Suite, suite.user)
}

func (suite *UserRestTestSuite) TestGetNonExistingUser() {
    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 500
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.UserDoesNotExist))
    })

    TestGetNonExistingUser(&suite.Suite, suite.user)
}

func (suite *UserRestTestSuite) TestDeleteUser() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewUserWithID("userTestId", "userTest", "9999 6666",
            "userTest@email.com", testUserCreationTime, testUserExpirationTime)
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 200
        result := entities.NewSuccessfulOperation("deleteUser")
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    TestDeleteUser(&suite.Suite, suite.user)
}

func (suite *UserRestTestSuite) TestDeleteNonExistingUser() {
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 500
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.UserDoesNotExist))
    })

    TestDeleteNonExistingUser(&suite.Suite, suite.user)
}

func (suite *UserRestTestSuite) TestUpdateUser() {

    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 200
        result := entities.NewUserWithID("userTestId", "userTest", "9999 6666",
            "userTest@email.com", testUserCreationTime, testUserExpirationTime)
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewUserWithID("userTestId", "1", "2",
            "3", testUserCreationTime, testUserExpirationTime)
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    TestUpdateUser(&suite.Suite, suite.user)
}

func (suite *UserRestTestSuite) TestListUsers() {
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewUserWithID("a", "name", "phone", "email",
            testUserCreationTime, testUserExpirationTime)
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        result := entities.NewUserWithID("b", "name", "phone", "email",
            testUserCreationTime, testUserExpirationTime)
        return dhttp.NewResponse(result, &statusCode, nil)
    })

    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 200
        result := []entities.UserExtended{
            entities.NewUserExtended("a", "name", "phone","email", testUserCreationTime,
                testUserExpirationTime,[]entities.RoleType{entities.GlobalAdmin}),
            entities.NewUserExtended("b", "name", "phone","email", testUserCreationTime, testUserExpirationTime,
                []entities.RoleType{entities.GlobalAdmin}),
        }
        return dhttp.NewResponse(&result, &statusCode, nil)
    })

    TestListUsers(&suite.Suite, suite.user)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserRest(t *testing.T) {
    suite.Run(t, new(UserRestTestSuite))
}
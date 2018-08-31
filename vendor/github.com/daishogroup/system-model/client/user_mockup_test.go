//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the User Client Mockup tests

package client

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/suite"
    "testing"
)

type UserMockupTestSuite struct {
    suite.Suite
    user User
}

func (suite *UserMockupTestSuite) SetupSuite() {
    suite.user = NewUserMockup()
}

func (suite *UserMockupTestSuite) SetupTest() {
    suite.user.(*UserMockup).InitMockup()
}

func (suite *UserMockupTestSuite) TestAddUser() {
    TestAddUser(&suite.Suite, suite.user)
}

func (suite *UserMockupTestSuite) TestGetUser() {
    TestGetUser(&suite.Suite, suite.user)
}

func (suite *UserMockupTestSuite) TestGetNonExistingUser() {
    TestGetNonExistingUser(&suite.Suite, suite.user)
}

func (suite *UserMockupTestSuite) TestDeleteUser() {
    TestDeleteUser(&suite.Suite, suite.user)
}

func (suite *UserMockupTestSuite) TestDeleteNonExistingUser() {
    TestDeleteNonExistingUser(&suite.Suite, suite.user)
}

func ListUsersMockup(suite *suite.Suite, user User) {
    // Add test user
    _, err := user.Add(*entities.NewAddUserRequest("a", "name", "phone",
        "email", testUserCreationTime, testUserExpirationTime))
    suite.Nil(err, "error must be Nil preparing testing environment")

    _, err = user.Add(*entities.NewAddUserRequest("b", "name", "phone",
        "email", testUserCreationTime, testUserExpirationTime))
    suite.Nil(err, "error must be Nil preparing testing environment")

    err = user.(*UserMockup).FixUserAccess("a", "pa", entities.GlobalAdmin)
    suite.Nil(err, "should set the role")
    err = user.(*UserMockup).FixUserAccess("b", "pb", entities.GlobalAdmin)
    suite.Nil(err, "should set the role")

    // Add roles
    users, err := user.ListUsers()
    suite.Nil(err, "unexpected nil list")
    suite.Equal(3, len(users), "unexpected length of entries")
}

func (suite *UserMockupTestSuite) TestList() {
    ListUsersMockup(&suite.Suite, suite.user)
}


// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserMockup(t *testing.T) {
    suite.Run(t, new(UserMockupTestSuite))
}
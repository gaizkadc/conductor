//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the Node Client TestSuite

package client

import (
    "time"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/entities"
)

var testUserCreationTime = time.Now()
var testUserExpirationTime = testUserCreationTime.Add(time.Hour)

func TestAddUser(suite *suite.Suite, user User) {
    n, err := user.Add(*entities.NewAddUserRequest("userTestId", "userTest", "9999 6666",
        "userTest@email.com", testUserCreationTime, testUserExpirationTime))
    suite.Nil(err, "error must be Nil")
    suite.NotNil(n, "node must not nil")
    suite.Equal("userTestId", n.ID, "user Id must match")
    suite.Equal("userTest", n.Name, "user Id must match")
    suite.Equal("9999 6666", n.Phone, "user Id must match")
    suite.Equal("userTest@email.com", n.Email, "user Id must match")
}

func TestAddAlreadyExistingUser(suite *suite.Suite, user User) {
    _, err := user.Add(*entities.NewAddUserRequest("userDefaultId", "userTest", "9999 6666",
        "userTest@email.com", testUserCreationTime, testUserExpirationTime))
    suite.NotNil(err,  "an error was expected")
}

func TestGetUser(suite *suite.Suite, user User) {
    // Add test user
    _, err := user.Add(*entities.NewAddUserRequest("userTestId", "userTest", "9999 6666",
        "userTest@email.com", testUserCreationTime, testUserExpirationTime))
    suite.Nil(err, "error must be Nil preparing testing environment")
    retrieved, err := user.Get("userTestId")
    suite.Nil(err, "unexpected error")
    suite.NotNil(retrieved, "retrieved user is nil")

    suite.Equal("userTestId", retrieved.ID, "user Id must match")
    suite.Equal("userTest", retrieved.Name, "user Id must match")
    suite.Equal("9999 6666", retrieved.Phone, "user Id must match")
    suite.Equal("userTest@email.com", retrieved.Email, "user Id must match")

}

func TestGetNonExistingUser(suite *suite.Suite, user User) {
    n, err := user.Get("nonexisting")
    suite.NotNil(err, "an  error must be returned")
    suite.Nil(n, "empty user must be returned")
}

func TestDeleteUser(suite *suite.Suite, user User) {
    // Add test user
    _, err := user.Add(*entities.NewAddUserRequest("userTestId", "userTest", "9999 6666",
        "userTest@email.com", testUserCreationTime, testUserExpirationTime))
    suite.Nil(err, "error must be Nil preparing testing environment")

    err = user.Delete("userTestId")
    suite.Nil(err, "unexpected error")
}

func TestDeleteNonExistingUser(suite *suite.Suite, user User) {
    err := user.Delete("nonexisting")
    suite.NotNil(err, "an error must be returned")
}

func TestUpdateUser(suite *suite.Suite, user User) {
    // Add test user
    _, err := user.Add(*entities.NewAddUserRequest("userTestId", "userTest", "9999 6666",
        "userTest@email.com", testUserCreationTime, testUserExpirationTime))
    suite.Nil(err, "error must be Nil preparing testing environment")

    updateRequest := entities.NewUpdateUserRequest().WithName("1").
        WithPhone("2").WithEmail("3")

    retrieved, err := user.Update("userTestId", *updateRequest)

    suite.Nil(err, "unexpected error")
    suite.NotNil(retrieved, "retrieved user is nil")

    suite.Equal("userTestId", retrieved.ID, "user Id must match")
    suite.Equal("1", retrieved.Name, "user Id must match")
    suite.Equal("2", retrieved.Phone, "user Id must match")
    suite.Equal("3", retrieved.Email, "user Id must match")
}

func TestListUsers(suite *suite.Suite, user User) {
    // Add test user
    _, err := user.Add(*entities.NewAddUserRequest("a", "name", "phone",
        "email", testUserCreationTime, testUserExpirationTime))
    suite.Nil(err, "error must be Nil preparing testing environment")

    _, err = user.Add(*entities.NewAddUserRequest("b", "name", "phone",
        "email", testUserCreationTime, testUserExpirationTime))
    suite.Nil(err, "error must be Nil preparing testing environment")

    // Add roles
    users, err := user.ListUsers()
    suite.Nil(err, "unexpected nil list")
    suite.Equal(2, len(users), "unexpected length of entries")
}

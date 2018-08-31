//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the Node Client TestSuite

package client

import (
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/entities"
)

func TestSetPassword(suite *suite.Suite, password Password) {
    p := "anotherPassword"
    toSet, err := entities.NewPassword("userDefault", &p)
    suite.Nil(err, "unexpected error")
    err = password.SetPassword(*toSet)
    suite.Nil(err, "unexpected error")

    // retrieve and check it was changed
    ret, err := password.GetPassword("userDefault")
    suite.Nil(err, "unexpected error")
    suite.Equal("userDefault", ret.UserID, "non matching entry")
    suite.True(ret.CompareWith("anotherPassword"), "passwords do not match")
}

func TestGetPassword(suite *suite.Suite, password Password) {
    p := "apassword"
    ret, err := password.GetPassword("userDefault")
    suite.Nil(err, "unexpected error")
    suite.Equal("userDefault", ret.UserID, "non matching entry")
    suite.True(ret.CompareWith(p), "passwords do not match")
}

func TestDeletePassword(suite *suite.Suite, password Password) {
    err := password.DeletePassword("userDefault")
    suite.Nil(err, "unexpected error")
    // try to get the password
    ret, err := password.GetPassword("userDefault")
    suite.NotNil(err, "an error was expected")
    suite.Nil(ret, "nil object was expected")
}

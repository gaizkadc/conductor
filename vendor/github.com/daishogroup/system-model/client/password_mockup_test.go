//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the User Client Rest Test

package client

import (
    "testing"

    "github.com/stretchr/testify/suite"
)

type PasswordMockupTestSuite struct {
    suite.Suite
    password Password
}

func (suite *PasswordMockupTestSuite) SetupSuite() {
    suite.password = NewPasswordMockup()
}

func (suite *PasswordMockupTestSuite) SetupTest() {
    suite.password.(*PasswordMockup).ClearMockup()
    suite.password.(*PasswordMockup).InitMockup()
}

func (suite *PasswordMockupTestSuite) TestSetPassword() {
    TestSetPassword(&suite.Suite, suite.password)
}

func (suite *PasswordMockupTestSuite) TestGetPassword() {
    TestGetPassword(&suite.Suite, suite.password)
}

func (suite *PasswordMockupTestSuite) TestDeletePassword() {
    TestDeletePassword(&suite.Suite, suite.password)
}


// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPasswordMockup(t *testing.T) {
    suite.Run(t, new(PasswordMockupTestSuite))
}
//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the Access Client Mockup tests

package client

import (
    "github.com/stretchr/testify/suite"
    "testing"
)

type AccessMockupTestSuite struct {
    suite.Suite
    access Access
}

func (suite *AccessMockupTestSuite) SetupSuite() {
    suite.access = NewAccessMockup()
}
func (suite * AccessMockupTestSuite) SetupTest() {
    suite.access.(*AccessMockup).ClearAccessMockup()
}

func (suite * AccessMockupTestSuite) TestAddUserAccess(){
    TestAddUserAccess(&suite.Suite, suite.access)
}

func (suite * AccessMockupTestSuite) TestGetUserAccess(){
    TestGetUserAccess(&suite.Suite, suite.access)
}

func (suite * AccessMockupTestSuite) TestDeleteUserAccess(){
    TestDeleteUserAccess(&suite.Suite, suite.access)
}
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAccessMockup(t *testing.T) {
    suite.Run(t, new(NodeMockupTestSuite))
}
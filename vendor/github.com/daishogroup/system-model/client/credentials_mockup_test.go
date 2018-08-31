//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the User Client Mockup tests

package client

import (
    "github.com/stretchr/testify/suite"
    "testing"
)

type CredentialsMockupTestSuite struct {
    suite.Suite
    credentials Credentials
}

func (suite *CredentialsMockupTestSuite) SetupSuite() {
    suite.credentials = NewCredentialsMockup()
}

func (suite *CredentialsMockupTestSuite) SetupTest() {
    suite.credentials.(*CredentialMockup).ClearMockup()
}

func (suite *CredentialsMockupTestSuite) TestAddUser() {
    TestAddCredentials(&suite.Suite, suite.credentials)
}

func (suite *CredentialsMockupTestSuite) TestGetUser() {
    TestGetCredentials(&suite.Suite, suite.credentials)
}

func (suite *CredentialsMockupTestSuite) TestDeleteCredentials() {
    TestDeleteCredentials(&suite.Suite, suite.credentials)
}


// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCredentialsMockup(t *testing.T) {
    suite.Run(t, new(CredentialsMockupTestSuite))
}
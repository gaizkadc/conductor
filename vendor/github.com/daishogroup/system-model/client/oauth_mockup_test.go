//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the User Client Rest Test

package client

import (
    "testing"

    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/entities"
)

type OAuthMockupTestSuite struct {
    suite.Suite
    oauth OAuth
}

func (suite *OAuthMockupTestSuite) SetupSuite() {
    suite.oauth = NewOAuthMockup()
}

func (suite *OAuthMockupTestSuite) SetupTest() {
    suite.oauth.(*OAuthMockup).ClearMockup()
}

func (suite *OAuthMockupTestSuite) TestSetPassword() {
    entrySet := entities.NewOAuthSecrets(OAuthTestUserID)
    entrySet.AddEntry("app1", entities.NewOAuthEntry("clientID1","secret1"))
    suite.oauth.(*OAuthMockup).oauthProvider.Add(entrySet)
    TestGetSecrets(&suite.Suite, suite.oauth)
}

func (suite *OAuthMockupTestSuite) TestSetSecrets() {
    entrySet := entities.NewOAuthSecrets(OAuthTestUserID)
    suite.oauth.(*OAuthMockup).oauthProvider.Add(entrySet)

    TestSetSecrets(&suite.Suite, suite.oauth)
}

func (suite *OAuthMockupTestSuite) TestDeleteSecret() {
    entrySet := entities.NewOAuthSecrets(OAuthTestUserID)
    suite.oauth.(*OAuthMockup).oauthProvider.Add(entrySet)

    TestDeleteSecrets(&suite.Suite, suite.oauth)
}


// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOAuthMockup(t *testing.T) {
    suite.Run(t, new(OAuthMockupTestSuite))
}
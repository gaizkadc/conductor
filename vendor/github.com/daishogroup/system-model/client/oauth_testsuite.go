//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the Oauth Client TestSuite

package client

import (
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/entities"
)

const OAuthTestUserID = "testUserID"

func TestSetSecrets(suite *suite.Suite, oauth OAuth) {
    secrets := entities.NewOAuthAddEntryRequest("app1","clientID1", "secret1")
    err := oauth.SetSecret(OAuthTestUserID,secrets)
    suite.Nil(err, "unexpected error")

    // Try to get
    returned, err := oauth.GetSecrets(OAuthTestUserID)
    suite.Nil(err, "unexpected error")
    suite.Equal(OAuthTestUserID, returned.UserID, "unexpected user id")
}

func TestGetSecrets(suite *suite.Suite, oauth OAuth) {
    // Try to get
    returned, err := oauth.GetSecrets(OAuthTestUserID)
    suite.Nil(err, "unexpected error")
    suite.Equal(OAuthTestUserID, returned.UserID, "unexpected user id")
    _,isThere := returned.Entries["app1"]
    suite.True(isThere, "entry was not found")
}

func TestDeleteSecrets(suite *suite.Suite, oauth OAuth) {
    // Delete
    err := oauth.DeleteSecrets(OAuthTestUserID)
    suite.Nil(err, "unexpected error")
    // Try to get
    _, err = oauth.GetSecrets(OAuthTestUserID)
    suite.NotNil(err, "unexpected error")

}
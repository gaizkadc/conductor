//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the credentials.
package client

import (
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/entities"
)



func TestAddCredentials(suite *suite.Suite, cred Credentials) {
    toAdd := entities.NewAddCredentialsRequest("uuid", "pubkey", "privatekey", "desc", "type")
    err := cred.Add(*toAdd)
    suite.Nil(err, "unexpected error")
}


func TestGetCredentials(suite *suite.Suite, cred Credentials) {
    toAdd := entities.NewAddCredentialsRequest("uuid", "pubkey", "privatekey", "desc", "type")
    err := cred.Add(*toAdd)
    suite.Nil(err, "unexpected error")

    // get it
    returned, err2:= cred.Get("uuid")
    suite.Nil(err2, "unexpected error")
    suite.NotNil(returned, "empty response")
}


func TestDeleteCredentials(suite *suite.Suite, cred Credentials) {

    toAdd := entities.NewAddCredentialsRequest("uuid", "pubkey", "privatekey", "desc", "type")
    err := cred.Add(*toAdd)
    suite.Nil(err, "unexpected error")

    // get it
    returned, err2:= cred.Get("uuid")
    suite.Nil(err2, "unexpected error")
    suite.NotNil(returned, "empty response")

    // delete it
    err3 := cred.Delete("uuid")
    suite.Nil(err3, "unexpected error when deleting")

    // Try to retrieve it
    returned2, err4 := cred.Get("uuid")
    suite.Nil(returned2, "the returned was expected to be nil")
    suite.NotNil(err4, "an error was expected")

}
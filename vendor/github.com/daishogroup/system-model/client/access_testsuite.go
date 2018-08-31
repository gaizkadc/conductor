//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the Access Client TestSuite

package client

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/suite"
)

// TestAddUserAccess checks that new roles can be added.
func TestAddUserAccess(suite * suite.Suite, client Access) {
    toAdd := entities.NewAddUserAccessRequest(
        []entities.RoleType{entities.GlobalAdmin})

    added, err := client.AddAccess("user", *toAdd)
    suite.Nil(err, "descriptor should be added")
    suite.NotNil(added, "expecting new descriptor")
}

// TestGetUserAccess checks that a newly added role can be retrieved.
func TestGetUserAccess(suite * suite.Suite, client Access) {
    // add new entry
    toAdd := entities.NewAddUserAccessRequest(
        []entities.RoleType{entities.GlobalAdmin})
    client.AddAccess("user", *toAdd)

    result, err := client.GetAccess("user")
    suite.Nil(err, "error must be nil")
    suite.NotNil(result, "unexpected nil entry")
}

// TestDeleteUserAccess checks that a role can be deleted.
func TestDeleteUserAccess(suite * suite.Suite, client Access) {

    // add new entry
    toAdd := entities.NewAddUserAccessRequest(
        []entities.RoleType{entities.GlobalAdmin})
    client.AddAccess("user", *toAdd)

    err := client.DeleteAccess("user")
    suite.Nil(err, "error must be nil")

    // Check we have already deleted the entry
    result, err := client.GetAccess("user")
    suite.Nil(result, "no entry must be returned")
    suite.NotNil(err, "error message must be returned")

}

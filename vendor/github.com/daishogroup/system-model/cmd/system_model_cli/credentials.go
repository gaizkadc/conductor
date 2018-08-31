//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Credentials line operations.

package main

import (
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

type CredentialsCommand struct {
    GlobalCommand

    addCredUUID *string
    addCredPublicKey *string
    addCredPrivateKey *string
    addCredDescription *string
    addCredTypeKey *string

    getCredUserId *string

    delCredUserId *string

}

func NewCredentialsCommand(app * kingpin.Application, global GlobalCommand) * CredentialsCommand {

    c := &CredentialsCommand{
        GlobalCommand: global,
    }

    credsCommand := app.Command("credentials", "User credentials")

    cmd := credsCommand.Command("add", "Add user credentials.").Action(c.addCredentials)
    c.addCredUUID = cmd.Arg("uuid", "Credential UUID").Required().String()
    c.addCredPublicKey = cmd.Arg("public", "Public user key").Required().String()
    c.addCredPrivateKey = cmd.Arg("private", "Private user key").Required().String()
    c.addCredDescription = cmd.Arg("description", "Credential description").Required().String()
    c.addCredTypeKey = cmd.Arg("typeKey", "Credential type key").Required().String()

    cmd = credsCommand.Command("get", "Get user credentials.").Action(c.getCredentials)
    c.getCredUserId = cmd.Arg("id", "User id").Required().String()

    cmd = credsCommand.Command("delete", "Delete user credentials.").Action(c.deleteCredentials)
    c.delCredUserId = cmd.Arg("id", "User id").Required().String()

    return c
}



func (cmd * CredentialsCommand) addCredentials(c *kingpin.ParseContext) error {
    credRest := client.NewCredentialsClientRest(cmd.IP.String(), *cmd.Port)
    // Build the request
    req := entities.NewAddCredentialsRequest(*cmd.addCredUUID, *cmd.addCredPublicKey, *cmd.addCredPrivateKey,
        *cmd.addCredDescription, *cmd.addCredTypeKey)
    err := credRest.Add(*req)
    return cmd.printResultOrError(entities.NewSuccessfulOperation(errors.CredentialsAdded),err)
}

func (cmd * CredentialsCommand) getCredentials(c *kingpin.ParseContext) error {
    credRest := client.NewCredentialsClientRest(cmd.IP.String(), *cmd.Port)
    creds, err := credRest.Get(*cmd.getCredUserId)
    return cmd.printResultOrError(creds,err)
}

func (cmd * CredentialsCommand) deleteCredentials(c *kingpin.ParseContext) error {
    credRest := client.NewCredentialsClientRest(cmd.IP.String(), *cmd.Port)
    err := credRest.Delete(*cmd.delCredUserId)
    return cmd.printResultOrError(entities.NewSuccessfulOperation(errors.CredentialsRemoved),err)

}
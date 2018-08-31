//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Password related commands.

package main

import (
    "gopkg.in/alecthomas/kingpin.v2"

    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

type PasswordCommand struct {
    GlobalCommand

    setPasswordUserID      * string
    setPasswordPass        * string

    deletePasswordUserID    * string

    getPasswordUserID     *string

}

func NewPassword(app * kingpin.Application, global GlobalCommand) * PasswordCommand {

    c := &PasswordCommand{
        GlobalCommand : global,
    }
    passCommand := app.Command("password", "Password commands")

    cmd := passCommand.Command("set", "Set an already existing user password").Action(c.setPassword)
    c.setPasswordUserID = cmd.Arg("userId", "User identifier").Required().String()
    c.setPasswordPass = cmd.Arg("password", "User password").Required().String()

    cmd = passCommand.Command("delete", "Delete existing password entry").Action(c.deletePassword)
    c.deletePasswordUserID = cmd.Arg("userId", "User identifier").Required().String()

    cmd = passCommand.Command("get", "Get existing password entry").Action(c.getPassword)
    c.getPasswordUserID = cmd.Arg("userId", "User identifier").Required().String()


    return c
}

func (cmd * PasswordCommand) setPassword(c *kingpin.ParseContext) error {
    passRest := client.NewPasswordClientRest(cmd.IP.String(), *cmd.Port)
    toSet, err := entities.NewPassword(*cmd.setPasswordUserID, cmd.setPasswordPass)

    if err!= nil {
        return cmd.printResultOrError("", err)
    }
    err = passRest.SetPassword(*toSet)
    return cmd.printResultOrError(entities.NewSuccessfulOperation(errors.PasswordSet), err)
}

func (cmd * PasswordCommand) getPassword(c *kingpin.ParseContext) error {
    passRest := client.NewPasswordClientRest(cmd.IP.String(), *cmd.Port)
    result, err := passRest.GetPassword(*cmd.getPasswordUserID)

    return cmd.printResultOrError(result, err)
}

func (cmd * PasswordCommand) deletePassword(c *kingpin.ParseContext) error {
    passRest := client.NewPasswordClientRest(cmd.IP.String(), *cmd.Port)
    err := passRest.DeletePassword(*cmd.deletePasswordUserID)
    return cmd.printResultOrError(entities.NewSuccessfulOperation(errors.PasswordDeleted), err)
}

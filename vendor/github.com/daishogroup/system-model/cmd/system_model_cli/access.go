//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Access related commands.

package main

import (
    "gopkg.in/alecthomas/kingpin.v2"

    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

type AccessCommand struct {
    GlobalCommand

    addAccessUserID      * string
    addAccessRole        * string

    setAccessUserID      * string
    setAccessRole        * string

    getAccessUserID     *string

    deleteAccessUserID  *string

}

func NewAccess(app * kingpin.Application, global GlobalCommand) * AccessCommand {

    c := &AccessCommand{
        GlobalCommand : global,
    }
    accCommand := app.Command("access", "Access commands")

    cmd := accCommand.Command("add", "Add new access").Action(c.addAccess)
    c.addAccessUserID = cmd.Arg("userId", "User identifier").Required().String()
    c.addAccessRole = cmd.Arg("role", "User role").Required().Enum(
        entities.AvailableRolesString()...
    )

    cmd = accCommand.Command("set", "Set user roles").Action(c.setAccess)
    c.setAccessUserID = cmd.Arg("userId", "User identifier").Required().String()
    c.setAccessRole = cmd.Arg("role", "User role").Required().Enum(
        //entities.AvailableRolesString()...
    )

    cmd = accCommand.Command("get", "Get user access entry").Action(c.getAccess)
    c.getAccessUserID = cmd.Arg("userId", "User identifier").Required().String()

    cmd = accCommand.Command("delete", "Delete user access entry").Action(c.deleteAccess)
    c.deleteAccessUserID = cmd.Arg("userId", "User identifier").Required().String()

    cmd = accCommand.Command("list", "List all the users with their privileges").Action(c.listAccess)

    return c
}

func (cmd * AccessCommand) addAccess(c *kingpin.ParseContext) error {
    accRest := client.NewAccessClientRest(cmd.IP.String(), *cmd.Port)
    // TODO this should be a list of access privileges
    toAdd := entities.NewAddUserAccessRequest([]entities.RoleType{(entities.RoleType)(*cmd.addAccessRole)})

    result, err := accRest.AddAccess(*cmd.addAccessUserID,*toAdd)

    return cmd.printResultOrError(result, err)
}

func (cmd * AccessCommand) setAccess(c *kingpin.ParseContext) error {
    accRest := client.NewAccessClientRest(cmd.IP.String(), *cmd.Port)
    // TODO this should be a list of access privileges
    toAdd := entities.NewAddUserAccessRequest([]entities.RoleType{(entities.RoleType)(*cmd.setAccessRole)})

    result, err := accRest.SetAccess(*cmd.setAccessUserID,*toAdd)

    return cmd.printResultOrError(result, err)
}


func (cmd * AccessCommand) getAccess(c *kingpin.ParseContext) error {
    accRest := client.NewAccessClientRest(cmd.IP.String(), *cmd.Port)

    result, err := accRest.GetAccess(*cmd.getAccessUserID)

    return cmd.printResultOrError(result, err)
}

func (cmd * AccessCommand) deleteAccess(c *kingpin.ParseContext) error {
    accRest := client.NewAccessClientRest(cmd.IP.String(), *cmd.Port)

    err := accRest.DeleteAccess(*cmd.deleteAccessUserID)

    return cmd.printResultOrError(entities.NewSuccessfulOperation(errors.AccessDeleted), err)
}

func (cmd * AccessCommand) listAccess(c *kingpin.ParseContext) error {
    accRest := client.NewAccessClientRest(cmd.IP.String(), *cmd.Port)

    result,err := accRest.ListAccess()

    return cmd.printResultOrError(result, err)
}

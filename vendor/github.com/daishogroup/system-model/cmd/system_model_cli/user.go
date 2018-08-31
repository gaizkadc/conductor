//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// User command line operations.

package main

import (
    "fmt"
    "time"
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

type UserCommand struct {
    GlobalCommand

    addUserID       *string
    addUserName     *string
    addUserEmail    *string
    addUserPhone    *string
    addUserCreation *string
    addUserExpiration *string

    getUserID       *string

    deleteUserID    *string

    updateUserID       *string
    updateUserName     *string
    updateUserEmail    *string
    updateUserPhone    *string
}

func NewUserCommand(app * kingpin.Application, global GlobalCommand) * UserCommand {

    c := &UserCommand{
        GlobalCommand: global,
    }

    userCommand := app.Command("user", "User commands")

    cmd := userCommand.Command("add", "Add a new user").Action(c.addUser)
    c.addUserID = cmd.Arg("id", "User id").Required().String()
    c.addUserName = cmd.Arg("name", "User name").Required().String()
    c.addUserPhone = cmd.Arg("phone", "User phone").Required().String()
    c.addUserEmail = cmd.Arg("email", "User email").Required().String()
    c.addUserCreation = cmd.Arg("creation", "User creation time using format "+time.RFC3339).Required().String()
    c.addUserExpiration = cmd.Arg("expiration", "User expiration time using format "+time.RFC3339).Required().String()

    cmd = userCommand.Command("get", "Get an existing user").Action(c.getUser)
    c.getUserID = cmd.Arg("id", "User id").Required().String()

    cmd = userCommand.Command("delete", "Delete an existing user").Action(c.deleteUser)
    c.deleteUserID = cmd.Arg("id", "Delete id").Required().String()

    cmd = userCommand.Command("update", "Update an existing user").Action(c.updateUser)
    c.updateUserID = cmd.Arg("id", "User id").Required().String()
    c.updateUserName = cmd.Flag("name", "User name").Default(NotAssigned).String()
    c.updateUserPhone = cmd.Flag("phone", "User phone").Default(NotAssigned).String()
    c.updateUserEmail = cmd.Flag("email", "User email").Default(NotAssigned).String()

    cmd = userCommand.Command("list", "List users").Action(c.listUsers)

    return c
}

func (cmd * UserCommand) addUser(c *kingpin.ParseContext) error {
    userRest := client.NewUserClientRest(cmd.IP.String(),*cmd.Port)
    cTime, err := time.Parse(time.RFC3339,*cmd.addUserCreation)
    if err != nil {
        fmt.Println("Creation time format is not correct.")
        return err
    }
    eTime, err := time.Parse(time.RFC3339,*cmd.addUserExpiration)
    if err != nil {
        fmt.Println("Expiration time format is not correct.")
        return err
    }
    toAdd := entities.NewAddUserRequest(
        *cmd.addUserID, *cmd.addUserName, *cmd.addUserPhone, *cmd.addUserEmail, cTime, eTime)
    added, err2 := userRest.Add(* toAdd)
    return cmd.printResultOrError(added, err2)
}

func (cmd * UserCommand) getUser(c *kingpin.ParseContext) error {
    userRest := client.NewUserClientRest(cmd.IP.String(),*cmd.Port)
    retrieved, err := userRest.Get(*cmd.getUserID)
    return cmd.printResultOrError(retrieved, err)
}

func (cmd * UserCommand) deleteUser(c *kingpin.ParseContext) error {
    userRest := client.NewUserClientRest(cmd.IP.String(),*cmd.Port)
    err := userRest.Delete(*cmd.deleteUserID)

    return cmd.printResultOrError(entities.NewSuccessfulOperation(errors.UserDeleted),err)

}

func (cmd * UserCommand) updateUser(c *kingpin.ParseContext) error {
    userRest := client.NewUserClientRest(cmd.IP.String(),*cmd.Port)
    toUpdate := entities.NewUpdateUserRequest()
    if *cmd.updateUserName != NotAssigned && cmd.updateUserName != nil{
        toUpdate.WithName(*cmd.updateUserName)
    }
    if *cmd.updateUserPhone != NotAssigned && cmd.updateUserPhone != nil{
        toUpdate.WithPhone(*cmd.updateUserPhone)
    }
    if *cmd.updateUserEmail != NotAssigned && cmd.updateUserEmail != nil{
        toUpdate.WithEmail(*cmd.updateUserEmail)
    }

    updated, err := userRest.Update(*cmd.updateUserID, *toUpdate)
    return cmd.printResultOrError(updated, err)
}

func (cmd * UserCommand) listUsers(c *kingpin.ParseContext) error {
    userRest := client.NewUserClientRest(cmd.IP.String(),*cmd.Port)
    users, err := userRest.ListUsers()

    return cmd.printResultOrError(users,err)

}

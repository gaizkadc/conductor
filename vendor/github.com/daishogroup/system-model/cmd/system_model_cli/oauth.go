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

type OauthCommand struct {
    GlobalCommand

    setOAuthUserID      * string
    setOAuthApp         * string
    setOAuthClientID    * string
    setOAuthSecret      * string

    deleteOAuthUserID   * string

    getOauthUserID      * string

}

func NewOauth(app * kingpin.Application, global GlobalCommand) *OauthCommand {

    c := &OauthCommand{
        GlobalCommand : global,
    }
    oauthCommand := app.Command("oauth", "Oauth commands")

    cmd := oauthCommand.Command("set", "Set a new application secret for an existing user").Action(c.setSecret)
    c.setOAuthUserID = cmd.Arg("userId", "User ID").Required().String()
    c.setOAuthApp = cmd.Arg("app", "Application name").Required().String()
    c.setOAuthClientID = cmd.Arg("clientId", "Application client id").Required().String()
    c.setOAuthSecret = cmd.Arg("secret", "Application client secret").Required().String()

    cmd = oauthCommand.Command("delete", "Delete user oauth entries").Action(c.deleteSecrets)
    c.deleteOAuthUserID = cmd.Arg("userId", "User identifier").Required().String()

    cmd = oauthCommand.Command("get", "Get existing user oauth entries").Action(c.getSecrets)
    c.getOauthUserID = cmd.Arg("userId", "User identifier").Required().String()


    return c
}

func (cmd *OauthCommand) setSecret(c *kingpin.ParseContext) error {
    rest := client.NewOAuthClientRest(cmd.IP.String(), *cmd.Port)
    toSet := entities.NewOAuthAddEntryRequest(*cmd.setOAuthApp, *cmd.setOAuthClientID, *cmd.setOAuthSecret)
    err := rest.SetSecret(*cmd.setOAuthUserID, toSet)
    return cmd.printResultOrError(entities.NewSuccessfulOperation(errors.OAuthEntrySet), err)
}

func (cmd *OauthCommand) getSecrets(c *kingpin.ParseContext) error {
    rest := client.NewOAuthClientRest(cmd.IP.String(), *cmd.Port)
    result, err := rest.GetSecrets(*cmd.getOauthUserID)

    return cmd.printResultOrError(result, err)
}

func (cmd *OauthCommand) deleteSecrets(c *kingpin.ParseContext) error {
    rest := client.NewOAuthClientRest(cmd.IP.String(), *cmd.Port)
    err := rest.DeleteSecrets(*cmd.deleteOAuthUserID)

    return cmd.printResultOrError(entities.NewSuccessfulOperation(errors.OAuthUserDeleted), err)
}


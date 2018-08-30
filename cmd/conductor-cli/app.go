//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Source code to operate with conductor service using command line

package main

import (

    "github.com/daishogroup/conductor/client"
    "github.com/daishogroup/system-model/entities"
    "gopkg.in/alecthomas/kingpin.v2"
    entitiesConductor "github.com/daishogroup/conductor/entities"
    "github.com/daishogroup/conductor/errors"
)

type orchestratorCommand struct {
    globalCommand

    deployNetworkId      *string
    deployAppDescId      *string
    deployAppLabel       *string
    deployAppArguments   *string
    deployAppDescription *string
    deployAppName        *string
    deployAppPersistence *string
    deployAppStorageType *string

    undeployNetworkId  *string
    undeployInstanceId *string

    logsNetworkId  *string
    logsInstanceId *string
}

// Initialize deploy command
func newDeployCommand(app *kingpin.Application, global globalCommand) *orchestratorCommand {
    c := &orchestratorCommand{
        globalCommand: global,
    }
    orchCommand := app.Command("orchestrator", "Application deployment commands")

    // Add new application commands here ...
    // app deploy
    cmdDeploy := orchCommand.Command("deploy", "Deploy an application on a given network").Action(c.deploy)
    c.deployNetworkId = cmdDeploy.Arg("network.id", "Network ID").Required().String()
    c.deployAppDescId = cmdDeploy.Arg("app.descriptor", "Application descriptor id").Required().String()
    c.deployAppLabel = cmdDeploy.Arg("app.label", "Application label. Typically edge, gateway or cloud").Required().String()
    c.deployAppName = cmdDeploy.Arg("app.name", "Application name").Required().String()
    c.deployAppArguments = cmdDeploy.Arg("app.arguments", "Application arguments").Default("").String()
    c.deployAppDescription = cmdDeploy.Arg("app.description", "Application description").String()
    c.deployAppPersistence = cmdDeploy.Arg("app.persistenceSize", "Persistence size").Default("1Gb").String()
    c.deployAppStorageType = cmdDeploy.Arg("app.storage", "Storage type to be used").
        Default("default").String()

    cmdUndeploy := orchCommand.Command("undeploy", "Undeploy and application on a given network").Action(c.undeploy)
    c.undeployNetworkId = cmdUndeploy.Arg("network.id", "Network ID").Required().String()
    c.undeployInstanceId = cmdUndeploy.Arg("instance.id", "Application instance ID").Required().String()

    cmdLogs := orchCommand.Command("logs", "Get the list of logs.").Action(c.logs)
    c.logsNetworkId = cmdLogs.Arg("network.id", "Network ID").Required().String()
    c.logsInstanceId = cmdLogs.Arg("instance.id", "Application instance ID").Required().String()
    return c
}

func (cmd *orchestratorCommand) deploy(c *kingpin.ParseContext) error {
    conductorRest := client.NewConductorRest(cmd.Ip.String(),*cmd.Port)

    // storage type
    // TODO Control the values in this app storage type to be under the defined range
    storageType := entities.AppStorageType(*cmd.deployAppStorageType)

    deployAppRequest := entitiesConductor.NewDeployAppRequest(
        *cmd.deployAppName, *cmd.deployAppDescId, *cmd.deployAppDescription,
        *cmd.deployAppLabel, make(map[string]string, 0),
        *cmd.deployAppArguments, *cmd.deployAppPersistence, storageType)

    response, err := conductorRest.Deploy(*cmd.deployNetworkId, deployAppRequest)
    return cmd.printResultOrError(response, err)
}

func (cmd *orchestratorCommand) undeploy(c *kingpin.ParseContext) error {
    conductorRest := client.NewConductorRest(cmd.Ip.String(),*cmd.Port)

    err := conductorRest.Undeploy(*cmd.undeployNetworkId, *cmd.undeployInstanceId)
    return cmd.printResultOrError(entities.NewSuccessfulOperation(errors.UndeploySuccess),err)
}

func (cmd *orchestratorCommand) logs(c *kingpin.ParseContext) error {
    conductorRest := client.NewConductorRest(cmd.Ip.String(),*cmd.Port)
    return cmd.printResultOrError(conductorRest.Logs(*cmd.logsNetworkId, *cmd.logsInstanceId))

}

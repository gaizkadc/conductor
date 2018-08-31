//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package main

import (
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/daishogroup/system-model/client"
)

// InfoCommand contains all parameters and arguments for the info commands.
type InfoCommand struct {
    GlobalCommand

    reducedByNetworkID *string
}

// NewInfoCommand is the method to generate the commands.
func NewInfoCommand(app *kingpin.Application, global GlobalCommand) *InfoCommand {
    c := &InfoCommand{
        GlobalCommand: global,
    }

    infoCmd := app.Command("info", "Info commands")
    infoCmd.Command("reduced", "Get the essential info of the system model").Action(c.reduced)

    infoCmd.Command("summary", "Get the number of stored elements").Action(c.summary)

    cmd := infoCmd.Command("reduced-by-network", "Get the essential info of the system model by network").
        Action(c.reducedByNetwork)
    c.reducedByNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    return c
}

func (cmd *InfoCommand) reduced(c *kingpin.ParseContext) error {
    infoRest := client.NewInfoClientRest(cmd.IP.String(), *cmd.Port)
    info, err := infoRest.ReducedInfo()
    return cmd.printResultOrError(info,err)
}

func (cmd *InfoCommand) summary(c *kingpin.ParseContext) error {
    infoRest := client.NewInfoClientRest(cmd.IP.String(), *cmd.Port)
    info, err := infoRest.SummaryInfo()
    return cmd.printResultOrError(info,err)
}

func (cmd *InfoCommand) reducedByNetwork(c *kingpin.ParseContext) error {
    infoRest := client.NewInfoClientRest(cmd.IP.String(), *cmd.Port)
    info, err := infoRest.ReducedInfoByNetwork(*cmd.reducedByNetworkID)
    return cmd.printResultOrError(info,err)
}

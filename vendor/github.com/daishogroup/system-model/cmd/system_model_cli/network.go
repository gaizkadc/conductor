//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Network related commands.

package main

import (
    "gopkg.in/alecthomas/kingpin.v2"

    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
)

type NetworkCommand struct {
    GlobalCommand

    addNetworkName * string
    addNetworkDescription * string
    addNetworkAdminName * string
    addNetworkAdminPhone * string
    addNetworkAdminEmail * string

    getNetworkID * string

    deleteNetworkID * string
}

func NewNetworkCommand(app * kingpin.Application, global GlobalCommand) * NetworkCommand {

    c := &NetworkCommand{
        GlobalCommand : global,
    }

    networkCmd := app.Command("network", "Network commands")

    cmd := networkCmd.Command("add", "Add a new network").Action(c.addNetwork)
    c.addNetworkName = cmd.Arg("name", "Name of the network to add").Required().String()
    c.addNetworkDescription= cmd.Flag("desc", "Description of the network to add").Default("").String()
    c.addNetworkAdminName  = cmd.Flag("admin.name", "Administrator's name").Default("admin").String()
    c.addNetworkAdminPhone = cmd.Flag("admin.phone", "Administrator's phone number").Default("").String()
    c.addNetworkAdminEmail = cmd.Flag("admin.email", "Administrator's email").Default("").String()

    cmd = networkCmd.Command("list", "List all networks").Action(c.listNetworks)

    cmd = networkCmd.Command("get", "Get a network").Action(c.getNetwork)
    c.getNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()

    cmd = networkCmd.Command("delete", "Delete a network").Action(c.deleteNetwork)
    c.deleteNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()

    return c
}

func (cmd * NetworkCommand) addNetwork(c *kingpin.ParseContext) error {
    netRest := client.NewNetworkClientRest(cmd.IP.String(), *cmd.Port)
    toAdd := entities.NewAddNetworkRequest(* cmd.addNetworkName, * cmd.addNetworkDescription,
        * cmd.addNetworkAdminName, * cmd.addNetworkAdminPhone, * cmd.addNetworkAdminEmail)
    added, err := netRest.Add(* toAdd)
    return cmd.printResultOrError(added, err)
}

func (cmd * NetworkCommand) listNetworks(c *kingpin.ParseContext) error {
    netRest := client.NewNetworkClientRest(cmd.IP.String(), *cmd.Port)
    networks, err := netRest.List()
    return cmd.printResultOrError(networks, err)
}

func (cmd * NetworkCommand) getNetwork(c *kingpin.ParseContext) error {
    netRest := client.NewNetworkClientRest(cmd.IP.String(), *cmd.Port)
    network, err := netRest.Get(* cmd.getNetworkID)
    return cmd.printResultOrError(network, err)
}

func (cmd * NetworkCommand) deleteNetwork(c *kingpin.ParseContext) error {
    netRest := client.NewNetworkClientRest(cmd.IP.String(), *cmd.Port)
    err := netRest.Delete(* cmd.deleteNetworkID)
    return cmd.printResultOrError(entities.NewSuccessfulOperation("Network deleted"), err)
}
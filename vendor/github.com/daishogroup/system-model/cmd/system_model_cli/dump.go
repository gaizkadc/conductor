//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Dump command.

package main

import (
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/daishogroup/system-model/client"
)

// DumpCommand contains all parameters and arguments for the dump commands.
type DumpCommand struct {
    GlobalCommand
}

func NewDumpCommand(app *kingpin.Application, global GlobalCommand) * DumpCommand {
    c := &DumpCommand{
        GlobalCommand: global,
    }

    dumpCmd := app.Command("dump", "Dump commands")
    dumpCmd.Command("export", "Export all system model entities").Action(c.export)

    return c
}

func (cmd * DumpCommand) export(c *kingpin.ParseContext) error {
    dumpRest := client.NewDumpClientRest(cmd.IP.String(), *cmd.Port)
    dump, err := dumpRest.Export()
    return cmd.printResultOrError(dump, err)
}
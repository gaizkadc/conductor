//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Main entry command line to run the integration for conductor.

package main

import (
    "os"
    log "github.com/sirupsen/logrus"
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/nalej/conductor/integration"
)



//Main application command.
var app = kingpin.New("conductor-integration", "Conductor integration analysis.")

// Demo API command.
var (
    startCommand = app.Command("start", "Launch the service API")

    FactoryTargetPath = startCommand.Flag(
        "factoryTargetPath",
        "Path were the daisho factory targets are available.",
    ).Default("~/daisho/factory/target").String()

    AsmPath = startCommand.Flag(
        "asmPath",
        "ASM path folder",
    ).Default("~/daisho/appmgr/bazel-bin/asmcli").String()

    ColonyPath = startCommand.Flag(
        "colonyPath",
        "Colony folder were the binary is allocated",
    ).Default("~/daisho/colony/colonize").String()

    AppdevkitPath = startCommand.Flag(
        "appdevkitPath",
        "Folder were we can find appdevkit",
    ).Default("~/daisho/appdevkit").String()

    debugFlag = startCommand.Flag("debug", "Activate debug logging").Default("false").Bool()
)


// Enable debug log level.
func enableDebug() {
    log.SetLevel(log.DebugLevel)
}

func main() {
    var err error = nil
    switch kingpin.MustParse(app.Parse(os.Args[1:])) {
    case startCommand.FullCommand():
        if *debugFlag {
            enableDebug()
        }
        configuration := integration.Config{FactoryTargetPath: *FactoryTargetPath,
        AsmPath: *AsmPath, ColonyPath: *ColonyPath, AppdevKitPath: *AppdevkitPath}

        integrator := integration.NewIntegration(configuration)
        integrator.RunIntegration()
    }
    if err != nil {
        os.Exit(1)
    }
}

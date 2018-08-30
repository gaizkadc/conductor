//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Main entry to run the service using command line.

package main

import (
    "os"
    log "github.com/sirupsen/logrus"
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/daishogroup/conductor/server"
    "github.com/daishogroup/service"
)

// Logger for the main package.
var logger = log.WithField("package", "main").WithField("file", "conductor.go")

//AppDevKit application command.
var app = kingpin.New("conductor", "Conductor for deploying applications on top of Daisho.")

// Demo API command.
var (
    apiCommand = app.Command("api", "Launch the service API")

    portFlag = apiCommand.Flag(
        "port",
        "Local port where the CONDUCTOR API listens to incoming requests.",
    ).Default("9000").Uint16()

    systemModelAddressFlag = apiCommand.Flag(
        "system-model-address",
        "Address and port where the system model API server listens to requests.",
    ).Default("http://localhost:8800").String()

    loggerAddressFlag = apiCommand.Flag(
        "logger-address",
        "Address and port where the logger API server listens to requests.",
    ).Default("http://localhost:8083").String()

    debugFlag = apiCommand.Flag("debug", "Activate debug logging").Default("false").Bool()
)

// Version command.
var (
    versionCommand = app.Command("version", "View version of the code")
)

// Function to launch the Demo API.
func launchApi() error {
    configApi := server.Config{Port: *portFlag,
        SystemModelAddress: *systemModelAddressFlag,
        LoggerAddress: *loggerAddressFlag}
    serviceApi := server.Service{Configuration: configApi}
    return service.Launch(&serviceApi)
}

// Enable debug log level.
func enableDebug() {
    log.SetLevel(log.DebugLevel)
}

func main() {
    var err error = nil
    switch kingpin.MustParse(app.Parse(os.Args[1:])) {
    case apiCommand.FullCommand():
        if *debugFlag {
            enableDebug()
        }
        logger.Info("Launch Api Command.")
        err = launchApi()
    case versionCommand.FullCommand():
        logger.Info("Version 0.1.0.")
        err = nil
    }
    if err != nil {
        logger.Error("Unexpected error in the application.", err)
        os.Exit(1)
    }
}

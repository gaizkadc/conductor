//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the System Model server app

package main

import (
    log "github.com/sirupsen/logrus"
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/daishogroup/system-model/server"
    "os"
    "github.com/daishogroup/service"
)

// Logger for the main package.
var loggerSystemModel = log.WithField("package", "main").WithField("file", "system-model.go")

//AppDevKit application command.
var app = kingpin.New("system-model", "System Model application.")

// Demo API command.
var (
    apiCommand = app.Command("api", "Launch the System Model Service")

    portFlag = apiCommand.Flag(
        "port",
        "Local port where the System model listens to incoming requests.",
    ).Default("8800").Uint16()

    inMemoryPersistenceFlag       = apiCommand.Flag("in-memory-persistence", "Use in-memory providers").Default("true").Bool()
    fileSystemPersistenceFlag     = apiCommand.Flag("filesystem-persistence", "Use file backed persistence providers").Default("false").Bool()
    fileSystemPersistenceBasePath = apiCommand.Flag("filesystem-basePath", "Base path to store the system model data").Default("/tmp/").String()
    defaultAdminUser = apiCommand.Flag("default-admin-user", "Name of the basic default admin user").Default("admin").String()
    defaultAdminPassword = apiCommand.Flag("default-admin-password", "Password for the default admin user").Default("admin").String()

    debugFlag = apiCommand.Flag("debug", "Activate debug logging").Default("false").Bool()
)

// Version command.
var (
    versionCommand = app.Command("version", "View version of the code")
)

// Function to launch the Demo API.
func launchAPI() error {
    configAPI := server.Config{
        Port:                  *portFlag,
        UseInMemoryProviders:  * inMemoryPersistenceFlag,
        UseFileSystemProvider: * fileSystemPersistenceFlag,
        FileSystemBasePath:    * fileSystemPersistenceBasePath,
        DefaultAdminUser:      * defaultAdminUser,
        DefaultAdminPassword:  * defaultAdminPassword}

    err := configAPI.Validate()
    if err != nil {
        return err
    }
    serviceAPI := server.Service{Configuration: configAPI}
    return service.Launch(&serviceAPI)
}

// Enable debug log level.
func enableDebug() {
    log.SetLevel(log.DebugLevel)
}

// ApiDevKit entry point.
func main() {
    var err error
    switch kingpin.MustParse(app.Parse(os.Args[1:])) {
    case apiCommand.FullCommand():
        if *debugFlag {
            enableDebug()
        }
        loggerSystemModel.Info("Launch Api Command.")
        err = launchAPI()
    case versionCommand.FullCommand():
        loggerSystemModel.Info("Version 1.0.1")
        err = nil
    }
    if err != nil {
        loggerSystemModel.Error("Unexpected error in the application.", err)
        os.Exit(1)
    }
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Main file for the system model command line interface.

package main

import (
log "github.com/sirupsen/logrus"
"gopkg.in/alecthomas/kingpin.v2"
"os"
)

func main() {

    var (
        app = kingpin.New("system-model-cli", "Command line tool for Daisho System Model").DefaultEnvars()
    )

    globalCfg := NewGlobalCommand(app)
    NewNetworkCommand(app, * globalCfg)
    NewClusterCommand(app, * globalCfg)
    NewNodeCommand(app, * globalCfg)
    NewApplicationCommand(app, * globalCfg)
    NewDumpCommand(app, * globalCfg)
    NewInfoCommand(app, * globalCfg)
    NewUserCommand(app, * globalCfg)
    NewAccess(app, * globalCfg)
    NewPassword(app, * globalCfg)
    NewOauth(app, * globalCfg)
    NewCredentialsCommand(app, * globalCfg)
    kingpin.MustParse(app.Parse(os.Args[1:]))
}

func init() {
    log.SetOutput(os.Stdout)
    log.SetLevel(log.DebugLevel)
}

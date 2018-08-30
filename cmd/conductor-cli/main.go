//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)


// Logger for the main package.
var logger = log.WithField("package", "main").WithField("file","main.go")


func main() {

	var (
		app = kingpin.New("conductor-cli", "Command line tool for Daisho conductor.").DefaultEnvars()
	)

	globalCfg := newGlobalCommand(app)
	newDeployCommand(app, *globalCfg)


	kingpin.MustParse(app.Parse(os.Args[1:]))

}


func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

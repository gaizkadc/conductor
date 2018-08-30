//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/daishogroup/derrors"
	"gopkg.in/alecthomas/kingpin.v2"
)

type globalCommand struct {
	Ip     *net.IP
	Port   *int
	// Whether detailed responses should be printed.
	Debug * bool
}

func newGlobalCommand(app *kingpin.Application) *globalCommand {
	c := &globalCommand{}

	c.Ip   = app.Flag("ip", "Ip address of conductor service.").Default("127.0.0.1").IP()
	c.Port = app.Flag("port", "Port number of conductor service.").Default("9000").Int()
	c.Debug = app.Flag("debug", "Print detailed responses").Default("false").Bool()

	return c
}

func (c *globalCommand) getBasePath() string{
	base :=  fmt.Sprintf("http://%s:%d/", *c.Ip, *c.Port)
	return base
}

func (cmd * globalCommand) printResultOrError(result interface {}, err derrors.DaishoError) error {
	if err != nil {
		if * cmd.Debug {
			fmt.Println(err.DebugReport())
		}else{
			fmt.Println(err.Error())
		}
		return nil
	}else{
		return cmd.printResult(result)
	}
}

// Output the command result.
//   params:
//     result A JSON object.
//   returns:
//     An error if the JSON processing fails.
func (cmd * globalCommand) printResult(result interface{}) error {
	//Print descriptors
	res, err := json.MarshalIndent(result,"","  ")
	if err == nil {
		fmt.Println(string(res))
	}else{
		fmt.Println("Error found in printResult: " + err.Error())
	}
	return err
}
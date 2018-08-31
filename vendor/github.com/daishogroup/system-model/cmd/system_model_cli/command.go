//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Generic command to interact with the system model.

package main

import (
    "encoding/json"
    "fmt"
    "github.com/daishogroup/derrors"
    "gopkg.in/alecthomas/kingpin.v2"
    "net"
)

// NotAssigned represents the nil value for update commands.
const NotAssigned = "#!@NotAssigned"

const (
    NetworkID        = "networkId"
    NetworkIDDesc    = "ID of the network"
    ClusterID        = "clusterId"
    ClusterIDDesc    = "ID of the cluster"
    NodeID           = "nodeId"
    NodeIDDesc       = "ID of the node"
    DescriptorID     = "descriptorId"
    DescriptorIDDesc = "Descriptor identifier"
    DeployedID       = "deployedId"
    DeployedIDDesc   = "Instance identifier"
)

// Struct that defines the elements required to connect to the system model.
type GlobalCommand struct {
    // System model IP.
    IP * net.IP
    // System model port
    Port * int
    // Whether detailed responses should be printed.
    Debug * bool
}

// Generate a new command.
//   params:
//     app The cli App.
//   returns:
//     A global command.
func NewGlobalCommand(app * kingpin.Application) * GlobalCommand {
    c := &GlobalCommand{}

    c.IP = app.Flag("ip", "IP address of Daisho System Model").Default("127.0.0.1").IP()
    c.Port = app.Flag("port", "Port of Daisho System Model").Default("8800").Int()
    c.Debug = app.Flag("debug", "Print detailed responses").Default("false").Bool()

    return c
}



func (cmd * GlobalCommand) printResultOrError(result interface {}, err derrors.DaishoError) error {
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
func (cmd * GlobalCommand) printResult(result interface{}) error {
    //Print descriptors
    res, err := json.MarshalIndent(result,"","  ")
    if err == nil {
        fmt.Println(string(res))
    }else{
        fmt.Println("Error found in printResult: " + err.Error())
    }
    return err
}


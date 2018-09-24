//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package main

import (
    "github.com/nalej/conductor/cmd/conductor/cmd"
    "github.com/nalej/golang-template/version"
)


var MainVersion string
var MainCommit string


func main() {
    version.AppVersion = MainVersion
    version.Commit = MainCommit
    cmd.Execute()
}
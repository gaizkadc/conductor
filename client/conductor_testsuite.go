//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Testsuite for conductor

package client

import (
	"github.com/stretchr/testify/suite"
    "github.com/daishogroup/conductor/entities"
)


func TestDeploy(suite * suite.Suite, client Conductor){
    app, err := client.Deploy("1", entities.DeployAppRequest{
        Label: "label",
        Arguments: "arguments",
        AppDescriptorId: "1",
        Description: "description",
        Name: "name",
    })
    suite.Nil(err, "There was an error deploying the cluster")
    suite.NotNil(app, "Deployed instance was nil")
}

func TestUndeploy(suite * suite.Suite, client Conductor){
    err := client.Undeploy("1","1")
    suite.Nil(err, "Error undeploying instance")
}
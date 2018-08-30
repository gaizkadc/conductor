//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package entities

import (
	"strings"
	"testing"

	"github.com/daishogroup/system-model/entities"
	"github.com/stretchr/testify/assert"
)

func getRequest(name string) DeployAppRequest {
	return NewDeployAppRequest(name, "appDescriptorId", "description", "",
		make(map[string]string, 0), "args", "1Gb", entities.AppStorageDefault)
}

func TestValid(t *testing.T) {
	deployRequest := getRequest("proper-name")
	assert.Nil(t, deployRequest.IsValid())
}

func TestMaxName(t *testing.T) {
	name := strings.Repeat("n", MaxDeployNameLength+1)
	deployRequest := getRequest(name)
	assert.NotNil(t, deployRequest.IsValid())
}

func TestNoUpperCase(t *testing.T) {
	deployRequest := getRequest("Name")
	assert.NotNil(t, deployRequest.IsValid())
}

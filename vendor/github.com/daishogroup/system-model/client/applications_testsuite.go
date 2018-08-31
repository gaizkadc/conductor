//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Application Client TestSuite

package client

import (
    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/entities"
)

const (
    // TestDescriptorID with the descriptor identifier.
    TestDescriptorID     = "testDescriptorId"
    // TestAppNetworkID with the network identifier.
    TestAppNetworkID     = "networkId"
    // TestDescription with a test value for description.
    TestDescription      = "description"
    // TestName with a test value for the name.
    TestName             = "name"
    // TestLabel with a test value for label.
    TestLabel            = "label"
    // TestArguments with a test value for arguments.
    TestArguments        = "argument"
    // TestAppDescName with a test descriptor name.
    TestAppDescName      = "descName"
    // TestServiceName with a test value for service names.
    TestServiceName      = "service"
    // TestServiceVersion with a test value for service versions.
    TestServiceVersion   = "version"
    // TestPersistenceSize with a test persistence size.
    TestPersistenceSize  = "1Gb"
    // TestStorageType with a default storage.
    TestStorageType      = entities.AppStorageDefault
    // TestInstanceID with a test application instance identifier.
    TestInstanceID       = "instanceId"
    // TestInstanceIDUpdate with the ID to be used on update tests.
    TestInstanceIDUpdate = "instanceToUpdate"
    // TestNetworkIDDelete with the ID to be used on delete tests.
    TestNetworkIDDelete  = "networkDelete"
    // TestInstanceIDDelete with the ID to be used on delete tests.
    TestInstanceIDDelete = "instanceToDelete"
    // TestPort with a test value for port.
    TestPort             = 0
    TestImage            = "nginx:1.12-alpine"
)

// TestAddDescriptor checks that a new descriptor can be added.
func TestAddDescriptor(suite * suite.Suite, client Applications) {
    toAdd := entities.NewAddAppDescriptorRequest(
        TestAppDescName, TestDescription, TestServiceName, TestServiceVersion, TestLabel, TestPort, []string{TestImage})
    added, err := client.AddApplicationDescriptor(TestAppNetworkID, *toAdd)
    suite.Nil(err, "descriptor should be added")
    suite.NotNil(added, "expecting new descriptor")
}

// TestListDescriptors checks that existing descriptors can be listed.
func TestListDescriptors(suite * suite.Suite, client Applications) {
    retrieved, err := client.ListDescriptors(TestAppNetworkID)
    suite.Nil(err, "list should be returned")
    suite.Equal(1, len(retrieved), "expecting 1 descriptor")
}

// TestGetDescriptor checks that an existing descriptor can be retrieved.
func TestGetDescriptor(suite * suite.Suite, client Applications) {
    retrieved, err := client.GetDescriptor(TestAppNetworkID, TestDescriptorID)
    suite.Nil(err, "descriptor should be returned")
    suite.NotNil(retrieved, "descriptor should be returned")
}

func TestDeleteDescriptor(suite * suite.Suite, client Applications){
    err := client.DeleteDescriptor(TestAppNetworkID, TestDescriptorID)
    suite.Nil(err, "descriptor should be deleted")
}

// TestAddInstance checks that a new application instance can be added.
func TestAddInstance(suite * suite.Suite, client Applications) {
    toAdd := entities.NewAddAppInstanceRequest(
        TestDescriptorID, TestName, TestDescription, TestLabel, TestArguments,
        TestPersistenceSize, TestStorageType)
    added, err := client.AddApplicationInstance(TestAppNetworkID, *toAdd)
    suite.Nil(err, "instance should be added")
    suite.NotNil(added, "expecting new instance")
}

// TestListInstances checks that existing instances can be listed.
func TestListInstances(suite * suite.Suite, client Applications){
    retrieved, err := client.ListInstances(TestAppNetworkID)
    suite.Nil(err, "list should be returned")
    suite.Equal(2, len(retrieved), "expecting 2 instances")
}

// TestGetInstance checks that an existing application instance can be retrieved.
func TestGetInstance(suite * suite.Suite, client Applications){
    retrieved, err := client.GetInstance(TestAppNetworkID, TestInstanceID)
    suite.Nil(err, "instance should be returned")
    suite.NotNil(retrieved, "instance should be returned")
}

// TestUpdateInstance checks that an existing application instance can be updated.
func TestUpdateInstance(suite * suite.Suite, client Applications){
    retrieved, err := client.GetInstance(TestAppNetworkID, TestInstanceIDUpdate)
    suite.Nil(err, "instance should be returned")
    updateRequest := entities.NewUpdateAppInstRequest().WithDescription("new description")
    updated, err := client.UpdateInstance(TestAppNetworkID, retrieved.DeployedID, * updateRequest)
    suite.Nil(err, "instance should be returned")
    suite.NotNil(updated, "instance should be returned")
    suite.Equal("new description", updated.Description)
}

// TestDeleteInstance checks that an existing instance can be updated.
func TestDeleteInstance(suite * suite.Suite, client Applications){
    err := client.DeleteInstance(TestNetworkIDDelete, TestInstanceIDDelete)
    suite.Nil(err, "instance should be deleted")
}

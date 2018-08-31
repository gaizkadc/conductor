//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package app

import (

    "strconv"
    "testing"

    "github.com/stretchr/testify/suite"

    "github.com/daishogroup/system-model/entities"

    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/networkstorage"

    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/stretchr/testify/assert"
)

const(
    testInstName = "instName"
    testLabel = "label"
    testArguments = "--argument"
    testNetworkName    = "networkName"
    testDescription    = "description"
    testAdminName      = "adminName"
    testAdminPhone     = "adminPhone"
    testAdminEmail     = "adminEmail"
    testAppDescName    = "descName"
    testServiceName    = "service"
    testServiceVersion = "version"
    testPort           = 0
    testImage          = "nginx:1.12-alpine"
)

type ManagerHelper struct {
    networkProvider * networkstorage.MockupNetworkProvider
    appDescProvider * appdescstorage.MockupAppDescProvider
    appInstanceProvider * appinststorage.MockupAppInstProvider
    appManager Manager
}

func NewManagerHelper() ManagerHelper {
    networkProvider := networkstorage.NewMockupNetworkProvider()
    appInstanceProvider := appinststorage.NewMockupAppInstProvider()
    appDescriptorProvider := appdescstorage.NewMockupAppDescProvider()
    var instanceManager = NewManager(networkProvider, appDescriptorProvider, appInstanceProvider)
    return ManagerHelper{networkProvider, appDescriptorProvider, appInstanceProvider, instanceManager}
}

type TestHelper struct {
    suite.Suite
    manager ManagerHelper
}

func (helper * TestHelper) SetupSuite(){
    managerHelper := NewManagerHelper()
    helper.manager = managerHelper
}

func (helper * TestHelper) SetupTest() {
    helper.manager.networkProvider.Clear()
    helper.manager.appDescProvider.Clear()
    helper.manager.appInstanceProvider.Clear()
}

func TestManagerSuite(t *testing.T) {
    suite.Run(t, new(TestHelper))
}

func (helper *TestHelper) addTestingNetwork(id string) {
    var toAdd= entities.NewNetworkWithID(
        id, testNetworkName, testDescription, testAdminName, testAdminPhone, testAdminEmail)
    helper.manager.networkProvider.Add(*toAdd)
}

func (helper *TestHelper) getTestAppDescriptor() *entities.AddAppDescriptorRequest {
    return entities.NewAddAppDescriptorRequest(
        testAppDescName, testDescription, testServiceName, testServiceVersion, testLabel, testPort, []string{testImage})
}

func (helper * TestHelper) getTestAppInstance(networkID string) * entities.AddAppInstanceRequest{
    return entities.NewAddAppInstanceRequest(networkID,
        testInstName, testDescription, testLabel, testArguments,
        "1Gb", entities.AppStorageDefault)
}

func (helper *TestHelper) TestAddDescriptor() {
    networkID := "testAddDescriptor"
    helper.addTestingNetwork(networkID)
    newDescriptor := helper.getTestAppDescriptor()
    added, err := helper.manager.appManager.AddApplicationDescriptor(networkID, *newDescriptor)
    helper.Nil(err, "descriptor should be added")
    helper.NotNil(added, "descriptor should be returned")
}

func (helper *TestHelper) TestRetrieveDescriptor() {
    networkID := "testRetrieveDescriptor"
    helper.addTestingNetwork(networkID)
    newDescriptor := helper.getTestAppDescriptor()
    added, err := helper.manager.appManager.AddApplicationDescriptor(networkID, *newDescriptor)
    helper.Nil(err, "descriptor should be added")
    helper.NotNil(added, "descriptor should be returned")
    retrieved, err := helper.manager.appManager.GetDescriptor(networkID, added.ID)
    helper.Nil(err, "descriptor should be retrieved")
    helper.EqualValues(added, retrieved, "structs should match")
}

func (helper *TestHelper) TestDeleteDescriptor() {
    networkID := "testDeleteDescriptor"
    helper.addTestingNetwork(networkID)
    newDescriptor := helper.getTestAppDescriptor()
    added, err := helper.manager.appManager.AddApplicationDescriptor(networkID, *newDescriptor)
    helper.Nil(err, "descriptor should be added")
    helper.NotNil(added, "descriptor should be returned")
    err = helper.manager.appManager.DeleteDescriptor(networkID, added.ID)
    helper.Nil(err, "descriptor should be deleted")
}

func (helper *TestHelper) TestListDescriptors() {
    networkID := "testListDescriptors"
    helper.addTestingNetwork(networkID)
    numberDescriptors := 5
    for i := 0; i < numberDescriptors; i++ {
        newDescriptor := helper.getTestAppDescriptor()
        _, err := helper.manager.appManager.AddApplicationDescriptor(networkID, *newDescriptor)
        assert.Nil(helper.T(), err, "descriptor should be added")
    }
    descriptors, err := helper.manager.appManager.ListDescriptors(networkID)
    assert.Nil(helper.T(), err, "list should be retrieved")
    assert.Equal(helper.T(), numberDescriptors, len(descriptors),
        "expecting "+strconv.Itoa(numberDescriptors)+" descriptors")
}

func (helper *TestHelper) TestListDescEmptyNetwork() {
    networkID := "testListDescEmptyNetwork"
    helper.addTestingNetwork(networkID)
    descriptors, err := helper.manager.appManager.ListDescriptors(networkID)
    helper.Nil(err, "list should be retrieved")
    helper.Equal(0, len(descriptors), "Expecting empty list")
}

func (helper * TestHelper) TestAddInstance() {
    networkID := "testAddInstance"
    helper.addTestingNetwork(networkID)
    newInstance := helper.getTestAppInstance(networkID)
    added, err := helper.manager.appManager.AddApplicationInstance(networkID, *newInstance)
    helper.Nil(err, "instance should be added")
    helper.NotNil(added, "instance should be returned")
}

func (helper * TestHelper) TestRetrieveInstance() {
    networkID := "testRetrieveInstance"
    helper.addTestingNetwork(networkID)
    newInstance := helper.getTestAppInstance(networkID)
    added, err := helper.manager.appManager.AddApplicationInstance(networkID, *newInstance)
    helper.Nil(err, "instance should be added")
    helper.NotNil(added, "instance should be returned")
    retrieved, err := helper.manager.appManager.GetInstance(networkID, added.DeployedID)
    helper.Nil(err, "instance should be retrieved")
    helper.NotNil(retrieved, "expecting instance")
    helper.EqualValues(added, retrieved, "structs should match")
}

func (helper * TestHelper) TestListInstance() {
    networkID := "testListInstance"
    helper.addTestingNetwork(networkID)
    numberInstances := 5
    for i:= 0; i < numberInstances; i++ {
        newInstance :=  helper.getTestAppInstance(networkID)
        _, err := helper.manager.appManager.AddApplicationInstance(networkID, *newInstance)
        helper.Nil(err, "instances should be added")
    }
    instances, err := helper.manager.appManager.ListInstances(networkID)
    helper.Nil(err, "list should be retrieved")
    helper.Equal(numberInstances, len(instances), "expecting " + strconv.Itoa(numberInstances) + " instances")
}

func (helper * TestHelper) TestListInstEmptyNetwork(){
    networkID := "testListInstEmptyNetwork"
    helper.addTestingNetwork(networkID)
    descriptors, err := helper.manager.appManager.ListInstances(networkID)
    helper.Nil(err, "list should be retrieved")
    helper.Equal(0, len(descriptors), "Expecting empty list")
}

func (helper * TestHelper) TestUpdateInstance() {
    networkID := "TestUpdateInstance"
    helper.addTestingNetwork(networkID)
    newInstance := helper.getTestAppInstance(networkID)
    added, err := helper.manager.appManager.AddApplicationInstance(networkID, * newInstance)
    helper.Nil(err, "instance should be added")
    helper.NotNil(added, "instance should be returned")

    update := entities.NewUpdateAppInstRequest().
        WithDescription("newDescription").WithClusterID("newCluster")
    updated, err := helper.manager.appManager.UpdateInstance(networkID, added.DeployedID, * update)
    helper.Nil(err, "instance should be updated")
    helper.NotNil(updated, "expecting updated instance")
    helper.Equal(added.DeployedID, updated.DeployedID)
    helper.Equal("newDescription", updated.Description)
    helper.Equal("newCluster", updated.ClusterID)
}

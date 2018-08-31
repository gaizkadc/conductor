//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Applications client mockup implementation

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/server/app"
)

// ApplicationsMockup structure with the required providers and managers.
type ApplicationsMockup struct {
    networkProvider * networkstorage.MockupNetworkProvider
    appDescriptorProvider * appdescstorage.MockupAppDescProvider
    appInstanceProvider * appinststorage.MockupAppInstProvider
    appManager app.Manager
}

// NewApplicationsMockup creates a new mockup Applications client.
func NewApplicationsMockup() Applications {
    networkProvider := networkstorage.NewMockupNetworkProvider()
    appDescriptorProvider := appdescstorage.NewMockupAppDescProvider()
    appInstanceProvider := appinststorage.NewMockupAppInstProvider()
    appManager := app.NewManager(networkProvider, appDescriptorProvider, appInstanceProvider)

    return &ApplicationsMockup{
        networkProvider,
        appDescriptorProvider,
        appInstanceProvider,
        appManager}
}

// AddTestNetwork adds a test network to the mockup.
//   params:
//     networkId The network identifier.
//     An error if the network cannot be added.
func (mockup * ApplicationsMockup) AddTestNetwork(networkID string) derrors.DaishoError {
    toAdd := entities.NewNetworkWithID(
        networkID, "testNetwork", "network description",
        "Admin", "123 123 123", "admin@admins.com")
    return mockup.networkProvider.Add(* toAdd)
}

// AddTestDescriptor adds a test descriptor to the mockup.
//   params:
//     networkId The network identifier.
//     descriptorId The descriptor identifier
//   returns
//     An error if the descriptor cannot be added.
func (mockup * ApplicationsMockup) AddTestDescriptor(networkID string, descriptorID string) derrors.DaishoError {
    toAdd := entities.NewAppDescriptorWithID(networkID,
        descriptorID, "App Name", "description",
        "serviceName", "0.1.0", "label", 0, []string {"nginx:1.12-alpine"})
    err := mockup.appDescriptorProvider.Add(* toAdd)
    if err == nil {
        err = mockup.networkProvider.RegisterAppDesc(networkID, descriptorID)
    }
    return err
}

// AddTestInstance adds a test instance to the mockup.
//   params:
//     networkId The network identifier.
//     descriptorId The descriptor identifier
//     deployedId The instance identifier
//   returns
//     An error if the instance cannot be added.
func (mockup * ApplicationsMockup) AddTestInstance(networkID string, descriptorID string, deployedID string) derrors.DaishoError{
    toAdd := entities.NewAppInstanceWithID(networkID,
        deployedID, descriptorID, "", "name-"+deployedID, "description",
        "", "", entities.AppInstReady, "1Gb", entities.AppStorageDefault,
        make([]entities.ApplicationPort, 0),0, "")
    err := mockup.appInstanceProvider.Add(* toAdd)
    if err == nil{
        err = mockup.networkProvider.RegisterAppInst(networkID, deployedID)
    }
    return err
}

// AddApplicationDescriptor adds a new application descriptor.
//   params:
//     networkID The network identifier.
//     descriptor The application descriptor.
//   returns:
//     The added application descriptor.
//     An error if the application cannot be added.
func (mockup * ApplicationsMockup) AddApplicationDescriptor(networkID string,
    descriptor entities.AddAppDescriptorRequest) (*entities.AppDescriptor, derrors.DaishoError){
    if descriptor.IsValid(){
        return mockup.appManager.AddApplicationDescriptor(networkID, descriptor)
    }
    return nil, derrors.NewEntityError(descriptor, errors.InvalidEntity).WithParams(networkID)
}

// ListDescriptors lists all the application descriptors available for a given network.
//   params:
//     networkID The network identifier.
//   returns:
//     An array of application descriptors.
//     An error in case the list cannot be retrieved.
func (mockup * ApplicationsMockup) ListDescriptors(networkID string) ([] entities.AppDescriptor, derrors.DaishoError){
    return mockup.appManager.ListDescriptors(networkID)
}

// GetDescriptor gets an application descriptor
//   params:
//     networkID The network identifier.
//     appDescriptorID The application descriptor identifier.
//   returns:
//     An application descriptor
//     An error if the descriptor cannot be retrieved or is not associated with the network.
func (mockup * ApplicationsMockup) GetDescriptor(networkID string, appDescriptorID string) (* entities.AppDescriptor, derrors.DaishoError){
    return mockup.appManager.GetDescriptor(networkID, appDescriptorID)
}

// DeleteDescriptor deletes an application descriptor
//   params:
//     networkID The network identifier.
//     appDescriptorID The application descriptor identifier.
//   returns:
//     An error if the descriptor cannot be retrieved or is not associated with the network.
func (mockup * ApplicationsMockup) DeleteDescriptor(networkID string, appDescriptorID string) derrors.DaishoError {
    return mockup.appManager.DeleteDescriptor(networkID, appDescriptorID)
}

// AddApplicationInstance adds a new application instance.
//   params:
//     networkID The network identifier.
//     instance The application instance.
//   returns:
//     The added instance.
//     An error if the application instance cannot be added.
func (mockup * ApplicationsMockup) AddApplicationInstance(networkID string,
    instance entities.AddAppInstanceRequest) (* entities.AppInstance, derrors.DaishoError){
    if instance.IsValid() {
        return mockup.appManager.AddApplicationInstance(networkID, instance)
    }
    return nil, derrors.NewEntityError(instance, errors.InvalidEntity).WithParams(networkID)
}

// ListInstances obtains the list of all the application instances inside a network.
//   params:
//     networkID The network identifier.
//   returns:
//     An array of application instances.
//     An error if the list cannot be retrieved.
func (mockup * ApplicationsMockup) ListInstances(networkID string) ([] entities.AppInstance, derrors.DaishoError){
    return mockup.appManager.ListInstances(networkID)
}

// GetInstance retrieves an application instance.
//   params:
//     networkID The network identifier.
//     appInstanceID The application instance identifier.
//   returns:
//     An application instance.
//     An error if the instance cannot be retrieved.
func (mockup * ApplicationsMockup) GetInstance(networkID string, appInstanceID string) (* entities.AppInstance, derrors.DaishoError){
    return mockup.appManager.GetInstance(networkID, appInstanceID)
}

// UpdateInstance updates an application instance.
//   params:
//     networkID The network identifier.
//     appInstanceID The application instance identifier.
//     request The application update request.
//   returns:
//     The updated application instance.
//     An error if the instance cannot be retrieved.
func (mockup * ApplicationsMockup) UpdateInstance(networkID string, appInstanceID string,
    request entities.UpdateAppInstanceRequest) (* entities.AppInstance, derrors.DaishoError){
    return mockup.appManager.UpdateInstance(networkID, appInstanceID, request)
}

// DeleteInstance deletes an application instance.
//   params:
//     networkID The network identifier.
//     appInstanceID The application instance identifier.
//   returns:
//     An error if the instance cannot be removed.
func (mockup * ApplicationsMockup) DeleteInstance(networkID string, appInstanceID string) derrors.DaishoError {
    return mockup.appManager.DeleteInstance(networkID, appInstanceID)
}

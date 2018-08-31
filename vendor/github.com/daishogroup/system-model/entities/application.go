//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the specification of the applications inside the system model. The concept
// of application is treated in two forms inside Daisho. Application descriptors contain the definition
// of an app including the internal services it is composed of. Application instances on the other hand
// represent the deployment of an app inside a network.

package entities

import (
    "fmt"
)

// AppDescriptorPrefix is the application identifiers.
const AppDescriptorPrefix = "app-"

// AppInstancePrefix is the instance identifiers.
const AppInstancePrefix = "inst-"

// AppStatus is an enum definition of the application instance status.
type AppStatus string

// AppInstInit is the application init status.
const AppInstInit AppStatus = "init"

// AppInstOk is the application ok status.
// Deprecated: Use AppInstReady instead.
const AppInstOk AppStatus = "ok"

// AppInstReady indicates the application is deployed in kubernetes and everything is running.
const AppInstReady AppStatus = "Ready"

// AppInstNotReady indicates the application is deployed in kubernetes but some element (service, deployment, etc.)
// is not ready.
const AppInstNotReady AppStatus = "NotReady"

// AppInstError is the application error status.
const AppInstError AppStatus = "error"

// ValidAppStatus checks the status enum to determine if the string belongs to the enumeration.
//   params:
//     status The status to be checked
//   returns:
//     Whether it is contained in the enum.
func ValidAppStatus(status AppStatus) bool {
    switch status {
    case "":
        return false
    case AppInstInit:
        return true
    case AppInstOk:
        return true
    case AppInstReady:
        return true
    case AppInstNotReady:
        return true
    case AppInstError:
        return true
    default:
        return false
    }
}

// AppStorageType is the enum definition of the application storage type.
type AppStorageType string

// AppStorageDefault is the default storage.
const AppStorageDefault AppStorageType = "default"

// AppStoragePersistent  is the persistent storage..
const AppStoragePersistent AppStorageType = "persistent"

// AppStorageNetPersistent is the network storage.
const AppStorageNetPersistent AppStorageType = "networkPersistent"

// ValidAppStorage checks the storage enum to determine if the string belongs to the enumeration.
//   params:
//     storage The storage to be checked
//   returns:
//     Whether it is contained in the enum.
func ValidAppStorage(storage AppStorageType) bool {
    switch storage {
    case "":
        return false
    case AppStorageDefault:
        return true
    case AppStoragePersistent:
        return true
    case AppStorageNetPersistent:
        return true
    default:
        return false
    }
}

// AppDescriptor represents the concept of application from the point of view of the user. An
// application is composed of several services and the connectivity among them. Information about the
// services need to match those found inside the corresponding ASM packages.
type AppDescriptor struct {
    // Parent network ID.
    NetworkID string `json:"networkId, omitempty"`
    // Application descriptor identifier.
    ID string `json:"id, omitempty"`
    // Application name name.
    Name string `json:"name"`
    // Application description.
    Description string `json:"description, omitempty"`
    // TODO Considering just one service at this point.
    // Name of the service.
    ServiceName string `json:"serviceName"`
    // Version of the service phone.
    ServiceVersion string `json:"serviceVersion"`
    // Label for the super-orchestration system.
    Label string `json:"label"`
    // Listening port for the APP UI.
    // Deprecated: This field will be removed as ports are automatically discovered.
    Port int `json:"port"`
    // List of the required images
    Images []string `json:"images"`
}

func (appDesc *AppDescriptor) String() string {
    return fmt.Sprintf("%#v", appDesc)
}

// NewAppDescriptorWithID generates a new application descriptor.
//   params:
//     id The application descriptor identifier.
//     name The application name.
//     description The application description.
//     serviceName The name of the service.
//     serviceVersion The version of the service.
//     label The application label.
//     port The application exposed port.
//   returns:
//     A new application descriptor.
func NewAppDescriptorWithID(
    networkID string,
    id string,
    name string,
    description string,
    serviceName string,
    serviceVersion string,
    label string,
    port int,
    images []string) *AppDescriptor {
    descriptor := &AppDescriptor{
        networkID,
        id,
        name, description,
        serviceName, serviceVersion,
        label, port, images}
    return descriptor
}

// NewAppDescriptor generates a new application descriptor.
//   params:
//     id The application descriptor identifier.
//     name The application name.
//     description The application description.
//     serviceName The name of the service.
//     serviceVersion The version of the service.
//     label The application label.
//     port The application exposed port.
//   returns:
//     A new application descriptor.
func NewAppDescriptor(
    networkID string,
    name string,
    description string,
    serviceName string,
    serviceVersion string,
    label string,
    port int,
    images []string) *AppDescriptor {
    uuid := GenerateUUID(AppDescriptorPrefix)
    return NewAppDescriptorWithID(networkID, uuid, name, description, serviceName, serviceVersion, label, port, images)
}

// AppDescriptorReducedInfo is the struct to represent the minimal information of AppDescriptor.
type AppDescriptorReducedInfo struct {
    // Parent network ID.
    NetworkID string `json:"networkId, omitempty"`
    // Application descriptor identifier.
    ID string `json:"id, omitempty"`
    // Application name name.
    Name string `json:"name"`
}

// NewAppDescriptorReducedInfo is the main builder of AppDescriptorReducedInfo.
//   params:
//     networkID The network identifier.
//     id The application descriptor identifier.
//     name The application name.
//   returns:
//     A new application descriptor with the minimal information.
func NewAppDescriptorReducedInfo(networkID string, id string,
    name string) *AppDescriptorReducedInfo {
    return &AppDescriptorReducedInfo{NetworkID: networkID, ID: id, Name: name}
}

// AddAppDescriptorRequest to add a new app descriptor
type AddAppDescriptorRequest struct {
    // Application descriptor name.
    Name string `json:"name"`
    // Application description.
    Description string `json:"description, omitempty"`
    // TODO Considering just one service at this point.
    // Name of the service.
    ServiceName string `json:"serviceName"`
    // Version of the service phone.
    ServiceVersion string `json:"serviceVersion"`
    // Label for the super-orchestration system.
    Label string `json:"label"`
    // Listening port for the APP UI.
    Port int `json:"port"`
    // List of required images
    Images []string `json:"images"`
}

func (appDesc *AddAppDescriptorRequest) String() string {
    return fmt.Sprintf("%#v", appDesc)
}

// IsValid checks if the request is valid by checking the name, service name and service version.
func (appDesc *AddAppDescriptorRequest) IsValid() bool {
    return appDesc.Name != "" && appDesc.ServiceName != "" && appDesc.ServiceVersion != "" && appDesc.Label != ""
}

// NewAddAppDescriptorRequest generates a new add application descriptor request.
//   params:
//     name The application name.
//     description The application description.
//     serviceName The name of the service.
//     serviceVersion The version of the service.
//     label The application label.
//     port The application exposed port.
//   returns:
//     A new application descriptor.
func NewAddAppDescriptorRequest(
    name string,
    description string,
    serviceName string,
    serviceVersion string,
    label string,
    port int,
    images []string) *AddAppDescriptorRequest {
    descriptor := &AddAppDescriptorRequest{
        name,
        description,
        serviceName,
        serviceVersion,
        label,
        port,
        images}
    return descriptor
}

// ToAppDescriptor transforms a network request into a network by adding the UUID.
//   params:
//     request The add network request.
//   returns:
//     A network with an UUID.
func ToAppDescriptor(networkID string, request AddAppDescriptorRequest) *AppDescriptor {
    uuid := GenerateUUID(AppDescriptorPrefix)
    return NewAppDescriptorWithID(
        networkID, uuid, request.Name, request.Description,
        request.ServiceName, request.ServiceVersion, request.Label, request.Port, request.Images)
}

// AppInstance represents a deployed application in the system including all associated services.
type AppInstance struct {
    //Parent Network ID
    NetworkID string `json:"networkId"`
    // Application instance identifier.
    DeployedID string `json:"deployedId"`
    // Application descriptor identifier.
    AppDescriptorID string `json:"appDescriptorId"`
    // Cluster identifier where the application is deployed.
    ClusterID string `json:"clusterId, omitempty"`
    // Application instance name.
    Name string `json:"name"`
    // Application instance description.
    Description string `json:"description, omitempty"`
    // TODO Considering just one service at this point.
    // Labels added to the instance. Label can be used by the super orchestration system to decide where an instance
    // should be launched.
    Label string `json:"label"`
    // Arguments passed to the application.
    Arguments string `json:"arguments, omitempty"`
    // Application status
    Status AppStatus `json:"status"`
    // Persistence size required by the app
    PersistenceSize string `json:"persistentSize"`
    // Storage type
    StorageType AppStorageType `json:"storageType"`
    // Ports contains the list of exposed ports.
    Ports [] ApplicationPort `json:"ports"`
    // Listening port for the APP UI.
    // Deprecated: Use Ports instead.
    Port int `json:"port"`
    // Address where the port is exposed.
    ClusterAddress string `json:"clusterAddress, omitempty"`
}

func (appInstance *AppInstance) String() string {
    return fmt.Sprintf("%#v", appInstance)
}

// IsValid checks if the instance is valid by checking the application required fields.
func (appInstance *AppInstance) IsValid() bool {
    return appInstance.DeployedID != "" &&
        appInstance.AppDescriptorID != "" &&
        appInstance.Name != "" &&
        appInstance.Label != "" &&
        appInstance.PersistenceSize != "" &&
        ValidAppStatus(appInstance.Status) &&
        ValidAppStorage(appInstance.StorageType)
}

// DefaultPort returns the default port that will be used in the UI to open to application dashboard.
func (appInstance * AppInstance) DefaultPort() *ApplicationPort {
    var result * ApplicationPort = nil
    found := false
    for i := 0; i < len(appInstance.Ports) && !found ; i++ {
        if appInstance.Ports[i].NodePort != 0 {
            result = &appInstance.Ports[i]
            found = true
        }
    }
    return result
}

// Clone create a copy from the original object.
func (appInstance *AppInstance) Clone() *AppInstance {
    ports := make([]ApplicationPort, 0)
    for _, p := range appInstance.Ports {
        toAdd := NewApplicationPort(p.Name, p.Protocol, p.Port, p.TargetPort, p.NodePort)
        ports = append(ports, *toAdd)
    }
    return NewAppInstanceWithID(appInstance.NetworkID, appInstance.DeployedID, appInstance.AppDescriptorID,
        appInstance.ClusterID,
        appInstance.Name, appInstance.Description, appInstance.Label, appInstance.Arguments, appInstance.Status,
        appInstance.PersistenceSize, appInstance.StorageType, ports, appInstance.Port, appInstance.ClusterAddress)
}

// Merge an update request with the information of the current instance.
func (appInstance *AppInstance) Merge(update UpdateAppInstanceRequest) *AppInstance {
    updated := appInstance.Clone()
    if update.ClusterID != nil {
        updated.ClusterID = * update.ClusterID
    }
    if update.Description != nil {
        updated.Description = * update.Description
    }
    if update.Status != nil {
        updated.Status = * update.Status
    }
    if update.ClusterAddress != nil {
        updated.ClusterAddress = * update.ClusterAddress
    }
    if update.Ports != nil {
        updated.Ports = * update.Ports
    }
    return updated
}

// NewAppInstanceWithID generate a new application instance.
//   params:
//     deployedId The application instance identifier.
//     appDescriptorId The application descriptor identifier.
//     clusterId The cluster where the application is deployed.
//     name The application instance name.
//     description The application instance description.
//     label The label for the instance.
//     arguments The arguments for launching the app.
//     status The application status.
//     persistentSize The persistent size required by the instance.
//     storageType The storage type.
//     port Application port.
//     clusterAddress Address where the port is exposed.
//   returns:
//     A new application instance.
func NewAppInstanceWithID(
    networkID string,
    deployedID string,
    appDescriptorID string,
    clusterID string,
    name string,
    description string,
    label string,
    arguments string,
    status AppStatus,
    persistentSize string,
    storageType AppStorageType,
    ports []ApplicationPort,
    port int,
    clusterAddress string) *AppInstance {
    instance := &AppInstance{
        networkID,
        deployedID,
        appDescriptorID,
        clusterID,
        name,
        description,
        label,
        arguments,
        status,
        persistentSize,
        storageType,
        ports,
        port,
        clusterAddress}
    return instance
}

// NewAppInstance generate a new application instance.
//   params:
//     networkID The parent network identification.
//     appDescriptorId The application descriptor identifier.
//     clusterId The cluster where the application is deployed.
//     name The application instance name.
//     description The application instance description.
//     label The label for the instance.
//     arguments The arguments for launching the app.
//     persistentSize The persistent size required by the instance.
//     storageType The storage type.
//     port Application port.
//     clusterAddress Address where the port is exposed.
//   returns:
//     A new application instance.
func NewAppInstance(
    networkID string,
    appDescriptorID string,
    clusterID string,
    name string,
    description string,
    label string,
    arguments string,
    persistentSize string,
    storageType AppStorageType,
    ports [] ApplicationPort,
    port int,
    clusterAddress string) *AppInstance {
    uuid := GenerateUUID(AppInstancePrefix)
    return NewAppInstanceWithID(networkID, uuid, appDescriptorID, clusterID,
        name, description, label, arguments, AppInstInit, persistentSize, storageType, ports, port, clusterAddress)
}

// AppInstanceReducedInfo is the struct to represent the minimal information of AppInstance.
type AppInstanceReducedInfo struct {
    // Parent network ID.
    NetworkID string `json:"networkId, omitempty"`
    // Cluster identifier where the application is deployed.
    ClusterID string `json:"clusterId, omitempty"`
    // Application descriptor identifier.
    AppDescriptorID string `json:"appDescriptorId"`
    // Application instance identifier.
    DeployedID string `json:"deployedId"`
    // Application name name.
    Name string `json:"name"`
    // Application description.
    Description string `json:"description, omitempty"`
    // Ports contains the list of exposed ports.
    Ports [] ApplicationPort `json:"ports"`
    // Listening port for the APP UI.
    // Deprecated: Use Ports instead.
    Port int `json:"port"`
    // Address where the port is exposed.
    ClusterAddress string `json:"clusterAddress, omitempty"`
}

// NewAppInstanceReducedInfo is the main builder of AppInstanceReducedInfo.
//   params:
//     networkID The network identifier.
//     clusterID The cluster identifier.
//     appDescriptorID The app descriptor identifier.
//     deployedID The app instance identifier.
//     name The application name.
//     description The application description.
//     ports The list of exposed ports.
//     port The deployed port.
//     clusterAddress The current ip.
//   returns:
//     A new application descriptor with the minimal information.
func NewAppInstanceReducedInfo(networkID string, clusterID string, appDescriptorID, deployedID string,
    name string, description string, ports [] ApplicationPort, port int, clusterAddress string) *AppInstanceReducedInfo {
    return &AppInstanceReducedInfo{networkID, clusterID,
        appDescriptorID, deployedID,
        name, description,
        ports, port, clusterAddress}
}

// AddAppInstanceRequest to add a new app instance.
type AddAppInstanceRequest struct {
    // Application descriptor identifier.
    AppDescriptorID string `json:"appDescriptorId, omitempty"`
    // Application instance name.
    Name string `json:"name, omitempty"`
    // Application description.
    Description string `json:"description, omitempty"`
    // TODO Considering just one service at this point.
    // Instance label.
    Label string `json:"label, omitempty"`
    // Arguments for launching the app.
    Arguments string `json:"arguments, omitempty"`
    // Persistence size required by the app
    PersistenceSize string `json:"persistentSize"`
    // Storage type
    StorageType AppStorageType `json:"storageType"`
}

func (appInstanceReq *AddAppInstanceRequest) String() string {
    return fmt.Sprintf("%#v", appInstanceReq)
}

// IsValid checks if the request is valid by checking the application descriptor identifier.
func (appInstanceReq *AddAppInstanceRequest) IsValid() bool {
    return appInstanceReq.AppDescriptorID != "" && appInstanceReq.Name != ""
}

// NewAddAppInstanceRequest generates a new add application instance request.
//   params:
//     appDescriptorId The application descriptor identifier.
//     name The application instance name.
//     description The application instance description.
//     label The label for the instance.
//     arguments The arguments for launching the app.
//     persistentSize The persistent size required by the instance.
//     storageType The storage type.
//   returns:
//     A new add application instance request.
func NewAddAppInstanceRequest(
    appDescriptorID string,
    name string,
    description string,
    label string,
    arguments string,
    persistentSize string,
    storageType AppStorageType) *AddAppInstanceRequest {
    request := &AddAppInstanceRequest{
        appDescriptorID,
        name, description, label,
        arguments, persistentSize, storageType}
    return request
}

// ToAppInstance transforms a application instance request into a instance by adding the UUID.
//   params:
//     request The add application instance request.
//   returns:
//     An instance with an UUID.
func ToAppInstance(networkID string, request AddAppInstanceRequest) *AppInstance {
    uuid := GenerateUUID(AppInstancePrefix)
    return NewAppInstanceWithID(networkID, uuid, request.AppDescriptorID, "", request.Name, request.Description,
        request.Label, request.Arguments, AppInstInit, request.PersistenceSize, request.StorageType,
        make([]ApplicationPort, 0), -1, "")
}

// UpdateAppInstanceRequest is the struct to update an instance.
type UpdateAppInstanceRequest struct {
    // Cluster identifier where the application is deployed.
    ClusterID *string `json:"clusterId, omitempty"`
    // Application instance description.
    Description *string `json:"description, omitempty"`
    // Application status
    Status *AppStatus `json:"status"`
    // Address where the port is exposed.
    ClusterAddress *string `json:"clusterAddress, omitempty"`
    // Ports exposed by the application.
    Ports * []ApplicationPort `json:"ports, omitempty"`
}

// NewUpdateAppInstRequest creates a new update node request.
func NewUpdateAppInstRequest() *UpdateAppInstanceRequest {
    return &UpdateAppInstanceRequest{}
}

// IsValid checks if the request is valid by checking the application update status.
func (update *UpdateAppInstanceRequest) IsValid() bool {
    if update.Status != nil {
        return ValidAppStatus(* update.Status)
    }
    return true
}

func (update *UpdateAppInstanceRequest) String() string {
    return fmt.Sprintf("%#v", update)
}

// WithClusterID updates the request with a new cluster id.
//   params:
//     clusterID The new cluster identifier.
//   returns:
//     An update app instance request.
func (update *UpdateAppInstanceRequest) WithClusterID(clusterID string) *UpdateAppInstanceRequest {
    update.ClusterID = &clusterID
    return update
}

// WithDescription updates the request with a new description.
//   params:
//     description The new name.
//   returns:
//     An update app instance request.
func (update *UpdateAppInstanceRequest) WithDescription(description string) *UpdateAppInstanceRequest {
    update.Description = &description
    return update
}

// WithStatus updates the request with a new status.
//   params:
//     status The new status.
//   returns:
//     An update app instance request.
func (update *UpdateAppInstanceRequest) WithStatus(status AppStatus) *UpdateAppInstanceRequest {
    update.Status = &status
    return update
}

// WithClusterAddress updates the request with a new cluster address.
//   params:
//     address The new cluster address where the app is exposed.
//   returns:
//     An update app instance request.
func (update *UpdateAppInstanceRequest) WithClusterAddress(address string) *UpdateAppInstanceRequest {
    update.ClusterAddress = &address
    return update
}

// WithPorts updates the request with a new set of application ports.
func (update * UpdateAppInstanceRequest) WithPorts(ports []ApplicationPort) *UpdateAppInstanceRequest {
    update.Ports = &ports
    return update
}
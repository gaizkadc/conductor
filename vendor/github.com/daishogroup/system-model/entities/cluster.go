//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package entities

import (
    "fmt"
    "strconv"
)

// Definition of the different types of clusters in the system.
type ClusterType string

// Cluster deployed on a cloud infraestructure (e.g., Amazon, Azure, etc.).
const CloudType ClusterType = "cloud"

// Cluster deployed at the edge.
const EdgeType ClusterType = "edge"

// Cluster deployed near the data source.
const GatewayType ClusterType = "gateway"

// ValidClusterType checks the type enum to determine if the string belongs to the enumeration.
//   params:
//     clusterType The type to be checked
//   returns:
//     Whether it is contained in the enum.
func ValidClusterType(clusterType ClusterType) bool {
    switch clusterType {
    case "" : return false
    case CloudType : return true
    case EdgeType : return true
    case GatewayType : return true
    default: return false
    }
}

type ClusterStatus string

const ClusterCreated ClusterStatus = "Created"
const ClusterReadyToInstall ClusterStatus = "ReadyToInstall"
const ClusterInstalling ClusterStatus = "Installing"
const ClusterInstalled ClusterStatus = "Installed"
const ClusterUninstalling ClusterStatus = "Uninstalling"
const ClusterError ClusterStatus = "Error"

// ValidClusterStatus checks the type enum to determine if the string belongs to the enumeration.
//   params:
//     status The type to be checked
//   returns:
//     Whether it is contained in the enum.
func ValidClusterStatus(status ClusterStatus) bool {
    switch status {
    case "" : return false
    case ClusterCreated : return true
    case ClusterReadyToInstall : return true
    case ClusterInstalling : return true
    case ClusterInstalled : return true
    case ClusterUninstalling : return true
    case ClusterError : return true
    default: return false
    }
}

// Cluster contains a set of nodes that make the computing system where the user can deploy apps on.
type Cluster struct {
    // NetworkID with the network identifier.
    NetworkID string `json:"networkId"`
    // The cluster identifier.
    ID string `json:"id"`
    // Cluster name.
    Name string `json:"name, omitempty"`
    // Cluster description.
    Description string `json:"description, omitempty"`
    // Cluster type.
    Type ClusterType `json:"type, omitempty"`
    // Cluster location.
    Location string `json:"location, omitempty"`
    // The admin email of the cluster.
    Email string `json:"email, omitempty"`
    // The current status of the cluster.
    Status ClusterStatus `json:"status"`
    // The cluster don't accept applications and is empty.
    Drain bool `json:"drain"`
    // The cluster don't accept applications but it can contains launched applications previously.
    Cordon bool `json:"cordon"`
}

func (c *Cluster) String() string {
    return fmt.Sprintf("%#v", c)
}

// Create a clone of the object.
func (c *Cluster) Clone() *Cluster {
    return NewClusterWithID(c.NetworkID, c.ID, c.Name, c.Description, c.Type,
        c.Location, c.Email, c.Status, c.Drain, c.Cordon)
}

// Create a new Cluster merge the information with a UpdateClusterRequest object.
//     params:
//         update The list of updates.
//     returns:
//         A new Cluster with the applied changes.
func (c *Cluster) Merge(update UpdateClusterRequest) *Cluster {
    newCluster := c.Clone()
    if update.Name != nil {
        newCluster.Name = * update.Name
    }
    if update.Description != nil {
        newCluster.Description = * update.Description
    }
    if update.Type != nil {
        newCluster.Type = * update.Type
    }
    if update.Location != nil {
        newCluster.Location = * update.Location
    }
    if update.Email != nil {
        newCluster.Email = * update.Email
    }
    if update.Status != nil {
        newCluster.Status = * update.Status
    }
    if update.Drain != nil {
        newCluster.Drain = * update.Drain
    }
    if update.Cordon != nil {
        newCluster.Cordon = * update.Cordon
    }
    return newCluster
}

// Generate a new cluster.
//   params:
//     name The cluster name.
//     description The cluster description.
//     clusterType The type of cluster.
//     location The cluster location.
//     canRemove Whether the cluster can be removed.
//   returns:
//     A new cluster.
func NewCluster(networkID string, name string, description string,
    clusterType ClusterType, location string,
    email string, status ClusterStatus, drain bool, cordon bool) *Cluster {
    uuid := GenerateUUID(ClusterPrefix)
    return NewClusterWithID(networkID, uuid, name, description, clusterType,
        location, email, status, drain, cordon)
}

// Generate a new cluster.
//   params:
//     networkID The network identifier.
//     id The cluster identifier.
//     name The cluster name.
//     description The cluster description.
//     clusterType The type of cluster.
//     location The cluster location.
//     canRemove Whether the cluster can be removed.
//   returns:
//     A new cluster.
func NewClusterWithID(networkID string, id string, name string, description string, clusterType ClusterType,
    location string, email string, status ClusterStatus, drain bool, cordon bool) *Cluster {
    cluster := &Cluster{networkID, id, name, description, clusterType,
        location, email, status, drain, cordon,}
    return cluster
}

// Cluster contains a set of nodes that make the computing system where the user can deploy apps on.
type ClusterReducedInfo struct {
    //The network identifier
    NetworkID string `json:"networkId"`
    // The cluster identifier.
    ID string `json:"id"`
    // Cluster name.
    Name string `json:"name, omitempty"`
    // Cluster type.
    Type ClusterType `json:"type, omitempty"`
}

// NewClusterReducedInfo generates a new cluster.
//   params:
//     networkID The network identifier.
//     id The cluster identifier.
//     name The name of the cluster.
//     clusterType The cluster type.
//   returns:
//     A new cluster with reduced info.
func NewClusterReducedInfo(networkID string, id string, name string,
    clusterType ClusterType) *ClusterReducedInfo {
    return &ClusterReducedInfo{
        NetworkID: networkID,
        ID:        id,
        Name:      name,
        Type:      clusterType,
    }
}

// Structure used by the REST endpoints to create a new cluster.
type AddClusterRequest struct {
    // Cluster name.
    Name string `json:"name, omitempty"`
    // Cluster description.
    Description string `json:"description, omitempty"`
    // Cluster type.
    Type ClusterType `json:"type, omitempty"`
    // Cluster location.
    Location string `json:"location, omitempty"`
    // The admin email of the cluster.
    Email string `json:"email, omitempty"`
}

// Create a new add cluster request.
//   params:
//     name The cluster name.
//     description The cluster description.
//     clusterType The type of cluster.
//     location The cluster location.
//   returns:
//     A new add cluster request.
func NewAddClusterRequest(name string, description string, clusterType ClusterType,
    location string, email string) *AddClusterRequest {
    request := &AddClusterRequest{
        name, description, clusterType, location, email,
    }
    return request
}

// Check if the request is valid
//   params:
//     request The add network request.
//   returns:
//     Whether the request is valid.
func (request *AddClusterRequest) IsValid() bool {
    // TODO Add validation based on required fields
    return request.Name != "" && ValidClusterType(request.Type)
}

// Transform a cluster request into a cluster by adding the UUID.
//   params:
//     request The add cluster request.
//     networkID The parent network ID.
//   returns:
//     A cluster with an UUID.
func ToCluster(networkID string, request AddClusterRequest) *Cluster {
    uuid := GenerateUUID(ClusterPrefix)
    return NewClusterWithID(networkID, uuid, request.Name, request.Description, request.Type, request.Location,
        request.Email, ClusterCreated, false, false)
}

func (c *AddClusterRequest) String() string {
    return fmt.Sprintf("%#v", c)
}

// Object for updating the cluster information.
type UpdateClusterRequest struct {
    // Cluster name.
    Name *string `json:"name, omitempty"`
    // Cluster description.
    Description *string `json:"description, omitempty"`
    // Cluster type.
    Type *ClusterType `json:"type, omitempty"`
    // Cluster location.
    Location *string `json:"location, omitempty"`
    // The admin email of the cluster.
    Email *string `json:"email, omitempty"`
    // The current status of the cluster.
    Status *ClusterStatus `json:"status"`
    // The cluster don't accept applications and is empty.
    Drain *bool `json:"drain"`
    // The cluster don't accept applications but it can contains launched applications previously.
    Cordon *bool `json:"cordon"`
}

func NewUpdateClusterRequest() *UpdateClusterRequest {
    return &UpdateClusterRequest{}
}

// IsValid checks the enum types are valid when present.
func (update * UpdateClusterRequest) IsValid() bool {
    result := true
    if update.Status != nil {
        result = result && ValidClusterStatus(* update.Status)
    }
    if update.Type != nil {
        result = result && ValidClusterType(* update.Type)
    }
    return result
}

// PrettyString generates a string representation showing nil values.
func (update *UpdateClusterRequest) PrettyString() string {
    name := "nil"
    description := "nil"
    cType := "nil"
    location := "nil"
    email := "nil"
    status := "nil"
    drain := "nil"
    cordon := "nil"
    if update.Name != nil {
        name = * update.Name
    }
    if update.Description != nil {
        description = * update.Description
    }
    if update.Type != nil {
        cType = string(* update.Type)
    }
    if update.Location != nil {
        location = * update.Location
    }
    if update.Email != nil {
        email = * update.Email
    }
    if update.Status != nil {
        status = string(* update.Status)
    }
    if update.Drain != nil {
        drain = strconv.FormatBool(* update.Drain)
    }
    if update.Cordon != nil {
        cordon = strconv.FormatBool(* update.Cordon)
    }
    return "UpdateClusterRequest{name:" + name + ", description:" + description + ", Type:" + cType + ", Location:" + location + ", Email:" + email + ", Status:" + status + ", Drain:" + drain + ", Cordon:" + cordon + "}"
}

func (update *UpdateClusterRequest) String() string {
    return fmt.Sprintf("%#v", update)
}

// Update the request with a new name.
//   params:
//     value The new name.
//   returns:
//     An update cluster request.
func (update *UpdateClusterRequest) WithName(value string) *UpdateClusterRequest {
    update.Name = &value
    return update
}

// Update the request with a new description.
//   params:
//     value The new description.
//   returns:
//     An update cluster request.
func (update *UpdateClusterRequest) WithDescription(value string) *UpdateClusterRequest {
    update.Description = &value
    return update
}

// Update the request with a new cluster type.
//   params:
//     value The new cluster type.
//   returns:
//     An update cluster request.
func (update *UpdateClusterRequest) WithType(value ClusterType) *UpdateClusterRequest {
    update.Type = &value
    return update
}

// Update the request with a new location.
//   params:
//     value The new location.
//   returns:
//     An update cluster request.
func (update *UpdateClusterRequest) WithLocation(value string) *UpdateClusterRequest {
    update.Location = &value
    return update
}

// Update the request with a new email.
//   params:
//     value The new email.
//   returns:
//     An update cluster request.
func (update *UpdateClusterRequest) WithEmail(value string) *UpdateClusterRequest {
    update.Email = &value
    return update
}

// Update the request with a new cluster status.
//   params:
//     value The new cluster status.
//   returns:
//     An update cluster request.
func (update *UpdateClusterRequest) WithClusterStatus(value ClusterStatus) *UpdateClusterRequest {
    update.Status = &value
    return update
}

// Update the request with a new drain flag.
//   params:
//     value The new drain flag.
//   returns:
//     An update cluster request.
func (update *UpdateClusterRequest) WithDrain(value bool) *UpdateClusterRequest {
    update.Drain = &value
    return update
}

// Update the request with a new cordon flag.
//   params:
//     value The new cordon flag.
//   returns:
//     An update cluster request.
func (update *UpdateClusterRequest) WithCordon(value bool) *UpdateClusterRequest {
    update.Cordon = &value
    return update
}

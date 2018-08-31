//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the specification of the node entities inside the system model.

package entities

import (
    "fmt"
)

// Definition of the different types of node statuses in the system. The list of status and the expected transitions
// can be found at: https://www.lucidchart.com/invitations/accept/bdb4d21f-f32d-4da4-9f06-0ba7a80fc8f7
type NodeStatus string

// NodeUnchecked corresponds to the state a node enters into the system.
const NodeUnchecked NodeStatus = "Unchecked"

// NodeReadyToInstall represents a node with proper credentials and configuration that can be installed.
const NodeReadyToInstall NodeStatus = "ReadyToInstall"

// NodeInstalling represents a node that is being installed at the moment.
const NodeInstalling NodeStatus = "Installing"

// NodeInstalled represents a node with the Daisho platform components.
const NodeInstalled NodeStatus = "Installed"

// NodeUninstalling represents a node that is being uninstalled.
const NodeUninstalling NodeStatus = "Uninstalling"

// NodePrecheckError represents a node that failed the precheck before installing.
const NodePrecheckError NodeStatus = "PrecheckError"

// NodeError represents a node in a generic error state
const NodeError NodeStatus = "Error"

// ValidNodeStatus checks the status enum to determine if the string belongs to the enumeration.
//   params:
//     status The status to be checked
//   returns:
//     Whether it is contained in the enum.
func ValidNodeStatus(status NodeStatus) bool {
    switch status {
    case "" : return false
    case NodeUnchecked : return true
    case NodeReadyToInstall : return true
    case NodeInstalling : return true
    case NodeInstalled : return true
    case NodeUninstalling : return true
    case NodePrecheckError : return true
    case NodeError : return true
    default: return false
    }
}

// A node represents a Kubernetes node. Nodes compose a cluster, and each node has the information
// to connect to the physical element.
type Node struct {
    // The parent network.
    NetworkID string `json:"networkId"`
    // The edgenet network ID this node is part of. This is filled in during
    // the add node operation in the system model node manager, based on what
    // is set in the network.
    EdgenetNetworkID string `json:"edgenetNetwork,omitempty"`
    // The parent cluster.
    ClusterID string `json:"clusterId"`
    // The node identifier.
    ID string `json:"id"`
    // Node name.
    Name string `json:"name"`
    // Node description.
    Description string `json:"description,omitempty"`
    // Labels of the current node. This information will be used by upper layers to filter target nodes.
    Labels []string
    // The public IP of the node.
    PublicIP string `json:"publicIP"`
    // The private IP of the node.
    PrivateIP string `json:"privateIP"`
    // The edgenet address of the node
    EdgenetAddress string `json:"edgenetAddress,omitempty"`
    // The edgenet IP of the node (TODO: should this replace public/private?)
    EdgenetIP string `json:"edgenetIP,omitempty"`
    // A flag to indicate if the current node is deployed Daisho platform.
    Installed bool `json:"installed"`
    // Username required to connect through SSH.
    Username string `json:"username"`
    // Password used to connect through SSH.
    Password string `json:"password,omitempty"`
    // SSH key used to connect through SSH.
    SSHKey string `json:"sshKey,omitempty"`
    // Node status.
    Status NodeStatus `json:"status"`
}

// Generate a new node.
//   params:
//      name            The node name.
//      description     The node description.
//      publicIP        The public IP of the node.
//      privateIP       The private IP of the node.
//      installed       A flag to indicate if the current node is deployed Daisho platform.
//      username        The username used to connect through SSH.
//      password        The password used to connect through SSH.
//      sshKey          The sshKey used to connect through SSH.
//   returns:
//     A new node.
func NewNode(
    networkID string, clusterID string, name string, description string, labels [] string,
    publicIP string, privateIP string, installed bool,
    username string, password string, sshKey string) *Node {
    uuid := GenerateUUID(NodePrefix)
    return NewNodeWithID(networkID, clusterID,
        uuid, name, description, labels,
        publicIP, privateIP, installed,
        username, password, sshKey, NodeUnchecked)
}

// Generate a new node.
//   params:
//      id              The node identifier.
//      name            The node name.
//      description     The node description.
//      labels          The set of node labels.
//      publicIP        The public IP of the node.
//      privateIP       The private IP of the node.
//      installed       A flag to indicate if the current node is deployed Daisho platform.
//      username        The username used to connect through SSH.
//      password        The password used to connect through SSH.
//      sshKey          The sshKey used to connect through SSH.
//      status          The node status.
//   returns:
//     A new node.
func NewNodeWithID(networkID string, clusterID string,
    id string, name string, description string, labels []string,
    publicIP string, privateIP string, installed bool,
    username string, password string, sshKey string, status NodeStatus) *Node {
    node := &Node{
        NetworkID: networkID,
        ClusterID: clusterID,
        ID: id,
        Name: name,
        Description: description,
        Labels: labels,
        PublicIP: publicIP,
        PrivateIP: privateIP,
        Installed: installed,
        Username: username,
        Password: password,
        SSHKey: sshKey,
        Status: status,
    }
    return node
}

func (node *Node) String() string {
    return fmt.Sprintf("%#v", node)
}

func (node *Node) Clone() *Node {
    newNode := &Node{}
    *newNode = *node
    return newNode
}

func (node *Node) Merge(update UpdateNodeRequest) *Node {
    updatedNode := node.Clone()
    if update.Name != nil {
        updatedNode.Name = * update.Name
    }
    if update.Description != nil {
        updatedNode.Description = * update.Description
    }
    if update.Labels != nil {
        updatedNode.Labels = *update.Labels
    }
    if update.PublicIP != nil {
        updatedNode.PublicIP = * update.PublicIP
    }
    if update.PrivateIP != nil {
        updatedNode.PrivateIP = * update.PrivateIP
    }
    if update.Installed != nil {
        updatedNode.Installed = * update.Installed
    }
    if update.Username != nil {
        updatedNode.Username = * update.Username
    }
    if update.Password != nil {
        updatedNode.Password = * update.Password
    }
    if update.SSHKey != nil {
        updatedNode.SSHKey = * update.SSHKey
    }
    if update.Status != nil {
        updatedNode.Status = * update.Status
    }
    if update.EdgenetAddress != nil {
        updatedNode.EdgenetAddress = *update.EdgenetAddress
    }
    if update.EdgenetIP != nil {
        updatedNode.EdgenetIP = *update.EdgenetIP
    }
    return updatedNode
}

// NodeReducedInfo represents a Kubernetes node. Nodes compose a cluster, and each node has the information
// to connect to the physical element.
type NodeReducedInfo struct {
    // The parent network.
    NetworkID string `json:"networkId"`
    // The parent cluster.
    ClusterID string `json:"clusterId"`
    // The node identifier.
    ID string `json:"id"`
    // Node name.
    Name string `json:"name"`
    // Node status.
    Status NodeStatus `json:"status"`
    // PublicIP is the public ip address of a node.
    PublicIP string `json:"publicIP"`
}

// NewNodeReducedInfo is a basic builder.
//   params:
//      networkID       The selected network.
//      id              The node identifier.
//      name            The node name.
//      status          The node status.
//   returns:
//     A new node.
func NewNodeReducedInfo(networkID string, clusterID string,
    id string, name string, status NodeStatus, publicIP string) *NodeReducedInfo {
    node := &NodeReducedInfo{
        networkID,
        clusterID,
        id,
        name,
        status,
        publicIP,
    }
    return node
}

// Structure used by the REST endpoints to create a new node.
type AddNodeRequest struct {
    // Node name.
    Name string `json:"name"`
    // Node description.
    Description string `json:"description,omitempty"`
    // Labels of the node to facilitate filtering.
    Labels []string `json:"labels"`
    //The public IP of the node.
    PublicIP string `json:"publicIP"`
    //The private IP of the node.
    PrivateIP string `json:"privateIP"`
    // The edgenet address of the node
    EdgenetAddress string `json:"edgenetAddress,omitempty"`
    // The edgenet IP of the node (TODO: should this replace public/private?)
    EdgenetIP string `json:"edgenetIP,omitempty"`
    // A flag to indicate if the current node is deployed Daisho platform.
    Installed bool `json:"installed,omitempty"`
    // Username required to connect through SSH.
    Username string `json:"username"`
    // Password used to connect through SSH.
    Password string `json:"password,omitempty"`
    // SSH key used to connect through SSH.
    SSHKey string `json:"sshKey,omitempty"`
}

// Check if the request is valid
//   params:
//     request The add node request.
//   returns:
//     Whether the request is valid.
func (request *AddNodeRequest) IsValid() bool {
    return request.Name != "" && request.PublicIP != "" && request.PrivateIP != "" && request.Username != ""
}

func (request *AddNodeRequest) String() string {
    return fmt.Sprintf("%#v", request)
}

// Create a new add cluster request.
//   params:
//      name            The node name.
//      description     The node description.
//      labels Set of labels.
//      publicIP        The public IP of the node.
//      privateIP       The private IP of the node.
//      installed       A flag to indicate if the current node is deployed Daisho platform.
//      username        The username used to connect through SSH.
//      password        The password used to connect through SSH.
//      sshKey          The sshKey used to connect through SSH.
//   returns:
//     A new add node request.
func NewAddNodeRequest(
    name string, description string, labels [] string,
    publicIP string, privateIP string, installed bool,
    username string, password string, sshKey string) *AddNodeRequest {
    request := &AddNodeRequest{
        Name: name,
        Description: description,
        Labels: labels,
        PublicIP: publicIP,
        PrivateIP: privateIP,
        Installed: installed,
        Username: username,
        Password: password,
        SSHKey: sshKey,
    }
    return request
}

// Transform a node request into a node by adding the UUID.
//   params:
//     request The add cluster request.
//   returns:
//     A node with an UUID.
func ToNode(networkID string, clusterID string, request AddNodeRequest) *Node {
    node := NewNode(networkID, clusterID,
        request.Name, request.Description, request.Labels,
        request.PublicIP, request.PrivateIP, request.Installed,
        request.Username, request.Password, request.SSHKey)

    // We do want to have a node that reflects all values in the request.
    // However, we don't want all (optional) request fields in the NewNode
    // functions, as it makes it
    // a) a lot of work to add new fields, and
    // b) the functional call is not very descriptive.
    //
    // Setting fields explicitly by name makes it very clear what is happening.
    // In general, New* functions should return minimially initialized struct,
    // i.e., one that is valid (has all mandatory fields set). After creation,
    // it can be modified by the calling function directly. Fields that should
    // not be modified after creation should be made private.
    // We should at some point reflect this in our NewNode function by
    // removing everything but ID and Name from the arguments.
    node.EdgenetAddress = request.EdgenetAddress
    node.EdgenetIP = request.EdgenetIP

    return node
}

// The update node request contains all the parameters that can be updated from a node instance.
type UpdateNodeRequest struct {
    // Node name.
    Name *string `json:"name,omitempty"`
    // Node description.
    Description *string `json:"description,omitempty"`
    // Labels of the node to facilitate filtering.
    Labels * []string `json:"labels"`
    // The public IP of the node.
    PublicIP *string `json:"publicIP,omitempty"`
    // The private IP of the node.
    PrivateIP *string `json:"privateIP,omitempty"`
    // The edgenet address of the node
    EdgenetAddress *string `json:"edgenetAddress,omitempty"`
    // The edgenet IP of the node (TODO: should this replace public/private?)
    EdgenetIP *string `json:"edgenetIP,omitempty"`
    // A flag to indicate if the current node is deployed Daisho platform.
    Installed *bool `json:"installed,omitempty"`
    // Username required to connect through SSH.
    Username *string `json:"username,omitempty"`
    // Password used to connect through SSH.
    Password *string `json:"password,omitempty"`
    // SSH key used to connect through SSH.
    SSHKey *string `json:"sshKey,omitempty"`
    // Node status.
    Status *NodeStatus `json:"status,omitempty"`
}

// Create a new update node request.
func NewUpdateNodeRequest() *UpdateNodeRequest {
    return &UpdateNodeRequest{}
}

func (update * UpdateNodeRequest) IsValid() bool {
    if update.Status != nil {
        return ValidNodeStatus(* update.Status)
    }
    return true
}

func (update *UpdateNodeRequest) String() string {
    return fmt.Sprintf("%#v", update)
}

// Update the request with a new name.
//   params:
//     name The new name.
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithName(name string) *UpdateNodeRequest {
    update.Name = &name
    return update
}

// Update the request with a new description.
//   params:
//     description The new description.
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithDescription(description string) *UpdateNodeRequest {
    update.Description = &description
    return update
}

// Update the request with a new set of labels.
//   params:
//     description The new set of labels.
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithLabels(labels [] string) *UpdateNodeRequest {
    update.Labels = &labels
    return update
}

// Update the request with a new public IP.
//   params:
//     ip The new IP address.
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithPublicIP(ip string) *UpdateNodeRequest {
    update.PublicIP = &ip
    return update
}

// Update the request with a new private IP.
//   params:
//     ip The new IP address.
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithPrivateIP(ip string) *UpdateNodeRequest {
    update.PrivateIP = &ip
    return update
}

// Update the request with a new installed state.
//   params:
//     installed Whether the node is installed.
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithInstalled(installed bool) *UpdateNodeRequest {
    update.Installed = &installed
    return update
}

// Update the request with a new username.
//   params:
//     username The new username to connect through SSH.
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithUsername(username string) *UpdateNodeRequest {
    update.Username = &username
    return update
}

// Update the request with a new password.
//   params:
//     password The new password to connect through SSH.
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithPassword(password string) *UpdateNodeRequest {
    update.Password = &password
    return update
}

// Update the request with a new SSH key.
//   params:
//     key The new key to connect through SSH.
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithSSHKey(key string) *UpdateNodeRequest {
    update.SSHKey = &key
    return update
}

// Update the request with a new status.
//   params:
//     status The new node status.
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithStatus(status NodeStatus) *UpdateNodeRequest {
    update.Status = &status
    return update
}

// Update the request with new edgenet info.
//   params:
//     edgenetAddress The new node edgenet address.
//     edgenetIP The new node edgenet IP
//   returns:
//     An update node request.
func (update *UpdateNodeRequest) WithEdgenet(address, ip string) *UpdateNodeRequest {
    update.EdgenetAddress = &address
    update.EdgenetIP = &ip
    return update
}

type FilterNodesRequest struct {
    // Labels required on the node.
    Labels * []string `json:"labels"`
}

// Create a new update node request.
func NewFilterNodesRequest() *FilterNodesRequest {
    return &FilterNodesRequest{}
}

func (filter *FilterNodesRequest) String() string {
    return fmt.Sprintf("%#v", filter)
}

func (filter * FilterNodesRequest) ByLabel(label string) *FilterNodesRequest {
    labels := []string{label}
    filter.Labels = &labels
    return filter
}

func (filter * FilterNodesRequest) ByLabels(labels []string) *FilterNodesRequest {
    filter.Labels = &labels
    return filter
}

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package entities

import "fmt"

// A network represents the higher entity in the hierarchy and contains all clusters the user can deploy
// apps on.
type Network struct {
    // Network identifier.
    ID string `json:"id, omitempty"`
    // Identifier of the edge overlay network
    EdgenetID string `json:"edgenetID, omitempty"`
    // Network name.
    Name string `json:"name, omitempty"`
    // Network description.
    Description string `json:"description, omitempty"`
    // Administrator name.
    AdminName string `json:"adminName, omitempty"`
    // Administrator phone.
    AdminPhone string `json:"adminPhone, omitempty"`
    // Administrator email.
    AdminEmail string `json:"adminEmail, omitempty"`
    // The network operator.
    Operator * User `json:"operator, omitempty"`
}

// Generate a new network.
//   params:
//     name The network name.
//     description The network description.
//     adminName The administrator name.
//     adminPhone The administrator phone.
//     adminEmail The administrator email.
//   returns:
//     A new network.
func NewNetwork(name string, description string, adminName string, adminPhone string, adminEmail string) * Network {
    uuid := GenerateUUID(NetworkPrefix)
    return NewNetworkWithID(uuid, name, description, adminName, adminPhone, adminEmail)
}

// Generate a new network.
//   params:
//     id The network identifier.
//     name The network name.
//     description The network description.
//     adminName The administrator name.
//     adminPhone The administrator phone.
//     adminEmail The administrator email.
//   returns:
//     A new network.
func NewNetworkWithID(id string, name string, description string, adminName string, adminPhone string, adminEmail string) * Network {
    network := &Network{
        ID: id,
        Name: name,
        Description: description,
        AdminName: adminName,
        AdminPhone: adminPhone,
        AdminEmail: adminEmail,
    }
    return network
}

// NetworkReducedInfo represents the higher entity in the hierarchy and contains all clusters the user can deploy
// apps on.
type NetworkReducedInfo struct {
    // Network identifier.
    ID string `json:"id, omitempty"`
    // Network name.
    Name string `json:"name, omitempty"`
}

// Generate a new network.
//   params:
//     id The network identifier.
//     name The network name.
//   returns:
//     A new network.
func NewNetworkReducedInfo(id string, name string) * NetworkReducedInfo {
    network := &NetworkReducedInfo{
        id,
        name,
    }
    return network
}

// Copy the data of a network adding an identifier.
//   params:
//     source The source network to be copied.
//   returns:
//     A pointer to the new network.
func AddNetworkIdentifier(source Network) * Network {
    // Passing source by value already makes us a copy
    uuid := GenerateUUID(NetworkPrefix)
    source.ID = uuid
    return &source
}

// Structure received by the REST endpoints.
type AddNetworkRequest struct {
    // Network name.
    Name string `json:"name, omitempty"`
    // Network description.
    Description string `json:"description, omitempty"`
    // Administrator name.
    AdminName string `json:"adminName, omitempty"`
    // Administrator phone.
    AdminPhone string `json:"adminPhone, omitempty"`
    // Administrator email.
    AdminEmail string `json:"adminEmail, omitempty"`
    // Identifier of the edge overlay network
    EdgenetID string `json:"edgenetID, omitempty"`
}

// Generate a new network request.
//   params:
//     name The network name.
//     description The network description.
//     adminName The administrator name.
//     adminPhone The administrator phone.
//     adminEmail The administrator email.
//     edgenetID The edge network identifier, if available
//   returns:
//     A new network request.
func NewAddNetworkRequest(name string, description string, adminName string, adminPhone string, adminEmail string) *AddNetworkRequest {
    request := &AddNetworkRequest{
        Name: name,
        Description: description,
        AdminName: adminName,
        AdminPhone: adminPhone,
        AdminEmail: adminEmail,
    }
    return request
}

// Check if the request is valid
//   params:
//     request The add network request.
//   returns:
//     Whether the request is valid.
func (request * AddNetworkRequest) IsValid() bool {
    return request.Name != ""
}

// Transform a network request into a network by adding the UUID.
//   params:
//     request The add network request.
//   returns:
//     A network with an UUID.
func ToNetwork(request AddNetworkRequest) *Network {
    uuid := GenerateUUID(NetworkPrefix)
    network := NewNetworkWithID(uuid, request.Name, request.Description, request.AdminName, request.AdminPhone, request.AdminEmail)

    // We do want to have a network that reflects all values in the request.
    // However, we don't want all (optional) request fields in the NewNetwork
    // functions, as it makes it
    // a) a lot of work to add new fields, and 
    // b) the functional call is not very descriptive.
    //
    // Setting fields explicitly by name makes it very clear what is happening.
    // In general, New* functions should return minimially initialized struct,
    // i.e., one that is valid (has all mandatory fields set). After creation,
    // it can be modified by the calling function directly. Fields that should
    // not be modified after creation should be made private.
    // We should at some point reflect this in our NewNetwork function by
    // removing everything but ID and Name from the arguments.
    network.EdgenetID = request.EdgenetID

    return network
}

func (request * AddNetworkRequest) String() string {
    return fmt.Sprintf("%#v", request)
}

func (request * Network) String() string {
    return fmt.Sprintf("%#v", request)
}

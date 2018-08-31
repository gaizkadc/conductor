//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
// User credentials.
//

package entities

// Data structure for user credentials.
type Credentials struct {
    // UUID
    UUID string `json:"uuid, omitempty"`
    // Public key
    PublicKey string `json:"publicKey, omitempty"`
    // Private key
    PrivateKey string `json:"privateKey, omitempty"`
    // Credentials description
    Description string `json:"description, omitempty"`
    // Type of the key pair
    TypeKey string `json:"typeKey, omitempty"`
}

// Return a new user credentials entry.
//  params:
//   UUID User identifier.
//   publicKey User public key.
//   privateKey User private key.
//   description Entry description.
//   typeKey: type of key we are storing
//  return:
//   New user credentials entry.
func NewCredentials(uuID string, publicKey string, privateKey string, description string, typeKey string) *Credentials {
    return &Credentials{UUID: uuID, PublicKey: publicKey, PrivateKey: privateKey, Description: description, TypeKey: typeKey}
}

// Data structure to request the addition of new credentials.
type AddCredentialsRequest struct {
    // UUID
    UUID string `json:"uuid, omitempty"`
    // Public key
    PublicKey string `json:"publicKey, omitempty"`
    // Private key
    PrivateKey string `json:"privateKey, omitempty"`
    // Credentials description
    Description string `json:"description, omitempty"`
    // Type of the key pair
    TypeKey string `json:"typeKey, omitempty"`
}

// Return a new user request for adding entries.
//  params:
//   UUID User identifier.
//   publicKey User public key.
//   privateKey User private key.
//   description Entry description.
//   typeKey: type of key we are storing
//  return:
//   New user credentials entry.
func NewAddCredentialsRequest(uuID string, publicKey string, privateKey string, description string, typeKey string) *AddCredentialsRequest {
    return &AddCredentialsRequest{UUID: uuID, PublicKey: publicKey, PrivateKey: privateKey, Description: description, TypeKey: typeKey}
}

func(req *AddCredentialsRequest) IsValid() bool {
    if req.UUID != "" && req.PrivateKey != "" && req.PublicKey != "" && req.Description != "" && req.TypeKey != "" {
        return true
    }
    return false
}
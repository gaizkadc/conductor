//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// API REST client for access entities

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/accessstorage"
    "github.com/daishogroup/system-model/server/access"
)

// AccessMockup for Access resources.
type AccessMockup struct {
    accessProvider * accessstorage.MockupUserAccessProvider
    accessManager  access.Manager
}

// NewAccessMockup creates a new Access mockup
func NewAccessMockup() Access {
    accessProvider := accessstorage.NewMockupUserAccessProvider()

    accessManager := access.NewManager(accessProvider)

    return &AccessMockup{ accessProvider, accessManager}
}

// ClearAccessMockup cleans the mockup.
func (m *AccessMockup) ClearAccessMockup(){
    m.accessProvider.Clear()
}

func (m *AccessMockup) InitAccessMockup() {
    roles := make([]entities.RoleType, 0)
    m.accessProvider.Add(*entities.NewUserAccess("userDefault", roles))
}


// AddAccess adds a new user access
//   params:
//     userID  The user identifier.
//     request The add user request.
//   returns:
//     The user access entity.
//     An error if the application cannot be added.
func(m *AccessMockup) AddAccess(userID string, request entities.AddUserAccessRequest) (*entities.UserAccess, derrors.DaishoError){
    return m.accessManager.AddAccess(userID, request)
}

// SetAccess sets the access roles for an existing user.
//   params:
//     userID  The user identifier.
//     request The add user request.
//   returns:
//     The user access entity.
//     An error if the application cannot be added.
func(m *AccessMockup) SetAccess(userID string, request entities.AddUserAccessRequest) (*entities.UserAccess, derrors.DaishoError){
    return m.accessManager.SetAccess(userID, request)
}

// GetAccess gets user access privilege entry for an existing user.
//   params:
//     userID The user identifier.
//   returns:
//     UserAccess values
//     An error in case the list cannot be retrieved.
func(m *AccessMockup) GetAccess(userID string) (* entities.UserAccess, derrors.DaishoError){
    return m.accessManager.GetAccess(userID)
}

// DeleteAccess deletes user access privilege.
//   params:
//     userID The user id.
//   returns:
//     Error if any.
func(m *AccessMockup) DeleteAccess(networkID string) derrors.DaishoError {
    return m.accessManager.DeleteAccess(networkID)
}

// ListAccess gets a list of user privileges.
//   returns:
//     Complete list of users with their access roles.
//     An error if the user does not exist.
func(m *AccessMockup) ListAccess() ([]entities.UserAccessReducedInfo, derrors.DaishoError){
    return m.accessManager.ListAccess()
}

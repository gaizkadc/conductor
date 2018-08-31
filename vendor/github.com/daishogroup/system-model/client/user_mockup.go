//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the user Client Mockup

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/accessstorage"
    "github.com/daishogroup/system-model/provider/oauthstorage"
    "github.com/daishogroup/system-model/provider/passwordstorage"
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/daishogroup/system-model/server/user"
)

type UserMockup struct {
    userProvider *userstorage.MockupUserProvider
    accessProvider *accessstorage.MockupUserAccessProvider
    passwordProvider *passwordstorage.MockupPasswordProvider
    oauthProvider   *oauthstorage.MockupOAuthProvider
    userMgr         user.Manager
}

func NewUserMockup() User {
    var userProvider = userstorage.NewMockupUserProvider()
    var accessProvider = accessstorage.NewMockupUserAccessProvider()
    var passwordProvider = passwordstorage.NewMockupPasswordProvider()
    var oauthProvider = oauthstorage.NewMockupOAuthProvider()
    var userMgr = user.NewManager(userProvider, accessProvider, passwordProvider, oauthProvider)

    return &UserMockup{userProvider, accessProvider, passwordProvider,
        oauthProvider, userMgr}
}

func (mockup *UserMockup) ClearMockup() {
    mockup.userProvider.Clear()
    mockup.accessProvider.Clear()
    mockup.passwordProvider.Clear()
    mockup.oauthProvider.Clear()
}

func (mockup *UserMockup) InitMockup() {
    mockup.ClearMockup()

    testId := "userDefault"
    testName := "userDefaultName"
    testEmail := "userDefault@email.com"
    testPhone := "999 999 999"
    testPassword := "apassword"

    mockup.userProvider.Add(* entities.NewUserWithID(testId, testName, testPhone, testEmail, testUserCreationTime, testUserExpirationTime))
    mockup.accessProvider.Add(*entities.NewUserAccess(testId, []entities.RoleType{entities.GlobalAdmin}))
    pwd, _ := entities.NewPassword(testId, &testPassword)
    mockup.passwordProvider.Add(*pwd)

}

func (mockup *UserMockup) FixUserAccess(userID string, password string, userRole entities.RoleType) derrors.DaishoError {
    err := mockup.accessProvider.Update(*entities.NewUserAccess(userID, []entities.RoleType{userRole}))
    if err != nil {
        return err
    }
    pwd, err := entities.NewPassword(userID, &password)
    if err != nil {
        return err
    }
    return mockup.passwordProvider.Update(*pwd)
}

// Add a new user
// params:
//    user User request.
// return:
//    Generated user.
//    Error if any.
func (mockup *UserMockup) Add(user entities.AddUserRequest) (*entities.User, derrors.DaishoError) {
    return mockup.userMgr.AddUser(user)
}

// Get an existing user data.
// params:
//    userId User identifier.
// return:
//    Found user data.
//    Error if any.
func (mockup *UserMockup) Get(userId string) (*entities.User, derrors.DaishoError) {
    return mockup.userMgr.GetUser(userId)
}

// Delete an existing user entry.
// params:
//    userId User identifier.
// return:
//    Error if any.
func (mockup *UserMockup) Delete(userId string) derrors.DaishoError {
    return mockup.userMgr.DeleteUser(userId)
}

// Update an existing user entry.
// params:
//   userId User identifier.
//   updateRequest Update request.
// return:
//   Error if any.
func (mockup *UserMockup) Update(userId string, update entities.UpdateUserRequest) (*entities.User, derrors.DaishoError) {
    return mockup.userMgr.UpdateUser(userId, update)
}

// Return the list of users with their corresponding role/privileges.
//  return:
//   List of users with their roles.
func (mockup *UserMockup) ListUsers() ([]entities.UserExtended, derrors.DaishoError) {
    return mockup.userMgr.ListUsers()
}
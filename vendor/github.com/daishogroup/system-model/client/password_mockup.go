//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the oauth Client interface

package client

import (
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/passwordstorage"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/server/password"
)


type PasswordMockup struct {
    passwordProvider *passwordstorage.MockupPasswordProvider
    passwordMgr      password.Manager
}


func NewPasswordMockup() Password{
    var passwordProvider = passwordstorage.NewMockupPasswordProvider()
    var passwordMgr = password.NewManager(passwordProvider)
    return &PasswordMockup{passwordProvider, passwordMgr}
}

func (mockup *PasswordMockup) ClearMockup() {
    mockup.passwordProvider.Clear()
}

func (mockup *PasswordMockup) InitMockup() {
    mockup.ClearMockup()

    testId := "userDefault"
    testPassword := "apassword"

    pwd, _ := entities.NewPassword(testId, &testPassword)
    mockup.passwordProvider.Add(*pwd)

}

// Get a user oauth
// params:
//  userID user identifier.
// return:
//  Password entity.
//  Error if any.
func(m *PasswordMockup)  GetPassword(userID string) (*entities.Password, derrors.DaishoError){
    return m.passwordMgr.GetPassword(userID)

}

// Set a oauth entry with a new value.
//  params:
//   oauth New oauth entity.
//  return:
//   Error if any.
func(m *PasswordMockup) SetPassword(password entities.Password) derrors.DaishoError{
    return m.passwordMgr.SetPassword(password)
}

// Delete an existing oauth.
//  params:
//   userID the user identifier
//  return:
//   Error if any.
func(m *PasswordMockup) DeletePassword(userID string) derrors.DaishoError{
    return m.passwordMgr.DeletePassword(userID)
}

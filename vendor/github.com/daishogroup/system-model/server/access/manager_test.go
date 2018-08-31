//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Access manager tests.

package access

import (
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/entities"
    "testing"
    "github.com/daishogroup/system-model/provider/accessstorage"
)

const (
    testUserId = "userId"
    testUserRole = entities.GlobalAdmin
)

type ManagerHelper struct {
    accessProvider *accessstorage.MockupUserAccessProvider
    accessMgr      Manager
}

func NewManagerHelper() ManagerHelper {
    var accessProvider = accessstorage.NewMockupUserAccessProvider()
    var nodeMgr = NewManager(accessProvider)
    return ManagerHelper{accessProvider,
        nodeMgr}
}

type ManagerTestSuite struct {
    suite.Suite
    manager ManagerHelper
}

func (helper *ManagerTestSuite) SetupSuite() {
    managerHelper := NewManagerHelper()
    helper.manager = managerHelper
}

// The SetupTest method is called before every test on the suite.
func (helper *ManagerTestSuite) SetupTest() {
    helper.manager.accessProvider.Clear()
    helper.manager.accessProvider.Add(*entities.NewUserAccess(testUserId,[]entities.RoleType{entities.DeveloperType}))
}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(ManagerTestSuite))
}


func (helper *ManagerTestSuite) TestAddUserAccess() {
    userRequest := entities.NewAddUserAccessRequest([]entities.RoleType{testUserRole})
    added, err := helper.manager.accessMgr.AddAccess(testUserId,*userRequest)
    helper.Nil(err, "user must be added")
    helper.NotNil(added, "user must be returned")
    helper.Equal(testUserId, added.UserID, "user Id does not match")
    helper.Equal(2, len(added.Roles))
    helper.Equal(testUserRole, added.Roles[1], "user role does not match")
}

func (helper *ManagerTestSuite) TestSetUserAccess() {
    userRequest := entities.NewAddUserAccessRequest([]entities.RoleType{entities.OperatorType})
    added, err := helper.manager.accessMgr.SetAccess(testUserId,*userRequest)
    helper.Nil(err, "user must be set")
    helper.NotNil(added, "user must be returned")
    helper.Equal(testUserId, added.UserID, "user Id does not match")
    helper.Equal(1, len(added.Roles))
    helper.Equal(entities.OperatorType, added.Roles[0], "user role does not match")
}


func (helper *ManagerTestSuite) TestGetUserAccess() {
    // Retrieve the user from the mockup
    helper.manager.accessProvider.Add(*entities.NewUserAccess(testUserId, []entities.RoleType{testUserRole}))

    retrieved, err := helper.manager.accessMgr.GetAccess(testUserId)
    helper.Nil(err, "unexpected error")
    helper.Equal(testUserId, retrieved.UserID, "unexpected user id")
    helper.Equal(testUserId, retrieved.UserID, "user Id does not match")
}

func (helper *ManagerTestSuite) TestGetNonExistingUser() {
    retrieved, err := helper.manager.accessMgr.GetAccess("i dont exist")
    helper.NotNil(err, "an error must be returned")
    helper.Nil(retrieved,"returned user must be nil")
}

func (helper *ManagerTestSuite) TestDeleteUser() {
    // Retrieve the user from the mockup
    helper.manager.accessProvider.Add(*entities.NewUserAccess(testUserId, []entities.RoleType{testUserRole}))

    err := helper.manager.accessMgr.DeleteAccess(testUserId)
    helper.Nil(err, "unexpected error")

    //try to retrieve it
    retrieved, err := helper.manager.accessMgr.GetAccess(testUserId)
    helper.NotNil(err, "expected errors")
    helper.Nil(retrieved, "user was retrieved after deletion")
}

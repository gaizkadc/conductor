//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// User manager tests.

package user

import (
    "testing"
    "time"
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/accessstorage"
    "github.com/daishogroup/system-model/provider/passwordstorage"
    "github.com/daishogroup/system-model/provider/oauthstorage"
)

const (
    testUserId = "userId"
    testUserName = "userName"
    testUserPhone = "99999999"
    testUserEmail = "user@email.com"
    testUserPassword = "thepassword"
)

var testCreationTime = time.Date(2010, time.January,1,1,1,1,0, time.UTC)
var testExpirationTime = testCreationTime.Add(time.Hour)

type ManagerHelper struct {
    userProvider *userstorage.MockupUserProvider
    accessProvider *accessstorage.MockupUserAccessProvider
    passwordProvider *passwordstorage.MockupPasswordProvider
    oauthProvider    *oauthstorage.MockupOAuthProvider
    userMgr      Manager
}

func NewManagerHelper() ManagerHelper {
    var userProvider = userstorage.NewMockupUserProvider()
    var accessProvider = accessstorage.NewMockupUserAccessProvider()
    var passwordProvider = passwordstorage.NewMockupPasswordProvider()
    var oauthProvider = oauthstorage.NewMockupOAuthProvider()
    var nodeMgr = NewManager(userProvider, accessProvider, passwordProvider, oauthProvider)
    return ManagerHelper{userProvider, accessProvider,
        passwordProvider, oauthProvider,nodeMgr}
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
    helper.manager.userProvider.Clear()
    helper.manager.accessProvider.Clear()
    helper.manager.passwordProvider.Clear()
    helper.manager.oauthProvider.Clear()
}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(ManagerTestSuite))
}


func (helper *ManagerTestSuite) TestAddUser() {
    userRequest := entities.NewAddUserRequest(testUserId, testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime)
    added, err := helper.manager.userMgr.AddUser(*userRequest)
    helper.Nil(err, "user must be added")
    helper.NotNil(added, "user must be returned")
    helper.Equal(testUserId, added.ID, "user Id does not match")
    helper.Equal(testUserName, added.Name, "user name does not match")
    helper.Equal(testUserEmail, added.Email, "user email does not match")
    helper.Equal(testUserPhone, added.Phone, "user phone does not match")
    helper.Equal(testCreationTime, added.CreationTime.Time, "creation time does not match")
    helper.Equal(testExpirationTime, added.ExpirationTime.Time, "expiration time does not match")

    // There must be an empty access
    access, err := helper.manager.accessProvider.RetrieveAccess(testUserId)
    helper.Equal(0, len(access.Roles), "unexpected number of roles")
    // There must be an empty password
    password, err := helper.manager.passwordProvider.RetrievePassword(testUserId)
    helper.Equal(testUserId, password.UserID)
    // There must be an empty oauth entry.
    secrets, err := helper.manager.oauthProvider.Retrieve(testUserId)
    helper.Equal(testUserId, secrets.UserID)
}


func (helper *ManagerTestSuite) TestGetUser() {
    // Retrieve the user from the mockup
    helper.manager.userProvider.Add(*entities.NewUserWithID(testUserId, testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime))

    retrieved, err := helper.manager.userMgr.GetUser(testUserId)
    helper.Nil(err, "unexpected error")
    helper.Equal(testUserId, retrieved.ID, "unexpected user id")
    helper.Equal(testUserId, retrieved.ID, "user Id does not match")
    helper.Equal(testUserName, retrieved.Name, "user name does not match")
    helper.Equal(testUserEmail, retrieved.Email, "user email does not match")
    helper.Equal(testUserPhone, retrieved.Phone, "user phone does not match")
    helper.Equal(testCreationTime, retrieved.CreationTime.Time, "creation time does not match")
    helper.Equal(testExpirationTime, retrieved.ExpirationTime.Time, "expiration time does not match")
}

func (helper *ManagerTestSuite) TestGetNonExistingUser() {
    retrieved, err := helper.manager.userMgr.GetUser("i dont exist")
    helper.NotNil(err, "an error must be returned")
    helper.Nil(retrieved,"returned user must be nil")
}

func (helper *ManagerTestSuite) TestDeleteUser() {
    // Retrieve the user from the mockup
    helper.manager.userProvider.Add(*entities.NewUserWithID(testUserId, testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime))
    helper.manager.accessProvider.Add(*entities.NewUserAccess(testUserId,[]entities.RoleType{entities.DeveloperType}))
    pass, _ := entities.NewPassword(testUserId,nil)
    helper.manager.passwordProvider.Add(*pass)
    helper.manager.oauthProvider.Add(entities.NewOAuthSecrets(testUserId))

    err := helper.manager.userMgr.DeleteUser(testUserId)
    helper.Nil(err, "unexpected error")

    //try to retrieve it
    retrieved, err := helper.manager.userMgr.GetUser(testUserId)
    helper.NotNil(err, "expected errors")
    helper.Nil(retrieved, "user was retrieved after deletion")
    // Check there is no access entry
    helper.False(helper.manager.accessProvider.Exists(testUserId), "access still exists")
    // Check there is no password entry
    helper.False(helper.manager.passwordProvider.Exists(testUserId), "password still exists")
    // Check there is no secrets entry
    helper.False(helper.manager.oauthProvider.Exists(testUserId), "oauth secret still exists")
}

func (helper *ManagerTestSuite) TestUpdateUser() {
    // Retrieve the user from the mockup
    helper.manager.userProvider.Add(*entities.NewUserWithID(testUserId, testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime))

    // try to update it
    updateRequest := entities.NewUpdateUserRequest().WithName("1").
        WithPhone("2").WithEmail("3")
    updated, err := helper.manager.userMgr.UpdateUser(testUserId, *updateRequest)
    helper.Suite.Nil(err, "Unexpected error")
    helper.Equal(testUserId, updated.ID, "unexpected user id")
    helper.Equal(testUserId, updated.ID, "user Id does not match")
    helper.Equal("1", updated.Name, "user name does not match")
    helper.Equal("2", updated.Phone, "user phone does not match")
    helper.Equal("3", updated.Email, "user email does not match")
}

func (helper *ManagerTestSuite) TestListUsers() {

    // Inititate the whole thing
    helper.manager.userProvider.Add(*entities.NewUserWithID(testUserId, testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime))
    helper.manager.accessProvider.Add(*entities.NewUserAccess(testUserId, []entities.RoleType{entities.GlobalAdmin}))

    helper.manager.userProvider.Add(*entities.NewUserWithID("a", testUserName, testUserPhone, testUserEmail,
        testCreationTime, testExpirationTime))
    helper.manager.accessProvider.Add(*entities.NewUserAccess("a", []entities.RoleType{entities.GlobalAdmin}))


    entries, err := helper.manager.userMgr.ListUsers()

    helper.Nil(err, "unexpected error")
    helper.NotNil(entries, "unexpected nil entries")
    helper.Equal(2, len(entries), "unexpected number of returned elements")
}



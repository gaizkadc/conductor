//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Password manager testing.


package oauth

import (
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/provider/oauthstorage"
    "github.com/daishogroup/system-model/entities"
)

const(
    TestUserID   = "testUser"
    TestPassword = "testPassword123"
)

type ManagerHelper struct {
    oauthProvider *oauthstorage.MockupOAuthProvider
    oauthMgr      Manager
}

func NewManagerHelper() ManagerHelper {
    var oauthProvider = oauthstorage.NewMockupOAuthProvider()
    var oauthMgr = NewManager(oauthProvider)
    return ManagerHelper{oauthProvider, oauthMgr}
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
    helper.manager.oauthProvider.Clear()
}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(ManagerTestSuite))
}

func (helper *ManagerTestSuite) TestAddSecret() {
    // Add initial empty user
    p := entities.NewOAuthSecrets(TestUserID)
    err := helper.manager.oauthProvider.Add(p)
    helper.Nil(err, "unexpected error")

    // Add secrets
    request := entities.NewOAuthAddEntryRequest("app1", "client1", "secret1")
    err = helper.manager.oauthMgr.SetSecret(TestUserID, request)
    helper.Nil(err, "unexpected error")

    // Now try to retrieve it
    retrieved, err := helper.manager.oauthMgr.GetSecrets(TestUserID)
    helper.Nil(err, "unexpected error")
    helper.Equal(TestUserID, retrieved.UserID, "unexpected userID")
    helper.Equal("client1", retrieved.Entries["app1"].ClientID, "unexpected clientID")
    helper.Equal("secret1", retrieved.Entries["app1"].Secret, "unexpected secret")
}

func (helper *ManagerTestSuite) TestDeleteSecret() {
    // Add initial empty user
    p := entities.NewOAuthSecrets(TestUserID)
    err := helper.manager.oauthProvider.Add(p)
    helper.Nil(err, "unexpected error")

    // Remove
    err = helper.manager.oauthProvider.Delete(TestUserID)
    helper.Nil(err, "unexpected error when deleting entries")

    // Try to retrieve it
    retrieved, err := helper.manager.oauthMgr.GetSecrets(TestUserID)
    helper.NotNil(err, "unexpected error")
    helper.Nil(retrieved, "unexpected entry values")
}


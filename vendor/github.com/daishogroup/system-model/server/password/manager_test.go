//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Password manager testing.


package password

import (
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/provider/passwordstorage"
    "github.com/daishogroup/system-model/entities"
)

const(
    TestUserID   = "testUser"
    TestPassword = "testPassword123"
)

type ManagerHelper struct {
    passwordProvider *passwordstorage.MockupPasswordProvider
    passwordMgr      Manager
}

func NewManagerHelper() ManagerHelper {
    var passwordProvider = passwordstorage.NewMockupPasswordProvider()
    var passwordMgr = NewManager(passwordProvider)
    return ManagerHelper{passwordProvider, passwordMgr}
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
    helper.manager.passwordProvider.Clear()
}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(ManagerTestSuite))
}

func (helper *ManagerTestSuite) TestSetPassword() {
    p, err := entities.NewPassword(TestUserID, nil)
    helper.Nil(err, "unexpected error")
    helper.manager.passwordProvider.Add(*p)

    thePassword := TestPassword
    p, err = entities.NewPassword(TestUserID, &thePassword)
    err = helper.manager.passwordMgr.SetPassword(*p)
    helper.Nil(err, "unexpected error")
    // now try to retrieve it
    retrieved, err := helper.manager.passwordMgr.GetPassword(TestUserID)
    helper.Nil(err, "unexpected error")
    helper.Equal(TestUserID, retrieved.UserID, "unexpected user id")
    helper.True(retrieved.CompareWith(TestPassword), "passwords do not match")
    helper.False(retrieved.CompareWith("error"), "passwrods should not match")
}


func (helper *ManagerTestSuite) TestDeletePassword() {
    thePassword := TestPassword
    p, err := entities.NewPassword(TestUserID, &thePassword)
    helper.Nil(err, "unexpected error")
    err = helper.manager.passwordProvider.Add(*p)
    helper.Nil(err, "unexpected error")
    // now try to retrieve it
    retrieved, err := helper.manager.passwordMgr.GetPassword(TestUserID)
    helper.Nil(err, "unexpected error")
    helper.Equal(TestUserID, retrieved.UserID, "unexpected user id")
    helper.True(retrieved.CompareWith(TestPassword), "passwords do not match")
    helper.False(retrieved.CompareWith("error"), "password should not match")
    // Now delete it
    err = helper.manager.passwordMgr.DeletePassword(TestUserID)
    helper.Nil(err, "unexpected error")
    // Try to retrieve it
    retrieved, err = helper.manager.passwordMgr.GetPassword(TestUserID)
    helper.NotNil(err, "unexpected error")
    helper.Nil(retrieved, "unexpected returned entry")

}

//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Credentials manager test.


package credentials

import (
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/provider/credentialsstorage"
    "github.com/daishogroup/system-model/entities"
)

type ManagerHelper struct {
    credentialsProvider *credentialsstorage.MockupCredentialsProvider
    manager         Manager
}

type ManagerTestSuite struct {
    suite.Suite
    manager ManagerHelper
}

func NewManagerHelper() ManagerHelper {
    var credentialsProvider = credentialsstorage.NewMockupCredentialsProvider()
    var credentialsManager = Manager{credentialsProvider: credentialsProvider}
    return ManagerHelper{credentialsProvider, credentialsManager}
}

func (helper *ManagerTestSuite) SetupSuite() {
    managerHelper := NewManagerHelper()
    helper.manager = managerHelper
}

// The SetupTest method is called before every test on the suite.
func (helper *ManagerTestSuite) SetupTest() {
    helper.manager.credentialsProvider.Clear()
}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(ManagerTestSuite))
}

func (helper *ManagerTestSuite) TestAddCredentials() {
    newCredential := entities.NewAddCredentialsRequest(Testuuid, TestPublicKey, TestPrivateKey, TestDescription, TestTypeKey)
    err := helper.manager.manager.AddCredentials(*newCredential)
    helper.Nil(err, "unexpected error")

    // Get it
    returned, err := helper.manager.manager.GetCredentials(Testuuid)
    helper.Nil(err, "unexpected error")
    helper.Equal(Testuuid, returned.UUID)
    helper.Equal(TestPublicKey, returned.PublicKey)
    helper.Equal(TestPrivateKey, returned.PrivateKey)
    helper.Equal(TestDescription, returned.Description)
    helper.Equal(TestTypeKey, returned.TypeKey)

    // Delete it
    err = helper.manager.manager.DeleteCredentials(Testuuid)
    helper.Nil(err, "unexpected error when trying to delete")

    // Try to get it again
    returned2, err := helper.manager.manager.GetCredentials(Testuuid)
    helper.NotNil(err, "we expected an error")
    helper.Nil(returned2, "empty object was expected")
}

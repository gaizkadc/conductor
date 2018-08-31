//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// User manager tests.

package session

import (
    "testing"
    "time"
    "net/http"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/sessionstorage"
)

const (
    testUserId = "userId"
    testDomain = "mygooddomain.org"
    testCookieName = "testCookie"

)

var testCreationTime = time.Now()
var testExpirationTime = testCreationTime.Add(time.Hour)

type ManagerHelper struct {
    sessionProvider *sessionstorage.MockupSessionProvider
    sessionMgr      Manager
}

func NewManagerHelper() ManagerHelper {
    var sessionProvider = sessionstorage.NewMockupSessionProvider()
    var sessionMgr = NewManager(sessionProvider)
    return ManagerHelper{sessionProvider, sessionMgr}
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
    helper.manager.sessionProvider.Clear()
}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(ManagerTestSuite))
}


func (helper *ManagerTestSuite) TestAddSessionUser() {
    testSession := entities.NewSession(testUserId, testExpirationTime)
    testCookie := http.Cookie{
        Domain: "mygooddomain.org",
    }
    testSession.AddCookie(testCookieName, testCookie)
    testSessionRequest := entities.NewAddSessionRequest(*testSession)

    added, err := helper.manager.sessionMgr.AddSession(*testSessionRequest)
    helper.Nil(err, "session must be added")
    helper.NotNil(added, "session must be returned")
    helper.NotNil(added.ID)
    helper.Equal(testUserId, added.UserID)
    helper.Equal(added.ExpirationDate, testExpirationTime)
}


func (helper *ManagerTestSuite) TestGetSession() {
    // Add a session to be retrieved.
    // Retrieve the user from the mockup
    testSession := entities.NewSession(testUserId, testExpirationTime)
    testCookie := http.Cookie{
        Domain: "mygooddomain.org",
    }
    testSession.AddCookie(testCookieName, testCookie)
    err := helper.manager.sessionProvider.Add(*testSession)
    helper.Nil(nil, "unexpected error")

    retrieved, err := helper.manager.sessionMgr.GetSession(testSession.ID)
    helper.Nil(err, "unexpected error")
    helper.Equal(testUserId, retrieved.UserID, "unexpected user id")
    helper.Equal(testExpirationTime.UTC(), retrieved.ExpirationDate.UTC(), "creation time does not match")
    helper.Equal("mygooddomain.org", retrieved.Cookies[testCookieName].Domain)
}


func (helper *ManagerTestSuite) TestDeleteSession() {
    // Add a session to be retrieved.
    // Retrieve the user from the mockup
    testSession := entities.NewSession(testUserId, testExpirationTime)
    testCookie := http.Cookie{
        Domain: "mygooddomain.org",
    }
    testSession.AddCookie(testCookieName, testCookie)
    err := helper.manager.sessionProvider.Add(*testSession)
    helper.Nil(err, "unexpected error")

    // delete it
    err = helper.manager.sessionMgr.DeleteSession(testSession.ID)
    helper.Nil(err, "expected errors")

    // try to retrieve it
    // Check there is no access entry
    helper.False(helper.manager.sessionProvider.Exists(testSession.ID), "session still exists")

}

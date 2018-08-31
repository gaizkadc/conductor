//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the Session Mockup

package client

import (
    "time"
    "net/http"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/provider/sessionstorage"
    "github.com/daishogroup/system-model/server/session"
    "github.com/daishogroup/system-model/entities"
)

type SessionMockup struct {
    sessionProvider *sessionstorage.MockupSessionProvider
    sessionMgr session.Manager
}


func NewSessionMockup() Session{
    var sessionProvider = sessionstorage.NewMockupSessionProvider()
    var sessionMgr = session.NewManager(sessionProvider)

        //session.NewManager(sessionProvider)
    return &SessionMockup{sessionProvider, sessionMgr}
}

func (mockup *SessionMockup) ClearMockup() {
    mockup.sessionProvider.Clear()
}

func (mockup *SessionMockup) InitMockup() {
    mockup.ClearMockup()

    s := entities.NewSession("testUser", time.Now().Add(time.Hour))
    s.AddCookie("theCookie", http.Cookie{Domain: "localhost"})

    mockup.sessionProvider.Add(*s)

}


// Add a new user
// params:
//    session session request.
// return:
//    Generated session.
//    Error if any.
func (mockup *SessionMockup) Add(session entities.AddSessionRequest) (*entities.Session, derrors.DaishoError) {
    return mockup.sessionMgr.AddSession(session)
}

// Get an existing session.
// params:
//    sessionID Session identifier.
// return:
//    Found session data.
//    Error if any.
func (mockup *SessionMockup) Get(sessionID string) (*entities.Session, derrors.DaishoError) {
    return mockup.sessionMgr.GetSession(sessionID)
}

// Delete an existing session entry.
// params:
//    sessionID User identifier.
// return:
//    Error if any.
func (mockup *SessionMockup) Delete(sessionID string) derrors.DaishoError {
    return mockup.sessionMgr.DeleteSession(sessionID)
}

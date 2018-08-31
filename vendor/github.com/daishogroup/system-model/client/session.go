//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the Session Client interface

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Session interface that defines the operations available on Users.
type Session interface {

    // Add a new user
    // params:
    //    session session request.
    // return:
    //    Generated session.
    //    Error if any.
    Add(session entities.AddSessionRequest) (*entities.Session, derrors.DaishoError)

    // Get an existing session.
    // params:
    //    sessionID Session identifier.
    // return:
    //    Found session data.
    //    Error if any.
    Get(sessionID string) (*entities.Session, derrors.DaishoError)

    // Delete an existing session entry.
    // params:
    //    sessionID User identifier.
    // return:
    //    Error if any.
    Delete(sessionID string) derrors.DaishoError

}
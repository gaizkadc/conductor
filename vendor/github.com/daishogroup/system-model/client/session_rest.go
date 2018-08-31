//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the Session Rest Client

package client

import (
    "fmt"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/dhttp"
)

const (
    SessionAddURI    = "/api/v0/session/add"
    SessionGetURI    = "/api/v0/session/%s/get"
    SessionDeleteURI = "/api/v0/session/%s/delete"
)

// Client Rest for user resources.
type SessionRest struct {
    client dhttp.Client
}

// Deprecated: Use NewUserClientRest
func NewSessionRest(basePath string) Session {
    return NewSessionClientRest(ParseHostPort(basePath))
}

func NewSessionClientRest(host string, port int) Session {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host, port))
    return &SessionRest{rest}
}


// Add a new user
// params:
//    session session request.
// return:
//    Generated session.
//    Error if any.
func (rest *SessionRest) Add(session entities.AddSessionRequest) (*entities.Session, derrors.DaishoError) {
    response := rest.client.Post(SessionAddURI, session, new(entities.Session))
    if response.Error != nil {
        return nil, response.Error
    } else {
        result := response.Result.(*entities.Session)
        return result, nil
    }
}

// Get an existing session.
// params:
//    sessionID Session identifier.
// return:
//    Found session data.
//    Error if any.
func (rest *SessionRest) Get(sessionID string) (*entities.Session, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(SessionGetURI, sessionID), new(entities.Session))
    if response.Error != nil {
        return nil, response.Error
    } else {
        result := response.Result.(*entities.Session)
        return result, nil
    }
}


// Delete an existing session entry.
// params:
//    sessionID User identifier.
// return:
//    Error if any.
func (rest *SessionRest) Delete(sessionID string) derrors.DaishoError {
    response := rest.client.Delete(fmt.Sprintf(SessionDeleteURI, sessionID), new(entities.SuccessfulOperation))
    if response.Error != nil {
        return response.Error
    } else {
        return nil
    }
}


//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Definition of a user session.

package entities

import (
    "time"
    "net/http"
    "github.com/satori/go.uuid"
)

// User session structure. Right now this is an abstraction of gorilla sessions implementation.
type Session struct{
    // User session ID.
    ID string `json:"id, omitempty"`
    // User id to be used in this session.
    UserID string `json:"userId, omitempty"`
    // Expiration date.
    ExpirationDate time.Time `json:"expirationDate, omitempty"`
    // Session cookies
    Cookies map[string]http.Cookie `json:"cookies, omitempty"`
}

// Creates a new session with a predefined session id, with empty cookies.
// params:
//  userID user identifier.
//  expirationDate Date this session will expire
func NewSession(userID string, expirationDate time.Time) *Session {
    return &Session{uuid.NewV4().String(), userID,
    expirationDate,make(map[string]http.Cookie) }
}

// Add a cookie to the session.
// params:
//  cookieName name of the cookie
//  newCookie cookie to be added
func (s *Session) AddCookie(cookieName string, newCookie http.Cookie) {
    s.Cookies[cookieName] = newCookie
}

// It deletes a cookie from the existing list. If the cookie does not exist, nothing is done.
func (s *Session) DeleteCookie(cookieName string) {
    delete(s.Cookies, cookieName)
}

// Indicate if this session has already expired.
//  return:
//   Boolean value indicating if this session has expired or not.
func (s *Session) IsExpired() bool {
    return s.ExpirationDate.Before(time.Now())
}


// Request to include a new session
type AddSessionRequest struct {
    NewSession Session `json:"newSession, omitempty"`
}

func NewAddSessionRequest(session Session) *AddSessionRequest {
    return &AddSessionRequest{NewSession: session}
}

// Check if the request is already valid or not.
//  return:
//   True if the request is valid.
func (asr *AddSessionRequest) IsValid() bool {
    return !asr.NewSession.IsExpired()
}
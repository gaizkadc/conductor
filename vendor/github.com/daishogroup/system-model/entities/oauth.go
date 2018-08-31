//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Entities required to use OAuth2 user authentication.

package entities

import (
   "github.com/daishogroup/derrors"
)

// Data structure containing all the user OAuth entries.
type OAuthSecrets struct {
    UserID      string  `json:"userID"`
    Entries     map[string]OAuthEntry `json:"entries"`
}

type OAuthEntry struct {
    ClientID    string  `json:"clientID,omitempty"`
    Secret      string  `json:"secret,omitempty"`
}


// Create a new OAuthSecrets structure with empty entries.
//  params:
//   userID User identifier.
//  result:
//   OAuthSecrets struct with an empty entry set.
func NewOAuthSecrets (userID string) OAuthSecrets {
    return OAuthSecrets{
        UserID: userID,
        Entries: make(map[string]OAuthEntry,0),
    }
}

// Create a new OAuthEntry
//  params:
//   clientID client identifer.
//   secret OAuth secret for the app and client
//  return:
//   New entry instance.
func NewOAuthEntry (clientID string, secret string) OAuthEntry {
    return OAuthEntry{
        Secret: secret,
        ClientID: clientID,
    }
}

// Add a new OAuthEntry to the secrets set.
//  params:
//   appName The application name.
//   entry The new oauth entry.
func(s *OAuthSecrets) AddEntry(appName string, entry OAuthEntry) derrors.DaishoError{
    _, exists := s.Entries[appName]
    if exists {
        return derrors.NewOperationError("OAuth entry already exists")
    }
    s.Entries[appName] = entry
    return nil
}

// Remove an existing OAuthEntry.
//  params:
//   appName Name of the application to remove.
//  return:
//   Error if any.
func(s *OAuthSecrets) DeleteEntry(appName string) derrors.DaishoError{
    _, exists := s.Entries[appName]
    if !exists {
        return derrors.NewOperationError("OAuth entry does not exist")
    }
    delete(s.Entries, appName)
    return nil
}

// Request to add a new OAuth secrets entry for an app.
type OAuthAddEntryRequest struct {
    AppName     string  `json:"appName"`
    ClientID    string  `json:"clientID"`
    Secret      string  `json:"secret"`
}

func NewOAuthAddEntryRequest(appName string, clientID string, secret string) OAuthAddEntryRequest {
    return OAuthAddEntryRequest{
        AppName:    appName,
        ClientID:   clientID,
        Secret:     secret,
    }
}
package sessionstorage

// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the session provider specification.

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Provider definition of the node persistence-related methods.
type Provider interface {

    // Add a new session to the system.
    //   params:
    //     session The session to be added
    //   returns:
    //     An error if the session cannot be added.
    Add(session entities.Session) derrors.DaishoError

    // Check if a session exists.
    //   params:
    //     sessionID The session identifier.
    //   returns:
    //     Whether the session exists or not.
    Exists(sessionID string) bool

    // Retrieve a given session.
    //   params:
    //     sessionID The session identifier.
    //   returns:
    //     The session.
    //     An error if the session cannot be retrieved.
    Retrieve(sessionID string) (* entities.Session, derrors.DaishoError)

    // Delete a given session.
    //   params:
    //     sessionID The session identifier.
    //   returns:
    //     An error if the session cannot be removed.
    Delete(sessionID string) derrors.DaishoError

    // Update a session in the system.
    //   params:
    //     node The new session information.
    //   returns:
    //     An error if the session cannot be updated.
    Update(session entities.Session) derrors.DaishoError

    // Dump obtains the list of all session in the system.
    //   returns:
    //     The list of session.
    //     An error if the list cannot be retrieved.
    Dump() ([] entities.Session, derrors.DaishoError)
}


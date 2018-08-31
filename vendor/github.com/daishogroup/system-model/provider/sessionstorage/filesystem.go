//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
//

package sessionstorage

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

const dirCreatePerm = 0777
const fileCreatePerm = 0666

// FileSystemProvider implements a file-backed provider.
type FileSystemProvider struct {
    sync.Mutex
    BasePath string
}

// NewFileSystemProvider creates a new FileSystemProvider.
//   params:
//     basePath The base path.
//   returns:
//     A file system provider.
func NewFileSystemProvider(basePath string) *FileSystemProvider {
    os.MkdirAll(basePath+"/sessions/", dirCreatePerm)
    return &FileSystemProvider{BasePath: basePath}
}

func (fs *FileSystemProvider) getPath(path string) string {
    newPath := fmt.Sprintf("%s/sessions/%s", fs.BasePath, path)
    return newPath
}

func (fs *FileSystemProvider) getSubPath(parent string, sessionID string) string {
    newPath := fmt.Sprintf("%s/sessions/%s/%s", fs.BasePath, parent, sessionID)
    return newPath
}

// Add a new session to the system.
//   params:
//     session The session to be added
//   returns:
//     An error if the session cannot be added.
func (fs *FileSystemProvider) Add(session entities.Session) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if ! fs.unsafeExists(session.ID) {
        toWrite, err := json.Marshal(session)
        if err == nil {
            return derrors.AsDaishoError(ioutil.WriteFile(fs.getPath(session.ID), toWrite, fileCreatePerm), errors.IOError)
        }
        return derrors.NewEntityError(session, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.SessionAlreadyExists).WithParams(session)
}

// Exists checks if a session exists in the system.
//   params:
//     sessionID The session identifier.
//   returns:
//     Whether the session exists or not.
func (fs *FileSystemProvider) Exists(sessionID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExists(sessionID)
}

func (fs *FileSystemProvider) unsafeExists(sessionID string) bool {
    _, err := os.Stat(fs.getPath(sessionID))
    return err == nil
}

// Retrievesession retrieves a given session.
//   params:
//     sessionID The session identifier.
//   returns:
//     The session.
//     An error if the session cannot be retrieved.
func (fs *FileSystemProvider) Retrieve(sessionID string) (*entities.Session, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeRetrieve(sessionID)
}

func (fs *FileSystemProvider) unsafeRetrieve(sessionID string) (*entities.Session, derrors.DaishoError) {
    if fs.unsafeExists(sessionID) {
        content, err := ioutil.ReadFile(fs.getPath(sessionID))
        if err == nil {
            session := entities.Session{}
            err = json.Unmarshal(content, &session)
            if err == nil {
                return &session, nil
            }
            return nil, derrors.NewOperationError(errors.UnmarshalError).WithParams(sessionID)
        }
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(sessionID)
    }
    return nil, derrors.NewOperationError(errors.SessionDoesNotExist).WithParams(sessionID)
}

// Delete a given session.
//   params:
//     sessionID The session identifier.
//   returns:
//     An error if the session cannot be removed.
func (fs *FileSystemProvider) Delete(sessionID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(sessionID) {
        return derrors.AsDaishoError(os.Remove(fs.getPath(sessionID)), errors.IOError)
    }
    return derrors.NewOperationError(errors.SessionDoesNotExist).WithParams(sessionID)
}

// Update a session in the system.
//   params:
//     node The new session information.
//   returns:
//     An error if the session cannot be updated.
func (fs *FileSystemProvider) Update(session entities.Session) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(session.ID) {
        toWrite, err := json.Marshal(session)
        if err == nil {
            return derrors.AsDaishoError(ioutil.WriteFile(fs.getPath(session.ID), toWrite, fileCreatePerm), errors.IOError)
        }
        return derrors.NewEntityError(session, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.SessionDoesNotExist).WithParams(session)
}

// Dump obtains the list of all sessions in the system.
//   returns:
//     The list of sessions.
//     An error if the list cannot be retrieved.
func (fs *FileSystemProvider) Dump() ([] entities.Session, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err)
    }
    result := make([] entities.Session, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieve(f.Name())
            if err == nil {
                result = append(result, * toAdd)
            } else {
                return nil, err
            }
        }
    }
    return result, nil
}

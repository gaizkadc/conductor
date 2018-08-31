//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
//

package credentialsstorage

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
    os.MkdirAll(basePath+"/credentials/", dirCreatePerm)
    return &FileSystemProvider{BasePath: basePath}
}

func (fs *FileSystemProvider) getPath(path string) string {
    newPath := fmt.Sprintf("%s/credentials/%s", fs.BasePath, path)
    return newPath
}

func (fs *FileSystemProvider) getSubPath(parent string, nodeID string, path string) string {
    newPath := fmt.Sprintf("%s/credentials/%s/%s/%s", fs.BasePath, parent, nodeID, path)
    return newPath
}

// Add a new user to the system.
//   params:
//     user The user to be added
//   returns:
//     An error if the user cannot be added.
func (fs *FileSystemProvider) Add(credentials entities.Credentials) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if ! fs.unsafeExists(credentials.UUID) {
        toWrite, err := json.Marshal(credentials)
        if err == nil {
            return derrors.AsDaishoError(ioutil.WriteFile(fs.getPath(credentials.UUID), toWrite, fileCreatePerm), errors.IOError)
        }
        return derrors.NewEntityError(credentials, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.CredentialsAlreadyExist).WithParams(credentials)
}

// Exists checks if a user exists in the system.
//   params:
//     userID The user identifier.
//   returns:
//     Whether the user exists or not.
func (fs *FileSystemProvider) Exists(uuID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExists(uuID)
}

func (fs *FileSystemProvider) unsafeExists(uuID string) bool {
    _, err := os.Stat(fs.getPath(uuID))
    return err == nil
}

// RetrieveUser retrieves a given user.
//   params:
//     uuID The user identifier.
//   returns:
//     The user.
//     An error if the user cannot be retrieved.
func (fs *FileSystemProvider) Retrieve(uuID string) (*entities.Credentials, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeRetrieve(uuID)
}

func (fs *FileSystemProvider) unsafeRetrieve(uuID string) (*entities.Credentials, derrors.DaishoError) {
    if fs.unsafeExists(uuID) {
        content, err := ioutil.ReadFile(fs.getPath(uuID))
        if err == nil {
            credentials := entities.Credentials{}
            err = json.Unmarshal(content, &credentials)
            if err == nil {
                return &credentials, nil
            }
            return nil, derrors.NewOperationError(errors.UnmarshalError).WithParams(uuID)
        }
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(uuID)
    }
    return nil, derrors.NewOperationError(errors.CredentialsDoNotExist).WithParams(uuID)
}

// Delete a given user.
//   params:
//     uuID The user identifier.
//   returns:
//     An error if the user cannot be removed.
func (fs *FileSystemProvider) Delete(uuID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(uuID) {
        return derrors.AsDaishoError(os.Remove(fs.getPath(uuID)), errors.IOError)
    }
    return derrors.NewOperationError(errors.CredentialsDoNotExist).WithParams(uuID)
}

// Update a user in the system.
//   params:
//     node The new user information.
//   returns:
//     An error if the user cannot be updated.
func (fs *FileSystemProvider) Update(credentials entities.Credentials) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(credentials.UUID) {
        toWrite, err := json.Marshal(credentials)
        if err == nil {
            return derrors.AsDaishoError(ioutil.WriteFile(fs.getPath(credentials.UUID), toWrite, fileCreatePerm), errors.IOError)
        }
        return derrors.NewEntityError(credentials, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.CredentialsDoNotExist).WithParams(credentials)
}


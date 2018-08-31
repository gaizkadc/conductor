//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
//

package oauthstorage

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
    os.MkdirAll(basePath+"/oauth/", dirCreatePerm)
    return &FileSystemProvider{BasePath: basePath}
}

func (fs *FileSystemProvider) getPath(path string) string {
    newPath := fmt.Sprintf("%s/oauth/%s", fs.BasePath, path)
    return newPath
}

func (fs *FileSystemProvider) getSubPath(parent string, userID string, path string) string {
    newPath := fmt.Sprintf("%s/oauth/%s/%s/%s", fs.BasePath, parent, userID, path)
    return newPath
}


// Add new secrets to the system.
//   params:
//     secrets The secrets to be added.
//   returns:
//     An error if the user cannot be added.
func (fs *FileSystemProvider) Add(secrets entities.OAuthSecrets) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if ! fs.unsafeExists(secrets.UserID) {
        toWrite, err := json.Marshal(secrets)
        if err == nil {
            return derrors.AsDaishoError(ioutil.WriteFile(fs.getPath(secrets.UserID), toWrite, fileCreatePerm), errors.IOError)
        }
        return derrors.NewEntityError(secrets, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.UserAlreadyExists).WithParams(secrets.UserID)
}


// Check if a user exists in the system.
//   params:
//     userID The user identifier.
//   returns:
//     Whether the user exists or not.
func (fs *FileSystemProvider) Exists(userID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExists(userID)
}

func (fs *FileSystemProvider) unsafeExists(userID string) bool {
    _, err := os.Stat(fs.getPath(userID))
    return err == nil
}

// Retrieve a given user secrets.
//   params:
//     userID The user identifier.
//   returns:
//     The secrets.
//     An error if the user cannot be retrieved.
func (fs *FileSystemProvider) Retrieve(userID string) (* entities.OAuthSecrets, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeRetrieve(userID)
}

func (fs *FileSystemProvider) unsafeRetrieve(userID string) (*entities.OAuthSecrets, derrors.DaishoError) {
    if fs.unsafeExists(userID) {
        content, err := ioutil.ReadFile(fs.getPath(userID))
        if err == nil {
            secrets := entities.OAuthSecrets{}
            err = json.Unmarshal(content, &secrets)
            if err == nil {
                return &secrets, nil
            }
            return nil, derrors.NewOperationError(errors.UnmarshalError).WithParams(userID)
        }
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(userID)
    }
    return nil, derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}


// Delete an existing secrets collection.
//   params:
//     userID The user identifier.
//   returns:
//     An error if the user cannot be removed.
func (fs *FileSystemProvider) Delete(userID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(userID) {
        return derrors.AsDaishoError(os.Remove(fs.getPath(userID)), errors.IOError)
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}



// Update a secrets collection in the system.
//   params:
//     node The new user information.
//   returns:
//     An error if the user cannot be updated.
func (fs *FileSystemProvider) Update(secrets entities.OAuthSecrets) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(secrets.UserID) {
        toWrite, err := json.Marshal(secrets)
        if err == nil {
            return derrors.AsDaishoError(ioutil.WriteFile(fs.getPath(secrets.UserID), toWrite, fileCreatePerm), errors.IOError)
        }
        return derrors.NewEntityError(secrets, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(secrets)
}




// Dump obtains the list of all secrets in the system.
//   returns:
//     The list of user.
//     An error if the list cannot be retrieved.
func (fs *FileSystemProvider) Dump() ([] entities.OAuthSecrets, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err)
    }
    result := make([] entities.OAuthSecrets, 0)
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

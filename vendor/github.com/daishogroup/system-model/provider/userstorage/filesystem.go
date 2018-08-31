//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
//

package userstorage

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
    os.MkdirAll(basePath+"/users/", dirCreatePerm)
    return &FileSystemProvider{BasePath: basePath}
}

func (fs *FileSystemProvider) getPath(path string) string {
    newPath := fmt.Sprintf("%s/users/%s", fs.BasePath, path)
    return newPath
}

func (fs *FileSystemProvider) getSubPath(parent string, nodeID string, path string) string {
    newPath := fmt.Sprintf("%s/users/%s/%s/%s", fs.BasePath, parent, nodeID, path)
    return newPath
}

// Add a new user to the system.
//   params:
//     user The user to be added
//   returns:
//     An error if the user cannot be added.
func (fs *FileSystemProvider) Add(user entities.User) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if ! fs.unsafeExists(user.ID) {
        toWrite, err := json.Marshal(user)
        if err == nil {
            return derrors.AsDaishoError(ioutil.WriteFile(fs.getPath(user.ID), toWrite, fileCreatePerm), errors.IOError)
        }
        return derrors.NewEntityError(user, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.UserAlreadyExists).WithParams(user)
}

// Exists checks if a user exists in the system.
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

// RetrieveUser retrieves a given user.
//   params:
//     userID The user identifier.
//   returns:
//     The user.
//     An error if the user cannot be retrieved.
func (fs *FileSystemProvider) RetrieveUser(userID string) (*entities.User, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeRetrieveUser(userID)
}

func (fs *FileSystemProvider) unsafeRetrieveUser(userID string) (*entities.User, derrors.DaishoError) {
    if fs.unsafeExists(userID) {
        content, err := ioutil.ReadFile(fs.getPath(userID))
        if err == nil {
            user := entities.User{}
            err = json.Unmarshal(content, &user)
            if err == nil {
                return &user, nil
            }
            return nil, derrors.NewOperationError(errors.UnmarshalError).WithParams(userID)
        }
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(userID)
    }
    return nil, derrors.NewOperationError(errors.UserDoesNotExist).WithParams(userID)
}

// Delete a given user.
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

// Update a user in the system.
//   params:
//     node The new user information.
//   returns:
//     An error if the user cannot be updated.
func (fs *FileSystemProvider) Update(user entities.User) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(user.ID) {
        toWrite, err := json.Marshal(user)
        if err == nil {
            return derrors.AsDaishoError(ioutil.WriteFile(fs.getPath(user.ID), toWrite, fileCreatePerm), errors.IOError)
        }
        return derrors.NewEntityError(user, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.UserDoesNotExist).WithParams(user)
}

// Dump obtains the list of all users in the system.
//   returns:
//     The list of users.
//     An error if the list cannot be retrieved.
func (fs *FileSystemProvider) Dump() ([] entities.User, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err)
    }
    result := make([] entities.User, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveUser(f.Name())
            if err == nil {
                result = append(result, * toAdd)
            } else {
                return nil, err
            }
        }
    }
    return result, nil
}

// ReducedInfoList get a list with the reduced info.
//   returns:
//     List of the reduced info.
//     An error if the info cannot be retrieved.
func (fs *FileSystemProvider) ReducedInfoList() ([] entities.UserReducedInfo, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err)
    }
    result := make([] entities.UserReducedInfo, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveUser(f.Name())
            if err == nil {
                reduced := entities.NewUserReducedInfo(toAdd.ID, toAdd.Email)
                result = append(result, *reduced)
            } else {
                return nil, err
            }
        }
    }
    return result, nil
}

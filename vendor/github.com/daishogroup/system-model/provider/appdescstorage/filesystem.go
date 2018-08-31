//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package appdescstorage

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

// FileSystemProvider structure that provides a file system backed storage for applications.
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
    os.MkdirAll(basePath+"/appdesc/", dirCreatePerm)
    return &FileSystemProvider{BasePath: basePath}
}

func (fs *FileSystemProvider) getPath(path string) string {
    newPath := fmt.Sprintf("%s/appdesc/%s", fs.BasePath, path)
    return newPath
}

func (fs *FileSystemProvider) getSubPath(parent string, appDescriptorID string, path string) string {
    newPath := fmt.Sprintf("%s/appdesc/%s/%s/%s", fs.BasePath, parent, appDescriptorID, path)
    return newPath
}

// Add a new application descriptor to the system.
//   params:
//     descriptor The application descriptor to be added
//   returns:
//     An error if the descriptor cannot be added.
func (fs *FileSystemProvider) Add(descriptor entities.AppDescriptor) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if ! fs.unsafeExists(descriptor.ID) {
        toWrite, err := json.Marshal(descriptor)
        if err == nil {
            err = ioutil.WriteFile(fs.getPath(descriptor.ID), toWrite, fileCreatePerm)
            return derrors.AsDaishoError(err, errors.IOError)
        }
        return derrors.NewEntityError(descriptor, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.AppDescAlreadyExists).WithParams(descriptor)
}

// Exists checks if an application descriptor exists in the system.
//   params:
//     descriptorID The application descriptor identifier.
//   returns:
//     Whether the descriptor exists or not.
func (fs *FileSystemProvider) Exists(descriptorID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExists(descriptorID)
}

func (fs *FileSystemProvider) unsafeExists(descriptorID string) bool {
    _, err := os.Stat(fs.getPath(descriptorID))
    return err == nil
}

// RetrieveDescriptor retrieves a given application descriptor.
//   params:
//     descriptorID The application descriptor identifier.
//   returns:
//     The application descriptor.
//     An error if the descriptor cannot be retrieved.
func (fs *FileSystemProvider) RetrieveDescriptor(descriptorID string) (*entities.AppDescriptor, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeRetrieveDescriptor(descriptorID)
}

func (fs * FileSystemProvider) unsafeRetrieveDescriptor(descriptorID string) (*entities.AppDescriptor, derrors.DaishoError) {
    _, err := os.Stat(fs.getPath(descriptorID))
    if err == nil {
        content, err := ioutil.ReadFile(fs.getPath(descriptorID))
        if err == nil {
            desc := entities.AppDescriptor{}
            err = json.Unmarshal(content, &desc)
            if err == nil {
                return &desc, nil
            }
            return nil, derrors.NewOperationError(errors.UnmarshalError, err).WithParams(descriptorID)
        }
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(descriptorID)
    }
    return nil, derrors.NewOperationError(errors.AppDescDoesNotExists).WithParams(descriptorID)
}

// Delete a given application descriptor.
//   params:
//     instanceID The application descriptor identifier.
//   returns:
//     An error if the application descriptor cannot be removed.
func (fs *FileSystemProvider) Delete(descriptorID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(descriptorID) {
        err := os.Remove(fs.getPath(descriptorID))
        return derrors.AsDaishoError(err, errors.IOError)
    }
    return derrors.NewOperationError(errors.AppDescDoesNotExists).WithParams(descriptorID)
}

// Dump obtains the list of all app descriptors in the system.
//   returns:
//     The list of AppDescriptors.
//     An error if the list cannot be retrieved.
func (fs *FileSystemProvider) Dump() ([] entities.AppDescriptor, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewGenericError(errors.IOError, err)
    }
    result := make([] entities.AppDescriptor, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveDescriptor(f.Name())
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
//     List of the reduced app info.
//     An error if the descriptor cannot be retrieved.
func (fs *FileSystemProvider) ReducedInfoList() ([] entities.AppDescriptorReducedInfo, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewGenericError(errors.IOError, err)
    }
    result := make([] entities.AppDescriptorReducedInfo, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveDescriptor(f.Name())
            if err == nil {
                reduced := entities.NewAppDescriptorReducedInfo(toAdd.NetworkID, toAdd.ID, toAdd.Name)
                result = append(result, *reduced)
            } else {
                return nil, err
            }
        }
    }
    return result, nil
}

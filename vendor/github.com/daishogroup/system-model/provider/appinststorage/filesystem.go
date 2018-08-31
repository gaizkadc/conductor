//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package appinststorage

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

// FileSystemProvider implements the AppInstProvider using a file system as backend.
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
    os.MkdirAll(basePath+"/appinst/", dirCreatePerm)
    return &FileSystemProvider{BasePath: basePath}
}

func (fs *FileSystemProvider) getPath(path string) string {
    newPath := fmt.Sprintf("%s/appinst/%s", fs.BasePath, path)
    return newPath
}

func (fs *FileSystemProvider) getSubPath(parent string, deployedID string, path string) string {
    newPath := fmt.Sprintf("%s/appinst/%s/%s/%s", fs.BasePath, parent, deployedID, path)
    return newPath
}

// Add a new application instance to the system.
//   params:
//     instance The application instance to be added
//   returns:
//     An error if the instance cannot be added.
func (fs *FileSystemProvider) Add(instance entities.AppInstance) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if ! fs.unsafeExists(instance.DeployedID) {
        toWrite, err := json.Marshal(instance)
        if err == nil {
            err = ioutil.WriteFile(fs.getPath(instance.DeployedID), toWrite, fileCreatePerm)
            return derrors.AsDaishoError(err, errors.IOError)
        }
        return derrors.NewEntityError(instance, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.AppInstAlreadyExists).WithParams(instance)
}

// Update an instance in the system.
//   params:
//     instance The new instance information. The instance identifier will be used and cannot be modified.
//   returns:
//     An error if the instance cannot be updated.
func (fs *FileSystemProvider) Update(instance entities.AppInstance) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(instance.DeployedID) {
        toWrite, err := json.Marshal(instance)
        if err == nil {
            err = ioutil.WriteFile(fs.getPath(instance.DeployedID), toWrite, fileCreatePerm)
            return derrors.AsDaishoError(err, errors.IOError)
        }
        return derrors.NewEntityError(instance, errors.MarshalError, err)
    }
    return derrors.NewEntityError(instance, errors.AppInstDoesNotExists)
}

// Exists checks if an application instance exists in the system.
//   params:
//     instanceID The application instance identifier.
//   returns:
//     Whether the instance exists or not.
func (fs *FileSystemProvider) Exists(instanceID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExists(instanceID)
}

func (fs *FileSystemProvider) unsafeExists(instanceID string) bool {
    _, err := os.Stat(fs.getPath(instanceID))
    return err == nil
}

// RetrieveInstance retrieves a given application instance.
//   params:
//     instanceID The application instance identifier.
//   returns:
//     The application instance.
//     An error if the instance cannot be retrieved.
func (fs *FileSystemProvider) RetrieveInstance(instanceID string) (*entities.AppInstance, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeRetrieveInstance(instanceID)
}

func (fs *FileSystemProvider) unsafeRetrieveInstance(instanceID string) (*entities.AppInstance, derrors.DaishoError) {
    if fs.unsafeExists(instanceID) {
        content, err := ioutil.ReadFile(fs.getPath(instanceID))
        if err == nil {
            instance := entities.AppInstance{}
            err = json.Unmarshal(content, &instance)
            if err == nil {
                return &instance, nil
            }
            return nil, derrors.NewOperationError(errors.UnmarshalError, err).WithParams(instanceID)
        }
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(instanceID)
    }
    return nil, derrors.NewOperationError(errors.AppInstDoesNotExists).WithParams(instanceID)
}

// Delete a given instance.
//   params:
//     instanceID The application instance identifier.
//   returns:
//     An error if the instance cannot be removed.
func (fs *FileSystemProvider) Delete(instanceID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(instanceID) {
        err := os.Remove(fs.getPath(instanceID))
        return derrors.AsDaishoError(err, errors.IOError)
    }
    return derrors.NewOperationError(errors.AppInstDoesNotExists).WithParams(instanceID)
}

// Dump obtains the list of all application instances in the system.
//   returns:
//     The list of AppInstance.
//     An error if the list cannot be retrieved.
func (fs *FileSystemProvider) Dump() ([] entities.AppInstance, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewGenericError(errors.IOError, err)
    }
    result := make([] entities.AppInstance, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveInstance(f.Name())
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
func (fs *FileSystemProvider) ReducedInfoList() ([] entities.AppInstanceReducedInfo, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewGenericError(errors.IOError, err)
    }
    result := make([] entities.AppInstanceReducedInfo, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveInstance(f.Name())
            if err == nil {
                reduced := entities.NewAppInstanceReducedInfo(toAdd.NetworkID, toAdd.ClusterID, toAdd.AppDescriptorID,
                    toAdd.DeployedID, toAdd.Name, toAdd.Description, toAdd.Ports, toAdd.Port, toAdd.ClusterAddress)
                result = append(result, *reduced)
            } else {
                return nil, err
            }
        }
    }
    return result, nil
}

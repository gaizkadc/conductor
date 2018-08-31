//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// File system implementation of the network provider.

package networkstorage

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

// FileSystemProvider that implements a file-backed provider.
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
    os.MkdirAll(basePath+"/networks/", dirCreatePerm)
    return &FileSystemProvider{BasePath: basePath}
}

func (fs *FileSystemProvider) getPath(path string) string {
    newPath := fmt.Sprintf("%s/networks/%s", fs.BasePath, path)
    return newPath
}

func (fs *FileSystemProvider) getSubPath(parent string, networkID string, path string) string {
    newPath := fmt.Sprintf("%s/networks/%s/%s/%s", fs.BasePath, parent, networkID, path)
    return newPath
}

// Add a new network to the system.
//   params:
//     network The Network to be added
//   returns:
//     An error if the network cannot be added.
func (fs *FileSystemProvider) Add(network entities.Network) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if ! fs.unsafeExists(network.ID) {
        toWrite, err := json.Marshal(network)
        if err == nil {
            ioError := ioutil.WriteFile(fs.getPath(network.ID), toWrite, fileCreatePerm)
            return derrors.AsDaishoError(ioError, errors.IOError)
        }
        return derrors.NewEntityError(network, errors.MarshalError, err).WithParams(network)
    }
    return derrors.NewOperationError(errors.NetworkAlreadyExists).WithParams(network)
}

// Exists checks if a network exists in the system.
//   params:
//     networkID The network identifier.
//   returns:
//     Whether the network exists or not.
func (fs *FileSystemProvider) Exists(networkID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExists(networkID)
}

func (fs *FileSystemProvider) unsafeExists(networkID string) bool {
    _, err := os.Stat(fs.getPath(networkID))
    return err == nil
}

// RetrieveNetwork retrieves a given network.
//   params:
//     networkID The network identifier.
//   returns:
//     The network.
//     An error if the network cannot be retrieved.
func (fs *FileSystemProvider) RetrieveNetwork(networkID string) (*entities.Network, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeRetrieveNetwork(networkID)
}

func (fs *FileSystemProvider) unsafeRetrieveNetwork(networkID string) (*entities.Network, derrors.DaishoError) {
    if fs.unsafeExists(networkID) {
        content, err := ioutil.ReadFile(fs.getPath(networkID))
        if err == nil {
            network := entities.Network{}
            err = json.Unmarshal(content, &network)
            if err == nil {
                return &network, nil
            }
            return nil, derrors.NewEntityError(networkID, errors.UnmarshalError, err)
        }
        return nil, derrors.NewOperationError(errors.IOError, err)
    }
    return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}

// ListNetworks retrieves all the networks in the system.
//   returns:
//     An array of networks.
//     An error if the networks cannot be retrieved.
func (fs *FileSystemProvider) ListNetworks() ([]entities.Network, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err)
    }
    result := make([] entities.Network, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveNetwork(f.Name())
            if err == nil {
                result = append(result, * toAdd)
            } else {
                return nil, err
            }
        }
    }
    return result, nil
}

// AttachCluster attaches a cluster to an existing network.
//   params:
//     networkID The network identifier.
//     clusterID The cluster identifier.
//   returns:
//     An error if the cluster cannot be attached.
func (fs *FileSystemProvider) AttachCluster(networkID string, clusterID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(networkID) {
        if !fs.unsafeExistsCluster(networkID, clusterID) {
            err := os.MkdirAll(fs.getSubPath("clusters", networkID, ""), dirCreatePerm)
            if err != nil {
                return derrors.NewOperationError(errors.IOError, err)
            }
            newFile, err := os.Create(fs.getSubPath("clusters", networkID, clusterID))
            if err == nil {
                return derrors.AsDaishoError(newFile.Close(), errors.IOError)
            }
            return derrors.NewOperationError(errors.IOError, err)
        }
        return derrors.NewOperationError(errors.ClusterAlreadyAttached).WithParams(networkID, clusterID)
    }
    return derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}

// ListClusters lists the clusters of a network.
//   params:
//     networkID The network identifier.
//   returns:
//     An array of cluster identifiers.
//     An error if the clusters cannot be retrieved.
func (fs *FileSystemProvider) ListClusters(networkID string) ([]string, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    if !fs.unsafeExists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    result := make([] string, 0)

    path := fs.getSubPath("clusters", networkID, "")
    _, err := os.Stat(path)
    if err != nil {
        // Network exists but may have no clusters
        return result, nil
    }
    // List existing files
    files, err := ioutil.ReadDir(path)
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(path)
    }
    for _, f := range files {
        if !f.IsDir() {
            result = append(result, f.Name())
        }
    }
    return result, nil
}

// ExistsCluster checks if a cluster is associated with a given network.
//   params:
//     networkID The network identifier.
//     clusterID The cluster identifier.
//   returns:
//     Whether the cluster is associated to the network.
func (fs *FileSystemProvider) ExistsCluster(networkID string, clusterID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExistsCluster(networkID, clusterID)
}

func (fs *FileSystemProvider) unsafeExistsCluster(networkID string, clusterID string) bool {
    _, err := os.Stat(fs.getSubPath("clusters", networkID, clusterID))
    return err == nil
}

// DeleteCluster deletes a cluster from an existing network.
//   params:
//     networkID The network identifier.
//     clusterID The cluster identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (fs *FileSystemProvider) DeleteCluster(networkID string, clusterID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExistsCluster(networkID, clusterID) {
        return derrors.AsDaishoError(os.Remove(fs.getSubPath("clusters", networkID, clusterID)), errors.IOError)
    }
    return derrors.NewOperationError(errors.ClusterNotAttachedToNetwork).WithParams(networkID, clusterID)
}

// RegisterAppDesc registers a new application descriptor inside a given network.
//   params:
//     networkID The network identifier.
//     appDescriptorID The application descriptor identifier.
//   returns:
//     An error if the descriptor cannot be registered.
func (fs *FileSystemProvider) RegisterAppDesc(networkID string, appDescriptorID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(networkID) {
        if !fs.unsafeExistsAppDesc(networkID, appDescriptorID) {
            err := os.MkdirAll(fs.getSubPath("appdesc", networkID, ""), dirCreatePerm)
            if err != nil {
                return derrors.NewOperationError(errors.IOError, err)
            }
            newFile, err := os.Create(fs.getSubPath("appdesc", networkID, appDescriptorID))
            if err == nil {
                return derrors.AsDaishoError(newFile.Close(), errors.IOError)
            }
            return derrors.NewOperationError(errors.IOError, err)
        }
        return derrors.NewOperationError(errors.AppDescAlreadyAttached).WithParams(networkID, appDescriptorID)
    }
    return derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}

// ListAppDesc lists all the application descriptors in a given network.
//   params:
//     networkID The network identifier.
//   returns:
//     An array of application descriptor identifiers.
//     An error if the list cannot be retrieved.
func (fs *FileSystemProvider) ListAppDesc(networkID string) ([] string, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    if !fs.unsafeExists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    result := make([] string, 0)

    path := fs.getSubPath("appdesc", networkID, "")
    _, err := os.Stat(path)
    if err != nil {
        // Network exists but may have no clusters
        return result, nil
    }
    // List existing files
    files, err := ioutil.ReadDir(path)
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(path)
    }
    for _, f := range files {
        if !f.IsDir() {
            result = append(result, f.Name())
        }
    }
    return result, nil
}

// ExistsAppDesc checks if an application descriptor exists in a network.
//   params:
//     networkID The network identifier.
//     appDescriptorID The application descriptor identifier.
//   returns:
//     Whether the application exists in the given network.
func (fs *FileSystemProvider) ExistsAppDesc(networkID string, appDescriptorID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExistsAppDesc(networkID, appDescriptorID)
}

func (fs *FileSystemProvider) unsafeExistsAppDesc(networkID string, appDescriptorID string) bool {
    _, err := os.Stat(fs.getSubPath("appdesc", networkID, appDescriptorID))
    return err == nil
}


// DeleteAppDescriptor deletes an application descriptor from a network.
//   params:
//     networkID The network identifier.
//     appDescriptorID The application descriptor identifier.
//   returns:
//     An error if the application descriptor cannot be removed.
func (fs *FileSystemProvider) DeleteAppDescriptor(networkID string, appDescriptorID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExistsAppDesc(networkID, appDescriptorID) {
        return derrors.AsDaishoError(os.Remove(fs.getSubPath("appdesc", networkID, appDescriptorID)), errors.IOError)
    }
    return derrors.NewOperationError(errors.AppDescNotAttached).WithParams(networkID, appDescriptorID)
}


// RegisterAppInst registers a new application instance inside a given network.
//   params:
//     networkID The network identifier.
//     appInstanceID The application instance identifier.
//   returns:
//     An error if the descriptor cannot be registered.
func (fs *FileSystemProvider) RegisterAppInst(networkID string, appInstanceID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(networkID) {
        if !fs.unsafeExistsAppDesc(networkID, appInstanceID) {
            err := os.MkdirAll(fs.getSubPath("appinst", networkID, ""), dirCreatePerm)
            if err != nil {
                return derrors.NewOperationError(errors.IOError).WithParams(networkID)
            }
            newFile, err := os.Create(fs.getSubPath("appinst", networkID, appInstanceID))
            if err == nil {
                return derrors.AsDaishoError(newFile.Close(), errors.IOError)
            }
            return derrors.AsDaishoError(err, errors.IOError)
        }
        return derrors.NewOperationError(errors.AppInstAlreadyAttached).WithParams(networkID, appInstanceID)
    }
    return derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}

// ListAppInst list all the application instances in a given network.
//   params:
//     networkID The network identifier.
//   returns:
//     An array of application descriptor identifiers.
//     An error if the list cannot be retrieved.
func (fs *FileSystemProvider) ListAppInst(networkID string) ([] string, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    if !fs.unsafeExists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    result := make([] string, 0)

    path := fs.getSubPath("appinst", networkID, "")
    _, err := os.Stat(path)
    if err != nil {
        // Network exists but may have no clusters
        return result, nil
    }
    // List existing files
    files, err := ioutil.ReadDir(path)
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(path)
    }
    for _, f := range files {
        if !f.IsDir() {
            result = append(result, f.Name())
        }
    }
    return result, nil
}

// ExistsAppInst checks if an application instance exists in a network.
//   params:
//     networkID The network identifier.
//     appInstanceID The application instance identifier.
//   returns:
//     Whether the application exists in the given network.
func (fs *FileSystemProvider) ExistsAppInst(networkID string, appInstanceID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExistsAppInst(networkID, appInstanceID)
}

func (fs *FileSystemProvider) unsafeExistsAppInst(networkID string, appInstanceID string) bool {
    _, err := os.Stat(fs.getSubPath("appinst", networkID, appInstanceID))
    return err == nil
}

// DeleteAppInstance deletes an application instance from a network.
//   params:
//     networkID The network identifier.
//     appInstanceID The application instance identifier.
//   returns:
//     An error if the application cannot be removed.
func (fs *FileSystemProvider) DeleteAppInstance(networkID string, appInstanceID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExistsAppInst(networkID, appInstanceID) {
        return derrors.AsDaishoError(os.Remove(fs.getSubPath("appinst", networkID, appInstanceID)), errors.IOError)
    }
    return derrors.NewOperationError(errors.AppInstNotAttachedToNetwork).WithParams(networkID, appInstanceID)
}

// ReducedInfoList get a list with the reduced info.
//   returns:
//     List of the reduced info.
//     An error if the info cannot be retrieved.
func (fs *FileSystemProvider) ReducedInfoList() ([] entities.NetworkReducedInfo, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err)
    }
    result := make([] entities.NetworkReducedInfo, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveNetwork(f.Name())
            if err == nil {
                reduced := entities.NewNetworkReducedInfo(toAdd.ID, toAdd.Name)
                result = append(result, *reduced)
            } else {
                return nil, err
            }
        }
    }
    return result, nil
}

// DeleteNetwork deletes a given network
// returns:
//  Error if any
func (fs * FileSystemProvider) DeleteNetwork(networkID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if !fs.unsafeExists(networkID) {
        return derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    return derrors.AsDaishoError(os.Remove(fs.getPath(networkID)), errors.IOError)
}
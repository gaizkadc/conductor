//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package clusterstorage

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

// FileSystemProvider that provides a file system backed storage.
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
    os.MkdirAll(basePath+"/clusters/", dirCreatePerm)
    return &FileSystemProvider{BasePath: basePath}
}

func (fs *FileSystemProvider) getPath(path string) string {
    newPath := fmt.Sprintf("%s/clusters/%s", fs.BasePath, path)
    return newPath
}

func (fs *FileSystemProvider) getSubPath(parent string, clusterID string, path string) string {
    newPath := fmt.Sprintf("%s/clusters/%s/%s/%s", fs.BasePath, parent, clusterID, path)
    return newPath
}

// Add a new cluster to the system.
//   params:
//     cluster The Cluster to be added
//   returns:
//     An error if the cluster cannot be added.
func (fs *FileSystemProvider) Add(cluster entities.Cluster) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if ! fs.unsafeExists(cluster.ID) {
        toWrite, err := json.Marshal(cluster)
        if err == nil {
            err = ioutil.WriteFile(fs.getPath(cluster.ID), toWrite, fileCreatePerm)
            return derrors.AsDaishoError(err, errors.IOError)
        }
        return derrors.NewEntityError(cluster, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.ClusterAlreadyExists).WithParams(cluster)
}

// Update a existing cluster in the provider.
//   params:
//     cluster The Cluster to be updated, the id of the cluster must be exist.
//   returns:
//     An error if the cluster cannot be edited.
func (fs *FileSystemProvider) Update(cluster entities.Cluster) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(cluster.ID) {
        toWrite, err := json.Marshal(cluster)
        if err == nil {
            err = ioutil.WriteFile(fs.getPath(cluster.ID), toWrite, fileCreatePerm)
            return derrors.AsDaishoError(err, errors.IOError)
        }
        return derrors.NewEntityError(cluster, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(cluster)
}

// Exists checks if a cluster exists in the system.
//   params:
//     clusterID The cluster identifier.
//   returns:
//     Whether the cluster exists or not.
func (fs *FileSystemProvider) Exists(clusterID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExists(clusterID)
}

func (fs *FileSystemProvider) unsafeExists(clusterID string) bool {
    _, err := os.Stat(fs.getPath(clusterID))
    return err == nil
}

// RetrieveCluster retrieves a given cluster.
//   params:
//     clusterID The cluster identifier.
//   returns:
//     The cluster.
//     An error if the cluster cannot be retrieved.
func (fs *FileSystemProvider) RetrieveCluster(clusterID string) (*entities.Cluster, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeRetrieveCluster(clusterID)
}

func (fs *FileSystemProvider) unsafeRetrieveCluster(clusterID string) (*entities.Cluster, derrors.DaishoError) {
    if fs.unsafeExists(clusterID) {
        content, err := ioutil.ReadFile(fs.getPath(clusterID))
        if err == nil {
            cluster := entities.Cluster{}
            err = json.Unmarshal(content, &cluster)
            if err == nil {
                return &cluster, nil
            }
            return nil, derrors.NewOperationError(errors.UnmarshalError).WithParams(clusterID)
        }
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(clusterID)
    }
    return nil, derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(clusterID)
}

// Delete a given cluster.
//   params:
//     clusterID The cluster identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (fs *FileSystemProvider) Delete(clusterID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(clusterID) {
        err := os.Remove(fs.getPath(clusterID))
        return derrors.AsDaishoError(err, errors.IOError)
    }
    return derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(clusterID)
}

// AttachNode links a node to an existing node.
//   params:
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     An error if the node cannot be attached.
func (fs *FileSystemProvider) AttachNode(clusterID string, nodeID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(clusterID) {
        if !fs.unsafeExistsNode(clusterID, nodeID) {
            err := os.MkdirAll(fs.getSubPath("nodes", clusterID, ""), dirCreatePerm)
            if err != nil {
                return derrors.NewOperationError(errors.IOError, err).WithParams(clusterID, nodeID)
            }
            newFile, err := os.Create(fs.getSubPath("nodes", clusterID, nodeID))
            if err == nil {
                err = newFile.Close()
                return derrors.AsDaishoError(err, errors.IOError)
            }
            return derrors.AsDaishoError(err, errors.IOError)
        }
        return derrors.NewOperationError(errors.NodeAlreadyAttachedToCluster).WithParams(clusterID, nodeID)
    }
    return derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(clusterID)
}

// ListNodes lists the nodes of a cluster.
//   params:
//     clusterID The cluster identifier.
//   returns:
//     An array of node identifiers.
//     An error if the nodes cannot be retrieved.
func (fs *FileSystemProvider) ListNodes(clusterID string) ([]string, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    if !fs.unsafeExists(clusterID) {
        return nil, derrors.NewOperationError(errors.ClusterDoesNotExists).WithParams(clusterID)
    }
    result := make([] string, 0)

    path := fs.getSubPath("nodes", clusterID, "")
    _, err := os.Stat(path)
    if err != nil {
        // Network exists but may have no clusters
        return result, nil
    }
    // List existing files
    files, err := ioutil.ReadDir(path)
    if err != nil {
        return nil, derrors.AsDaishoError(err, errors.IOError)
    }
    for _, f := range files {
        if !f.IsDir() {
            result = append(result, f.Name())
        }
    }
    return result, nil
}

// ExistsNode checks if a node is associated with a given cluster.
//   params:
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     Whether the node is associated to the cluster.
func (fs *FileSystemProvider) ExistsNode(clusterID string, nodeID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExistsNode(clusterID, nodeID)
}

func (fs *FileSystemProvider) unsafeExistsNode(clusterID string, nodeID string) bool {
    _, err := os.Stat(fs.getSubPath("nodes", clusterID, nodeID))
    return err == nil
}

// DeleteNode deletes a node from an existing cluster.
//   params:
//     clusterID    The cluster identifier.
//     nodeID       The node identifier.
//   returns:
//     An error if the node cannot be removed.
func (fs *FileSystemProvider) DeleteNode(clusterID string, nodeID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExistsNode(clusterID, nodeID) {
        err := os.Remove(fs.getSubPath("nodes", clusterID, nodeID))
        return derrors.AsDaishoError(err, errors.IOError)
    }
    return derrors.NewOperationError(errors.NodeNotAttachedToCluster).WithParams(clusterID, nodeID)
}

// Dump obtains the list of all clusters in the system.
//   returns:
//     The list of clusters.
//     An error if the list cannot be retrieved.
func (fs *FileSystemProvider) Dump() ([] entities.Cluster, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewGenericError(errors.IOError, err)
    }
    result := make([] entities.Cluster, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveCluster(f.Name())
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
func (fs *FileSystemProvider) ReducedInfoList() ([] entities.ClusterReducedInfo, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewGenericError(errors.IOError, err)
    }
    result := make([] entities.ClusterReducedInfo, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveCluster(f.Name())
            if err == nil {
                reduced := entities.NewClusterReducedInfo(toAdd.NetworkID, toAdd.ID, toAdd.Name, toAdd.Type)
                result = append(result, *reduced)
            } else {
                return nil, err
            }
        }
    }
    return result, nil
}

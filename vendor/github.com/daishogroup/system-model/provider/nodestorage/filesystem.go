//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package nodestorage

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
    os.MkdirAll(basePath+"/nodes/", dirCreatePerm)
    return &FileSystemProvider{BasePath: basePath}
}

func (fs *FileSystemProvider) getPath(path string) string {
    newPath := fmt.Sprintf("%s/nodes/%s", fs.BasePath, path)
    return newPath
}

func (fs *FileSystemProvider) getSubPath(parent string, nodeID string, path string) string {
    newPath := fmt.Sprintf("%s/nodes/%s/%s/%s", fs.BasePath, parent, nodeID, path)
    return newPath
}

// Add a new node to the system.
//   params:
//     node The Node to be added
//   returns:
//     An error if the node cannot be added.
func (fs *FileSystemProvider) Add(node entities.Node) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if ! fs.unsafeExists(node.ID) {
        toWrite, err := json.Marshal(node)
        if err == nil {
            return derrors.AsDaishoError(ioutil.WriteFile(fs.getPath(node.ID), toWrite, fileCreatePerm), errors.IOError)
        }
        return derrors.NewEntityError(node, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.NodeAlreadyExists).WithParams(node)
}

// Exists checks if a node exists in the system.
//   params:
//     nodeID The node identifier.
//   returns:
//     Whether the node exists or not.
func (fs *FileSystemProvider) Exists(nodeID string) bool {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeExists(nodeID)
}

func (fs *FileSystemProvider) unsafeExists(nodeID string) bool {
    _, err := os.Stat(fs.getPath(nodeID))
    return err == nil
}

// RetrieveNode retrieves a given node.
//   params:
//     nodeID The node identifier.
//   returns:
//     The cluster.
//     An error if the node cannot be retrieved.
func (fs *FileSystemProvider) RetrieveNode(nodeID string) (*entities.Node, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    return fs.unsafeRetrieveNode(nodeID)
}

func (fs *FileSystemProvider) unsafeRetrieveNode(nodeID string) (*entities.Node, derrors.DaishoError) {
    if fs.unsafeExists(nodeID) {
        content, err := ioutil.ReadFile(fs.getPath(nodeID))
        if err == nil {
            node := entities.Node{}
            err = json.Unmarshal(content, &node)
            if err == nil {
                return &node, nil
            }
            return nil, derrors.NewOperationError(errors.UnmarshalError).WithParams(nodeID)
        }
        return nil, derrors.NewOperationError(errors.IOError, err).WithParams(nodeID)
    }
    return nil, derrors.NewOperationError(errors.NodeDoesNotExists).WithParams(nodeID)
}

// Delete a given node.
//   params:
//     nodeID The node identifier.
//   returns:
//     An error if the node cannot be removed.
func (fs *FileSystemProvider) Delete(nodeID string) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(nodeID) {
        return derrors.AsDaishoError(os.Remove(fs.getPath(nodeID)), errors.IOError)
    }
    return derrors.NewOperationError(errors.NodeDoesNotExists).WithParams(nodeID)
}

// Update a node in the system.
//   params:
//     node The new node information.
//   returns:
//     An error if the node cannot be updated.
func (fs *FileSystemProvider) Update(node entities.Node) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    if fs.unsafeExists(node.ID) {
        toWrite, err := json.Marshal(node)
        if err == nil {
            return derrors.AsDaishoError(ioutil.WriteFile(fs.getPath(node.ID), toWrite, fileCreatePerm), errors.IOError)
        }
        return derrors.NewEntityError(node, errors.MarshalError, err)
    }
    return derrors.NewOperationError(errors.NodeDoesNotExists).WithParams(node)
}

// Dump obtains the list of all nodes in the system.
//   returns:
//     The list of nodes.
//     An error if the list cannot be retrieved.
func (fs *FileSystemProvider) Dump() ([] entities.Node, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err)
    }
    result := make([] entities.Node, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveNode(f.Name())
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
func (fs *FileSystemProvider) ReducedInfoList() ([] entities.NodeReducedInfo, derrors.DaishoError) {
    fs.Lock()
    defer fs.Unlock()
    files, err := ioutil.ReadDir(fs.getPath("/"))
    if err != nil {
        return nil, derrors.NewOperationError(errors.IOError, err)
    }
    result := make([] entities.NodeReducedInfo, 0)
    for _, f := range files {
        if !f.IsDir() {
            toAdd, err := fs.unsafeRetrieveNode(f.Name())
            if err == nil {
                reduced := entities.NewNodeReducedInfo(toAdd.NetworkID, toAdd.ClusterID,
                    toAdd.ID, toAdd.Name, toAdd.Status, toAdd.PublicIP)
                result = append(result, *reduced)
            } else {
                return nil, err
            }
        }
    }
    return result, nil
}

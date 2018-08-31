//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package configstorage

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
const fileName = "config.json"

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
    os.MkdirAll(basePath+"/config/", dirCreatePerm)
    return &FileSystemProvider{BasePath: basePath}
}

func (fs *FileSystemProvider) getPath(path string) string {
    newPath := fmt.Sprintf("%s/config/%s", fs.BasePath, path)
    return newPath
}

// Store the configuration.
//   params:
//     config The Config to be stored.
//   returns:
//     An error if the config cannot be added.
func (fs * FileSystemProvider) Store(config entities.Config) derrors.DaishoError {
    fs.Lock()
    defer fs.Unlock()
    toWrite, err := json.Marshal(config)
    if err == nil {
        err = ioutil.WriteFile(fs.getPath(fileName), toWrite, fileCreatePerm)
        return derrors.AsDaishoError(err, errors.IOError)
    }
    return derrors.NewEntityError(config, errors.MarshalError, err)
}

// Retrieve the current configuration.
//   returns:
//     The config.
//     An error if the config cannot be retrieved.
func (fs * FileSystemProvider)Get() (*entities.Config, derrors.DaishoError){
    fs.Lock()
    defer fs.Unlock()
    content, err := ioutil.ReadFile(fs.getPath(fileName))
    if err == nil {
        config := entities.Config{}
        err = json.Unmarshal(content, &config)
        if err == nil {
            return &config, nil
        }
        return nil, derrors.NewOperationError(errors.UnmarshalError)
    }
    return nil, derrors.NewOperationError(errors.IOError, err)
}
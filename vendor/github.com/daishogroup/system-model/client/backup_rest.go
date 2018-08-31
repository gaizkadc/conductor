package client

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// REST implementation of the backup client.

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "fmt"
    "github.com/daishogroup/dhttp"
)

const (
    BackupURI = "/api/v0/backup"
)

// type struct for paramerter passing
type Params struct {
    Componment string `url:"component"`
    Operation  string `url:"operation"`
}

// Client Rest for Network resources.
type BackupRest struct {
    client dhttp.Client
}

func NewBackupRestoreRest(basePath string) Backup {
    return NewBackupRestoreClientRest(ParseHostPort(basePath))
}

func NewBackupRestoreClientRest(host string, port int) Backup {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host, port))
    return &BackupRest{rest}
}

// Export all the information in the system model into a backup structure.
//   returns:
//     A backup structure with the system model information.
//     An error if the data cannot be obtained.
func (rest *BackupRest) Export(component string) (*entities.BackupRestore, derrors.DaishoError) {
    if component == "" {
        component = "all"
    }
    response := rest.client.Get(fmt.Sprintf("%s/%s/%s", BackupURI, component, "create"), new(entities.BackupRestore))
    if response.Error != nil {
        return nil, response.Error
    } else {
        n := response.Result.(*entities.BackupRestore)
        return n, nil
    }
}

// Import all the information in the system model into a backup structure.
//   returns:
//     A backup structure with the system model information.
//     An error if the data cannot be restored.
func (rest *BackupRest) Import(component string, entity *entities.BackupRestore) (derrors.DaishoError) {
    if component == "" {
        component = "all"
    }
    response := rest.client.Post(fmt.Sprintf("%s/%s/%s", BackupURI, component, "restore"), entity, nil)
    return response.Error
}

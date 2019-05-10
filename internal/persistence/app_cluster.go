/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package persistence

import (
    "bytes"
    "encoding/gob"
    "fmt"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/provider"
    "github.com/nalej/derrors"
)

// Manipulation of persistent entries storing information about applications running
// in the application clusters. This information is used by conductor to have a clear
// vision of running fragments.

const (
    AppClusterDB_Cluster_Bucket = "cluster"
)

type AppClusterDB struct {
    // provider to persist information
    db provider.KeyValueProvider
}

func NewAppClusterDB(db provider.KeyValueProvider) *AppClusterDB {
    return &AppClusterDB{
        db: db,
    }
}

func (a *AppClusterDB) AddDeploymentFragment(fragment entities.DeploymentFragment) derrors.Error{
    var buffer bytes.Buffer
    e := gob.NewEncoder(&buffer)
    if err := e.Encode(fragment); err!= nil {
        return derrors.NewInternalError("impossible to marshall deployment fragment", err)
    }

    return a.db.Put([]byte(AppClusterDB_Cluster_Bucket),[]byte(fragment.ClusterId),buffer.Bytes())
}

func (a *AppClusterDB) DeleteDeploymentFragment(fragment entities.DeploymentFragment) derrors.Error {
    var buffer bytes.Buffer
    e := gob.NewEncoder(&buffer)
    if err := e.Encode(fragment); err!= nil {
        return derrors.NewInternalError("impossible to marshall deployment fragment", err)
    }
    return a.db.Delete([]byte(AppClusterDB_Cluster_Bucket),buffer.Bytes())
}
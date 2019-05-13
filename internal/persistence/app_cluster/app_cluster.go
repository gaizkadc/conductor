/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package app_cluster

import (
    "bytes"
    "encoding/gob"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/provider"
    "github.com/nalej/derrors"
)

// Manipulation of persistent entries storing information about applications running
// in the application clusters. This information is used by conductor to have a clear
// vision of running fragments.


// The app cluster db creates one bucket for every cluster and stores
// the deployment fragments in the corresponding bucket.
// bucket    --> key                    --> value
// clusterId --> AppInstanceId_1 --> deploymentFragment
// clusterId --> AppInstanceId_2 --> deploymentFragment

type AppClusterDB struct {
    // provider to persist information
    db provider.KeyValueProvider
}

func NewAppClusterDB(db provider.KeyValueProvider) *AppClusterDB {
    return &AppClusterDB{
        db: db,
    }
}

func (a *AppClusterDB) AddDeploymentFragment(fragment *entities.DeploymentFragment) derrors.Error{
    var buffer bytes.Buffer
    e := gob.NewEncoder(&buffer)
    if err := e.Encode(fragment); err!= nil {
        return derrors.NewInternalError("impossible to marshall deployment fragment", err)
    }

    return a.db.Put([]byte(fragment.ClusterId),[]byte(fragment.AppInstanceId),buffer.Bytes())
}

func (a *AppClusterDB) GetDeploymentFragment(clusterId string, appInstanceId string) (*entities.DeploymentFragment, derrors.Error){
    retrieved, err := a.db.Get([]byte(clusterId),[]byte(appInstanceId))
    if err != nil {
        return nil, derrors.NewInternalError("impossible to get deployment fragment",err)
    }
    if retrieved == nil {
        return nil, nil
    }

    b:=bytes.NewReader(retrieved)
    d := gob.NewDecoder(b)
    var df entities.DeploymentFragment
    if err := d.Decode(&df); err!= nil {
        return nil,derrors.NewInternalError("impossible to unmarshall deployment fragment", err)
    }
    return &df, nil
}

func (a *AppClusterDB) DeleteDeploymentFragment(clusterId string, appInstanceId string) derrors.Error {
    err := a.db.Delete([]byte(clusterId), []byte(appInstanceId))
    if err != nil {
        return derrors.NewInternalError("impossible to delete deployment fragment", err)
    }
    return nil
}

func (a *AppClusterDB) GetAppsInCluster(clusterId string) ([]string, derrors.Error) {
    pairs, err := a.db.GetAllPairsInBucket([]byte(clusterId))
    if err != nil {
        return nil,derrors.NewInternalError("impossible to get deployments from cluster")
    }
    result := make([]string,len(pairs))
    for i, pair := range pairs {
        result[i] = string(pair.Key)
    }
    return result, nil
}
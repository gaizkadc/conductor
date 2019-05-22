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
    "github.com/rs/zerolog/log"
)

// Manipulation of persistent entries storing information about applications running
// in the application clusters. This information is used by conductor to have a clear
// vision of running fragments.


// The app cluster db creates one bucket for every cluster and stores
// the deployment fragments in the corresponding bucket.
// bucket    --> key                    --> value
// clusterId --> DeploymentFragmentId_1 --> deploymentFragment
// clusterId --> DeploymentFragmentId_2 --> deploymentFragment

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

    return a.db.Put([]byte(fragment.ClusterId),[]byte(fragment.FragmentId),buffer.Bytes())
}

func (a *AppClusterDB) GetDeploymentFragment(clusterId string, fragmentId string) (*entities.DeploymentFragment, derrors.Error){
    retrieved, err := a.db.Get([]byte(clusterId),[]byte(fragmentId))
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

func (a *AppClusterDB) DeleteDeploymentFragment(clusterId string, fragmentId string) derrors.Error {
    log.Debug().Str("clusterId",clusterId).Str("fragmentId",fragmentId).Msg("delete deployment fragment from db")
    err := a.db.Delete([]byte(clusterId), []byte(fragmentId))
    if err != nil {
        return derrors.NewInternalError("impossible to delete deployment fragment", err)
    }
    return nil
}


func (a *AppClusterDB) GetFragmentsInCluster(clusterId string) ([]entities.DeploymentFragment, derrors.Error) {
    pairs, err := a.db.GetAllPairsInBucket([]byte(clusterId))
    if err != nil {
        return nil,derrors.NewInternalError("impossible to get deployments from cluster")
    }
    result := make([]entities.DeploymentFragment,len(pairs))

    for i, pair := range pairs {
        b:=bytes.NewReader(pair.Value)
        d := gob.NewDecoder(b)
        var df entities.DeploymentFragment
        if err := d.Decode(&df); err!= nil {
            return nil, derrors.NewInternalError("impossible to unmarshall deployment fragment", err)
        }
        result[i] = df
    }
    return result, nil
}

// Return the deployment fragments for an application
func (a *AppClusterDB) GetFragmentsApp(clusterId string, appInstanceId string)([]entities.DeploymentFragment, derrors.Error) {
    fragmentsCluster, err := a.GetFragmentsInCluster(clusterId)
    if err != nil {
        return nil, err
    }
    toReturn := make([]entities.DeploymentFragment,0)
    for _, f := range fragmentsCluster {
        if f.AppInstanceId == appInstanceId {
            toReturn = append(toReturn, f)
        }
    }
    return toReturn, nil
}

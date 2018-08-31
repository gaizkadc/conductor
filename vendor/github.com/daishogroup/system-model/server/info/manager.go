//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Dump manager operations.

package info

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/provider/userstorage"
)

// Manager structure with all the required providers for the dump operations.
type Manager struct {
    networkProvider networkstorage.Provider
    clusterProvider clusterstorage.Provider
    nodesProvider   nodestorage.Provider
    appDescProvider appdescstorage.Provider
    appInstProvider appinststorage.Provider
    userProvider    userstorage.Provider
}

// NewManager creates a new dump manager.
func NewManager(
    networkProvider networkstorage.Provider,
    clusterProvider clusterstorage.Provider,
    nodeProvider nodestorage.Provider,
    appDescProvider appdescstorage.Provider,
    appInstProvider appinststorage.Provider,
    userProvider userstorage.Provider) Manager {
    return Manager{
        networkProvider, clusterProvider, nodeProvider,
        appDescProvider, appInstProvider, userProvider}
}

// ReducedInfo exports the basic information in the system model.
//   returns:
//     A dump structure with the system model information.
//     An error if the data cannot be obtained.
func (mgr *Manager) ReducedInfo() (*entities.ReducedInfo, derrors.DaishoError) {
    networks, err := mgr.networkProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }

    clusters, err := mgr.clusterProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }

    nodes, err := mgr.nodesProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }

    descriptors, err := mgr.appDescProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }

    instances, err := mgr.appInstProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }

    users, err := mgr.userProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }

    return entities.NewReducedInfo(networks, clusters, nodes, descriptors, instances, users), nil
}

// ReducedInfoByNetwork basic information in the system model filter by networkID.
//   params:
//     networkID The selected network.
//   returns:
//     A dump structure with the system model information.
//     An error if the data cannot be obtained.
func (mgr *Manager) ReducedInfoByNetwork(networkID string) (*entities.ReducedInfo, derrors.DaishoError) {
    networks, err := mgr.networkProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }
    filterNetwork := make([] entities.NetworkReducedInfo, 0)
    for _, n := range networks {
        if n.ID == networkID {
            filterNetwork = append(filterNetwork, n)
        }
    }

    clusters, err := mgr.clusterProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }
    filterClusters := make([] entities.ClusterReducedInfo, 0)
    for _, c := range clusters {
        if c.NetworkID == networkID {
            filterClusters = append(filterClusters, c)
        }
    }

    nodes, err := mgr.nodesProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }
    filterNodes := make([] entities.NodeReducedInfo, 0)
    for _, n := range nodes {
        if n.NetworkID == networkID {
            filterNodes = append(filterNodes, n)
        }
    }

    descriptors, err := mgr.appDescProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }

    filterDescriptors := make([] entities.AppDescriptorReducedInfo, 0)
    for _, d := range descriptors {
        if d.NetworkID == networkID {
            filterDescriptors = append(filterDescriptors, d)
        }
    }

    instances, err := mgr.appInstProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }
    filterInstances := make([] entities.AppInstanceReducedInfo, 0)
    for _, i := range instances {
        if i.NetworkID == networkID {
            filterInstances = append(filterInstances, i)
        }
    }

    users, err := mgr.userProvider.ReducedInfoList()
    if err != nil {
        return nil, err
    }

    return entities.NewReducedInfo(filterNetwork, filterClusters,
        filterNodes, filterDescriptors, filterInstances, users), nil
}

// SummaryInfo exports the basic information of the system model.
//   returns:
//     A summary with the counters of each entity.
//     An error if the data cannot be obtained.
func (mgr *Manager) SummaryInfo() (*entities.SummaryInfo, derrors.DaishoError) {
    info, err := mgr.ReducedInfo()
    if err != nil {
        return nil, err
    }
    return entities.NewSummaryInfo(len(info.Networks), len(info.Clusters), len(info.Nodes),
        len(info.Descriptors), len(info.Instances), len(info.Users)), nil
}

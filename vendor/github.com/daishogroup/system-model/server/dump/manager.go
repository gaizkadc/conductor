//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Dump manager operations.

package dump

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/daishogroup/system-model/provider/accessstorage"
)

// Manager structure with all the required providers for the dump operations.
type Manager struct {
    networkProvider networkstorage.Provider
    clusterProvider clusterstorage.Provider
    nodesProvider nodestorage.Provider
    appDescProvider appdescstorage.Provider
    appInstProvider appinststorage.Provider
    userProvider    userstorage.Provider
    accessProvider  accessstorage.Provider
}

// NewManager creates a new dump manager.
func NewManager(
    networkProvider networkstorage.Provider,
    clusterProvider clusterstorage.Provider,
    nodeProvider nodestorage.Provider,
    appDescProvider appdescstorage.Provider,
    appInstProvider appinststorage.Provider,
    userProvider userstorage.Provider,
    accessProvider accessstorage.Provider) Manager {
    return Manager{
        networkProvider, clusterProvider, nodeProvider,
        appDescProvider, appInstProvider, userProvider,
        accessProvider}
}

// Export all the information in the system model into a Dump structure.
//   returns:
//     A dump structure with the system model information.
//     An error if the data cannot be obtained.
func (mgr * Manager) Export() (* entities.Dump, derrors.DaishoError) {
    networks, err := mgr.networkProvider.ListNetworks()
    if err != nil {
        return nil, err
    }
    dump := entities.NewDumpWithNetworks(networks)

    clusters, err := mgr.clusterProvider.Dump()
    if err != nil {
        return nil, err
    }
    dump.AddClusters(clusters)

    nodes, err := mgr.nodesProvider.Dump()
    if err != nil {
        return nil, err
    }
    dump.AddNodes(nodes)

    descriptors, err := mgr.appDescProvider.Dump()
    if err != nil {
        return nil, err
    }
    dump.AddAppDescriptors(descriptors)

    instances, err := mgr.appInstProvider.Dump()
    if err != nil {
        return nil, err
    }
    dump.AddAppInstances(instances)

    users, err := mgr.userProvider.Dump()
    if err != nil {
        return nil, err
    }
    dump.AddUsers(users)

    access, err2 := mgr.accessProvider.Dump()
    if err2 != nil {
        return nil, err
    }
    dump.AddAccess(access)

    return dump, nil
}

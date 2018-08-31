package backup

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// BackupRestore manager operations.



import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/daishogroup/system-model/provider/accessstorage"
    "github.com/daishogroup/system-model/provider/passwordstorage"
    log "github.com/sirupsen/logrus"
)

// Manager structure with all the required providers for the dump operations.
type Manager struct {
    networkProvider networkstorage.Provider
    clusterProvider clusterstorage.Provider
    nodesProvider nodestorage.Provider
    appDescProvider appdescstorage.Provider
    userProvider    userstorage.Provider
    accessProvider  accessstorage.Provider
    passwordProvider passwordstorage.Provider
}

// NewManager creates a new dump manager.
func NewManager(
    networkProvider networkstorage.Provider,
    clusterProvider clusterstorage.Provider,
    nodeProvider nodestorage.Provider,
    appDescProvider appdescstorage.Provider,
    userProvider userstorage.Provider,
    accessProvider accessstorage.Provider, passwordProvider passwordstorage.Provider) Manager {
    return Manager{
        networkProvider, clusterProvider, nodeProvider,
        appDescProvider,  userProvider,
        accessProvider,  passwordProvider}
}

// Export all the information in the system model into a Dump structure.
//   returns:
//     A dump structure with the system model information.
//     An error if the data cannot be obtained.
func (mgr * Manager) Export(component string) (* entities.BackupRestore, derrors.DaishoError) {
    log.WithField("component", component).Debug("Exporting component")

    // Create list of components to export
    var components []string
    if component == "all" {
        components = []string{"networks", "clusters", "nodes", "appdesc", "users"}
    } else {
        components = []string{component}
    }

    backup := entities.NewBackup()

    for _, c := range components {
        switch c {
        case "networks":
            networks, err := mgr.networkProvider.ListNetworks()
            if err != nil {
                return nil, err
            }
            backup.AddNetworks(networks)
        case "clusters":
            clusters, err := mgr.clusterProvider.Dump()
            if err != nil {
                return nil, err
            }
            backup.AddClusters(clusters)
        case "nodes":
            nodes, err := mgr.nodesProvider.Dump()
            if err != nil {
                return nil, err
            }
            backup.AddNodes(nodes)
        case "appdesc":
            descriptors, err := mgr.appDescProvider.Dump()
            if err != nil {
                return nil, err
            }
            backup.AddAppDescriptors(descriptors)
        case "users":
            users, err := mgr.userProvider.Dump()
            if err != nil {
                return nil, err
            }
            var access *entities.UserAccess
            var password *entities.Password

            for _, user := range users {
                if mgr.accessProvider.Exists(user.ID) {
                    access, err = mgr.accessProvider.RetrieveAccess(user.ID)
                    if err != nil {
                        break
                    }

                    // Filter out internal users
                    if access.Roles[0].IsInternalUser() {
                        continue
                    }
                } else {
                    access = &entities.UserAccess{"", []entities.RoleType{}}
                }
                if mgr.passwordProvider.Exists(user.ID) {
                    password, err = mgr.passwordProvider.RetrievePassword(user.ID)
                    if err != nil {
                        // user has no access , not sure if this is an error
                        break
                    }
                } else {
                    password = &entities.Password{"", &[]byte{}}
                }

                backup.AddUsers(entities.BackupUser{User:user, Access:*access, Password: *password})
            }
            if err != nil {
                break
            }

        default:
            return nil, derrors.NewGenericError("Unknown component requested").WithParams(c)
        }
    }

    return backup, nil
}

func (mgr * Manager) Import(component string, restore * entities.BackupRestore ) ( derrors.DaishoError) {
    log.WithField("component", component).Debug("Importing component")

    // Create list of components to export
    var components []string
    if component == "all" {
        components = []string{"networks", "clusters", "nodes", "appdesc", "users"}
    } else {
        components = []string{component}
    }

    var err derrors.DaishoError = nil

    for _, c := range components {
        switch c {
        case "networks":
            for _, network := range restore.Networks {
                if mgr.networkProvider.Exists(network.ID) {
                    if err = mgr.networkProvider.DeleteNetwork(network.ID); err != nil {
                        break
                    }
                }
                err = mgr.networkProvider.Add(network)

                if err != nil {
                    break
                }
            }
        case "clusters":
            for _, cluster := range restore.Clusters {
                if mgr.clusterProvider.Exists(cluster.ID) {
                    if err = mgr.clusterProvider.Delete(cluster.ID); err != nil {
                        break
                    }
                }
                err = mgr.clusterProvider.Add(cluster)
                if err != nil {
                    break
                }
            }
        case "nodes":
            for _, node := range restore.Nodes {
                if mgr.nodesProvider.Exists(node.ID) {
                    if err = mgr.nodesProvider.Delete(node.ID); err != nil {
                        break
                    }
                }
                err = mgr.nodesProvider.Add(node)
                if err != nil {
                    break
                }
            }
        case "appdesc":
            for _, descriptor := range restore.AppDescriptors {

                if mgr.appDescProvider.Exists(descriptor.ID) {
                    if err = mgr.appDescProvider.Delete(descriptor.ID); err != nil {
                        break
                    }
                }

                // Fix up network ID
                networks, err := mgr.networkProvider.ListNetworks()
                if err != nil {
                    break
                }
                if len(networks) == 0 {
                    err = derrors.NewOperationError(errors.NoNetworkError)
                    break
                }

                descriptor.NetworkID = networks[0].ID

                err = mgr.appDescProvider.Add(descriptor)
                if err != nil {
                    break
                }

                // Register to network if needed. Note that we didn't unregister
                // when deleting the appdescriptor above, so it might already be
                // registered and we're fine
                if !mgr.networkProvider.ExistsAppDesc(networks[0].ID, descriptor.ID) {
                    err = mgr.networkProvider.RegisterAppDesc(networks[0].ID, descriptor.ID)
                    if err != nil {
                        break
                    }
                }
            }

        case "users":
            for _, user := range restore.Users {

                if mgr.userProvider.Exists(user.User.ID) {
                    if err = mgr.userProvider.Delete(user.User.ID); err != nil {
                        break
                    }
                }
                err = mgr.userProvider.Add(user.User)
                if err != nil {
                    break
                }

                if mgr.accessProvider.Exists(user.User.ID) {
                    if err = mgr.accessProvider.Delete(user.User.ID); err != nil {
                        return err
                    }
                }
                err = mgr.accessProvider.Add(user.Access)

                if mgr.passwordProvider.Exists(user.User.ID) {
                    if err = mgr.passwordProvider.Delete(user.User.ID); err != nil {
                        return err
                    }
                }
                err = mgr.passwordProvider.Add(user.Password)
                if err != nil {
                    break
                }
            }

        default:
            err = derrors.NewGenericError("Unknown component requested").WithParams(c)
        }

        if err != nil {
            log.WithField("component", c).WithField("error", err).Error("Returning error")
            return err
        }
    }

    return  nil
}

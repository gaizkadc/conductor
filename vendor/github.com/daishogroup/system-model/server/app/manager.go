//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package app

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
)

// Manager structure with all the providers required for application operations.
type Manager struct {
    networkProvider       networkstorage.Provider
    appDescriptorProvider appdescstorage.Provider
    appInstanceProvider   appinststorage.Provider
}

// NewManager creates a new application manager.
func NewManager(
    networkProvider networkstorage.Provider,
    appDescriptorProvider appdescstorage.Provider,
    appInstanceProvider appinststorage.Provider) Manager {
    return Manager{networkProvider, appDescriptorProvider, appInstanceProvider}
}

// AddApplicationDescriptor adds a new application descriptor.
//   params:
//     networkId The network identifier.
//     descriptor The application descriptor.
//   returns:
//     The added application descriptor.
//     An error if the application cannot be added.
func (mgr *Manager) AddApplicationDescriptor(networkID string,
    request entities.AddAppDescriptorRequest) (*entities.AppDescriptor, derrors.DaishoError) {
    descriptor := entities.ToAppDescriptor(networkID, request)
    if !mgr.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists)
    }

    if mgr.appDescriptorProvider.Exists(descriptor.ID) {
        return nil, derrors.NewOperationError(errors.AppDescAlreadyExists)
    }

    err := mgr.appDescriptorProvider.Add(*descriptor)
    if err == nil {
        err = mgr.networkProvider.RegisterAppDesc(networkID, descriptor.ID)
        if err == nil {
            return descriptor, nil
        }
    }

    return nil, err

}

// ListDescriptors list all the application descriptors available for a given network.
//   params:
//     networkId The network identifier.
//   returns:
//     An array of application descriptors.
//     An error in case the list cannot be retrieved.
func (mgr *Manager) ListDescriptors(networkID string) ([] entities.AppDescriptor, derrors.DaishoError) {
    if !mgr.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    appIds, err := mgr.networkProvider.ListAppDesc(networkID)
    if err != nil {
        return nil, err
    }
    result := make([] entities.AppDescriptor, 0, len(appIds))
    failed := false
    for index := 0; index < len(appIds) && !failed; index++ {
        toAdd, err := mgr.appDescriptorProvider.RetrieveDescriptor(appIds[index])
        if err == nil {
            result = append(result, *toAdd)
        } else {
            logger.Warn("Failed to retrieved associated descriptor " + appIds[index] + " on network " + networkID)
            failed = false
        }
    }

    if !failed {
        return result, nil
    }
    return nil, derrors.NewOperationError(errors.OpFail)
}

// GetDescriptor retrieves an application descriptor
//   params:
//     networkId The network identifier.
//     appDescriptorId The application descriptor identifier.
//   returns:
//     An application descriptor
//     An error if the descriptor cannot be retrieved or is not associated with the network.
func (mgr *Manager) GetDescriptor(networkID string, appDescriptorID string) (*entities.AppDescriptor, derrors.DaishoError) {
    if !mgr.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if mgr.networkProvider.ExistsAppDesc(networkID, appDescriptorID) {
        return mgr.appDescriptorProvider.RetrieveDescriptor(appDescriptorID)
    }
    return nil, derrors.NewOperationError(errors.AppDescNotAttached).WithParams(networkID, appDescriptorID)
}

// DeleteDescriptor deletes an application descriptor.
//   params:
//     networkId The network identifier.
//     appDescriptorID The application descriptor identifier.
//   returns:
//     An error if the descriptor cannot be removed.
func (mgr *Manager) DeleteDescriptor(networkID string, appDescriptorID string) derrors.DaishoError {
    if !mgr.networkProvider.ExistsAppDesc(networkID, appDescriptorID) {
        return derrors.NewOperationError(errors.AppDescNotAttached).WithParams(networkID, appDescriptorID)
    }

    err := mgr.appDescriptorProvider.Delete(appDescriptorID)
    if err == nil {
        err = mgr.networkProvider.DeleteAppDescriptor(networkID, appDescriptorID)
    }
    return err
}

// AddApplicationInstance adds a new application instance.
//   params:
//     networkId The network identifier.
//     instance The application instance.
//   returns:
//     The added instance.
//     An error if the application instance cannot be added.
func (mgr *Manager) AddApplicationInstance(networkID string,
    request entities.AddAppInstanceRequest) (*entities.AppInstance, derrors.DaishoError) {
    instance := entities.ToAppInstance(networkID, request)
    if !mgr.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }

    if mgr.appInstanceProvider.Exists(instance.DeployedID) {
        return nil, derrors.NewOperationError(errors.AppInstAlreadyExists).WithParams(networkID, request)
    }

    err := mgr.appInstanceProvider.Add(*instance)
    if err == nil {
        err := mgr.networkProvider.RegisterAppInst(networkID, instance.DeployedID)
        if err == nil {
            return instance, nil
        }
    }

    return nil, err

}

// ListInstances lists all the application instances inside a network.
//   params:
//     networkId The network identifier.
//   returns:
//     An array of application instances.
//     An error if the list cannot be retrieved.
func (mgr *Manager) ListInstances(networkID string) ([] entities.AppInstance, derrors.DaishoError) {
    if !mgr.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    appIds, err := mgr.networkProvider.ListAppInst(networkID)
    if err != nil {
        return make([] entities.AppInstance, 0, 0), err
    }
    result := make([] entities.AppInstance, 0, len(appIds))
    failed := false
    for index := 0; index < len(appIds) && !failed; index++ {
        toAdd, err := mgr.appInstanceProvider.RetrieveInstance(appIds[index])
        if err == nil {
            result = append(result, *toAdd)
        } else {
            logger.Warn("Failed to retrieved associated instance " + appIds[index] + " on network " + networkID)
            failed = false
        }
    }

    if !failed {
        return result, nil
    }
    return make([] entities.AppInstance, 0, 0), derrors.NewOperationError(errors.OpFail)
}

// GetInstance retrieves an application instance.
//   params:
//     networkId The network identifier.
//     appInstanceId The application instance identifier.
//   returns:
//     An application instance.
//     An error if the instance cannot be retrieved.
func (mgr *Manager) GetInstance(networkID string, appInstanceID string) (*entities.AppInstance, derrors.DaishoError) {
    if !mgr.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if !mgr.appInstanceProvider.Exists(appInstanceID) {
        return nil, derrors.NewOperationError(errors.AppInstDoesNotExists).WithParams(appInstanceID)
    }
    if !mgr.networkProvider.ExistsAppInst(networkID, appInstanceID) {
        return nil, derrors.NewOperationError(errors.AppInstNotAttachedToNetwork).WithParams(networkID, appInstanceID)
    }

    retrieved, err := mgr.appInstanceProvider.RetrieveInstance(appInstanceID)
    if err == nil {
        return retrieved, nil
    }
    return nil, err
}

// UpdateInstance updates an application instance.
//   params:
//     networkId The network identifier.
//     appInstanceId The application instance identifier.
//     update The new instance information.
//   returns:
//     The updated application instance.
//     An error if the instance cannot be retrieved.
func (mgr *Manager) UpdateInstance(networkID string, appInstanceID string, update entities.UpdateAppInstanceRequest) (*entities.AppInstance, derrors.DaishoError) {
    if !mgr.networkProvider.Exists(networkID) {
        return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }
    if !mgr.appInstanceProvider.Exists(appInstanceID) {
        return nil, derrors.NewOperationError(errors.AppInstDoesNotExists).WithParams(appInstanceID)
    }
    if !mgr.networkProvider.ExistsAppInst(networkID, appInstanceID) {
        return nil, derrors.NewOperationError(errors.AppInstNotAttachedToNetwork).WithParams(networkID, appInstanceID)
    }

    previous, err := mgr.appInstanceProvider.RetrieveInstance(appInstanceID)
    if err != nil {
        return nil, err
    }
    updated := previous.Merge(update)

    err = mgr.appInstanceProvider.Update(* updated)
    if err == nil {
        return updated, nil
    }
    return nil, err
}

// DeleteInstance deletes an application instance.
//   params:
//     networkId The network identifier.
//     appInstanceId The application instance identifier.
//   returns:
//     An error if the instance cannot be removed.
func (mgr *Manager) DeleteInstance(networkID string, appInstanceID string) derrors.DaishoError {
    if !mgr.networkProvider.ExistsAppInst(networkID, appInstanceID) {
        return derrors.NewOperationError(errors.AppInstNotAttachedToNetwork).WithParams(networkID, appInstanceID)
    }

    err := mgr.appInstanceProvider.Delete(appInstanceID)
    if err == nil {
        err = mgr.networkProvider.DeleteAppInstance(networkID, appInstanceID)
    }
    return err
}

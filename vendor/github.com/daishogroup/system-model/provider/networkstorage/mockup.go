//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the network provider mockup using in-memory storage.

package networkstorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupNetworkProvider structure that contains hashmaps to store network information in memory.
type MockupNetworkProvider struct {
    sync.Mutex

    // Networks indexed by network identifier.
    networks map[string]entities.Network

    // Array of cluster by network identifier.
    clusters map[string][] string

    // Array of applications descriptors by network identifier.
    applicationDescriptors map[string][]string

    // Array of applications instances by network identifier.
    applicationInstances map[string][]string
}

// NewMockupNetworkProvider creates a mockup provider for the network operations.
func NewMockupNetworkProvider() *MockupNetworkProvider {
    return &MockupNetworkProvider{
        networks: make(map[string]entities.Network, 0),
        clusters: make(map[string][] string, 0),
        applicationDescriptors: make(map[string][] string, 0),
        applicationInstances: make(map[string][] string, 0)}
}

// Add a new network to the system.
//   params:
//     network The Network to be added
//   returns:
//     An error if the network cannot be added.
func (mockup *MockupNetworkProvider) Add(network entities.Network) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(network.ID) {
        mockup.networks[network.ID] = network
        mockup.clusters[network.ID] = make([] string, 0)
        return nil
    }
    return derrors.NewOperationError(errors.NetworkAlreadyExists).WithParams(network)
}

// Exists checks if a network exists in the system.
//   params:
//     networkID The network identifier.
//   returns:
//     Whether the network exists or not.
func (mockup *MockupNetworkProvider) Exists(networkID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExists(networkID)
}

func (mockup *MockupNetworkProvider) unsafeExists(networkID string) bool {
    _, exists := mockup.networks[networkID]
    return exists
}


// RetrieveNetwork retrieves a given network.
//   params:
//     networkID The network identifier.
//   returns:
//     The network.
//     An error if the network cannot be retrieved.
func (mockup *MockupNetworkProvider) RetrieveNetwork(networkID string) (*entities.Network, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    network, exists := mockup.networks[networkID]
    if exists {
        return &network, nil
    }
    return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}

// ListNetworks retrieves all the networks in the system.
//   returns:
//     An array of networks.
//     An error if the networks cannot be retrieved.
func (mockup *MockupNetworkProvider) ListNetworks() ([]entities.Network, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    result := make([]entities.Network, 0, len(mockup.networks))
    for _, value := range mockup.networks {
        result = append(result, value)
    }
    return result, nil
}


// DeleteNetwork deletes a given network.
//   params:
//     networkID The network identifier.
//  returns:
//   Error if any
func (mockup *MockupNetworkProvider) DeleteNetwork(networkID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    _, exists := mockup.networks[networkID]
    if !exists {
        return derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
    }

    delete(mockup.networks, networkID)
    delete(mockup.clusters, networkID)
    delete(mockup.applicationInstances, networkID)
    return nil
}

// AttachCluster attaches a cluster to an existing network.
//   params:
//     networkID The network identifier.
//     clusterID The cluster identifier.
//   returns:
//     An error if the cluster cannot be attached.
func (mockup *MockupNetworkProvider) AttachCluster(networkID string, clusterID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(networkID) {
        if !mockup.unsafeExistsCluster(networkID, clusterID) {
            clusters, _ := mockup.clusters[networkID]
            mockup.clusters[networkID] = append(clusters, clusterID)
            return nil
        }
        return derrors.NewOperationError(errors.ClusterAlreadyAttached).WithParams(networkID, clusterID)
    }
    return derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}

// ListClusters lists the clusters of a network.
//   params:
//     networkID The network identifier.
//   returns:
//     An array of cluster identifiers.
//     An error if the clusters cannot be retrieved.
func (mockup *MockupNetworkProvider) ListClusters(networkID string) ([]string, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    clusters, ok := mockup.clusters[networkID]
    if ok {
        return clusters, nil
    }
    return nil, derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}

// ExistsCluster checks if a cluster is associated with a given network.
//   params:
//     networkID The network identifier.
//     clusterID The cluster identifier.
//   returns:
//     Whether the cluster is associated to the network.
func (mockup *MockupNetworkProvider) ExistsCluster(networkID string, clusterID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExistsCluster(networkID, clusterID)
}

func (mockup *MockupNetworkProvider) unsafeExistsCluster(networkID string, clusterID string) bool {
    clusters, ok := mockup.clusters[networkID]
    if ok {
        for _, cluster := range clusters {
            if cluster == clusterID {
                return true
            }
        }
        return false
    }
    return false
}

// DeleteCluster deletes a cluster from an existing network.
//   params:
//     networkID The network identifier.
//     clusterID The cluster identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (mockup *MockupNetworkProvider) DeleteCluster(networkID string, clusterID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExistsCluster(networkID, clusterID) {
        previousClusters := mockup.clusters[networkID]
        newClusters := make([] string, 0)
        for _, cID := range previousClusters {
            if cID != clusterID {
                newClusters = append(newClusters, cID)
            }
        }
        mockup.clusters[networkID] = newClusters
        return nil
    }
    return derrors.NewOperationError(errors.ClusterNotAttachedToNetwork).WithParams(networkID, clusterID)
}

// RegisterAppDesc registers a new application descriptor inside a given network.
//   params:
//     networkID The network identifier.
//     appDescriptorID The application descriptor identifier.
//   returns:
//     An error if the descriptor cannot be registered.
func (mockup *MockupNetworkProvider) RegisterAppDesc(networkID string, appDescriptorID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(networkID) {
        if !mockup.unsafeExistsAppDesc(networkID, appDescriptorID) {
            descriptors, _ := mockup.applicationDescriptors[networkID]
            mockup.applicationDescriptors[networkID] = append(descriptors, appDescriptorID)
            return nil
        }
        return derrors.NewOperationError(errors.AppDescAlreadyAttached).WithParams(networkID, appDescriptorID)
    }
    return derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}

// ListAppDesc lists all the application descriptors in a given network.
//   params:
//     networkID The network identifier.
//   returns:
//     An array of application descriptor identifiers.
//     An error if the list cannot be retrieved.
func (mockup *MockupNetworkProvider) ListAppDesc(networkID string) ([] string, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    descriptors, ok := mockup.applicationDescriptors[networkID]
    if ok {
        return descriptors, nil
    }
    return make([] string, 0), nil
}

// ExistsAppDesc checks if an application descriptor exists in a network.
//   params:
//     networkID The network identifier.
//     appDescriptorID The application descriptor identifier.
//   returns:
//     Whether the application exists in the given network.
func (mockup *MockupNetworkProvider) ExistsAppDesc(networkID string, appDescriptorID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExistsAppDesc(networkID, appDescriptorID)
}

func (mockup *MockupNetworkProvider) unsafeExistsAppDesc(networkID string, appDescriptorID string) bool {
    descriptors, ok := mockup.applicationDescriptors[networkID]
    if ok {
        for _, descriptor := range descriptors {
            if descriptor == appDescriptorID {
                return true
            }
        }
        return false
    }
    return false
}


// DeleteAppDescriptor deletes an application descriptor from a network.
//   params:
//     networkID The network identifier.
//     appDescriptorID The application descriptor identifier.
//   returns:
//     An error if the application descriptor cannot be removed.
func (mockup *MockupNetworkProvider) DeleteAppDescriptor(networkID string, appDescriptorID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExistsAppDesc(networkID, appDescriptorID) {
        previousDescriptors := mockup.applicationDescriptors[networkID]
        newDescriptors := make([] string, 0, len(previousDescriptors)-1)
        for _, appID := range previousDescriptors {
            if appID != appDescriptorID {
                newDescriptors = append(newDescriptors, appID)
            }
        }
        mockup.applicationDescriptors[networkID] = newDescriptors
        return nil
    }

    return derrors.NewOperationError(errors.AppDescNotAttached).WithParams(networkID, appDescriptorID)
}

// RegisterAppInst registers a new application instance inside a given network.
//   params:
//     networkID The network identifier.
//     appInstanceID The application instance identifier.
//   returns:
//     An error if the descriptor cannot be registered.
func (mockup *MockupNetworkProvider) RegisterAppInst(networkID string, appInstanceID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(networkID) {
        if !mockup.unsafeExistsAppDesc(networkID, appInstanceID) {
            instances, _ := mockup.applicationInstances[networkID]
            mockup.applicationInstances[networkID] = append(instances, appInstanceID)
            return nil
        }
        return derrors.NewOperationError(errors.AppInstAlreadyAttached).WithParams(networkID, appInstanceID)
    }

    return derrors.NewOperationError(errors.NetworkDoesNotExists).WithParams(networkID)
}

// ListAppInst list all the application instances in a given network.
//   params:
//     networkID The network identifier.
//   returns:
//     An array of application descriptor identifiers.
//     An error if the list cannot be retrieved.
func (mockup *MockupNetworkProvider) ListAppInst(networkID string) ([] string, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    instances, ok := mockup.applicationInstances[networkID]
    if ok {
        return instances, nil
    }
    return make([] string, 0), nil
}

// ExistsAppInst checks if an application instance exists in a network.
//   params:
//     networkID The network identifier.
//     appInstanceID The application instance identifier.
//   returns:
//     Whether the application exists in the given network.
func (mockup *MockupNetworkProvider) ExistsAppInst(networkID string, appInstanceID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExistsAppInst(networkID, appInstanceID)
}

func (mockup *MockupNetworkProvider) unsafeExistsAppInst(networkID string, appInstanceID string) bool {
    instances, ok := mockup.applicationInstances[networkID]
    if ok {
        for _, instance := range instances {
            if instance == appInstanceID {
                return true
            }
        }
    }
    return false
}

// DeleteAppInstance deletes an application instance from a network.
//   params:
//     networkID The network identifier.
//     appInstanceID The application instance identifier.
//   returns:
//     An error if the application cannot be removed.
func (mockup *MockupNetworkProvider) DeleteAppInstance(networkID string, appInstanceID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExistsAppInst(networkID, appInstanceID) {
        previousInstances := mockup.applicationInstances[networkID]
        newInstances := make([] string, 0, len(previousInstances)-1)
        for _, appID := range previousInstances {
            if appID != appInstanceID {
                newInstances = append(newInstances, appID)
            }
        }
        mockup.applicationInstances[networkID] = newInstances
        return nil
    }

    return derrors.NewOperationError(errors.AppInstNotAttachedToNetwork).WithParams(networkID, appInstanceID)
}

// ReducedInfoList get a list with the reduced info.
//   returns:
//     List of the reduced info.
//     An error if the info cannot be retrieved.
func (mockup *MockupNetworkProvider) ReducedInfoList() ([] entities.NetworkReducedInfo, derrors.DaishoError){
    mockup.Lock()
    defer mockup.Unlock()
    result := make([] entities.NetworkReducedInfo, 0, len(mockup.clusters))
    for _, n := range mockup.networks {
        reducedInfo := entities.NewNetworkReducedInfo(n.ID,n.Name)
        result = append(result, *reducedInfo)
    }
    return result, nil
}

// Clear util function to clean the contents of the mockup.
func (mockup *MockupNetworkProvider) Clear() {
    mockup.Lock()
    mockup.networks = make(map[string]entities.Network)
    mockup.clusters = make(map[string][] string)
    mockup.applicationDescriptors = make(map[string][] string)
    mockup.applicationInstances = make(map[string][] string)
    mockup.Unlock()
}

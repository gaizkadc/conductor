//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package nodestorage

import (
    "sync"

    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/errors"
)

// MockupNodeProvider is a mockup in-memory implementation of a node provider.
type MockupNodeProvider struct {
    sync.Mutex
    // Clusters indexed by cluster identifier.
    nodes map[string]entities.Node
}

// NewMockupNodeProvider creates a mockup provider for node operations.
func NewMockupNodeProvider() *MockupNodeProvider {
    return &MockupNodeProvider{nodes:make(map[string]entities.Node)}
}

// Add a new node to the system.
//   params:
//     node The Node to be added
//   returns:
//     An error if the node cannot be added.
func (mockup *MockupNodeProvider) Add(node entities.Node) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if !mockup.unsafeExists(node.ID){
        mockup.nodes[node.ID] = node
        return nil
    }
    return derrors.NewOperationError(errors.NodeAlreadyExists).WithParams(node)
}

// Exists checks if a node exists in the system.
//   params:
//     nodeID The Node identifier.
//   returns:
//     Whether the cluster exists or not.
func (mockup *MockupNodeProvider) Exists(nodeID string) bool {
    mockup.Lock()
    defer mockup.Unlock()
    return mockup.unsafeExists(nodeID)
}

func (mockup *MockupNodeProvider) unsafeExists(nodeID string) bool {
    _, exists := mockup.nodes[nodeID]
    return exists
}

// RetrieveNode retrieves a given node.
//   params:
//     nodeID The node identifier.
//   returns:
//     The cluster.
//     An error if the cluster cannot be retrieved.
func (mockup *MockupNodeProvider) RetrieveNode(nodeID string) (*entities.Node, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    node, exists := mockup.nodes[nodeID]
    if exists {
        return &node, nil
    }
    return nil, derrors.NewOperationError(errors.NodeDoesNotExists).WithParams(nodeID)
}

// Delete a given cluster.
//   params:
//     nodeID The node identifier.
//   returns:
//     An error if the cluster cannot be removed.
func (mockup *MockupNodeProvider) Delete(nodeID string) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    _, exists := mockup.nodes[nodeID]
    if exists {
        delete(mockup.nodes, nodeID)
        return nil
    }
    return derrors.NewOperationError(errors.NodeDoesNotExists).WithParams(nodeID)
}

// Update a node in the system.
//   params:
//     node The new node information.
//   returns:
//     An error if the instance cannot be updated.
func (mockup *MockupNodeProvider) Update(node entities.Node) derrors.DaishoError {
    mockup.Lock()
    defer mockup.Unlock()
    if mockup.unsafeExists(node.ID) {
        mockup.nodes[node.ID] = node
        return nil
    }
    return derrors.NewOperationError(errors.NodeDoesNotExists).WithParams(node)
}

// ReducedInfoList get a list with the reduced info.
//   returns:
//     List of the reduced info.
//     An error if the info cannot be retrieved.
func (mockup *MockupNodeProvider) ReducedInfoList() ([] entities.NodeReducedInfo, derrors.DaishoError) {
    mockup.Lock()
    defer mockup.Unlock()
    result := make([] entities.NodeReducedInfo, 0, len(mockup.nodes))
    for _, n := range mockup.nodes {
        reducedInfo := entities.NewNodeReducedInfo(n.NetworkID, n.ClusterID, n.ID, n.Name, n.Status, n.PublicIP)
        result = append(result, *reducedInfo)
    }
    return result, nil
}

// Dump obtains the list of all nodes in the system.
//   returns:
//     The list of nodes.
//     An error if the list cannot be retrieved.
func (mockup * MockupNodeProvider) Dump() ([] entities.Node, derrors.DaishoError) {
    result := make([] entities.Node, 0)
    mockup.Lock()
    defer mockup.Unlock()
    for _, node := range mockup.nodes {
        result = append(result, node)
    }
    return result, nil
}

// Clear is a util function to clear the contents of the mockup.
func (mockup *MockupNodeProvider) Clear() {
    mockup.Lock()
    mockup.nodes = make(map[string]entities.Node)
    mockup.Unlock()
}

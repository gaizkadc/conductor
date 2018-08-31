//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the specification of the entities inside the system model.

package entities

import (
    "github.com/satori/go.uuid"
)

// NetworkPrefix is the prefix for the network identifiers.
const NetworkPrefix = "n"

// ClusterPrefix is the prefix for the cluster identifiers.
const ClusterPrefix = "c"

// NodePrefix is the prefix for the node identifiers.
const NodePrefix = "o"

// UserPrefix is the prefix for the user identifiers.
const UserPrefix = "u"

// GenerateUUID generates a new UUID.
//   params:
//     prefix The UUID prefix.
//   returns:
//     A new UUID.
func GenerateUUID(prefix string) string {
    return prefix + uuid.NewV4().String()
}

// ReducedInfo is the object that contains the essential info stores in System Model.
type ReducedInfo struct {
    // List of networks
    Networks [] NetworkReducedInfo `json:"networks"`
    // List of clusters
    Clusters [] ClusterReducedInfo `json:"clusters"`
    // List of nodes
    Nodes [] NodeReducedInfo `json:"nodes"`
    // List of descriptors
    Descriptors [] AppDescriptorReducedInfo `json:"descriptors"`
    // List of instances
    Instances [] AppInstanceReducedInfo `json:"instances"`
    // List of users
    Users [] UserReducedInfo `json:"users"`
}

//NewReducedInfo is the basic constructor of ReducedInfo.
//   params:
//     networks     The current networks.
//     clusters     The current clusters.
//     nodes        The current nodes.
//     descriptors  The current descriptors.
//     instances    The current instances.
//     users      The current users.
//   returns:
//     A new ReducedInfo.
func NewReducedInfo(networks [] NetworkReducedInfo, clusters [] ClusterReducedInfo,
    nodes [] NodeReducedInfo, descriptors [] AppDescriptorReducedInfo,
    instances [] AppInstanceReducedInfo, users [] UserReducedInfo) *ReducedInfo {
    return &ReducedInfo{Networks: networks, Clusters: clusters, Nodes: nodes, Descriptors: descriptors,
        Instances: instances, Users: users}
}

// SummaryInfo is the object that contains basic information about the system.
type SummaryInfo struct {
    // Number of networks
    NumNetworks int `json:"numNetworks"`
    // Number of clusters
    NumClusters int `json:"numClusters"`
    // Number of nodes
    NumNodes int `json:"numNodes"`
    // Number of descriptors
    NumDescriptors int `json:"numDescriptors"`
    // Number of instances
    NumInstances int `json:"numApps"`
    // Number of users
    NumUsers int `json:"numUsers"`
    // Number of roles
    NumRoles int `json:"numRoles"`
}

//NewSummaryInfo is the basic constructor of SummaryInfo.
//   params:
//     numNetworks     The current networks.
//     numClusters     The current clusters.
//     numNodes        The current nodes.
//     numDescriptors  The current descriptors.
//     numInstances    The current instances.
//     numUsers        The current users.
//   returns:
//     A new ReducedInfo.
func NewSummaryInfo(numNetworks int, numClusters int, numNodes int,
    numDescriptors int, numInstances int, numUsers int) *SummaryInfo {
    return &SummaryInfo{NumClusters: numClusters, NumDescriptors: numDescriptors,
        NumInstances: numInstances, NumNetworks: numNetworks, NumNodes: numNodes,
        NumUsers: numUsers, NumRoles: 3}

}

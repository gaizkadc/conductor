//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Definition of the dump object that will allow export and in the future import of the system model.

package entities

import "fmt"

// Dump structure that contains the list of all elements in the system model. Expect this object to be large. Notice
// that this object contains lists of elements instead of generating a tree representation as this will speedup the
// processing of the dump by upper layers given that all entities in the system contain references to their parents.
// From the point of view of security, this structure contains all information assuming the user has already beeing
// authorized.
type Dump struct{
    // List of networks
    Networks [] Network `json:"networks"`
    // List of clusters
    Clusters [] Cluster `json:"clusters"`
    // List of nodes
    Nodes [] Node `json:"nodes"`
    // List of descriptors
    Descriptors [] AppDescriptor `json:"descriptors"`
    // List of instances
    Instances [] AppInstance `json:"instances"`
    // List of users
    Users [] User `json:"users"`
    // List of access
    UserAccesses [] UserAccess `json:"access"`
}

// NewDump creates a new dump structure.
//   params:
//     networks The array of networks.
//     clusters The array of clusters.
//     nodes The array of clusters.
//     descriptors The array of descriptors.
//     instances The array of instances.
//   returns:
//     A new Dump structure.
func NewDump(
    networks [] Network,
    clusters [] Cluster,
    nodes [] Node,
    descriptors [] AppDescriptor,
    instances [] AppInstance,
    users [] User,
    accesses [] UserAccess) * Dump {
    return & Dump{
        networks, clusters, nodes, descriptors, instances,
        users, accesses}
}
// NewDumpWithNetworks creates a new structure with networks and with all other structures empty. This method
// helps building a structure that will be dynamically modified.
//   params:
//     networks The array of networks.
//   returns:
//     A new Dump structure.
func NewDumpWithNetworks(networks [] Network) * Dump {
    return NewDump(networks, make([] Cluster, 0), make([] Node, 0), make([] AppDescriptor, 0), make([] AppInstance, 0),
        make([] User,0), make([] UserAccess, 0))
}

// AddClusters appends new clusters to the existing ones.
//   params: newClusters The clusters to be appended.
func (dump * Dump) AddClusters(newClusters [] Cluster) {
    dump.Clusters = append(dump.Clusters, newClusters...)
}

// AddCluster appends a new cluster to the existing ones.
//   params:
//     newCluster The cluster to be added.
func (dump * Dump) AddCluster(newCluster Cluster){
    dump.Clusters = append(dump.Clusters, newCluster)
}

// AddNodes appends new nodes to the existing ones.
//   params:
//     newNodes The nodes to be appended.
func (dump * Dump) AddNodes(newNodes [] Node) {
    dump.Nodes = append(dump.Nodes, newNodes...)
}

// AddNode appends a new node to the existing ones.
//   params:
//     newNode The node to be appended.
func (dump * Dump) AddNode(newNodes Node) {
    dump.Nodes = append(dump.Nodes, newNodes)
}

// AddAppDescriptors appends new AppDescriptors to the existing ones.
//   params:
//     newDescriptors The descriptors to be appended.
func (dump * Dump) AddAppDescriptors(newDescriptors [] AppDescriptor) {
    dump.Descriptors = append(dump.Descriptors, newDescriptors...)
}

// AddAppDescriptor appends a new AppDescriptor to the existing ones.
//   params:
//     newDescriptor The descriptor to be appended.
func (dump * Dump) AddAppDescriptor(newDescriptor AppDescriptor) {
    dump.Descriptors = append(dump.Descriptors, newDescriptor)
}

// AddAppInstances appends new AppInstances to the existing ones.
//   params:
//     newInstances The instances to be appended.
func (dump * Dump) AddAppInstances(newInstances [] AppInstance) {
    dump.Instances = append(dump.Instances, newInstances...)
}

// AddAppInstance appends a new AppInstance to the existing ones.
//   params:
//     newInstances The instance to be appended.
func (dump * Dump) AddAppInstance(newInstance AppInstance) {
    dump.Instances = append(dump.Instances, newInstance)
}

// AddUsers appends a list of existing users.
//   params:
//     newUsers The list of users to be appended.
func (dump * Dump) AddUsers(newUsers [] User) {
    dump.Users = append(dump.Users, newUsers...)
}

// AddAccess appends a list of existing accesses.
//   params:
//     newAccess The list of access to be appended.
func (dump * Dump) AddAccess(newAccess [] UserAccess) {
    dump.UserAccesses = append(dump.UserAccesses, newAccess...)
}

// String representation of the structure.
func (dump * Dump) String() string {
    return fmt.Sprintf("%#v", dump)
}
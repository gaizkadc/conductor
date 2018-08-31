//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Definition of the dump object that will allow export and in the future import of the system model.

package entities

import "fmt"

// Backup structure that contains the list of all elements in the system model. Expect this object to be large. Notice
// that this object contains lists of elements instead of generating a tree representation as this will speedup the
// processing of the dump by upper layers given that all entities in the system contain references to their parents.
// From the point of view of security, this structure contains all information assuming the user has already beeing
// authorized.
type BackupRestore struct{
    // List of networks
    Networks []Network `json:"networks,omitempty"`
    // List of clusters
    Clusters []Cluster `json:"clusters,omitempty"`
    // List of nodes
    Nodes []Node `json:"nodes,omitempty"`
    // List of descriptors
    AppDescriptors []AppDescriptor `json:"appdesc,omitempty"`
    // List of users
    Users  []BackupUser `json:"users,omitempty"`
}


// User data with password and access in one record

type BackupUser struct {
    User  User    `json:"user,omitempty"`
    Access UserAccess `json:"access,omitempty"`
    Password Password `json:"password,omitempty"`
}

// NewDump creates a new empty backup.
//   returns:
//     A new backup structure.
func NewBackup() *BackupRestore {
    return &BackupRestore{
        Networks: []Network{},
        Clusters: []Cluster{},
        Nodes: []Node{},
        AppDescriptors: []AppDescriptor{},
        Users: []BackupUser{},
    }
}

// AddNetworks appends new networks to the existing ones.
//   params: newNetworks The networks to be appended.
func (br *BackupRestore) AddNetworks(newNetworks []Network) {
    br.Networks = append(br.Networks, newNetworks...)
}

// AddClusters appends new clusters to the existing ones.
//   params: newClusters The clusters to be appended.
func (br *BackupRestore) AddClusters(newClusters [] Cluster) {
    br.Clusters = append(br.Clusters, newClusters...)
}

// AddCluster appends a new cluster to the existing ones.
//   params:
//     newCluster The cluster to be added.
func (br *BackupRestore) AddCluster(newCluster Cluster){
    fmt.Printf("Adding cluster ... :%v \n", newCluster)
    br.Clusters = append(br.Clusters, newCluster)
}

// AddNodes appends new nodes to the existing ones.
//   params:
//     newNodes The nodes to be appended.
func (br *BackupRestore) AddNodes(newNodes [] Node) {
    br.Nodes = append(br.Nodes, newNodes...)
}

// AddNode appends a new node to the existing ones.
//   params:
//     newNode The node to be appended.
func (br *BackupRestore) AddNode(newNodes Node) {
    br.Nodes = append(br.Nodes, newNodes)
}

// AddAppDescriptors appends new AppDescriptors to the existing ones.
//   params:
//     newDescriptors The descriptors to be appended.
func (br *BackupRestore) AddAppDescriptors(newDescriptors [] AppDescriptor) {
    br.AppDescriptors = append(br.AppDescriptors, newDescriptors...)
}

// AddAppDescriptor appends a new AppDescriptor to the existing ones.
//   params:
//     newDescriptor The descriptor to be appended.
func (br *BackupRestore) AddAppDescriptor(newDescriptor AppDescriptor) {
    br.AppDescriptors = append(br.AppDescriptors, newDescriptor)
}

// AddUsers appends a list of existing users.
//   params:
//     newUsers The list of users to be appended.
func (br *BackupRestore) AddUsers(newUsers  BackupUser) {
    br.Users = append(br.Users, newUsers)
}



// String representation of the structure.
func (br *BackupRestore) String() string {
    return fmt.Sprintf("%#v", br)

}

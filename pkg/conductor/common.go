/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package conductor

// Set of common routines for conductor components. A pool of already opened client connections is maintained
// for the components below and implemented in a singleton instance accessible by all the elements in this package.
// When running tests, this pool uses listening buffers.

import (
    pbInfrastructure "github.com/nalej/grpc-infrastructure-go"
    pbOrganization  "github.com/nalej/grpc-organization-go"
    "github.com/nalej/grpc-utils/pkg/tools"
    "github.com/nalej/conductor/pkg/utils"
    "google.golang.org/grpc"
    "github.com/rs/zerolog/log"
    "sync"
    "context"
    "fmt"
    "errors"
)

var (
    // Singleton instance of connections with cluster clients clients
    ClusterClients *tools.ConnectionsMap
    onceClusters   sync.Once
    // Singleton instance of connections with the system model
    SMClients *tools.ConnectionsMap
    onceSM sync.Once
    // Singleton instance of connections with the network client
    NetworkingClients *tools.ConnectionsMap
    onceNC sync.Once
    // Translation map between cluster ids and their ip addresses
    ClusterReference map[string]string
)

func GetSystemModelClients() *tools.ConnectionsMap {
    onceSM.Do(func(){
        // reuse the conductor factory
        SMClients = tools.NewConnectionsMap(systemModelClientFactory)
    })
    return SMClients
}


func GetClusterClients() *tools.ConnectionsMap {
    onceClusters.Do(func(){
        ClusterClients = tools.NewConnectionsMap(clusterClientFactory)
        if ClusterReference == nil {
            ClusterReference = make(map[string]string, 0)
        }
    })
    return ClusterClients
}


func GetNetworkingClients() *tools.ConnectionsMap {
    onceNC.Do(func() {
        NetworkingClients = tools.NewConnectionsMap(networkingClientFactory)
        if NetworkingClients == nil {
            NetworkingClients = tools.NewConnectionsMap(networkingClientFactory)
        }
    })
    return NetworkingClients
}


// Factory in charge of generating new basic connections with a grpc server.
//  params:
//   address the communication has to be done with
//  return:
//   client and error if any
func basicClientFactory(address string) (*grpc.ClientConn, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatal().Msgf("Failed to start gRPC connection: %v", err)
    }
    log.Info().Msgf("Connected to address at %s", address)
    return conn, err
}

// Factory in charge of generating new connections for Conductor->system model.
//  params:
//   address the communication has to be done with
//  return:
//   client and error if any
func systemModelClientFactory(address string) (*grpc.ClientConn, error) {
    return basicClientFactory(address)
}


// Factory in charge of generating new connections for Conductor->cluster communication.
//  params:
//   address the communication has to be done with
//  return:
//   client and error if any
func clusterClientFactory(address string) (*grpc.ClientConn, error) {
    return basicClientFactory(address)
}


// Factory in charge of generating new connections for Conductor->Networkcommunication.
//  params:
//   address the communication has to be done with
//  return:
//   client and error if any
func networkingClientFactory(address string) (*grpc.ClientConn, error) {
    return basicClientFactory(address)
}



// This is a common sharing function to check the system model and update the available clusters.
// Additionally, the function updates the available connections for musicians and deployment managers.
// The common ClusterReference object is updated with the cluster ids and the corresponding ip.
//  params:
//   organizationId
func UpdateClusterConnections(organizationId string) error{
    log.Debug().Msg("update cluster connections...")
    // Rebuild the map
    ClusterReference = make(map[string]string,0)

    cmClients := GetSystemModelClients()
    // no available system model client
    if cmClients.NumConnections() == 0 {
        log.Error().Msg("there are no available system model clients")
        return errors.New("there are no available system model clients")
    }

    // Get an infrastructure client and check the available clusters
    client := pbInfrastructure.NewClustersClient(cmClients.GetConnections()[0])
    req := pbOrganization.OrganizationId{OrganizationId:organizationId}
    clusterList, err := client.ListClusters(context.Background(), &req)
    if err != nil {
        msg := fmt.Sprintf("there was a problem getting the list of " +
            "available cluster for org %s",organizationId)
        log.Error().Err(err).Msg(msg)
        return errors.New(msg)
    }

    toReturn := make([]string,0)
    clusters := GetClusterClients()

    for _, cluster := range clusterList.Clusters {
        log.Debug().Msgf("add connection to cluster with id %s and hostname %s",cluster.ClusterId, cluster.Hostname)
        ClusterReference[cluster.ClusterId] = cluster.Hostname
        clusters.AddConnection(fmt.Sprintf("%s:%d",cluster.Hostname,utils.APP_CLUSTER_API_PORT))
        toReturn = append(toReturn, cluster.Hostname)
    }
    return nil
}


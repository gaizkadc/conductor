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
)

var (
    // Singleton instance of connections with musician clients
    MusicianClients *tools.ConnectionsMap
    onceMusicians   sync.Once
    // Singleton instance of connections with deployment managers
    DMClients *tools.ConnectionsMap
    onceDM sync.Once
    // Singleton instance of connections with the system model
    SMClients *tools.ConnectionsMap
    onceSM sync.Once

)

func GetSystemModelClients() *tools.ConnectionsMap {
    onceSM.Do(func(){
        // reuse the conductor factory
        SMClients = tools.NewConnectionsMap(conductorClientFactory)
    })
    return SMClients
}


func GetMusicianClients() *tools.ConnectionsMap {
    onceMusicians.Do(func(){
        MusicianClients = tools.NewConnectionsMap(conductorClientFactory)
    })
    return MusicianClients
}

func GetDMClients() *tools.ConnectionsMap {
    onceDM.Do(func() {
        DMClients = tools.NewConnectionsMap(dmClientFactory)
    })
    return DMClients
}

// Factory in charge of generating new connections for Conductor->Musician communication.
//  params:
//   address the communication has to be done with
//  return:
//   client and error if any
func conductorClientFactory(address string) (*grpc.ClientConn, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatal().Msgf("Failed to start gRPC connection: %v", err)
    }
    log.Info().Msgf("Connected to address at %s", address)
    return conn, err
}

// Factory in charge of generating new connections for Conductor->DM communication.
//  params:
//   address the communication has to be done with
//  return:
//   client and error if any
func dmClientFactory(address string) (*grpc.ClientConn, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatal().Msgf("Failed to start gRPC connection: %v", err)
    }
    log.Info().Msgf("Connected to address at %s", address)
    return conn, err
}

// This is a common sharing function to check the system model and update the available clusters.
// Additionally, the function updates the available connections for musicians and deployment managers.
//  params:
//   organizationId
//  returns:
//   list of available cluster hostnames
func UpdateClusterConnections(organizationId string) []string{
    log.Debug().Msg("update cluster connections...")
    cmClients := GetSystemModelClients()
    // no available system model client
    if cmClients.NumConnections() == 0 {
        log.Error().Msg("there are no available system model clients")
        return nil
    }

    // Get an infrastructure client and check the available clusters
    client := pbInfrastructure.NewClustersClient(cmClients.GetConnections()[0])
    req := pbOrganization.OrganizationId{OrganizationId:organizationId}
    clusterList, err := client.ListClusters(context.Background(), &req)
    if err != nil {
        log.Error().Err(err).Msgf("there was a problem getting the list of " +
            "available cluster for org %s",organizationId)
        return nil
    }

    toReturn := make([]string,0)
    musicians := GetMusicianClients()
    dms := GetDMClients()
    for _, cluster := range clusterList.Clusters {
        musicians.AddConnection(fmt.Sprintf("%s:%d",cluster.Hostname,utils.MUSICIAN_PORT))
        dms.AddConnection(fmt.Sprintf("%s:%d",cluster.Hostname,utils.DEPLOYMENT_MANAGER_PORT))
        toReturn = append(toReturn, cluster.Hostname)
    }
    return toReturn
}


/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package conductor

// Set of common routines for conductor components. A pool of already opened client connections is maintained
// for the components below and implemented in a singleton instance accessible by all the elements in this package.
// When running tests, this pool uses listening buffers.

import (
    "github.com/nalej/grpc-utils/pkg/tools"
    "google.golang.org/grpc"
    "github.com/rs/zerolog/log"
    "sync"
)

var (
    // Singleton instance of connections with musician clients
    MusicianClients *tools.ConnectionsMap
    onceMusicians   sync.Once
    // Singleton instance of connections with deployment managers
    DMClients *tools.ConnectionsMap
    onceDM sync.Once

)


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


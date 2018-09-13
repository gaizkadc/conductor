//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package musician

import (
    "github.com/rs/zerolog/log"
    "fmt"
    "net"
    "google.golang.org/grpc/reflection"
    "google.golang.org/grpc"
    "github.com/nalej/conductor/internal/musician"
    conductor "github.com/nalej/grpc-conductor-go"
    "github.com/phf/go-queue/queue"
)

type MusicianServer struct {
    // Port our service is waiting for incoming data.
    Port uint32
}

// Create a new Musician server.
//  params:
//   port where the server will be listening
//  return:
//   implementation of a Musician server
func NewMusicianServer(port uint32) *MusicianServer {
    return &MusicianServer{port}
}

func(c *MusicianServer) Run() {
    log.Info().Msg("Running musician server...")
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
    if err != nil {
        log.Fatal().Errs("failed to listen: %v", []error{err})
    }

    // Available handlers
    musician := musician.NewHandler(queue.New())

    // Server and registry
    grpcServer := grpc.NewServer()
    conductor.RegisterMusicianServer(grpcServer,musician)

    // Register reflection service on gRPC server.
    reflection.Register(grpcServer)
    log.Info().Uint32("port", c.Port).Msg("Launching gRPC server")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatal().Errs("failed to serve: %v", []error{err})
    }

}

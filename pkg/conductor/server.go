//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//
// Conductor server using gRPC.

package conductor

import (
    "github.com/rs/zerolog/log"
    "fmt"
    "net"
    "google.golang.org/grpc/reflection"
    "google.golang.org/grpc"
    "github.com/nalej/conductor/internal/conductor/deployment"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/phf/go-queue/queue"
)

type ConductorServer struct {
    // Port our service is waiting for incoming data.
    Port uint32
}

// Create a new Conductor server.
//  params:
//   port where the server will be listening
//  return:
//   implementation of a Conductor server
func NewConductorServer(port uint32) *ConductorServer {
    return &ConductorServer{port}
}

func(c *ConductorServer) Run() {
    log.Info().Msg("Running conductor server...")
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
    if err != nil {
        log.Fatal().Errs("failed to listen: %v", []error{err})
    }

    // Available handlers
    deployment := deployment.NewHandler(queue.New())

    // Server and registry
    grpcServer := grpc.NewServer()
    pbConductor.RegisterDeploymentServer(grpcServer,deployment)

    // Register reflection service on gRPC server.
    reflection.Register(grpcServer)
    log.Info().Uint32("port", c.Port).Msg("Launching gRPC server")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatal().Errs("failed to serve: %v", []error{err})
    }

}

//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//
// Conductor server using gRPC.

package conductor

import (
    "github.com/rs/zerolog/log"
    "fmt"
    "net"
    "google.golang.org/grpc"
)

type ConductorServer struct {
    Port uint32
    Server *grpc.Server
}

// Create a new Conductor server.
//  params:
//   port where the server will be listening
//  return:
//   implementation of a Conductor server
func NewConductorServer(port uint32) *ConductorServer {
    s := grpc.NewServer()
    return &ConductorServer{port, s}
}

func(c *ConductorServer) Run() {
    log.Info().Msg("Running conductor server...")
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
    if err != nil {
        log.Fatal().Errs("failed to listen: %v", []error{err})
    }

    log.Info().Uint32("port", c.Port).Msg("Launching gRPC server")
    if err := c.Server.Serve(lis); err != nil {
        log.Fatal().Errs("failed to serve: %v", []error{err})
    }

}

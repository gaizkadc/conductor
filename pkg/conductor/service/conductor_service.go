/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 */

package service

import (
    "github.com/nalej/conductor/pkg/conductor/handler"
    "github.com/phf/go-queue/queue"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/conductor/tools"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "google.golang.org/grpc/reflection"
    "github.com/rs/zerolog/log"
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/conductor/requirementscollector"
)




type ConductorService struct {
    // Conductor manager
    conductor *handler.Manager
    // Server for incoming requests
    server *tools.GenericGRPCServer
    // List of musicians to be queried
    musicians []string
    // Connections with musicians
    connections *tools.ConnectionsMap
}

func NewConductorService(port uint32, q *queue.Queue, s scorer.Scorer, reqCollector requirementscollector.RequirementsCollector,
    designer plandesigner.PlanDesigner) (*ConductorService, error) {

    c := handler.NewManager(q, s, reqCollector, designer,port)
    conductorServer := tools.NewGenericGRPCServer(port)
    // instance := ConductorService{c, conductorServer, make([]string, 0)},connections)
    instance := ConductorService{conductor: c,
                                server: conductorServer,
                                musicians: make([]string,0),
                                connections: conductor.GetMusicianClients()}
    return &instance, nil
}



func(c *ConductorService) Run() {
    // register services
    deployment := handler.NewHandler(c.conductor)

    // Server and registry
    //grpcServer := grpc.NewServer()
    pbConductor.RegisterConductorServer(c.server.Server,deployment)

    // Register reflection service on gRPC server.
    reflection.Register(c.server.Server)

    go c.conductor.Run()
    c.server.Run()

}

// Set the musicians to be queried
// TODO: this has to be removed and check the system model instead. This is only for initial testing.
func(c *ConductorService) SetMusicians(musicians []string) {
    c.musicians=musicians
    for _, target := range musicians {
        _,err := c.connections.AddConnection(target)
        if err != nil {
            log.Error().Err(err)
        } else {
            log.Info().Str("address",target).Msg("musician address correctly added")
            c.musicians = append(c.musicians, target)
        }
    }
}
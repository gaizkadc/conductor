/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package service

import (
    "github.com/nalej/conductor/pkg/conductor/handler"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/grpc-utils/pkg/tools"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "google.golang.org/grpc/reflection"
    "github.com/rs/zerolog/log"
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/conductor/requirementscollector"
)

type ConductorConfig struct {
    // incoming port
    Port uint32
    // URL where the system model is available
    SystemModelURL string
    // List of musicians to be queried
    Musicians []string
}


type ConductorService struct {
    // Conductor manager
    conductor *handler.Manager
    // Server for incoming requests
    server *tools.GenericGRPCServer
    // Connections with musicians
    connections *tools.ConnectionsMap
    // Configuration object
    configuration *ConductorConfig
}


func NewConductorService(config *ConductorConfig) (*ConductorService, error) {

    // Initialize connections pool with system model
    smPool := conductor.GetSystemModelClients()
    _, err := smPool.AddConnection(config.SystemModelURL)
    if err != nil {
        log.Error().Err(err).Msg("error creating connection with system model")
        return nil, err
    }

    // Confligure cluster entries
    // TODO get from the system model
    InitPool(config.Musicians, conductor.GetMusicianClients())
    //InitPool([]string{"localhost:5200"}, conductor.GetDMClients())
    //SetMusicians(config.Musicians)

    q := handler.NewMemoryRequestQueue()
    scr := scorer.NewSimpleScorer()
    reqColl := requirementscollector.NewSimpleRequirementsCollector()
    designer := plandesigner.NewSimplePlanDesigner()

    c := handler.NewManager(q, scr, reqColl, designer)
    conductorServer := tools.NewGenericGRPCServer(config.Port)
    instance := ConductorService{conductor: c,
                                server: conductorServer,
                                connections: conductor.GetMusicianClients(),
                                configuration: config}

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

// Initialize a connections pool with a set of addresses.
func InitPool(addresses []string, connections *tools.ConnectionsMap) {
    for _, target := range addresses {
        _,err := connections.AddConnection(target)
        if err != nil {
            log.Error().Err(err)
        } else {
            log.Info().Str("address",target).Msg("correctly added to the connections pool")
        }
    }
}
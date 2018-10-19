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
    "github.com/nalej/conductor/pkg/conductor/monitor"
)

type ConductorConfig struct {
    // incoming port
    Port uint32
    // URL where the system model is available
    SystemModelURL string
}

func (conf * ConductorConfig) Print() {
    log.Info().Uint32("port", conf.Port).Msg("gRPC port")
    log.Info().Str("URL", conf.SystemModelURL).Msg("System Model")
}


type ConductorService struct {
    // Conductor manager
    conductor *handler.Manager
    // Conductor monitor
    monitor *monitor.Manager
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

    //InitPool(config.Musicians, conductor.GetMusicianClients())

    q := handler.NewMemoryRequestQueue()
    scr := scorer.NewSimpleScorer()
    reqColl := requirementscollector.NewSimpleRequirementsCollector()
    designer := plandesigner.NewSimplePlanDesigner()
    monitor := monitor.NewManager()

    if monitor == nil {
        log.Panic().Msg("impossible to create monitor service")
        return nil, err
    }

    c := handler.NewManager(q, scr, reqColl, designer, *monitor)

    conductorServer := tools.NewGenericGRPCServer(config.Port)
    instance := ConductorService{conductor: c,
                                monitor: monitor,
                                server: conductorServer,
                                connections: conductor.GetMusicianClients(),
                                configuration: config}

    return &instance, nil
}



func(c *ConductorService) Run() {
    // register services
    conductorService := handler.NewHandler(c.conductor)
    monitorService := monitor.NewHandler(c.monitor)

    // Server and registry
    // -- conductor service
    pbConductor.RegisterConductorServer(c.server.Server, conductorService)
    // -- monitor service
    pbConductor.RegisterConductorMonitorServer(c.server.Server, monitorService)

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
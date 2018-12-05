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
    "net"
    "fmt"
)

type ConductorConfig struct {
    // incoming port
    Port uint32
    // URL where the system model is available
    SystemModelURL string
    // URL where the networking client is available
    NetworkingServiceURL string
}

func (conf * ConductorConfig) Print() {
    log.Info().Uint32("port", conf.Port).Msg("gRPC port")
    log.Info().Str("URL", conf.SystemModelURL).Msg("System Model")
    log.Info().Str("NetworkingServiceURL", conf.NetworkingServiceURL).Msg("Networking service URL")
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

    // Initialize connections pool with networking client
    cnPool := conductor.GetNetworkingClients()
    _, err = cnPool.AddConnection(config.NetworkingServiceURL)

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
                                connections: conductor.GetClusterClients(),
                                configuration: config}

    return &instance, nil
}



func(c *ConductorService) Run() {

    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.configuration.Port))
    if err != nil {
        log.Fatal().Errs("failed to listen: %v", []error{err})
    }

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

    // Launch the main deployment manager in a separate routine
    go c.conductor.Run()

    // Run
    log.Info().Uint32("port", c.configuration.Port).Msg("Launching gRPC server")
    if err := c.server.Server.Serve(lis); err != nil {
        log.Fatal().Errs("failed to serve: %v", []error{err})
    }

}

/*
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
*/
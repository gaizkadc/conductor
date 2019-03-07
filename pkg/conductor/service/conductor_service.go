/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package service

import (
    "errors"
    "github.com/nalej/conductor/internal/structures"
    "github.com/nalej/conductor/pkg/conductor/baton"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/grpc-utils/pkg/tools"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "google.golang.org/grpc/reflection"
    "github.com/rs/zerolog/log"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/conductor/requirementscollector"
    "github.com/nalej/conductor/pkg/conductor/monitor"
    "net"
    "fmt"
    "github.com/nalej/conductor/pkg/utils"
    "strconv"
)

type ConductorConfig struct {
    // incoming port
    Port uint32
    // URL where the system model is available
    SystemModelURL string
    // URL where the networking client is available
    NetworkingServiceURL string
    // AppClusterAPI port
    AppClusterApiPort uint32
    // UseTLSForClusterAPI defines if TLS should be used to connect to the cluster API.
    UseTLSForClusterAPI bool
    // Path for the certificate of the CA
    CACertPath string
    // Skip CA validation
    SkipCAValidation bool
    // URL where authx client is available
    AuthxURL string
    //UnifiedLogging client is available
    UnifiedLoggingURL string
}

func (conf * ConductorConfig) Print() {
    log.Info().Uint32("port", conf.Port).Msg("gRPC port")
    log.Info().Str("URL", conf.SystemModelURL).Msg("System Model")
    log.Info().Str("NetworkingServiceURL", conf.NetworkingServiceURL).Msg("Networking service URL")
    log.Info().Str("AuthxURL", conf.AuthxURL).Msg("Authx service URL")
    log.Info().Str("UnifiedLoggingURL", conf.UnifiedLoggingURL).Msg("UnifiedLogging service URL")
    log.Info().Uint32("appclusterport", conf.AppClusterApiPort).Msg("appClusterApi gRPC port")
    log.Info().Bool("useTLS", conf.UseTLSForClusterAPI).Msg("Use TLS to connect the the Application Cluster API")
}


type ConductorService struct {
    // Conductor manager
    conductor *baton.Manager
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

    // TODO review this global set
    // set global port
    utils.APP_CLUSTER_API_PORT = config.AppClusterApiPort

    connectionsHelper := utils.NewConnectionsHelper(config.UseTLSForClusterAPI,config.CACertPath,config.SkipCAValidation)

    // Initialize connections pool with system model
    smPool := connectionsHelper.GetSystemModelClients()
    host, port, err := net.SplitHostPort(config.SystemModelURL)
    if err != nil {
        log.Fatal().Err(err).Msg("error getting the system model url")
    }

    p, err := strconv.Atoi(port)
    if err != nil {
        log.Fatal().Err(err).Msg("error getting the system model port")
    }

    _, err = smPool.AddConnection(host, p)
    if err != nil {
        log.Error().Err(err).Msg("error creating connection with system model")
        return nil, err
    }

    // Initialize connections pool with networking client
    cnPool := connectionsHelper.GetNetworkingClients()

    netHost, netPort, err := net.SplitHostPort(config.NetworkingServiceURL)
    if err != nil {
        log.Fatal().Err(err).Msg("error getting the system model url")
    }

    netp, err := strconv.Atoi(netPort)
    if err != nil {
        log.Fatal().Err(err).Msg("error getting the system model port")
    }

    _, err = cnPool.AddConnection(netHost, netp)
    if err != nil {
        log.Error().Err(err).Msg("error creating connection with system model")
        return nil, err
    }

    // Initialize connections pool with authx client
    authxPool := connectionsHelper.GetAuthxClients()
    authxHost, authxPort, err := net.SplitHostPort(config.AuthxURL)
    if err != nil {
        log.Fatal().Err(err).Msg("error getting the authx url")
    }

    authxp, err := strconv.Atoi(authxPort)
    if err != nil {
        log.Fatal().Err(err).Msg("error getting the authx port")
    }

    _, err = authxPool.AddConnection(authxHost, authxp)
    if err != nil {
        log.Error().Err(err).Msg("error creating connection with authx")
        return nil, err
    }

    uLoggingPool := connectionsHelper.GetUnifiedLoggingClients()
    uLoggingHost, uLogginPort, err := net.SplitHostPort(config.UnifiedLoggingURL)
    if err != nil {
        log.Fatal().Err(err).Msg("error getting the unified logging url")
    }
    uLogginp, err := strconv.Atoi(uLogginPort)
    if err != nil {
        log.Fatal().Err(err).Msg("error getting the unified logging port")
    }

    _, err = uLoggingPool.AddConnection(uLoggingHost, uLogginp)
    if err != nil {
        log.Error().Err(err).Msg("error creating connection with unified logging")
        return nil, err
    }


    log.Info().Msg("instantiate local queue in memory...")
    q := structures.NewMemoryRequestQueue()
    log.Info().Msg("done")
    log.Info().Msg("instantiate local pending plans structure...")
    pendingPlans := structures.NewPendingPlans()
    log.Info().Msg("done")
    scr := scorer.NewSimpleScorer(connectionsHelper)
    reqColl := requirementscollector.NewSimpleRequirementsCollector()

    //designer := plandesigner.NewSimplePlanDesigner(connectionsHelper)
    designer := plandesigner.NewSimpleReplicaPlanDesigner(connectionsHelper)

    batonMgr := baton.NewManager(connectionsHelper, q, scr, reqColl, designer,pendingPlans)
    if batonMgr == nil {
        log.Panic().Msg("impossible to create baton service")
        return nil, errors.New("impossible to create baton service")
    }

    monitorMgr := monitor.NewManager(connectionsHelper,q,pendingPlans, batonMgr)
    if monitorMgr == nil {
        log.Panic().Msg("impossible to create monitorMgr service")
        return nil, err
    }


    conductorServer := tools.NewGenericGRPCServer(config.Port)
    instance := ConductorService{conductor: batonMgr,
                                monitor: monitorMgr,
                                server: conductorServer,
                                connections: connectionsHelper.GetClusterClients(),
                                configuration: config}



    return &instance, nil
}


func(c *ConductorService) Run() {

    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.configuration.Port))
    if err != nil {
        log.Fatal().Errs("failed to listen: %v", []error{err})
    }

    // register services
    conductorService := baton.NewHandler(c.conductor)
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

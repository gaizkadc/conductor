/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package service

import (
    "errors"
    "github.com/nalej/conductor/internal/persistence/app_cluster"
    "github.com/nalej/conductor/internal/structures"
    "github.com/nalej/conductor/pkg/conductor/baton"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/conductor/pkg/provider/kv"
    "github.com/nalej/grpc-utils/pkg/tools"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/nalej-bus/pkg/queue/network/ops"
    "google.golang.org/grpc"

    "google.golang.org/grpc/reflection"
    "github.com/rs/zerolog/log"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/conductor/requirementscollector"
    "github.com/nalej/conductor/pkg/conductor/monitor"
    "github.com/nalej/conductor/pkg/conductor/queue"
    "github.com/nalej/nalej-bus/pkg/bus/pulsar-comcast"
    queueAppOps "github.com/nalej/nalej-bus/pkg/queue/application/ops"
    queueInfrOps "github.com/nalej/nalej-bus/pkg/queue/infrastructure/ops"
    queueInfrEvents "github.com/nalej/nalej-bus/pkg/queue/infrastructure/events"
    queueNetOps "github.com/nalej/nalej-bus/pkg/queue/network/ops"
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
    // Client Cert Path
    ClientCertPath string
    // Skip Server validation
    SkipServerCertValidation bool
    // URL where authx client is available
    AuthxURL string
    // UnifiedLogging client is available
    UnifiedLoggingURL string
    // Queue service
    QueueURL string
    // Folder for the local database
    DBFolder string
    // Debugging flag
    Debug bool

}

func (conf * ConductorConfig) Print() {
    log.Info().Uint32("port", conf.Port).Msg("gRPC port")
    log.Info().Str("URL", conf.SystemModelURL).Msg("System Model")
    log.Info().Str("NetworkingServiceURL", conf.NetworkingServiceURL).Msg("Networking service URL")
    log.Info().Str("AuthxURL", conf.AuthxURL).Msg("Authx service URL")
    log.Info().Str("UnifiedLoggingURL", conf.UnifiedLoggingURL).Msg("UnifiedLogging service URL")
    log.Info().Str("QueueURL", conf.QueueURL).Msg("Queue service URL")
    log.Info().Uint32("appclusterport", conf.AppClusterApiPort).Msg("appClusterApi gRPC port")
    log.Info().Bool("useTLS", conf.UseTLSForClusterAPI).Msg("Use TLS to connect the the Application Cluster API")
    log.Info().Str("DBFolder", conf.DBFolder).Msg("Folder for the local database")
    log.Info().Bool("Debug", conf.Debug).Msg("Debug enabled")
    log.Info().Bool("SkipServerCertValidation", conf.SkipServerCertValidation).Msg("SkipServerCertValidation enabled")
    log.Info().Str("CACertPath", conf.CACertPath).Msg("CA cert path")
    log.Info().Str("ClientCertPath", conf.ClientCertPath).Msg("Client cert path")
}


type ConductorService struct {
    // Conductor manager
    conductor *baton.Manager
    // Conductor monitor
    monitor *monitor.Manager
    // Server for incoming requests
    server *grpc.Server
    // Connections with musicians
    connections *tools.ConnectionsMap
    // Configuration object
    configuration *ConductorConfig
    // Application ops consumer
    appOpsConsumer *queueAppOps.ApplicationOpsConsumer
    // infrastructure ops consumer
    infOpsConsumer *queueInfrOps.InfrastructureOpsConsumer
    // infrastructure events consumer
    infEventsConsumer *queueInfrEvents.InfrastructureEventsConsumer
    // network operations producer
    networkOpsProducer *queueNetOps.NetworkOpsProducer
}


func NewConductorService(config *ConductorConfig) (*ConductorService, error) {

    // TODO review this global set
    // set global port
    utils.APP_CLUSTER_API_PORT = config.AppClusterApiPort

    connectionsHelper := utils.NewConnectionsHelper(config.UseTLSForClusterAPI,config.ClientCertPath,config.CACertPath,config.SkipServerCertValidation)

    // Initialize connections pool with system model
    log.Info().Msg("initialize system model client...")
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
    log.Info().Msg("done")

    // Initialize connections pool with networking client
    log.Info().Msg("initialize network client...")
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
    log.Info().Msg("done")

    // Initialize connections pool with authx client
    log.Info().Msg("initialize authx client...")
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
    log.Info().Msg("done")

    // Initialize unified logging service
    log.Info().Msg("initialize unified logging client...")
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
    log.Info().Msg("done")


    // Instantiate pulsar client
    log.Info().Str("address", config.QueueURL).Msg("instantiate pulsar comcast client")
    // Instantiate message queue clients
    log.Info().Msg("initialize message queue client...")
    pulsarClient := pulsar_comcast.NewClient(config.QueueURL)
    log.Info().Msg("done")

    log.Info().Msg("initialize application ops client...")
    appOpsConfig := queueAppOps.NewConfigApplicationOpsConsumer(1,
        queueAppOps.ConsumableStructsApplicationOpsConsumer{true, true})
    appsOps,err := queueAppOps.NewApplicationOpsConsumer(pulsarClient, "conductor-app-ops", true, appOpsConfig)
    if err != nil {
        log.Panic().Err(err).Msg("impossible to initialize application ops queue client")
    }
    log.Info().Msg("done")

    log.Info().Msg("initialize infrastructure ops client...")
    infrOpsConfig := queueInfrOps.NewConfigInfrastructureOpsConsumer(1,
        queueInfrOps.ConsumableStructsInfrastructureOpsConsumer{DrainRequest: true})
    infrOps, err := queueInfrOps.NewInfrastructureOpsConsumer(pulsarClient, "conductor-infr-ops", true, infrOpsConfig)
    if err != nil {
        log.Panic().Err(err).Msg("impossible to initialize infrastructure ops queue client")
    }
    log.Info().Msg("done")

    log.Info().Msg("initialize infrastructure events client...")
    infrEventsConfig := queueInfrEvents.NewConfigInfrastructureEventsConsumer(1,
        queueInfrEvents.ConsumableStructsInfrastructureEventsConsumer{UpdateClusterRequest: true, SetClusterStatusRequest: true})
    infrEvents, err := queueInfrEvents.NewInfrastructureEventsConsumer(pulsarClient, "conductor-infr-events", true, infrEventsConfig)
    if err != nil {
        log.Panic().Err(err).Msg("impossible to initialize infrastructure events queue client")
    }
    log.Info().Msg("done")

    log.Info().Msg("initialize network ops producer...")
    netOpsProducer, err := ops.NewNetworkOpsProducer(pulsarClient, "conductor-network-ops")
    if err != nil {
        log.Panic().Err(err).Msg("impossible to initialize network ops queue client")
    }
    log.Info().Msg("done")


    // TODO replace this memory queue by the system queue solution
    log.Info().Msg("instantiate local queue in memory...")
    q := structures.NewMemoryRequestQueue()
    log.Info().Msg("done")
    log.Info().Msg("instantiate local pending plans structure...")
    pendingPlans := structures.NewPendingPlans()
    log.Info().Msg("done")
    scr := scorer.NewSimpleScorer(connectionsHelper)
    reqColl := requirementscollector.NewSimpleRequirementsCollector()

    log.Info().Msg("instantiate plan designer...")
    designer := plandesigner.NewSimpleReplicaPlanDesigner(connectionsHelper)
    log.Info().Msg("done")

    log.Info().Msg("instantiate local app cluster db...")

    boltProvider, err := kv.NewLocalDB(config.DBFolder+"/appcluster.db")
    if err != nil {
        log.Panic().Err(err).Msgf("impossible to instantiate bolt provider for appcluster in %s",config.DBFolder)
        return nil, err
    }
    appClusterDB := app_cluster.NewAppClusterDB(boltProvider)
    log.Info().Msg("done")

    batonMgr := baton.NewManager(connectionsHelper, q, scr, reqColl, designer,pendingPlans, appClusterDB,netOpsProducer)
    if batonMgr == nil {
        log.Panic().Msg("impossible to create baton service")
        return nil, errors.New("impossible to create baton service")
    }

    monitorMgr := monitor.NewManager(connectionsHelper,q,pendingPlans, batonMgr)
    if monitorMgr == nil {
        log.Panic().Msg("impossible to create monitorMgr service")
        return nil, err
    }


    conductorServer := grpc.NewServer()
    instance := ConductorService{conductor: batonMgr,
                                monitor:            monitorMgr,
                                server:             conductorServer,
                                connections:        connectionsHelper.GetClusterClients(),
                                configuration:      config,
                                appOpsConsumer:     appsOps,
                                infOpsConsumer:     infrOps,
                                infEventsConsumer:  infrEvents,
                                networkOpsProducer: netOpsProducer,
    }


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
    pbConductor.RegisterConductorServer(c.server, conductorService)
    // -- monitor service
    pbConductor.RegisterConductorMonitorServer(c.server, monitorService)


    // Register reflection service on gRPC server.
    if c.configuration.Debug {
        reflection.Register(c.server)
    }


    log.Info().Msg("run application ops handler...")
    appOpsQueue := queue.NewApplicationOpsHandler(c.conductor, c.appOpsConsumer)
    appOpsQueue.Run()
    log.Info().Msg("done")

    log.Info().Msg("run infrastructure ops handler...")
    infrOpsQueue := queue.NewInfrastructureOpsHandler(c.conductor, c.infOpsConsumer)
    infrOpsQueue.Run()
    log.Info().Msg("done")

    log.Info().Msg("run infrastructure events handler...")
    infrEventsQueue := queue.NewInfrastructureEventsHandler(c.conductor, c.infEventsConsumer)
    infrEventsQueue.Run()
    log.Info().Msg("done")




    // Launch the main deployment manager in a separate routine
    go c.conductor.Run()

    // Run
    log.Info().Uint32("port", c.configuration.Port).Msg("Launching gRPC server")
    if err := c.server.Serve(lis); err != nil {
        log.Fatal().Errs("failed to serve: %v", []error{err})
    }

}

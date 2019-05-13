/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package cmd

import (
    "github.com/spf13/cobra"
    "github.com/rs/zerolog/log"
    "github.com/spf13/viper"
    "github.com/nalej/conductor/pkg/conductor/service"
    "fmt"
    "github.com/nalej/conductor/pkg/utils"
)


var runCmd = &cobra.Command{
    Use: "run",
    Short: "Run conductor",
    Long: "Run conductor service with... and with...",
    Run: func(cmd *cobra.Command, args [] string) {
        SetupLogging()
        RunConductor()
    },
}


func init() {

    RootCmd.AddCommand(runCmd)

    runCmd.Flags().Uint32P("port", "c",utils.CONDUCTOR_PORT,"port where conductor listens to")
    runCmd.Flags().StringP("systemModelAddress","s",fmt.Sprintf("localhost:%d",utils.SYSTEM_MODEL_PORT),
        "host:port address for system model")
    runCmd.Flags().StringP("networkManagerAddress", "n", fmt.Sprintf("localhost:%d", utils.NETWORKING_SERVICE_PORT),
        "host:port address for networking manager")
    runCmd.Flags().StringP("authxAddress","a",fmt.Sprintf("localhost:%d",utils.AUTHX_PORT),
        "host:port address for authx")
    runCmd.Flags().Uint32P("appClusterPort","p",utils.APP_CLUSTER_API_PORT, "port where the application cluster api is listening")
    runCmd.Flags().Bool("useTLS", true, "Use TLS to connect to the application cluster API")
    runCmd.Flags().String("caCertPath", "", "Part for the CA certificate")
    runCmd.Flags().Bool("skipCAValidation", true, "Skip CA authentication validation")
    runCmd.Flags().StringP("unifiedLogging", "u",fmt.Sprintf("localhost:%d",utils.UNIFIED_LOGGING_PORT),
        "host:port address for unifiedLogging")
    runCmd.Flags().StringP("queueAddress", "q", fmt.Sprintf("localhost:%d", utils.QUEUE_PORT),
        "host:port address for the Nalej management queue")
    runCmd.Flags().StringP("dbFolder", "f", "/data/",
        "path for the folder used to store the local database")

    viper.BindPFlags(runCmd.Flags())
}

// Entrypoint for a musician service.
func RunConductor() {
    // Incoming requests port
    var port uint32
    // System model url
    var systemModel string
    // Networking service url
    var networkingService string
    // AppClusterAPI port
    var appClusterApiPort uint32
    // useTLS boolean
    var useTLS bool
    // CA cert path
    var caCertPath string
    // Skip CA validation
    var skipCAValidation bool
    // Authx url
    var authxService string
    // Unified Logging url
    var unifiedLoggingService string
    // Queue url
    var queueAddress string
    // Database folder path
    var dbFolder string

    port = uint32(viper.GetInt32("port"))
    systemModel = viper.GetString("systemModelAddress")
    networkingService = viper.GetString("networkManagerAddress")
    authxService = viper.GetString("authxAddress")
    appClusterApiPort = uint32(viper.GetInt32("appClusterPort"))
    useTLS = viper.GetBool("useTLS")
    caCertPath = viper.GetString("caCertPath")
    skipCAValidation = viper.GetBool("skipCAValidation")
    unifiedLoggingService = viper.GetString("unifiedLogging")
    queueAddress = viper.GetString("queueAddress")
    dbFolder = viper.GetString("dbFolder")


    log.Info().Msg("launching conductor...")


    config := service.ConductorConfig{
        Port: port,
        SystemModelURL: systemModel,
        NetworkingServiceURL: networkingService,
        AppClusterApiPort: appClusterApiPort,
        UseTLSForClusterAPI: useTLS,
        CACertPath: caCertPath,
        SkipCAValidation: skipCAValidation,
        AuthxURL:authxService,
        UnifiedLoggingURL:unifiedLoggingService,
        QueueURL: queueAddress,
        DBFolder: dbFolder,
    }
    config.Print()

    conductorService, err := service.NewConductorService(&config)
    if err != nil {
        log.Fatal().AnErr("err", err).Msg("impossible to initialize conductor service")
    }
    conductorService.Run()
}
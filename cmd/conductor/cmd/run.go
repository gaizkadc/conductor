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
    runCmd.Flags().Uint32P("appClusterPort","p",utils.APP_CLUSTER_API_PORT, "port where the application cluster api is listening")

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

    port = uint32(viper.GetInt32("port"))
    systemModel = viper.GetString("systemModelAddress")
    networkingService = viper.GetString("networkManagerAddress")
    appClusterApiPort = uint32(viper.GetInt32("appClusterPort"))

    log.Info().Msg("launching conductor...")



    config := service.ConductorConfig{
        Port: port,
        SystemModelURL: systemModel,
        NetworkingServiceURL: networkingService,
        AppClusterApiPort: appClusterApiPort,
    }
    config.Print()
    conductorService, err := service.NewConductorService(&config)
    if err != nil {
        log.Fatal().AnErr("err", err).Msg("impossible to initialize conductor service")
    }
    conductorService.Run()
}
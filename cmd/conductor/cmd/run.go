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

    runCmd.Flags().Uint32P("conductor-port", "c",utils.CONDUCTOR_PORT,"port where conductor listens to")
    runCmd.Flags().StringP("systemmodel","s",fmt.Sprintf("localhost:%d",utils.SYSTEM_MODEL_PORT),
        "host:port indicating where is available the system model")

    viper.BindPFlags(runCmd.Flags())
}

// Entrypoint for a musician service.
func RunConductor() {
    // Incoming requests port
    var port uint32
    // System model url
    var systemModel string

    port = uint32(viper.GetInt32("conductor-port"))
    systemModel = viper.GetString("systemmodel")

    log.Info().Msg("launching conductor...")


    config := service.ConductorConfig{
        Port: port,
        SystemModelURL: systemModel,
    }
    config.Print()
    conductorService, err := service.NewConductorService(&config)
    if err != nil {
        log.Fatal().AnErr("err", err).Msg("impossible to initialize conductor service")
    }
    conductorService.Run()
}
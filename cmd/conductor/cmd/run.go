//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package cmd

import (
    "github.com/spf13/cobra"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/spf13/viper"
    "github.com/nalej/conductor/pkg/conductor"
)

// Incoming requests port
var port uint32


var runCmd = &cobra.Command{
    Use: "run",
    Short: "Run conductor",
    Long: "Run conductor service with... and with...",
    Run: func(cmd *cobra.Command, args [] string) {
        RunConductor()
    },
}


func init() {

    RootCmd.AddCommand(runCmd)

    runCmd.Flags().Uint32P("port", "p",5000,"port where conductor listens to")

    viper.BindPFlags(runCmd.Flags())
}

// Entrypoint for a musician service.
func RunConductor() {
    // UNIX Time is faster and smaller than most timestamps
    // If you set zerolog.TimeFieldFormat to an empty string,
    // logs will write with UNIX time
    zerolog.TimeFieldFormat = ""

    port = uint32(viper.GetInt32("port"))

    log.Info().Msg("launching conductor...")

    conductorServer := conductor.NewConductorServer(port)
    conductorServer.Run()
}
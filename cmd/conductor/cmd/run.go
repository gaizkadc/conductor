//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package cmd

import (
    "github.com/spf13/cobra"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/spf13/viper"
    "github.com/nalej/conductor/pkg/conductor/service"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/phf/go-queue/queue"
)

// Incoming requests port
var port uint32
// Array of musician addresses
var musicians[]string


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
    runCmd.Flags().StringArrayP("musicians", "m", make([]string,0),"list of addresses for musicians (192.168.1.1:3000, 127.0.0.1:3000)")

    viper.BindPFlags(runCmd.Flags())
}

// Entrypoint for a musician service.
func RunConductor() {
    // UNIX Time is faster and smaller than most timestamps
    // If you set zerolog.TimeFieldFormat to an empty string,
    // logs will write with UNIX time
    zerolog.TimeFieldFormat = ""

    port = uint32(viper.GetInt32("port"))
    musicians = viper.GetStringSlice("musicians")

    log.Info().Msg("launching conductor...")

    q := queue.New()
    scr := scorer.NewSimpleScorer()
    conductorService, err := service.NewConductorService(port, q, scr)
    conductorService.SetMusicians(musicians)
    if err != nil {
        log.Fatal().AnErr("err", err).Msg("impossible to initialize conductor service")
    }

    conductorService.Run()
}
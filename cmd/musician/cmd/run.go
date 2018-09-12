//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package cmd

import (
    "github.com/spf13/cobra"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/spf13/viper"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
)

// Prometheus URL
var prometheus string
// Time to sleep between monitoring queries
var sleepTime uint32


var runCmd = &cobra.Command{
    Use: "run",
    Short: "Run a musician service",
    Long: "Run a musician service for the cluster this node belongs to",
    Run: func(cmd *cobra.Command, args [] string) {
        RunMusician()
    },
}


func init() {

    RootCmd.AddCommand(runCmd)

    runCmd.Flags().StringP("prometheus", "p", "", "prometheus endpoint")
    runCmd.Flags().Uint32P("sleep", "s",10000,"time to sleep between queries in milliseconds")

    viper.BindPFlags(runCmd.Flags())
}

// Entrypoint for a musician service.
func RunMusician() {
    // UNIX Time is faster and smaller than most timestamps
    // If you set zerolog.TimeFieldFormat to an empty string,
    // logs will write with UNIX time
    zerolog.TimeFieldFormat = ""

    prometheus = viper.GetString("prometheus")
    sleepTime = uint32(viper.GetInt32("sleep"))

    log.Info().Msg("launching status collector...")

    collector := statuscollector.NewPrometheusStatusCollector(prometheus, sleepTime)
    collector.Run()

}
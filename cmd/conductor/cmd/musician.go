/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/rs/zerolog/log"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
    "github.com/nalej/conductor/pkg/musician/service"
    "github.com/nalej/conductor/pkg/musician/scorer"
    "os"
    "github.com/nalej/conductor/pkg/utils"
)

var musicianCmd = &cobra.Command{
    Use: "musician",
    Short: "Run a musician service",
    Long: "Run a musician service for the cluster this node belongs to",
    Run: func(cmd *cobra.Command, args [] string) {
        SetupLogging()
        RunMusician()
    },
}


func init() {

    RootCmd.AddCommand(musicianCmd)

    musicianCmd.Flags().Uint32P("musician-port", "u",utils.MUSICIAN_PORT,"musician endpoint")
    musicianCmd.Flags().StringP("prometheus", "o", "", "prometheus endpoint")
    musicianCmd.Flags().Uint32P("sleep", "s",10000,"time to sleep between queries in milliseconds")

    viper.BindPFlags(musicianCmd.Flags())
}

// Entrypoint for a musician service.
func RunMusician() {
    // Prometheus URL
    var prometheus string
    // Time to sleep between monitoring queries
    var sleepTime uint32
    // Application port
    var port uint32


    port = uint32(viper.GetInt32("musician-port"))
    prometheus = viper.GetString("prometheus")
    sleepTime = uint32(viper.GetInt32("sleep"))

    log.Info().Msg("launching musician...")
    collector := statuscollector.NewPrometheusStatusCollector(prometheus, sleepTime)
    go collector.Run()

    scorer := scorer.NewSimpleScorer(collector)

    conf := &service.MusicianConfig{
        Port: port,
        Scorer: &scorer,
        Collector: &collector,
    }

    musicianService, err := service.NewMusicianService(conf)

    if err!=nil{
        log.Fatal().AnErr("error",err).Msg("impossible to start service")
        os.Exit(1)
    }

    musicianService.Run()

}

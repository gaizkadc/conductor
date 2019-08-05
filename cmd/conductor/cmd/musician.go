/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package cmd

import (
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/rs/zerolog/log"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
    "github.com/nalej/conductor/pkg/musician/service"
    "github.com/nalej/conductor/pkg/musician/scorer"
    "github.com/nalej/conductor/pkg/utils"

    "github.com/nalej/grpc-monitoring-go"
    "google.golang.org/grpc"
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
    musicianCmd.Flags().StringP("metrics", "m", "", "metrics api endpoint")
    // 60s is default Prometheus scrape time - no use in collecting status more often
    musicianCmd.Flags().Uint32P("sleep", "s",60000,"time to sleep between queries in milliseconds")

    viper.BindPFlags(musicianCmd.Flags())
}

// Entrypoint for a musician service.
func RunMusician() {
    // Prometheus URL
    var prometheus string
    // Metrics collector address
    var metrics string
    // Time to sleep between monitoring queries
    var sleepTime uint32
    // Application port
    var port uint32
    // Debug flag
    var debug bool


    port = uint32(viper.GetInt32("musician-port"))
    prometheus = viper.GetString("prometheus")
    metrics = viper.GetString("metrics")
    sleepTime = uint32(viper.GetInt32("sleep"))
    debug = viper.GetBool("debug")

    log.Info().Msg("launching musician...")

    if prometheus != "" && metrics != "" {
        log.Fatal().Msg("only one of 'prometheus' and 'metrics' can be set")
    }
    if prometheus == "" && metrics == "" {
        log.Fatal().Msg("one of 'prometheus' or 'metrics' should be set")
    }

    var collector statuscollector.StatusCollector

    if prometheus != "" {
        collector = statuscollector.NewPrometheusStatusCollector(prometheus, sleepTime)
    }

    if metrics != "" {
        metricsConn, err := grpc.Dial(metrics, grpc.WithInsecure())
        if err != nil {
            log.Fatal().Err(err).Str("metricsAddress", metrics).Msg("cannot create connection with the metrics collector")
        }
	metricsClient := grpc_monitoring_go.NewMetricsCollectorClient(metricsConn)

        organizationId := os.Getenv("ORGANIZATION_ID")
	clusterId := os.Getenv("CLUSTER_ID")
	if organizationId == "" || clusterId == "" {
            log.Fatal().Msg("ORGANIZATION_ID or CLUSTER_ID environment not set")
        }

        collector = statuscollector.NewMetricsAPICollector(metricsClient, organizationId, clusterId, sleepTime)
    }

    go collector.Run()

    scorer := scorer.NewSimpleScorer(collector)

    conf := &service.MusicianConfig{
        Port: port,
        Scorer: &scorer,
        Collector: &collector,
        Debug: debug,
    }

    musicianService, err := service.NewMusicianService(conf)

    if err!=nil{
        log.Fatal().AnErr("error",err).Msg("impossible to start service")
        os.Exit(1)
    }

    musicianService.Run()

}

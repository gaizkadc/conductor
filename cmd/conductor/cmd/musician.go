/*
 * Copyright 2018 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
    "github.com/rs/zerolog"
)

var musicianCmd = &cobra.Command{
    Use: "musician",
    Short: "Run a musician service",
    Long: "Run a musician service for the cluster this node belongs to",
    Run: func(cmd *cobra.Command, args [] string) {
        RunMusician()
    },
}


func init() {
    // UNIX Time is faster and smaller than most timestamps
    // If you set zerolog.TimeFieldFormat to an empty string,
    // logs will write with UNIX time
    zerolog.TimeFieldFormat = ""

    RootCmd.AddCommand(musicianCmd)

    musicianCmd.Flags().Uint32P("musician-port", "u",5001,"musician endpoint")
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
    musicianService,err := service.NewMusicianService(port, &collector, &scorer)

    if err!=nil{
        log.Fatal().AnErr("error",err).Msg("impossible to start service")
        os.Exit(1)
    }

    musicianService.Run()


}

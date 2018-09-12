//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package cmd

import (
    "github.com/spf13/cobra"
    "github.com/rs/zerolog/log"
    "github.com/rs/zerolog"
    "github.com/spf13/viper"
    "path/filepath"
    "strings"
    "os"
)

var RootCmd = &cobra.Command{
    Use: "musician",
    Short: "Musician data collector for clusters",
    Long: `Musicians collect information for cluster monitoring and interact with conductor to schedule deployments.`,
    TraverseChildren: true,
}

// Variables
// Path of the configuration file
var configFile string


func Execute() {
    if err := RootCmd.Execute(); err != nil {
        log.Error().Msg(err.Error())
    }
}


func initConfig() {
    // if --config is passed, attempt to parse the config file
    log.Info().Msg("Running init config")
    log.Info().Str("configfile",configFile).Msg("file")
    if configFile != "" {

        // get the filepath
        abs, err := filepath.Abs(configFile)
        if err != nil {
            log.Error().Msgf("Error reading filepath: ", err.Error())
        }

        // get the config name
        base := filepath.Base(abs)

        // get the path
        path := filepath.Dir(abs)

        //
        viper.SetConfigName(strings.Split(base, ".")[0])
        viper.AddConfigPath(path)

        viper.AutomaticEnv()

        // Find and read the config file; Handle errors reading the config file
        if err := viper.ReadInConfig(); err != nil {
            log.Fatal().Msgf("Failed to read config file: ", err.Error())
            os.Exit(1)
        }
    }
}

func init() {
    zerolog.TimeFieldFormat = ""
    cobra.OnInitialize(initConfig)
    // initialization file
    RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file path")

}

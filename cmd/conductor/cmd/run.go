/*
 * Copyright 2019 Nalej
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
 *
 */

package cmd

import (
	"fmt"
	"github.com/nalej/conductor/pkg/conductor/service"
	"github.com/nalej/conductor/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run conductor",
	Long:  "Run conductor service with... and with...",
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		RunConductor()
	},
}

func init() {

	RootCmd.AddCommand(runCmd)

	runCmd.Flags().Uint32P("port", "c", utils.CONDUCTOR_PORT, "port where conductor listens to")
	runCmd.Flags().StringP("systemModelAddress", "s", fmt.Sprintf("localhost:%d", utils.SYSTEM_MODEL_PORT),
		"host:port address for system model")
	runCmd.Flags().StringP("networkManagerAddress", "n", fmt.Sprintf("localhost:%d", utils.NETWORKING_SERVICE_PORT),
		"host:port address for networking manager")
	runCmd.Flags().StringP("authxAddress", "a", fmt.Sprintf("localhost:%d", utils.AUTHX_PORT),
		"host:port address for authx")
	runCmd.Flags().Uint32P("appClusterPort", "p", utils.APP_CLUSTER_API_PORT, "port where the application cluster api is listening")
	runCmd.Flags().Bool("useTLS", true, "Use TLS to connect to the application cluster API")
	runCmd.Flags().String("caCertPath", "", "Path for the CA certificate")
	runCmd.Flags().String("clientCertPath", "", "Path for the client certificate")
	runCmd.Flags().Bool("skipServerCertValidation", true, "Skip CA authentication validation")
	runCmd.Flags().StringP("unifiedLogging", "u", fmt.Sprintf("localhost:%d", utils.UNIFIED_LOGGING_PORT),
		"host:port address for unifiedLogging")
	runCmd.Flags().StringP("queueAddress", "q", fmt.Sprintf("localhost:%d", utils.QUEUE_PORT),
		"host:port address for the Nalej management queue")
	runCmd.Flags().StringP("dbFolder", "f", "/data/",
		"path for the folder used to store the local database")
	runCmd.Flags().StringP("networkMode","t", service.ConductorNetworkingModeZT,
		"Indicate the kind of networking solution conductor will work on top of (zt, istio)")

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
	// Client cert path
	var clientCertPath string
	// Skip CA validation
	var skipServerCertValidation bool
	// Authx url
	var authxService string
	// Unified Logging url
	var unifiedLoggingService string
	// Queue url
	var queueAddress string
	// Database folder path
	var dbFolder string
	// Networking mode
	var networkingMode string
	// Debug flag
	var debug bool

	port = uint32(viper.GetInt32("port"))
	systemModel = viper.GetString("systemModelAddress")
	networkingService = viper.GetString("networkManagerAddress")
	authxService = viper.GetString("authxAddress")
	appClusterApiPort = uint32(viper.GetInt32("appClusterPort"))
	useTLS = viper.GetBool("useTLS")
	caCertPath = viper.GetString("caCertPath")
	clientCertPath = viper.GetString("clientCertPath")
	skipServerCertValidation = viper.GetBool("skipServerCertValidation")
	unifiedLoggingService = viper.GetString("unifiedLogging")
	queueAddress = viper.GetString("queueAddress")
	dbFolder = viper.GetString("dbFolder")
	networkingMode = viper.GetString("networkMode")
	debug = viper.GetBool("debug")

	log.Info().Msg("launching conductor...")

	netMode, err := service.ConductorNetworkingModeFromString(networkingMode)
	if err != nil{
		log.Panic().Err(err).Msg("configuration error... leaving")
		return
	}

	config := service.ConductorConfig{
		Port:                     port,
		SystemModelURL:           systemModel,
		NetworkingServiceURL:     networkingService,
		AppClusterApiPort:        appClusterApiPort,
		UseTLSForClusterAPI:      useTLS,
		CACertPath:               caCertPath,
		ClientCertPath:           clientCertPath,
		SkipServerCertValidation: skipServerCertValidation,
		AuthxURL:                 authxService,
		UnifiedLoggingURL:        unifiedLoggingService,
		QueueURL:                 queueAddress,
		DBFolder:                 dbFolder,
		NetworkingMode:           netMode,
		Debug:                    debug,
	}
	config.Print()

	conductorService, err := service.NewConductorService(&config)
	if err != nil {
		log.Fatal().AnErr("err", err).Msg("impossible to initialize conductor service")
	}
	conductorService.Run()
}

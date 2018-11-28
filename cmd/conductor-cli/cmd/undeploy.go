/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package cmd

import (
	"context"
	pbConductor "github.com/nalej/grpc-conductor-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)


// Deployment Manager IP
var undeployConductorServer string

// Organization ID
var undeployOrgId string

// App Instance ID
var undeployAppId string

var undeployAppCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploy an application",
	Long:  `Undeploy an application`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		undeployApp()
	},
}

func init() {
	RootCmd.AddCommand(undeployAppCmd)
	undeployAppCmd.Flags().StringVar(&undeployConductorServer, "server", "localhost:5000", "Deployment Manager server URL")
	undeployAppCmd.Flags().StringVar(&undeployOrgId, "orgId", "", "Organization ID")
	undeployAppCmd.Flags().StringVar(&undeployAppId, "appId", "", "App Instance ID")
	undeployAppCmd.MarkFlagRequired("orgId")
	undeployAppCmd.MarkFlagRequired("appId")
}

func undeployApp() {

	conn, err := grpc.Dial(undeployConductorServer, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Err(err).Msgf("impossible to connect to server %s", undeployConductorServer)
	}

	client := pbConductor.NewConductorClient(conn)

	request := pbConductor.UndeployRequest{
		OrganizationId: undeployOrgId,
		AppInstanceId: undeployAppId,
	}

	_, err = client.Undeploy(context.Background(), &request)
	if err != nil {
		log.Error().Err(err).Msgf("error deleting app %s", undeployAppId)
		return
	}

	log.Info().Msg("Application successfully undeployed")
}
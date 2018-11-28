/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

// Service in charge of processing deployment gRPC requests.

package handler

import (
	"context"

	"github.com/nalej/conductor/internal/entities"
	pbCommon "github.com/nalej/grpc-common-go"
	pbConductor "github.com/nalej/grpc-conductor-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

//Handler struct with a related Manager
type Handler struct {
	c *Manager
}

//NewHandler creates a new Handler with its given Manager
func NewHandler(c *Manager) *Handler {
	return &Handler{c}
}

//Deploy validates and queues a deployment request
func (h *Handler) Deploy(ctx context.Context, request *pbConductor.DeploymentRequest) (*pbConductor.DeploymentResponse, error) {
	log.Debug().Interface("deploymentRequest", request).Msg("Deploy")
	if err := entities.ValidDeploymentRequest(request); err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	// Enqueue request for later processing
	log.Debug().Msgf("enqueue request %s", request.RequestId)
	instance, err := h.c.PushRequest(request)
	if err != nil {
		return nil, err
	}

	toReturn := pbConductor.DeploymentResponse{
		RequestId:     request.RequestId,
		AppInstanceId: instance.InstanceId,
		Status:        pbConductor.ApplicationStatus_QUEUED}
	log.Debug().Interface("deploymentResponse", toReturn).Msg("Response")
	return &toReturn, nil
}

func (h *Handler) Undeploy(ctx context.Context, request *pbConductor.UndeployRequest) (*pbCommon.Success, error) {
	log.Debug().Interface("undeployRequest", request).Msg("Undeploy")
	if err := entities.ValidUndeployRequest(request); err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	toUndeploy := entities.UndeployRequest{
		OrganizationId: request.OrganizationId,
		AppInstanceId: request.AppInstanceId,
	}
	err := h.c.Undeploy(&toUndeploy)
	if err != nil {
		log.Error().Msgf("Unable to undeploy application %s", request.AppInstanceId)
		return nil, err
	}
	log.Debug().Msgf("Application %s undeployed", request.AppInstanceId)
	return &pbCommon.Success{}, nil
}

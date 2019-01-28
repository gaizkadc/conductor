/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

// The baton monitor collects information from deployment fragments and updates the status of services.

package monitor

import (
	"context"

	"github.com/nalej/conductor/internal/entities"
	"github.com/nalej/grpc-utils/pkg/conversions"

	pbCommon "github.com/nalej/grpc-common-go"
	pbConductor "github.com/nalej/grpc-conductor-go"
	"github.com/rs/zerolog/log"
)

//Handler struct with a related Manager
type Handler struct {
	mng *Manager
}

//NewHandler creates a new Handler with its given Manager
func NewHandler(m *Manager) *Handler {
	return &Handler{mng: m}
}

//UpdateDeploymentFragmentStatus validates the requests and updates the fragment status
func (h *Handler) UpdateDeploymentFragmentStatus(ctx context.Context, request *pbConductor.DeploymentFragmentUpdateRequest) (*pbCommon.Success, error) {
	if err := entities.ValidDeploymentFragmentUpdateRequest(request); err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	err := h.mng.UpdateFragmentStatus(request)
	if err != nil {
		return nil, err
	}
	return &pbCommon.Success{}, nil
}

//UpdateServiceStatus validates the request and updates the service status
func (h *Handler) UpdateServiceStatus(ctx context.Context, request *pbConductor.DeploymentServiceUpdateRequest) (*pbCommon.Success, error) {
	if err := entities.ValidDeploymentFragmentID(request.FragmentId); err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	err := h.mng.UpdateServicesStatus(request)
	if err != nil {
		log.Error().Err(err).Msgf("error when updating service status in system model")
		return nil, err
	}

	return &pbCommon.Success{}, nil
}

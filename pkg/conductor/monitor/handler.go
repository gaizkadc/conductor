/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

// The handler monitor collects information from deployment fragments and updates the status of services.

package monitor

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbCommon "github.com/nalej/grpc-common-go"
    "github.com/rs/zerolog/log"
    "errors"
    "context"
)

type Handler struct {
    mng * Manager
}

func NewHandler(m * Manager) *Handler {
    return &Handler{mng: m}
}

func (h *Handler) UpdateDeploymentFragmentStatus(ctx context.Context, request *pbConductor.DeploymentFragmentUpdateRequest) (*pbCommon.Success, error) {
    if request.FragmentId == "" {
        err := errors.New("non valid empty fragment id in DeploymentFragmentUpdateRequest")
        return nil, err
    }

    if !h.ValidDeploymentFragmentUpdateRequest(request) {
        return nil, errors.New("missing mandatory fields")
    }

    err := h.mng.UpdateFragmentStatus(request)
    if err != nil {
        return nil, err
    }
    return &pbCommon.Success{}, nil
}


func (h *Handler) UpdateServiceStatus(ctx context.Context, request *pbConductor.DeploymentServiceUpdateRequest) (*pbCommon.Success, error) {
    if request.FragmentId == "" {
        err := errors.New("non valid empty fragment id in DeploymentServiceUpdateRequest")
        return nil, err
    }



    err := h.mng.UpdateServicesStatus(request)

    if err != nil {
        log.Error().Err(err).Msgf("error when updating service status in system model")
        return nil, err
    }

    return &pbCommon.Success{}, nil
}


func (h *Handler) ValidDeploymentFragmentUpdateRequest(request *pbConductor.DeploymentFragmentUpdateRequest) bool {
    if request.OrganizationId == "" || request.FragmentId == "" || request.AppInstanceId == "" ||
        request.ClusterId == "" {
        return false
    }
    return true
}

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

    // TODO finish this
    log.Debug().Msgf("UpdateServiceStatus receives %v", request)
    for _, serv := range request.List {
        log.Debug().Msgf("--> %s-%s", serv.ServiceInstanceId,serv.Status)
    }

    err := h.mng.UpdateServicesStatus(request)

    if err != nil {
        log.Error().Err(err).Msgf("error when updating service status in system model")
        return nil, err
    }

    return &pbCommon.Success{}, nil
}




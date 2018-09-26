/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 */

 // Service in charge of processing deployment gRPC requests.

package handler

import (
    "context"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "errors"
    "github.com/rs/zerolog/log"
)

type Handler struct{
    c *Manager
}

func NewHandler(c *Manager) *Handler {
    return &Handler{c}
}


func (h *Handler) Deploy(ctx context.Context, request *pbConductor.DeploymentRequest) (*pbConductor.DeploymentResponse, error) {
    log.Debug().Interface("deploymentRequest", request).Msg("Deploy")
    if request == nil {
        return nil, errors.New("invalid request")
    }

    // Enqueue request for later processing
    err := h.c.PushRequest(request)
    if err != nil {
        return nil, err
    }

    // TODO
    // Modify system model accordingly

    toReturn := pbConductor.DeploymentResponse{RequestId: request.RequestId, Status: pbConductor.ApplicationStatus_QUEUED}
    log.Debug().Interface("deploymentResponse", toReturn).Msg("Response")
    return &toReturn, nil
}
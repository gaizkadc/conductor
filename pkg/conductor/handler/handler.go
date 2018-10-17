/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


// Service in charge of processing deployment gRPC requests.

package handler

import (
    "context"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbCommon "github.com/nalej/grpc-common-go"
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

    if !h.ValidDeploymentRequest(request) {
        return nil, errors.New("mandatory parameters in this request are missing")
    }

    // Enqueue request for later processing
    log.Debug().Msgf("enqueue request %s",request.RequestId)
    instance, err := h.c.PushRequest(request)
    if err != nil {
        return nil, err
    }

    toReturn := pbConductor.DeploymentResponse{
        RequestId: request.RequestId,
        AppInstanceId: instance.InstanceID,
        Status: pbConductor.ApplicationStatus_QUEUED}
    log.Debug().Interface("deploymentResponse", toReturn).Msg("Response")
    return &toReturn, nil
}

func (h *Handler) Undeploy(ctx context.Context, request *pbConductor.UndeployRequest) (*pbCommon.Success, error) {
    panic("undeploy operation is not implemented yet")
}

func (h *Handler) ValidDeploymentRequest(request *pbConductor.DeploymentRequest) bool {
    if request.RequestId == "" || request.Name == ""  || request.AppId == nil || request.AppId.OrganizationId == "" ||
        request.AppId.AppDescriptorId == "" {
            return false
    }
    return true
}
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
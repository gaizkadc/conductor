//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//
// Service in charge of processing deployment gRPC requests.

package deployment

import (
    "context"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/rs/zerolog/log"
    "errors"
    "github.com/phf/go-queue/queue"
)

type Handler struct{
    queue *queue.Queue
}

func NewHandler(q *queue.Queue) *Handler {
    return &Handler{q}
}


func (h *Handler) Deploy(ctx context.Context, request *pbConductor.DeploymentRequest) (*pbConductor.DeploymentResponse, error) {
    if request == nil {
        return nil, errors.New("invalid request")
    }
    log.Debug().Interface("request", request).Msg("received deployment request")
    h.queue.PushBack(request)

    response := pbConductor.DeploymentResponse{"this is your answer"}
    return &response, nil
}
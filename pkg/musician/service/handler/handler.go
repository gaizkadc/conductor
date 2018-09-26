/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 */

package handler

import (
    "context"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/rs/zerolog/log"
    "errors"
)

type Handler struct{
    m *Manager
}

func NewHandler(m *Manager) *Handler {
    return &Handler{m}
}

func (h *Handler) Score(ctx context.Context, request *pbConductor.ClusterScoreRequest) (*pbConductor.ClusterScoreResponse, error) {
    log.Debug().Msg("musician score was requested")
    if request == nil {
        return nil, errors.New("empty request")
    }
    return h.m.Score(request)
}

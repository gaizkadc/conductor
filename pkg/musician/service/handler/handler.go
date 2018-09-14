//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package handler

import (
    "context"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "errors"
)

type Handler struct{
    m *Manager
}

func NewHandler(m *Manager) *Handler {
    return &Handler{m}
}

func (h *Handler) Score(ctx context.Context, request *pbConductor.ClusterScoreRequest) (*pbConductor.ClusterScoreResponse, error) {
    if request == nil {
        return nil, errors.New("empty request")
    }
    return h.m.Score(request)
}

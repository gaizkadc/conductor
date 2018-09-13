//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package musician

import (
    conductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/pkg/musician/scorer"
    "github.com/rs/zerolog/log"
    "context"
    "errors"
)

type Handler struct{
     scorer *scorer.Scorer
}

func NewHandler(s *scorer.Scorer) *Handler {
    return &Handler{s}
}


func (h *Handler) Score(ctx context.Context, request *conductor.ClusterScoreRequest) (*conductor.ClusterScoreResponse, error) {
    if request == nil {
        return nil, errors.New("invalid request")
    }

    response, err := h.scorer.Score(request)
    if err != nil {
        return nil, err
    }

    log.Debug().Interface("response", response).Msg("response to return...")

    return &response, nil
}
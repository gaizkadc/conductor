/*
 * Copyright (C) 2018 Nalej Group -All Rights Reserved
 */


package handler

import (
    "github.com/nalej/conductor/pkg/musician/scorer"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
)

type Manager struct {
    Status *statuscollector.StatusCollector
    ScorerMethod scorer.Scorer
}

func NewManager(collector *statuscollector.StatusCollector, serv scorer.Scorer) *Manager {
    return &Manager{collector, serv}
}


func (m *Manager) Score(request *pbConductor.ClusterScoreRequest) (*pbConductor.ClusterScoreResponse, error) {
    return m.ScorerMethod.Score(request)
}

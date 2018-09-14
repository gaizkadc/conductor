//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package handler

import (
    "github.com/phf/go-queue/queue"
    "github.com/rs/zerolog/log"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/pkg/conductor/scorer"

    "github.com/nalej/conductor/internal/entities"
)

type Manager struct {
    // Queue for incoming messages
    Queue *queue.Queue
    // ScorerMethod
    ScorerMethod scorer.Scorer
    // List of musicians to be queried
    Musicians []string
}

func NewManager(queue *queue.Queue, scorer scorer.Scorer, port uint32) *Manager {
    // instantiate a server
    return &Manager{queue, scorer, make([]string,0)}
}



func(c *Manager) ProcessDeploymentRequest(request *pbConductor.DeploymentRequest) (*pbConductor.DeploymentResponse, error) {
    // Empty queue process it.
    if c.Queue.Len() == 0 {
        log.Debug().Str("request_id",request.RequestId).Msg("It's time to process request")

        req:= entities.Requirements{Disk:0.1,CPU:0.2,Memory:0.3}

        returned,_ :=c.ScorerMethod.ScoreRequirements (&req, c.Musicians)
        log.Debug().Msgf("Returned %v",returned)

    } else {
        log.Debug().Str("request_id", request.RequestId).Msg("deployment request send to the queue")
        c.Queue.PushBack(request)
    }
    response := pbConductor.DeploymentResponse{RequestId: "this is a response"}
    return &response, nil
}









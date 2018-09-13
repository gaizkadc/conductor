//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package conductor

import (
    "github.com/phf/go-queue/queue"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "github.com/rs/zerolog/log"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/conductor/pkg/conductor/handler"
)

type Conductor struct {
    // Queue for incoming messages
    Queue *queue.Queue
    // Scorer
    Scorer *scorer.Scorer
    // grpc server
    server *ConductorServer
}

func NewConductor(queue *queue.Queue, scorer *scorer.Scorer, port uint32) *Conductor {
    // instantiate a server
    server := NewConductorServer(port)
    return &Conductor{queue, scorer, server}
}

func(c *Conductor) Run() {
    // register services
    deployment := handler.NewHandler(c)

    // Server and registry
    grpcServer := grpc.NewServer()
    pbConductor.RegisterConductorServer(grpcServer,deployment)

    // Register reflection service on gRPC server.
    reflection.Register(grpcServer)

    c.server.Run()
}

func(c *Conductor) ProcessDeploymentRequest(request *pbConductor.DeploymentRequest) (*pbConductor.DeploymentResponse, error) {
    // Empty queue process it.
    if h.queue.Len() == 0 {
        log.Info().Str("request_id",request.RequestId).Msg("It's time to process request")
    } else {
        log.Debug().Str("request_id", request.RequestId).Msg("deployment request send to the queue")
        h.queue.PushBack(request)
    }

    response := pbConductor.DeploymentResponse{RequestId: "this is a response"}
    return &response, nil
}







//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package service

import (
    "github.com/nalej/conductor/pkg/conductor/service/handler"
    "github.com/phf/go-queue/queue"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/conductor/tools"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "google.golang.org/grpc/reflection"
)

type ConductorService struct {
    conductor *handler.Manager
    server *tools.GenericGRPCServer
}

func NewConductorService(port uint32, q *queue.Queue, s scorer.Scorer) (*ConductorService, error) {
    c := handler.NewManager(q, s, port)
    conductorServer := tools.NewGenericGRPCServer(port)

    instance := ConductorService{c, conductorServer}
    return &instance, nil
}

func(c *ConductorService) Run() {
    // register services
    deployment := handler.NewHandler(c.conductor)

    // Server and registry
    //grpcServer := grpc.NewServer()
    pbConductor.RegisterConductorServer(c.server.Server,deployment)

    // Register reflection service on gRPC server.
    reflection.Register(c.server.Server)

    c.server.Run()
}

// Set the musicians to be queried
// TODO: this has to be removed and check the system model instead. This is only for initial testing.
func(c *ConductorService) SetMusicians(musicians []string) {
    c.conductor.Musicians=musicians
}
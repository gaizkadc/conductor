/*
 * Copyright (C) 2018 Nalej Group -All Rights Reserved
 */

package service

import (
    "github.com/nalej/conductor/pkg/musician/service/handler"
    "github.com/nalej/conductor/pkg/musician/scorer"
    "google.golang.org/grpc/reflection"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/grpc-utils/pkg/tools"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
)

type MusicianService struct {
    musician *handler.Manager
    server *tools.GenericGRPCServer
}

func NewMusicianService(port uint32, collector *statuscollector.StatusCollector, scor *scorer.Scorer) (*MusicianService, error) {
    c := handler.NewManager(collector, *scor)
    musicianServer := tools.NewGenericGRPCServer(port)
    instance := MusicianService{c, musicianServer}
    return &instance, nil
}

func(c *MusicianService) Run() {
    // register services
    deployment := handler.NewHandler(c.musician)

    // Server and registry
    pbConductor.RegisterMusicianServer(c.server.Server,deployment)

    // Register reflection service on gRPC server.
    reflection.Register(c.server.Server)

    c.server.Run()
}

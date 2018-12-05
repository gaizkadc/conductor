/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package service

import (
    "github.com/nalej/conductor/pkg/musician/service/handler"
    "github.com/nalej/conductor/pkg/musician/scorer"
    "google.golang.org/grpc/reflection"
    pbConductor "github.com/nalej/grpc-conductor-go"

    "github.com/nalej/grpc-utils/pkg/tools"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
    "os"
    "github.com/nalej/conductor/pkg/utils"
    "github.com/rs/zerolog/log"
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

    if os.Getenv(utils.MUSICIAN_CLUSTER_ID)==""{
        log.Panic().Msgf("%s variable has to be set before running the musician service", utils.MUSICIAN_CLUSTER_ID)
    }

    // register services
    deployment := handler.NewHandler(c.musician)

    // Server and registry
    pbConductor.RegisterMusicianServer(c.server.Server,deployment)

    // Register reflection service on gRPC server.
    reflection.Register(c.server.Server)

    c.server.Run()
}

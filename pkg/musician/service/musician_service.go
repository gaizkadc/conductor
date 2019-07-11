/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package service

import (
    "github.com/nalej/conductor/pkg/musician/service/handler"
    "github.com/nalej/conductor/pkg/musician/scorer"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
    "os"
    "github.com/nalej/conductor/pkg/utils"
    "github.com/rs/zerolog/log"
    "fmt"
    "net"
)

type MusicianConfig struct {
    // Musician server port
    Port uint32
    // Status collector
    Collector *statuscollector.StatusCollector
    // Scorer
    Scorer *scorer.Scorer
    // Debug enabled
    Debug bool
}

type MusicianService struct {
    musician *handler.Manager
    configuration *MusicianConfig
    server *grpc.Server
}

//func NewMusicianService(port uint32, collector *statuscollector.StatusCollector, scor *scorer.Scorer) (*MusicianService, error) {
func NewMusicianService(config *MusicianConfig) (*MusicianService, error) {
    musicianServer := grpc.NewServer()
    c := handler.NewManager(config.Collector, *config.Scorer)
    instance := MusicianService{c, config,musicianServer}
    return &instance, nil
}

func(m *MusicianService) Run() {

    if os.Getenv(utils.MUSICIAN_CLUSTER_ID)==""{
        log.Panic().Msgf("%s variable has to be set before running the musician service", utils.MUSICIAN_CLUSTER_ID)
    }

    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", m.configuration.Port))
    if err != nil {
        log.Fatal().Errs("failed to listen: %v", []error{err})
    }

    // register services
    deployment := handler.NewHandler(m.musician)

    // Server and registry
    pbConductor.RegisterMusicianServer(m.server,deployment)

    // Register reflection service on gRPC server.
    if m.configuration.Debug {
        reflection.Register(m.server)
    }


    // Run
    log.Info().Uint32("port", m.configuration.Port).Msg("Launching gRPC server")
    if err := m.server.Serve(lis); err != nil {
        log.Fatal().Errs("failed to serve: %v", []error{err})
    }

}

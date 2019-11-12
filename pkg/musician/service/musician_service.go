/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package service

import (
	"fmt"
	"github.com/nalej/conductor/pkg/musician/scorer"
	"github.com/nalej/conductor/pkg/musician/service/handler"
	"github.com/nalej/conductor/pkg/musician/statuscollector"
	"github.com/nalej/conductor/pkg/utils"
	pbConductor "github.com/nalej/grpc-conductor-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
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
	musician      *handler.Manager
	configuration *MusicianConfig
	server        *grpc.Server
}

//func NewMusicianService(port uint32, collector *statuscollector.StatusCollector, scor *scorer.Scorer) (*MusicianService, error) {
func NewMusicianService(config *MusicianConfig) (*MusicianService, error) {
	musicianServer := grpc.NewServer()
	c := handler.NewManager(config.Collector, *config.Scorer)
	instance := MusicianService{c, config, musicianServer}
	return &instance, nil
}

func (m *MusicianService) Run() {

	if os.Getenv(utils.MUSICIAN_CLUSTER_ID) == "" {
		log.Panic().Msgf("%s variable has to be set before running the musician service", utils.MUSICIAN_CLUSTER_ID)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", m.configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	// register services
	deployment := handler.NewHandler(m.musician)

	// Server and registry
	pbConductor.RegisterMusicianServer(m.server, deployment)

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

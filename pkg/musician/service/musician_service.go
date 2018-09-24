/*
 * Copyright 2018 Nalej
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
 */

package service

import (
    "github.com/nalej/conductor/pkg/musician/service/handler"
    "github.com/nalej/conductor/pkg/musician/scorer"
    "google.golang.org/grpc/reflection"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/tools"
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

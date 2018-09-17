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
    "github.com/nalej/conductor/pkg/conductor/handler"
    "github.com/phf/go-queue/queue"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/conductor/tools"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "google.golang.org/grpc/reflection"
    "github.com/rs/zerolog/log"
    "github.com/nalej/conductor/pkg/conductor"
)




type ConductorService struct {
    // Conductor manager
    conductor *handler.Manager
    // Server for incoming requests
    server *tools.GenericGRPCServer
    // List of musicians to be queried
    musicians []string
    // Connections with musicians
    connections *tools.ConnectionsMap
}

func NewConductorService(port uint32, q *queue.Queue, s scorer.Scorer) (*ConductorService, error) {
    c := handler.NewManager(q, s, port)
    conductorServer := tools.NewGenericGRPCServer(port)
    // instance := ConductorService{c, conductorServer, make([]string, 0)},connections)
    instance := ConductorService{conductor: c,
                                server: conductorServer,
                                musicians: make([]string,0),
                                connections: conductor.GetMusicianClients()}
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
    c.musicians=musicians
    for _, target := range musicians {
        _,err := c.connections.AddConnection(target)
        if err != nil {
            log.Error().Err(err)
        } else {
            log.Info().Str("address",target).Msg("musician address correctly added")
            c.musicians = append(c.musicians, target)
        }
    }
}
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

package scorer

import (
    . "github.com/onsi/ginkgo"
    "google.golang.org/grpc"
    "google.golang.org/grpc/test/bufconn"
    "github.com/nalej/conductor/tools"
    musicianHandler "github.com/nalej/conductor/pkg/musician/service/handler"
    musicianScorer "github.com/nalej/conductor/pkg/musician/scorer"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/nalej/conductor/internal/entities"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "fmt"
)



var _ = Describe("Simple scorer basic functionality with one musicianHandler", func() {
    // grpc server
    var server *grpc.Server
    // grpc test listener
    var listener *bufconn.Listener
    // scorer
    var scorerMethod Scorer
    // Musician handler
    var manager *musicianHandler.Manager
    // Musician clients
    var clients *tools.ConnectionsMap


    BeforeEach(func() {
        // instantiate musicianHandler server
        scorerMethod = NewSimpleScorer()
        collector := statuscollector.NewFakeCollector()
        manager = musicianHandler.NewManager(&collector, musicianScorer.NewSimpleScorer())

        listener = tools.GetDefaultListener()
        server = grpc.NewServer()

        tools.LaunchServer(server, listener)
        // Add the client
        pbConductor.RegisterMusicianServer(server, musicianHandler.NewHandler(manager))

        clients = conductor.GetMusicianClients()

        clients.AddConnection(listener.Addr().String())

    })

    Context("new score request arrives", func(){
        var requirements entities.Requirements

        BeforeEach(func(){

            requirements = entities.Requirements{Memory:0.1, Disk: 0.2, CPU: 0.3}
        })

        It("send a score request", func(){
            response, err := scorerMethod.ScoreRequirements(&requirements)
            fmt.Printf("response ---> %v\n",response)
            fmt.Printf("error ---> %v\n",err)
        })
    })

})

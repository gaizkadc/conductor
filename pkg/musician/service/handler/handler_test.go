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

package handler

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "google.golang.org/grpc/test/bufconn"
    "google.golang.org/grpc"
    "context"

    "github.com/nalej/conductor/tools"
    "github.com/nalej/conductor/pkg/musician/scorer"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
)


const TestPort=4000

var _ = Describe("Deployment server API", func() {
    // grpc server
    var server *grpc.Server
    // conductor object
    var mgr *Manager
    // grpc test listener
    var listener *bufconn.Listener

    BeforeEach(func(){
        collector := statuscollector.NewFakeCollector()
        listener = tools.GetDefaultListener()
        server = grpc.NewServer()
        scorerMethod := scorer.NewSimpleScorer(collector)
        mgr = NewManager(&collector, scorerMethod)
        tools.LaunchServer(server,listener)
    })

    Context("A new score requests arrives", func() {
        var request pbConductor.ClusterScoreRequest
        var response pbConductor.ClusterScoreResponse
        var client pbConductor.MusicianClient

        BeforeEach(func() {
            // Register the service.
            pbConductor.RegisterMusicianServer(server, NewHandler(mgr))

            request = pbConductor.ClusterScoreRequest{RequestId: "myrequestId"}
            response = pbConductor.ClusterScoreResponse{RequestId: "myrequestId", Score: 0.1}

            conn, err := tools.GetConn(*listener)
            Expect(err).ShouldNot(HaveOccurred())
            client = pbConductor.NewMusicianClient(conn)
        })

        It("receive an expected message", func() {
            resp, err := client.Score(context.Background(), &request)

            Expect(resp.RequestId).To(Equal(response.RequestId))
            Expect(err).ShouldNot(HaveOccurred())
        })
    })
})

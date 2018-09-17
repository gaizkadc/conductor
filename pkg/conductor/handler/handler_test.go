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
    "fmt"
    "github.com/nalej/conductor/tools"
    "github.com/phf/go-queue/queue"
    "github.com/nalej/conductor/pkg/conductor/scorer"
)


const TestPort=4000

var _ = Describe("Deployment server API", func() {
    // grpc server
    var server *grpc.Server
    // conductor object
    var cond *Manager
    // grpc test listener
    var listener *bufconn.Listener
    // Queue
    var q *queue.Queue

    BeforeEach(func(){
        listener = tools.GetDefaultListener()
        server = grpc.NewServer()
        scorerMethod := scorer.NewSimpleScorer()
        q = queue.New()
        cond = NewManager(q, scorerMethod, TestPort)
        tools.LaunchServer(server,listener)
    })

    Context("The queue is empty and a new request arrives", func() {
        var request pbConductor.DeploymentRequest
        var response pbConductor.DeploymentResponse
        var client pbConductor.ConductorClient

        BeforeEach(func() {
            // Register the service.
            q.Init()
            pbConductor.RegisterConductorServer(server, NewHandler(cond))

            request = pbConductor.DeploymentRequest{RequestId: "myrequestId"}
            response = pbConductor.DeploymentResponse{RequestId: "this is a response"}

            conn, err := tools.GetConn(*listener)
            Expect(err).ShouldNot(HaveOccurred())
            client = pbConductor.NewConductorClient(conn)
        })

        It("receive an expected message", func() {
            resp, err := client.Deploy(context.Background(), &request)

            Expect(resp.String()).To(Equal(response.String()))
            Expect(err).ShouldNot(HaveOccurred())
        })

        It("the queue was empty, the request is sent to process", func() {
            Expect(q.Len()).To(Equal(0))
            _, err := client.Deploy(context.Background(), &request)
            Expect(q.Len()).To(Equal(0))
            Expect(err).ShouldNot(HaveOccurred())
        })
    })

    Context("there was something in the queue and a new request arrives", func() {
        var request pbConductor.DeploymentRequest
        var client pbConductor.ConductorClient

        BeforeEach(func() {
            // Register the service.
            q.Init()
            // put something in the queue
            q.PushBack(pbConductor.DeploymentRequest{RequestId: "this was enqueued"})
            pbConductor.RegisterConductorServer(server, NewHandler(cond))

            request = pbConductor.DeploymentRequest{RequestId: "myrequestId2"}

            conn, err := tools.GetConn(*listener)
            Expect(err).ShouldNot(HaveOccurred())
            client = pbConductor.NewConductorClient(conn)
        })

        It("the new request is enqueued and the very first is processed", func() {
            Expect(q.Len()).To(Equal(1))
            _, err := client.Deploy(context.Background(), &request)
            Expect(err).ShouldNot(HaveOccurred())
            Expect(q.Len()).To(Equal(1))
            fmt.Printf("====%v\n",q.Front().(*pbConductor.DeploymentRequest).String())
            Expect(request.String()).To(Equal(q.PopFront().(*pbConductor.DeploymentRequest).String()))
        })
    })
})


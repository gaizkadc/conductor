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
    // queue
    var q *queue.Queue
    // Conductor client
    var client pbConductor.ConductorClient

    BeforeSuite(func(){
        listener = tools.GetDefaultListener()
        server = grpc.NewServer()
        scorerMethod := scorer.NewSimpleScorer()
        q = queue.New()
        cond = NewManager(q, scorerMethod, TestPort)
        tools.LaunchServer(server,listener)

        // Register the service.
        pbConductor.RegisterConductorServer(server, NewHandler(cond))

        conn, err := tools.GetConn(*listener)
        Expect(err).ShouldNot(HaveOccurred())
        client = pbConductor.NewConductorClient(conn)
    })

    AfterSuite(func(){
        server.Stop()
        listener.Close()
    })


    Context("The queue is empty and a new request arrives", func() {
        var request pbConductor.DeploymentRequest
        var response pbConductor.DeploymentResponse


        BeforeEach(func() {
            request = pbConductor.DeploymentRequest{RequestId: "myrequestId"}
            response = pbConductor.DeploymentResponse{RequestId: "myrequestId"}
            q.Init()
        })


        It("receive an expected message", func() {
            resp, err := client.Deploy(context.Background(), &request)

            Expect(resp.String()).To(Equal(response.String()))
            Expect(err).ShouldNot(HaveOccurred())
        })

    })

    Context("there was something in the queue and a new request arrives", func() {
        var request pbConductor.DeploymentRequest
        var toEnqueue pbConductor.DeploymentRequest

        BeforeEach(func() {
            q.Init()
            request = pbConductor.DeploymentRequest{RequestId: "myrequestId"}
            // put something in the queue
            toEnqueue = pbConductor.DeploymentRequest{RequestId: "this was enqueued"}
            cond.queue.PushBack(&toEnqueue)
        })

        It("the new request is enqueued and the very first is processed", func() {
            Expect(cond.queue.Len()).To(Equal(1))
            Expect(toEnqueue.String()).To(Equal(cond.queue.Front().(*pbConductor.DeploymentRequest).String()))
            _, err := client.Deploy(context.Background(), &request)
            Expect(err).ShouldNot(HaveOccurred())
            Expect(cond.queue.Len()).To(Equal(2))
            back:=cond.queue.Back()
            Expect(back.(*pbConductor.DeploymentRequest).String()).To(Equal(request.String()))
        })
    })
})


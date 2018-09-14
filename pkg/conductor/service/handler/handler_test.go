//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

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
    // Queue
    var q *queue.Queue

    BeforeEach(func(){
        listener = bufconn.Listen(tools.BufSize)
        server = grpc.NewServer()
        scorerMethod := scorer.NewSimpleScorer()
        q = queue.New()
        cond = NewManager(q, &scorerMethod, TestPort)
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

        It("the new request is enqueued", func() {
            Expect(q.Len()).To(Equal(1))
            _, err := client.Deploy(context.Background(), &request)
            Expect(err).ShouldNot(HaveOccurred())
            Expect(q.Len()).To(Equal(2))
        })
    })
})


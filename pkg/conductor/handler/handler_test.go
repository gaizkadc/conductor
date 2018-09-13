//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package handler

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    conductor "github.com/nalej/grpc-conductor-go"
    "google.golang.org/grpc/test/bufconn"
    "google.golang.org/grpc"
    "context"

    "github.com/nalej/conductor/tools"
    "github.com/phf/go-queue/queue"
)

const bufSize = 1024*1024
var listener *bufconn.Listener

var _ = Describe("Deployment server API", func() {
    // grpc server
    var server *grpc.Server


    BeforeEach(func(){
        listener = bufconn.Listen(bufSize)
        server = grpc.NewServer()
        tools.LaunchServer(server,listener)
    })

    Context("A new deployment is requested", func(){
        var request conductor.DeploymentRequest
        var response conductor.DeploymentResponse
        var client conductor.DeploymentClient
        var q *queue.Queue

        BeforeEach(func(){
            // Register the service.
            q = queue.New()
            conductor.RegisterDeploymentServer(server,NewHandler(q))

            request = conductor.DeploymentRequest{RequestId: "myrequestId"}
            response = conductor.DeploymentResponse{RequestId: "this is a response"}

            conn, err := tools.GetConn(*listener)
            Expect(err).ShouldNot(HaveOccurred())
            client = conductor.NewDeploymentClient(conn)
        })

        It("receive an expected message", func() {
            resp, err := client.Deploy(context.Background(), &request)

            Expect(resp.String()).To(Equal(response.String()))
            Expect(err).ShouldNot(HaveOccurred())
        })
        It("increases the size of the queue", func(){
            Expect(q.Len()).To(Equal(0))
            _, err := client.Deploy(context.Background(), &request)
            Expect(q.Len()).To(Equal(1))
            Expect(err).ShouldNot(HaveOccurred())
        })
        It("we can pop the request from the queue", func(){
            _, err := client.Deploy(context.Background(), &request)
            Expect(err).ShouldNot(HaveOccurred())
            pop := q.PopFront().(*conductor.DeploymentRequest)
            Expect(request.String()).To(Equal(pop.String()))
        })
    })
})


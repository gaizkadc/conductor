/*
 * Copyright (C) 2018 Nalej Group -All Rights Reserved
 */


package handler

import (
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "google.golang.org/grpc/test/bufconn"
    "google.golang.org/grpc"
    "context"
    "github.com/nalej/grpc-utils/pkg/test"
    "github.com/phf/go-queue/queue"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/conductor/requirementscollector"
)


const TestPort=4000

var _ = ginkgo.Describe("Deployment server API", func() {
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

    ginkgo.BeforeSuite(func(){
        listener = test.GetDefaultListener()
        server = grpc.NewServer()
        scorerMethod := scorer.NewSimpleScorer()
        designer := plandesigner.NewSimplePlanDesigner()
        reqcoll := requirementscollector.NewSimpleRequirementsCollector()
        q = queue.New()
        cond = NewManager(q, scorerMethod, reqcoll, designer,TestPort)
        test.LaunchServer(server,listener)

        // Register the service.
        pbConductor.RegisterConductorServer(server, NewHandler(cond))

        conn, err := test.GetConn(*listener)
        gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
        client = pbConductor.NewConductorClient(conn)
    })

    ginkgo.AfterSuite(func(){
        server.Stop()
        listener.Close()
    })


    ginkgo.Context("The queue is empty and a new request arrives", func() {
        var request pbConductor.DeploymentRequest
        var response pbConductor.DeploymentResponse


        ginkgo.BeforeEach(func() {
            request = pbConductor.DeploymentRequest{RequestId: "myrequestId"}
            response = pbConductor.DeploymentResponse{RequestId: "myrequestId"}
            q.Init()
        })


        ginkgo.It("receive an expected message", func() {
            resp, err := client.Deploy(context.Background(), &request)

            gomega.Expect(resp.String()).To(gomega.Equal(response.String()))
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
        })

    })

    ginkgo.Context("there was something in the queue and a new request arrives", func() {
        var request pbConductor.DeploymentRequest
        var toEnqueue pbConductor.DeploymentRequest

        ginkgo.BeforeEach(func() {
            q.Init()
            request = pbConductor.DeploymentRequest{RequestId: "myrequestId"}
            // put something in the queue
            toEnqueue = pbConductor.DeploymentRequest{RequestId: "this was enqueued"}
            cond.queue.PushBack(&toEnqueue)
        })

        ginkgo.It("the new request is enqueued and the very first is processed", func() {
            gomega.Expect(cond.queue.Len()).To(gomega.Equal(1))
            gomega.Expect(toEnqueue.String()).To(gomega.Equal(cond.queue.Front().(*pbConductor.DeploymentRequest).String()))
            _, err := client.Deploy(context.Background(), &request)
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
            gomega.Expect(cond.queue.Len()).To(gomega.Equal(2))
            back:=cond.queue.Back()
            gomega.Expect(back.(*pbConductor.DeploymentRequest).String()).To(gomega.Equal(request.String()))
        })
    })
})


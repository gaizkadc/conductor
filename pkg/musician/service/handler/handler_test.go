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
    "github.com/nalej/conductor/pkg/musician/scorer"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
)



var _ = ginkgo.Describe("Deployment server API", func() {
    // grpc server
    var server *grpc.Server
    // conductor object
    var mgr *Manager
    // grpc test listener
    var listener *bufconn.Listener

    ginkgo.BeforeEach(func(){
        collector := statuscollector.NewFakeCollector()
        listener = test.GetDefaultListener()
        server = grpc.NewServer()
        scorerMethod := scorer.NewSimpleScorer(collector)
        mgr = NewManager(&collector, scorerMethod)
        test.LaunchServer(server,listener)
    })

    ginkgo.Context("A new score requests arrives", func() {
        var request pbConductor.ClusterScoreRequest
        var response pbConductor.ClusterScoreResponse
        var client pbConductor.MusicianClient

        ginkgo.BeforeEach(func() {
            // Register the service.
            pbConductor.RegisterMusicianServer(server, NewHandler(mgr))

            request = pbConductor.ClusterScoreRequest{RequestId: "myrequestId"}
            response = pbConductor.ClusterScoreResponse{RequestId: "myrequestId", Score: 0.1}

            conn, err := test.GetConn(*listener)
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
            client = pbConductor.NewMusicianClient(conn)
        })

        ginkgo.It("receive an expected message", func() {
            resp, err := client.Score(context.Background(), &request)

            gomega.Expect(resp.RequestId).To(gomega.Equal(response.RequestId))
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
        })
    })
})

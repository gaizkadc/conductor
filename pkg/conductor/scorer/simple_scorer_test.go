/*
 * Copyright (C) 2018 Nalej Group -All Rights Reserved
 */


package scorer

import (
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    //"google.golang.org/grpc"
    //"google.golang.org/grpc/test/bufconn"
    "github.com/nalej/grpc-utils/pkg/tools"
    "github.com/nalej/grpc-utils/pkg/test"
    musicianHandler "github.com/nalej/conductor/pkg/musician/service/handler"
    musicianScorer "github.com/nalej/conductor/pkg/musician/scorer"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/nalej/conductor/internal/entities"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "fmt"

    "time"

)


var _ = ginkgo.Describe ("Simple scorer functionality with two musicians", func() {
    // grpc servers
    var servers []*tools.GenericGRPCServer
    // scorer
    var scorerMethod Scorer
    // Musician handler
    var managers []*musicianHandler.Manager
    // Fake status collectors
    var collectors []statuscollector.StatusCollector
    // Musician clients
    var clients *tools.ConnectionsMap

    ginkgo.BeforeEach(func() {
        // instantiate musicianHandler server
        scorerMethod = NewSimpleScorer()
        // instantiate collectors
        collectors = make([]statuscollector.StatusCollector,2)
        collectors[0] = statuscollector.NewFakeCollector()
        collectors[1] = statuscollector.NewFakeCollector()

        // instantiate musicians
        managers = make([]*musicianHandler.Manager,2)
        managers[0] = musicianHandler.NewManager(&collectors[0], musicianScorer.NewSimpleScorer(collectors[0]))
        managers[1] = musicianHandler.NewManager(&collectors[1], musicianScorer.NewSimpleScorer(collectors[1]))

        servers = make([]*tools.GenericGRPCServer,2)
        port1, _ := test.GetAvailablePort()
        servers[0] = tools.NewGenericGRPCServer(uint32(port1))
        port2, _ := test.GetAvailablePort()
        servers[1] = tools.NewGenericGRPCServer(uint32(port2))

        go servers[0].Run()
        go servers[1].Run()

        // Add the client
        pbConductor.RegisterMusicianServer(servers[0].Server, musicianHandler.NewHandler(managers[0]))
        pbConductor.RegisterMusicianServer(servers[1].Server, musicianHandler.NewHandler(managers[1]))

        clients = conductor.GetMusicianClients()

        // courtesy sleep to ensure all the grpc servers are up.
        time.Sleep(time.Second*2)
        clients.AddConnection(fmt.Sprintf("localhost:%d",servers[0].Port))
        clients.AddConnection(fmt.Sprintf("localhost:%d",servers[1].Port))

    })

    ginkgo.AfterEach(func(){
        for _,s := range servers {
            s.Server.Stop()
        }
    })

    ginkgo.Describe("sent requirements that only fit into one cluster", func(){
        var request entities.Requirements

        ginkgo.BeforeEach(func(){
            request = entities.Requirements{CPU:0.5,Memory:100, Disk:100}

            // collector 0 says overload
            overloaded_status := entities.Status{CPU: 0.87, Mem: 32000, Disk:100}
            collectors[0].(*statuscollector.FakeCollector).SetStatus(overloaded_status)
            // collector 1 says free
            free_status := entities.Status{CPU: 0.10, Mem: 5000, Disk: 200}
            collectors[1].(*statuscollector.FakeCollector).SetStatus(free_status)

        })

        ginkgo.Context("the cluster with lowest occupation is chosen", func(){
            ginkgo.It("second cluster has the highest score", func(){
                response, err := scorerMethod.ScoreRequirements(&request)
                gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
                gomega.Expect(response).NotTo(gomega.BeNil())
                gomega.Expect(response.TotalEvaluated).To(gomega.Equal(2))
                gomega.Expect(response.Score).To(gomega.Equal(float32(1.88)))
            })
        })
    })
})

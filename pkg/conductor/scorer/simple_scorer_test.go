/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package scorer

import (
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    "github.com/nalej/grpc-utils/pkg/tools"
    musicianHandler "github.com/nalej/conductor/pkg/musician/service/handler"
    musicianScorer "github.com/nalej/conductor/pkg/musician/scorer"
    pbInfrastructure "github.com/nalej/grpc-infrastructure-go"
    pbOrganization "github.com/nalej/grpc-organization-go"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
    "github.com/nalej/conductor/internal/entities"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "os"
    "time"
    "github.com/nalej/conductor/pkg/utils"
    "context"
    "fmt"
)


var _ = ginkgo.Describe ("Simple scorer functionality with two musicians", func() {
    var isReady bool
    // Connections helper
    var connHelper *utils.ConnectionsHelper
    // grpc servers
    var servers []*tools.GenericGRPCServer
    // scorer
    var scorerMethod Scorer
    // Musician baton
    var managers []*musicianHandler.Manager
    // Fake status collectors
    var collectors []statuscollector.StatusCollector
    // Musician clients
    var clients *tools.ConnectionsMap
    // System model address
    var smHostname string
    // organization
    organizationName := "testOrganization"
    var organizationId string



    ginkgo.BeforeSuite(func(){
        // Check these are integration tests
        isReady = false
        if utils.RunIntegrationTests() {
            smHostname = os.Getenv(utils.IT_SYSTEM_MODEL)
            if smHostname != "" {
                isReady=true
            }
        }

        if !isReady {
            return
        }

        connHelper = utils.NewConnectionsHelper(false, "", true)

        smAddress := fmt.Sprintf("%s:%d", utils.IT_SYSTEM_MODEL, utils.SYSTEM_MODEL_PORT)

        // initialize a system model
        sm := connHelper.GetSystemModelClients()
        sm.AddConnection(smHostname, int(utils.SYSTEM_MODEL_PORT))

        pool := connHelper.GetSystemModelClients()
        conn,err := pool.GetConnection(smAddress)
        gomega.Expect(err).ToNot(gomega.HaveOccurred())

        // Add an organization
        smOrganization := pbOrganization.NewOrganizationsClient(conn)

        orgReq := pbOrganization.AddOrganizationRequest{Name: organizationName}
        orgResp, err := smOrganization.AddOrganization(context.Background(), &orgReq)
        gomega.Expect(err).ToNot(gomega.HaveOccurred())
        organizationId = orgResp.OrganizationId

        // Add available clusters
        smClient := pbInfrastructure.NewClustersClient(conn)
        req := pbInfrastructure.AddClusterRequest{
            Hostname: "localhost",
            OrganizationId: orgResp.OrganizationId,
            Name: "testCluster",
            RequestId: "req001",
            Description: "A test description",
            Labels: map[string]string{"key1":"value1"},
        }
        _, err = smClient.AddCluster(context.Background(),&req)
        gomega.Expect(err).ToNot(gomega.HaveOccurred())

    })


    ginkgo.BeforeEach(func() {
        if !isReady {
            ginkgo.Skip("run integration test not configured")
        }

        // instantiate musicianHandler server
        scorerMethod = NewSimpleScorer(connHelper)
        // instantiate collectors
        collectors = make([]statuscollector.StatusCollector,2)
        collectors[0] = statuscollector.NewFakeCollector()
        collectors[1] = statuscollector.NewFakeCollector()

        // instantiate musicians
        managers = make([]*musicianHandler.Manager,2)
        managers[0] = musicianHandler.NewManager(&collectors[0], musicianScorer.NewSimpleScorer(collectors[0]))
        managers[1] = musicianHandler.NewManager(&collectors[1], musicianScorer.NewSimpleScorer(collectors[1]))

        servers = make([]*tools.GenericGRPCServer,1)
        //port1, _ := test.GetAvailablePort()
        port1 := utils.MUSICIAN_PORT
        servers[0] = tools.NewGenericGRPCServer(uint32(port1))
        // Only one cluster can be attempted until we have
        //port2, _ := test.GetAvailablePort()
        //servers[1] = tools.NewGenericGRPCServer(uint32(port2))

        go servers[0].Run()
        //go servers[1].Run()

        // Add the client
        pbConductor.RegisterMusicianServer(servers[0].Server, musicianHandler.NewHandler(managers[0]))
        //pbConductor.RegisterMusicianServer(servers[1].Server, musicianHandler.NewHandler(managers[1]))

        clients = connHelper.GetClusterClients()

        // courtesy sleep to ensure all the grpc servers are up.
        time.Sleep(time.Second*2)
        clients.AddConnection("localhost",int(servers[0].Port))
        //clients.AddConnection(fmt.Sprintf("localhost:%d",servers[1].Port))

    })

    ginkgo.AfterEach(func(){
        if !isReady {
            ginkgo.Skip("run integration test not configured")
        }
        for _,s := range servers {
            s.Server.Stop()
        }
    })

    ginkgo.Describe("sent requirements that only fit into one cluster", func(){
        var request entities.Requirements

        ginkgo.BeforeEach(func(){


            if !isReady {
                ginkgo.Skip("run integration test not configured")
            }

            request = entities.Requirements{[]entities.Requirement{
                {GroupServiceId: "serviceid",Replicas:1, CPU:50,Memory:100, Storage:100},
            }}

            // collector 0 says overload
            overloaded_status := entities.Status{CPU: 0.87, Mem: 32000, Disk:100}
            collectors[0].(*statuscollector.FakeCollector).SetStatus(overloaded_status)
            // collector 1 says free
            free_status := entities.Status{CPU: 0.10, Mem: 5000, Disk: 200}
            collectors[1].(*statuscollector.FakeCollector).SetStatus(free_status)

        })

        ginkgo.Context("the cluster with lowest occupation is chosen", func(){
            ginkgo.It("second cluster has the highest score", func(){
                if !isReady {
                    ginkgo.Skip("integration tests were not set")
                }
                response, err := scorerMethod.ScoreRequirements(organizationId,&request)
                gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
                gomega.Expect(response).NotTo(gomega.BeNil())
                gomega.Expect(response.NumEvaluatedClusters).To(gomega.Equal(1))
                gomega.Expect(len(response.DeploymentsScore)).To(gomega.Equal(1))
                gomega.Expect(response.DeploymentsScore[0].Score).To(gomega.Equal(float32(-48.87312316894531)))
            })
        })
    })
})

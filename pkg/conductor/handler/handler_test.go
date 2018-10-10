/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package handler

import (
    "github.com/nalej/conductor/pkg/utils"
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
    pbOrganization "github.com/nalej/grpc-organization-go"
    "github.com/rs/zerolog/log"
    "google.golang.org/grpc/test/bufconn"
    "google.golang.org/grpc"
    "context"
    "github.com/nalej/grpc-utils/pkg/test"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/conductor/requirementscollector"

    "github.com/nalej/conductor/pkg/conductor"
)


const (
    // TODO Set an ENV variable.
    SystemModelAddress="127.0.0.1:8800"
)

func InitializeEntries(orgClient pbOrganization.OrganizationsClient, appClient pbApplication.ApplicationsClient) *pbApplication.AppDescriptor{
    // add an organization
    orgRequest := pbOrganization.AddOrganizationRequest{Name: "org001"}
    resp, err := orgClient.AddOrganization(context.Background(),&orgRequest)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())


    port1 := pbApplication.Port{Name: "port1", ExposedPort: 3000}
    port2 := pbApplication.Port{Name: "port2", ExposedPort: 3001}
    credentials := pbApplication.ImageCredentials{Username: "user1", Password: "password1", Email: "email@email.com"}

    serv := pbApplication.Service{
        OrganizationId: resp.OrganizationId,
        ServiceId: "service_001",
        Name: "test-image-1",
        Image: "nginx:1.12",
        ExposedPorts: []*pbApplication.Port{&port1, &port2},
        Labels: map[string]string { "label1":"value1", "label2":"value2"},
        Specs: &pbApplication.DeploySpecs{Replicas: 1},
        AppDescriptorId: "app001",
        Description: "Test service",
        EnvironmentVariables: map[string]string{"var1":"var1"},
        Type: pbApplication.ServiceType_DOCKER,
        DeployAfter: []string{},
        Storage: []*pbApplication.Storage{&pbApplication.Storage{MountPath: "/tmp",}},
        Credentials: &credentials,
        Configs: []*pbApplication.ConfigFile{&pbApplication.ConfigFile{MountPath:"/tmp"}},
    }

    servGroup := pbApplication.ServiceGroup{
        OrganizationId:resp.OrganizationId,
        Description:"a service group",
        AppDescriptorId: "app001",
        Name: "group001",
        Services: []string{"test-image-1"},
        Policy: pbApplication.CollocationPolicy_SAME_CLUSTER,
        ServiceGroupId: "group-id",
    }

    secRule := pbApplication.SecurityRule{
        Name:"securityrule",
        AppDescriptorId: "app001",
        OrganizationId: resp.OrganizationId,
        Access: pbApplication.PortAccess_ALL_APP_SERVICES,
        AuthServices: []string{"auth"},
        DeviceGroups: []string{"devgroup"},
        RuleId: "rule001",
        SourcePort: 30000,
        SourceServiceId: "sourceserv001",
    }

    // add a desriptor
    appDescriptor := pbApplication.AddAppDescriptorRequest{
        RequestId: "req001",
        Name:"app_descriptor_test",
        Description: "app_descriptor_test description",
        OrganizationId: resp.OrganizationId,
        EnvironmentVariables: map[string]string{"var1":"var1_value", "var2":"var2_value"},
        Labels: map[string]string{"label1":"label1_value", "label2":"label2_value"},
        Services: []*pbApplication.Service{&serv},
        ConfigurationOptions: map[string]string{"conf1":"valueconf1", "conf2":"valueconf2"},
        Groups: []*pbApplication.ServiceGroup{&servGroup},
        Rules: []*pbApplication.SecurityRule{&secRule},
    }
    desc, err := appClient.AddAppDescriptor(context.Background(),&appDescriptor)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
    return desc

}


var _ = ginkgo.Describe("Deployment server API", func() {

    if ! utils.RunIntegrationTests() {
        log.Warn().Msg("Integration tests are skipped")
        return
    }

    // grpc server
    var server *grpc.Server
    // conductor object
    var cond *Manager
    // grpc test listener
    var listener *bufconn.Listener
    // queue
    var q RequestsQueue
    // Connection with system model
    var connSM *grpc.ClientConn
    // Applications client
    var appClient pbApplication.ApplicationsClient
    // Organizations client
    var orgClient pbOrganization.OrganizationsClient
    // Conductor client
    var client pbConductor.ConductorClient
    // Used application descriptor
    var appDescriptor *pbApplication.AppDescriptor

    ginkgo.BeforeSuite(func(){
        // connect with external system model using the pool
        pool := conductor.GetSystemModelClients()
        var err error
        connSM, err = pool.AddConnection(SystemModelAddress)
        gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

        listener = test.GetDefaultListener()
        server = grpc.NewServer()
        scorerMethod := scorer.NewSimpleScorer()
        designer := plandesigner.NewSimplePlanDesigner()
        reqcoll := requirementscollector.NewSimpleRequirementsCollector()
        q = NewMemoryRequestQueue()

        conn, err := test.GetConn(*listener)
        gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
        client = pbConductor.NewConductorClient(conn)

        // clients
        appClient = pbApplication.NewApplicationsClient(connSM)
        orgClient = pbOrganization.NewOrganizationsClient(connSM)

        cond = NewManager(q, scorerMethod, reqcoll, designer)
        test.LaunchServer(server,listener)

        // Register the service.
        pbConductor.RegisterConductorServer(server, NewHandler(cond))

        appDescriptor = InitializeEntries(orgClient, appClient)

    })

    ginkgo.AfterSuite(func(){
        listener.Close()
        server.Stop()
        connSM.Close()
    })


    ginkgo.Context("The queue is empty and a new request arrives", func() {
        var request pbConductor.DeploymentRequest
        var response pbConductor.DeploymentResponse

        ginkgo.BeforeEach(func() {
            request = pbConductor.DeploymentRequest{
                RequestId: "myrequestId",
                AppId: &pbApplication.AppDescriptorId{OrganizationId:appDescriptor.OrganizationId,AppDescriptorId: appDescriptor.AppDescriptorId},
                Description: "A single description",
                Name: "A testing application",
            }
            response = pbConductor.DeploymentResponse{
                RequestId: "myrequestId",
                Status: pbApplication.ApplicationStatus_QUEUED}
        })


        ginkgo.It("receive an expected message", func() {
            resp, err := client.Deploy(context.Background(), &request)
            //gomega.Expect(resp.String()).To(gomega.Equal(response.String()))
            gomega.Expect(resp.RequestId).To(gomega.Equal(response.RequestId))
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
        })
    })

})


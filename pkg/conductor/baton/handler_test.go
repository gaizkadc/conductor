/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package baton

import (
    "context"
    "github.com/google/uuid"
    "github.com/nalej/conductor/internal/structures"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/conductor/requirementscollector"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/conductor/pkg/utils"
    pbApplication "github.com/nalej/grpc-application-go"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbOrganization "github.com/nalej/grpc-organization-go"
    "github.com/nalej/grpc-utils/pkg/test"
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    "google.golang.org/grpc"
    "google.golang.org/grpc/test/bufconn"

    "os"
)


func InitializeEntries(orgClient pbOrganization.OrganizationsClient, appClient pbApplication.ApplicationsClient) *pbApplication.ParametrizedDescriptor{
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
        ServiceGroupId: "group001",
        ExposedPorts: []*pbApplication.Port{&port1, &port2},
        Labels: map[string]string { "label1":"value1", "label2":"value2"},
        Specs: &pbApplication.DeploySpecs{Replicas: 1},
        AppDescriptorId: "app001",
        EnvironmentVariables: map[string]string{"var1":"var1"},
        Type: pbApplication.ServiceType_DOCKER,
        DeployAfter: []string{},
        Storage: []*pbApplication.Storage{&pbApplication.Storage{MountPath: "/tmp",}},
        Credentials: &credentials,
        Configs: []*pbApplication.ConfigFile{&pbApplication.ConfigFile{MountPath:"/tmp"}},
    }

    servGroup := pbApplication.ServiceGroup{
        OrganizationId:resp.OrganizationId,
        AppDescriptorId: "app001",
        Name: "group001",
        Services: []*pbApplication.Service{&serv},
        Policy: pbApplication.CollocationPolicy_SAME_CLUSTER,
    }

    secRule := pbApplication.SecurityRule{
        Name:"securityrule",
        AppDescriptorId: "app001",
        OrganizationId: resp.OrganizationId,
        Access: pbApplication.PortAccess_ALL_APP_SERVICES,
        AuthServices: []string{"auth"},
        DeviceGroupNames: []string{"devgroup"},
        RuleId: "rule001",
        TargetServiceGroupName: "group001",
        TargetServiceName: "service_001",
        TargetPort: 30000,
        AuthServiceGroupName: "group001",
    }

    // add a desriptor
    appDescriptor := pbApplication.ParametrizedDescriptor{
        Name:"app_descriptor_test",
        OrganizationId: resp.OrganizationId,
        AppDescriptorId: uuid.New().String(),
        AppInstanceId: uuid.New().String(),
        EnvironmentVariables: map[string]string{"var1":"var1_value", "var2":"var2_value"},
        Labels: map[string]string{"label1":"label1_value", "label2":"label2_value"},
        ConfigurationOptions: map[string]string{"conf1":"valueconf1", "conf2":"valueconf2"},
        Groups: []*pbApplication.ServiceGroup{&servGroup},
        Rules: []*pbApplication.SecurityRule{&secRule},
    }
    _, err = appClient.AddParametrizedDescriptor(context.Background(),&appDescriptor)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())


    return &appDescriptor

}


var _ = ginkgo.Describe("Deployment server API", func() {
    var isReady bool
    // Connections helper
    var connHelper *utils.ConnectionsHelper
    // System model address
    var systemModelHostname string
    // grpc server
    var server *grpc.Server
    // conductor object
    var cond *Manager
    // grpc test listener
    var listener *bufconn.Listener
    // queue
    var q structures.RequestsQueue
    // pending plans controller
    var plans *structures.PendingPlans
    // Connection with system model
    var connSM *grpc.ClientConn
    // Applications client
    var appClient pbApplication.ApplicationsClient
    // Organizations client
    var orgClient pbOrganization.OrganizationsClient
    // Conductor client
    var client pbConductor.ConductorClient
    // Used application descriptor
    var appDescriptor *pbApplication.ParametrizedDescriptor


    ginkgo.BeforeSuite(func(){
        isReady = false

        if utils.RunIntegrationTests() {
            systemModelHostname = os.Getenv(utils.IT_SYSTEM_MODEL)
            if systemModelHostname != "" {
                isReady = true
            }
        }

        if !isReady {
            return
        }

        connHelper = utils.NewConnectionsHelper(false, "", true)


        // connect with external system model using the pool
        pool := connHelper.GetSystemModelClients()
        var err error
        connSM, err = pool.AddConnection(systemModelHostname, int(utils.SYSTEM_MODEL_PORT))
        gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

        listener = test.GetDefaultListener()
        server = grpc.NewServer()
        scorerMethod := scorer.NewSimpleScorer(connHelper)
        designer := plandesigner.NewSimpleReplicaPlanDesigner(connHelper)
        reqcoll := requirementscollector.NewSimpleRequirementsCollector()
        q = structures.NewMemoryRequestQueue()
        plans = structures.NewPendingPlans()


        conn, err := test.GetConn(*listener)
        gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
        client = pbConductor.NewConductorClient(conn)

        // clients
        appClient = pbApplication.NewApplicationsClient(connSM)
        orgClient = pbOrganization.NewOrganizationsClient(connSM)


        cond = NewManager(connHelper, q, scorerMethod, reqcoll, designer,plans,nil)
        test.LaunchServer(server,listener)

        // Register the service.
        pbConductor.RegisterConductorServer(server, NewHandler(cond))

        appDescriptor = InitializeEntries(orgClient, appClient)

    })

    ginkgo.AfterSuite(func(){
        if isReady {
            listener.Close()
            server.Stop()
            connSM.Close()
        }
    })


    ginkgo.Context("The queue is empty and a new request arrives", func() {
        var request pbConductor.DeploymentRequest

        ginkgo.BeforeEach(func() {
            if !isReady {
                return
            }
            request = pbConductor.DeploymentRequest{
                RequestId: "myrequestId",
                AppInstanceId: &pbApplication.AppInstanceId{OrganizationId:appDescriptor.OrganizationId,AppInstanceId: appDescriptor.AppInstanceId},
                Name: "A testing application",
            }
        })


        ginkgo.It("receive an expected message", func() {
            if !isReady {
                ginkgo.Skip("Integration environment was not set")
            }

            success, err := client.Deploy(context.Background(), &request)
            //gomega.Expect(resp.String()).To(gomega.Equal(response.String()))
            gomega.Expect(success).NotTo(gomega.BeNil())
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
        })
    })

})


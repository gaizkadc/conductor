/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package plandesigner

import (
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    "github.com/nalej/conductor/pkg/utils"
    "google.golang.org/grpc"
    pbApplication "github.com/nalej/grpc-application-go"
    pbOrganization "github.com/nalej/grpc-organization-go"
    "os"
    "context"
    "github.com/nalej/conductor/internal/entities"
)

var _ = ginkgo.Describe("Check plan designer", func(){

    var isReady bool

    // Connections helper
    var connHelper *utils.ConnectionsHelper

    var localPlanDesigner PlanDesigner

    // System model address
    var systemModelHost string
    // Connection with system model
    var connSM *grpc.ClientConn
    // Applications client
    var appClient pbApplication.ApplicationsClient
    // Organizations client
    var orgClient pbOrganization.OrganizationsClient


    ginkgo.BeforeSuite(func(){
        isReady = false

        if utils.RunIntegrationTests() {
            systemModelHost = os.Getenv(utils.IT_SYSTEM_MODEL)
            if systemModelHost != "" {
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
        connSM, err = pool.AddConnection(systemModelHost, int(utils.SYSTEM_MODEL_PORT))
        gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

        // clients
        appClient = pbApplication.NewApplicationsClient(connSM)
        orgClient = pbOrganization.NewOrganizationsClient(connSM)

    })

    ginkgo.Context("single fragment two stages", func(){
        ginkgo.It("create the expected deployment plan", func(){
            if !isReady {
                ginkgo.Skip("run integration test not configured")
            }
            appInstance := CreateApp1(orgClient, appClient)
            localPlanDesigner = NewSimplePlanDesigner(connHelper)
            score := entities.ClustersScore{TotalEvaluated: 1, Scoring: []entities.ClusterScore{{Score:0.99, ClusterId: "cluster1"}}}
            plan,err := localPlanDesigner.DesignPlan(appInstance, &score)
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
            // only one fragment
            gomega.Expect(len(plan.Fragments)).To(gomega.Equal(1))
            // two stages
            stages := plan.Fragments[0].Stages
            gomega.Expect(len(stages)).To(gomega.Equal(2))
            // only one service per stage
            gomega.Expect(len(stages[0].Services)).To(gomega.Equal(1))
            gomega.Expect(len(stages[1].Services)).To(gomega.Equal(1))
            // check the correctness of the service id
            gomega.Expect(stages[0].Services[0].ServiceId).To(gomega.Equal("service_001"))
            gomega.Expect(stages[1].Services[0].ServiceId).To(gomega.Equal("service_002"))
        })
    })

    ginkgo.Context("single fragment two stages with different services", func(){
        ginkgo.It("create the expected deployment plan", func(){
            if !isReady {
                ginkgo.Skip("run integration test not configured")
            }
            appInstance := CreateApp2(orgClient, appClient)
            localPlanDesigner = NewSimplePlanDesigner(connHelper)
            score := entities.ClustersScore{TotalEvaluated: 1, Scoring: []entities.ClusterScore{{Score:0.99, ClusterId: "cluster1"}}}
            plan,err := localPlanDesigner.DesignPlan(appInstance, &score)
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
            // only one fragment
            gomega.Expect(len(plan.Fragments)).To(gomega.Equal(1))
            // two stages
            stages := plan.Fragments[0].Stages
            gomega.Expect(len(stages)).To(gomega.Equal(2))
            // first stage two services serv1, serv3 and second stage with serv2
            gomega.Expect(len(stages[0].Services)).To(gomega.Equal(2))
            gomega.Expect(len(stages[1].Services)).To(gomega.Equal(1))
            // check the correctness of the service id
            gomega.Expect(stages[0].Services[0].ServiceId).To(gomega.Equal("service_001"))
            gomega.Expect(stages[0].Services[1].ServiceId).To(gomega.Equal("service_003"))
            gomega.Expect(stages[1].Services[0].ServiceId).To(gomega.Equal("service_002"))
        })
    })

    ginkgo.Context("single fragment twe stages with 1, 2 and 3 services", func(){
        ginkgo.It("create the expected deployment plan", func(){
            if !isReady {
                ginkgo.Skip("run integration test not configured")
            }
            appInstance := CreateApp3(orgClient, appClient)
            localPlanDesigner = NewSimplePlanDesigner(connHelper)
            score := entities.ClustersScore{TotalEvaluated: 1, Scoring: []entities.ClusterScore{{Score:0.99, ClusterId: "cluster1"}}}
            plan,err := localPlanDesigner.DesignPlan(appInstance, &score)
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
            // only one fragment
            gomega.Expect(len(plan.Fragments)).To(gomega.Equal(1))
            // two stages
            stages := plan.Fragments[0].Stages
            gomega.Expect(len(stages)).To(gomega.Equal(3))
            // first stage two services serv1, serv3 and second stage with serv2
            gomega.Expect(len(stages[0].Services)).To(gomega.Equal(1))
            gomega.Expect(len(stages[1].Services)).To(gomega.Equal(2))
            gomega.Expect(len(stages[2].Services)).To(gomega.Equal(1))
            // check the correctness of the service id
            gomega.Expect(stages[0].Services[0].ServiceId).To(gomega.Equal("service_001"))
            gomega.Expect(stages[1].Services[0].ServiceId).To(gomega.Equal("service_002"))
            gomega.Expect(stages[1].Services[1].ServiceId).To(gomega.Equal("service_004"))
            gomega.Expect(stages[2].Services[0].ServiceId).To(gomega.Equal("service_003"))
        })
    })

})

// CreateApp1 instantiates an application with two depending services.
func CreateApp1(orgClient pbOrganization.OrganizationsClient, appClient pbApplication.ApplicationsClient) *pbApplication.AppInstance{
    // add an organization
    orgRequest := pbOrganization.AddOrganizationRequest{Name: "org001"}
    resp, err := orgClient.AddOrganization(context.Background(),&orgRequest)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())


    port1 := pbApplication.Port{Name: "port1", ExposedPort: 3000}
    port2 := pbApplication.Port{Name: "port2", ExposedPort: 3001}
    credentials := pbApplication.ImageCredentials{Username: "user1", Password: "password1", Email: "email@email.com"}

    serv1 := pbApplication.Service{
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
    // serv2 must be deployed after serv1
    serv2 := pbApplication.Service{
        OrganizationId: resp.OrganizationId,
        ServiceId: "service_002",
        Name: "test-image-2",
        Image: "nginx:1.12",
        ExposedPorts: []*pbApplication.Port{&port1, &port2},
        Labels: map[string]string { "label1":"value1", "label2":"value2"},
        Specs: &pbApplication.DeploySpecs{Replicas: 1},
        AppDescriptorId: "app001",
        Description: "Test service",
        EnvironmentVariables: map[string]string{"var1":"var1"},
        Type: pbApplication.ServiceType_DOCKER,
        DeployAfter: []string{"service_001"},
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
        Services: []*pbApplication.Service{&serv1, &serv2},
        ConfigurationOptions: map[string]string{"conf1":"valueconf1", "conf2":"valueconf2"},
        Groups: []*pbApplication.ServiceGroup{&servGroup},
        Rules: []*pbApplication.SecurityRule{&secRule},
    }
    desc, err := appClient.AddAppDescriptor(context.Background(),&appDescriptor)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())


    // Add the application instance
    instanceReq := pbApplication.AddAppInstanceRequest{
        AppDescriptorId: desc.AppDescriptorId,
        Name:"app_descriptor_test",
        Description: "app_descriptor_test description",
        OrganizationId: resp.OrganizationId,
    }
    appInstance, err := appClient.AddAppInstance(context.Background(), &instanceReq)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

    return appInstance

}

// CreateApp2 instanciates a service with the following dependencies
// serv1 <- serv2, serv3
func CreateApp2(orgClient pbOrganization.OrganizationsClient, appClient pbApplication.ApplicationsClient) *pbApplication.AppInstance{
    // add an organization
    orgRequest := pbOrganization.AddOrganizationRequest{Name: "org001"}
    resp, err := orgClient.AddOrganization(context.Background(),&orgRequest)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())


    port1 := pbApplication.Port{Name: "port1", ExposedPort: 3000}
    port2 := pbApplication.Port{Name: "port2", ExposedPort: 3001}
    credentials := pbApplication.ImageCredentials{Username: "user1", Password: "password1", Email: "email@email.com"}

    serv1 := pbApplication.Service{
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
    // serv2 must be deployed after serv1
    serv2 := pbApplication.Service{
        OrganizationId: resp.OrganizationId,
        ServiceId: "service_002",
        Name: "test-image-2",
        Image: "nginx:1.12",
        ExposedPorts: []*pbApplication.Port{&port1, &port2},
        Labels: map[string]string { "label1":"value1", "label2":"value2"},
        Specs: &pbApplication.DeploySpecs{Replicas: 1},
        AppDescriptorId: "app001",
        Description: "Test service",
        EnvironmentVariables: map[string]string{"var1":"var1"},
        Type: pbApplication.ServiceType_DOCKER,
        DeployAfter: []string{"service_001"},
        Storage: []*pbApplication.Storage{&pbApplication.Storage{MountPath: "/tmp",}},
        Credentials: &credentials,
        Configs: []*pbApplication.ConfigFile{&pbApplication.ConfigFile{MountPath:"/tmp"}},
    }

    // serv2 must be deployed after serv1
    serv3 := pbApplication.Service{
        OrganizationId: resp.OrganizationId,
        ServiceId: "service_003",
        Name: "test-image-3",
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
        Services: []*pbApplication.Service{&serv1, &serv2, &serv3},
        ConfigurationOptions: map[string]string{"conf1":"valueconf1", "conf2":"valueconf2"},
        Groups: []*pbApplication.ServiceGroup{&servGroup},
        Rules: []*pbApplication.SecurityRule{&secRule},
    }
    desc, err := appClient.AddAppDescriptor(context.Background(),&appDescriptor)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())


    // Add the application instance
    instanceReq := pbApplication.AddAppInstanceRequest{
        AppDescriptorId: desc.AppDescriptorId,
        Name:"app_descriptor_test",
        Description: "app_descriptor_test description",
        OrganizationId: resp.OrganizationId,
    }
    appInstance, err := appClient.AddAppInstance(context.Background(), &instanceReq)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

    return appInstance

}

// CreateApp2 instanciates a service with the following dependencies
// serv1 <- serv2 <- serv3,
// serv1 <- serv4
func CreateApp3(orgClient pbOrganization.OrganizationsClient, appClient pbApplication.ApplicationsClient) *pbApplication.AppInstance{
    // add an organization
    orgRequest := pbOrganization.AddOrganizationRequest{Name: "org001"}
    resp, err := orgClient.AddOrganization(context.Background(),&orgRequest)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())


    port1 := pbApplication.Port{Name: "port1", ExposedPort: 3000}
    port2 := pbApplication.Port{Name: "port2", ExposedPort: 3001}
    credentials := pbApplication.ImageCredentials{Username: "user1", Password: "password1", Email: "email@email.com"}

    serv1 := pbApplication.Service{
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
    // serv2 must be deployed after serv1
    serv2 := pbApplication.Service{
        OrganizationId: resp.OrganizationId,
        ServiceId: "service_002",
        Name: "test-image-2",
        Image: "nginx:1.12",
        ExposedPorts: []*pbApplication.Port{&port1, &port2},
        Labels: map[string]string { "label1":"value1", "label2":"value2"},
        Specs: &pbApplication.DeploySpecs{Replicas: 1},
        AppDescriptorId: "app001",
        Description: "Test service",
        EnvironmentVariables: map[string]string{"var1":"var1"},
        Type: pbApplication.ServiceType_DOCKER,
        DeployAfter: []string{"service_001"},
        Storage: []*pbApplication.Storage{&pbApplication.Storage{MountPath: "/tmp",}},
        Credentials: &credentials,
        Configs: []*pbApplication.ConfigFile{&pbApplication.ConfigFile{MountPath:"/tmp"}},
    }

    // serv3 must be deployed after serv2
    serv3 := pbApplication.Service{
        OrganizationId: resp.OrganizationId,
        ServiceId: "service_003",
        Name: "test-image-3",
        Image: "nginx:1.12",
        ExposedPorts: []*pbApplication.Port{&port1, &port2},
        Labels: map[string]string { "label1":"value1", "label2":"value2"},
        Specs: &pbApplication.DeploySpecs{Replicas: 1},
        AppDescriptorId: "app001",
        Description: "Test service",
        EnvironmentVariables: map[string]string{"var1":"var1"},
        Type: pbApplication.ServiceType_DOCKER,
        DeployAfter: []string{"service_002"},
        Storage: []*pbApplication.Storage{&pbApplication.Storage{MountPath: "/tmp",}},
        Credentials: &credentials,
        Configs: []*pbApplication.ConfigFile{&pbApplication.ConfigFile{MountPath:"/tmp"}},
    }

    // serv4 must be deployed after serv1
    serv4 := pbApplication.Service{
        OrganizationId: resp.OrganizationId,
        ServiceId: "service_004",
        Name: "test-image-4",
        Image: "nginx:1.12",
        ExposedPorts: []*pbApplication.Port{&port1, &port2},
        Labels: map[string]string { "label1":"value1", "label2":"value2"},
        Specs: &pbApplication.DeploySpecs{Replicas: 1},
        AppDescriptorId: "app001",
        Description: "Test service",
        EnvironmentVariables: map[string]string{"var1":"var1"},
        Type: pbApplication.ServiceType_DOCKER,
        DeployAfter: []string{"service_001"},
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
        Services: []*pbApplication.Service{&serv1, &serv2, &serv3, &serv4},
        ConfigurationOptions: map[string]string{"conf1":"valueconf1", "conf2":"valueconf2"},
        Groups: []*pbApplication.ServiceGroup{&servGroup},
        Rules: []*pbApplication.SecurityRule{&secRule},
    }
    desc, err := appClient.AddAppDescriptor(context.Background(),&appDescriptor)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())


    // Add the application instance
    instanceReq := pbApplication.AddAppInstanceRequest{
        AppDescriptorId: desc.AppDescriptorId,
        Name:"app_descriptor_test",
        Description: "app_descriptor_test description",
        OrganizationId: resp.OrganizationId,
    }
    appInstance, err := appClient.AddAppInstance(context.Background(), &instanceReq)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

    return appInstance

}
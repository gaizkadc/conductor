/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package cmd

import (
    "fmt"
    pbApplication "github.com/nalej/grpc-application-go"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbInfrastructure "github.com/nalej/grpc-infrastructure-go"
    pbOrganization "github.com/nalej/grpc-organization-go"
    "github.com/rs/zerolog/log"
    "github.com/spf13/cobra"

    "bufio"
    "context"
    "github.com/nalej/conductor/pkg/utils"
    "google.golang.org/grpc"
    "os"
)


var runDemo = &cobra.Command{
    Use: "demo",
    Short: "Run a demo",
    Long: "Run a conductor demo with some service examples...",
    Run: func(cmd *cobra.Command, args [] string) {
        SetupLogging()
        RunExample()
    },
}

const msg =
` $$$$$$\   $$$$$$\  $$\   $$\ $$$$$$$\  $$\   $$\  $$$$$$\ $$$$$$$$\  $$$$$$\  $$$$$$$\  
$$  __$$\ $$  __$$\ $$$\  $$ |$$  __$$\ $$ |  $$ |$$  __$$\\__$$  __|$$  __$$\ $$  __$$\ 
$$ /  \__|$$ /  $$ |$$$$\ $$ |$$ |  $$ |$$ |  $$ |$$ /  \__|  $$ |   $$ /  $$ |$$ |  $$ |
$$ |      $$ |  $$ |$$ $$\$$ |$$ |  $$ |$$ |  $$ |$$ |        $$ |   $$ |  $$ |$$$$$$$  |
$$ |      $$ |  $$ |$$ \$$$$ |$$ |  $$ |$$ |  $$ |$$ |        $$ |   $$ |  $$ |$$  __$$< 
$$ |  $$\ $$ |  $$ |$$ |\$$$ |$$ |  $$ |$$ |  $$ |$$ |  $$\   $$ |   $$ |  $$ |$$ |  $$ |
\$$$$$$  | $$$$$$  |$$ | \$$ |$$$$$$$  |\$$$$$$  |\$$$$$$  |  $$ |    $$$$$$  |$$ |  $$ |
 \______/  \______/ \__|  \__|\_______/  \______/  \______/   \__|    \______/ \__|  \__|
                                                                                         
                                                                                         
                                                                                         
      $$\                                                                                
      $$ |                                                                               
 $$$$$$$ | $$$$$$\  $$$$$$\$$$$\   $$$$$$\                                               
$$  __$$ |$$  __$$\ $$  _$$  _$$\ $$  __$$\                                              
$$ /  $$ |$$$$$$$$ |$$ / $$ / $$ |$$ /  $$ |                                             
$$ |  $$ |$$   ____|$$ | $$ | $$ |$$ |  $$ |                                             
\$$$$$$$ |\$$$$$$$\ $$ | $$ | $$ |\$$$$$$  |                                             
 \_______| \_______|\__| \__| \__| \______/                                              
                                                                                         
`

func init() {
    RootCmd.AddCommand(runDemo)

}

// Entrypoint for a musician service.
func RunExample() {
    fmt.Println()
    fmt.Println(msg)

    conn2, err := grpc.Dial("localhost:8800", grpc.WithInsecure())
    if err != nil {
        log.Panic().Err(err).Msg("impossible to connect with system model at port 8800")
        return
    }

    applicationClient := pbApplication.NewApplicationsClient(conn2)
    organizationClient := pbOrganization.NewOrganizationsClient(conn2)
    clusterClient := pbInfrastructure.NewClustersClient(conn2)

    orgId, orgErr := InitializeInfrastructure(organizationClient, clusterClient)
    if orgErr != nil {
        return
    }

    desc := InitializeEntries(orgId,applicationClient)

    log.Info().Msgf("%v",desc)

    log.Info().Msg("\nNow start conductor and musician then, press any key to run a deployment request")
    bufio.NewReader(os.Stdin).ReadBytes('\n')



    log.Info().Msg("connect with conductor api at port ")

    conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
    if err != nil {
        log.Panic().Err(err).Msg("impossible to connect with conductor at port 5000")
        return
    }
    conductorClient := pbConductor.NewConductorClient(conn)



    request := pbConductor.DeploymentRequest{
        RequestId: "req0001",
        Name: "Conductor demo deployment",
        Description: "A Nalej demo deployment",
        AppId: &pbApplication.AppDescriptorId{OrganizationId: desc.OrganizationId, AppDescriptorId: desc.AppDescriptorId},
    }
    x, err := conductorClient.Deploy(context.Background(), &request)
    if err != nil {
        log.Panic().Err(err).Msg("impossible to connect with conductor for deployment")
    }

    //log.Info().Msgf("The output instance works with id: %s",x.AppInstanceId)

    log.Info().Msg("\nPress any key to delete the generated namespace")
    bufio.NewReader(os.Stdin).ReadBytes('\n')

    //targetNamespace := getNamespace(desc.OrganizationId,x.AppInstanceId)
    //
    //output, err := exec.Command("kubectl","delete","namespace",targetNamespace).CombinedOutput()
    //if err != nil {
    //    os.Stderr.WriteString(err.Error())
    //}
    //fmt.Println(string(output))

    // Namespace deletion with newly implemented Undeploy operation
    undeployRequest := pbConductor.UndeployRequest{
        OrganizationId: desc.OrganizationId,
        AppInstaceId: x.AppInstanceId,
    }
    _, err = conductorClient.Undeploy(context.Background(), &undeployRequest)
    if err != nil {
        log.Panic().Err(err).Msg("unable to undeploy demo namespace")
    }
    log.Debug().Msgf("app %s successfully undeployed", x.AppInstanceId)

}

func InitializeInfrastructure(orgClient pbOrganization.OrganizationsClient,
    clustersClient pbInfrastructure.ClustersClient) (string,error) {
    // add an organization
    orgRequest := pbOrganization.AddOrganizationRequest{Name: "org001"}
    orgResp, err := orgClient.AddOrganization(context.Background(),&orgRequest)
    if err != nil {
        log.Panic().Err(err).Msg("impossible to connect with system model to update organization")
        return "",err
    }

    // add a cluster
    clusterReq := pbInfrastructure.AddClusterRequest{
        Name: "cluster001",
        OrganizationId: orgResp.OrganizationId,
        Labels: map[string]string{"clusterlabel":"clustervalue"},
        Description: "This is a simple testing value",
        RequestId: "req001",
        Hostname: "localhost",
    }
    addedCluster, err := clustersClient.AddCluster(context.Background(), &clusterReq)

    if err != nil {
        log.Panic().Err(err).Msg("impossible to add cluster")
        return "",err
    }

    log.Info().Msgf("set export %s=%s",utils.MUSICIAN_CLUSTER_ID, addedCluster.ClusterId)
    os.Setenv(utils.MUSICIAN_CLUSTER_ID,addedCluster.ClusterId)

    return orgResp.OrganizationId, nil
}


func InitializeEntries(orgId string, appClient pbApplication.ApplicationsClient) *pbApplication.AppDescriptor{

    port1 := pbApplication.Port{Name: "webport", ExposedPort: 80}
    port2 := pbApplication.Port{Name: "mysqlport", ExposedPort: 3306}

    credentials := pbApplication.ImageCredentials{Username: "user1", Password: "password1", Email: "email@email.com"}


    serv1 := pbApplication.Service {
        OrganizationId: orgId,
        ServiceId: "service_001",
        Name: "demo-wordpress",
        Image: "wordpress:4.8-apache",
        ExposedPorts: []*pbApplication.Port{&port1},
        Labels: map[string]string { "app":"test-wordpress", "component":"my-component"},
        Specs: &pbApplication.DeploySpecs{Replicas: 3},
        AppDescriptorId: "app001",
        Description: "A wordpress demo",
        EnvironmentVariables: map[string]string{"WORDPRESS_DB_HOST":"demo-mysql",
                                                "WORDPRESS_DB_PASSWORD":"root"},
        Type: pbApplication.ServiceType_DOCKER,
        DeployAfter: []string{"service_002"},
        Storage: []*pbApplication.Storage{&pbApplication.Storage{MountPath: "/tmp",}},
        Credentials: &credentials,
        Configs: []*pbApplication.ConfigFile{&pbApplication.ConfigFile{MountPath:"/tmp"}},
    }

    serv2 := pbApplication.Service {
        OrganizationId: orgId,
        ServiceId: "service_002",
        Name: "demo-mysql",
        Image: "mysql:5.6",
        ExposedPorts: []*pbApplication.Port{&port2},
        Labels: map[string]string { "app":"test-mysql", "component":"my-component"},
        Specs: &pbApplication.DeploySpecs{Replicas: 1},
        AppDescriptorId: "app001",
        Description: "A mysql demo",
        EnvironmentVariables: map[string]string{"MYSQL_ROOT_PASSWORD":"root"},
        Type: pbApplication.ServiceType_DOCKER,
        DeployAfter: []string{},
        Storage: []*pbApplication.Storage{&pbApplication.Storage{MountPath: "/tmp",}},
        Credentials: &credentials,
        Configs: []*pbApplication.ConfigFile{&pbApplication.ConfigFile{MountPath:"/tmp"}},
    }

    servGroup := pbApplication.ServiceGroup{
        OrganizationId: orgId,
        Description:"a service group",
        AppDescriptorId: "app001",
        Name: "group001",
        Services: []string{serv1.Name, serv2.Name},
        Policy: pbApplication.CollocationPolicy_SAME_CLUSTER,
        ServiceGroupId: "group-id",
    }

    secRule := pbApplication.SecurityRule{
        Name:"securityrule",
        AppDescriptorId: "app001",
        OrganizationId: orgId,
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
        OrganizationId: orgId,
        EnvironmentVariables: map[string]string{"var1":"var1_value", "var2":"var2_value"},
        Labels: map[string]string{"label1":"label1_value", "label2":"label2_value"},
        Services: []*pbApplication.Service{&serv1,&serv2},
        //Services: []*pbApplication.Service{&serv3},
        ConfigurationOptions: map[string]string{"conf1":"valueconf1", "conf2":"valueconf2"},
        Groups: []*pbApplication.ServiceGroup{&servGroup},
        Rules: []*pbApplication.SecurityRule{&secRule},
    }
    desc, err := appClient.AddAppDescriptor(context.Background(),&appDescriptor)
    if err != nil {
        log.Error().Msg("impossible to add descriptor")
        panic(err)
    }
    return desc
}


func getNamespace(organizationId string, appInstanceId string) string {
    target := fmt.Sprintf("%s-%s", organizationId, appInstanceId)
    // check if the namespace is larger than the allowed k8s namespace length
    if len(target) > 63 {
        return target[:63]
    }
    return target
}
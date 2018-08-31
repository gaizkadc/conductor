//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Testing for the basic orchestrator.

package orchestrator

import (
    "testing"

    log "github.com/sirupsen/logrus"

    "github.com/nalej/conductor/asm"
    entitiesConductor "github.com/nalej/conductor/entities"
    smclient "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/suite"
)

type BasicOrchestratorTestSuite struct {
    suite.Suite
    orchestrator     Orchestrator
    clusterClient    smclient.Cluster
    appsClient       smclient.Applications
    nodesClient      smclient.Node
    asmClientFactory asm.ClientFactory
}

func (suite *BasicOrchestratorTestSuite) SetupSuite() {
    suite.clusterClient = smclient.NewClusterMockup()
    suite.nodesClient = smclient.NewNodeMockup()
    suite.appsClient = smclient.NewApplicationsMockup()
    suite.asmClientFactory = asm.NewMockupClientFactory()
}

func (suite *BasicOrchestratorTestSuite) SetupTest() {
    suite.nodesClient.(*smclient.NodeMockup).InitNodeMockup()
    suite.clusterClient.(*smclient.ClusterMockup).InitMockup()

     // force clusters to be installed
    clusters, _ := suite.clusterClient.ListByNetwork("1")
    for _, c := range clusters {
        suite.clusterClient.Update(c.NetworkID, c.ID,
            *entities.NewUpdateClusterRequest().WithClusterStatus(entities.ClusterInstalled).WithType(entities.EdgeType))
    }
    suite.appsClient.(*smclient.ApplicationsMockup).AddTestNetwork("1")
    // Recreate the orchestrator in every test
    suite.orchestrator = NewBasicOrchestrator(suite.clusterClient, suite.nodesClient, suite.appsClient, suite.asmClientFactory)
}

// we Deploy three applications and we check that the target cluster
// is the correct one
func (suite *BasicOrchestratorTestSuite) TestApplicationDeployment() {
    // test descriptor
    descriptor := entities.NewAppDescriptorWithID("1",
        "toBeDeployed", "toDeploy", "description",
        "toBeDeployed", "v1.0", string(entities.EdgeType), 5555, []string {"image:version"})

    request := entitiesConductor.NewDeployAppRequest(
        descriptor.Name, descriptor.ID, descriptor.Description,
        descriptor.Label, make(map[string]string, 0), "arguments", "1Gb", entities.AppStorageDefault)


    appInstance, err := suite.orchestrator.Deploy("1", * descriptor, request)
    suite.NotNil(appInstance, "unexpected nil result")
    suite.Nil(err, "There was an error assigning the first cluster")

    suite.Nil(err, "There was an error assigning the first cluster")
    suite.Equal("1", appInstance.ClusterID, "Not correctly assigned cluster")

    // Deploy one more and check that we deploy onto the next cluster
    appInstance2, err2 := suite.orchestrator.Deploy("1", * descriptor, request)
    suite.Nil(err2, "There was an error assigning the first cluster")
    suite.Equal("2", appInstance2.ClusterID, "Not correctly assigned second cluster")

    // Deploy a third one and check that we assign to the first cluster following a complete round robin cycle
    appInstance3, err3 := suite.orchestrator.Deploy("1", * descriptor, request)
    suite.Nil(err3, "There was an error assigning the first cluster")
    suite.Equal("1", appInstance3.ClusterID, "Not correctly assigned third cluster")

}


// Deploy to a non existing network
func (suite *BasicOrchestratorTestSuite) TestApplicationDeploymentNonExistingNetwork() {
    // test descriptor
    descriptor := entities.NewAppDescriptorWithID("networkID",
        "toBeDeployed", "toDeploy", "description",
        "toBeDeployed", "v1.0", string(entities.EdgeType), 5555, []string {"image:version"})

    request := entitiesConductor.NewDeployAppRequest(
        "name", descriptor.ID, descriptor.Description,
        descriptor.Label, make(map[string]string, 0), "arguments", "1Gb", entities.AppStorageDefault)

    _, err := suite.orchestrator.Deploy("fakeNetwork", * descriptor, request)
    suite.NotNil(err, "The orchestrator has deployed an application onto a non-existing network.")
}

// Deploy when there are no available clusters.
func (suite *BasicOrchestratorTestSuite) TestApplicationDeploymentNonAvailableClusters() {
    // test descriptor
    descriptor := entities.NewAppDescriptorWithID("1",
        "toBeDeployed", "toDeploy", "description",
        "toBeDeployed", "v1.0", string(entities.EdgeType), 5555, []string {"image:version"})

    request := entitiesConductor.NewDeployAppRequest(
        "name", descriptor.ID, descriptor.Description,
        descriptor.Label, make(map[string]string, 0), "arguments", "1Gb", entities.AppStorageDefault)

    // Create an empty set of clusters
    emptyClusters := smclient.NewClusterMockup()

    suite.orchestrator = NewBasicOrchestrator(emptyClusters, suite.nodesClient, suite.appsClient, suite.asmClientFactory)

    _, err := suite.orchestrator.Deploy("1", * descriptor, request)
    suite.NotNil(err, "The orchestrator deployed an application onto network 1 without clusters")

    suite.orchestrator = NewBasicOrchestrator(emptyClusters, suite.nodesClient, suite.appsClient, suite.asmClientFactory)
    _, err2 := suite.orchestrator.Deploy("2", * descriptor, request)
    suite.NotNil(err2, "The orchestrator deployed an application onto network 2 without clusters")
}

// Undeploy an already deployed application
func (suite *BasicOrchestratorTestSuite) TestApplicationUndeploy() {
    // test descriptor
    descriptor := entities.NewAppDescriptorWithID("1",
        "toBeDeployed", "toDeploy", "description",
        "toBeDeployed", "v1.0", string(entities.EdgeType), 5555, []string {"image:version"})

    request := entitiesConductor.NewDeployAppRequest(
        descriptor.Name, descriptor.ID, descriptor.Description,
        descriptor.Label, make(map[string]string, 0), "arguments", "1Gb", entities.AppStorageDefault)

    // Check number of instances before deployment
    instancesBefore, err := suite.appsClient.ListInstances("1")

    // Deploy the application
    appInstance, err := suite.orchestrator.Deploy("1", * descriptor, request)
    suite.Nil(err, "There was an error assigning the first cluster")
    suite.Equal("1", appInstance.ClusterID, "Not correctly assigned cluster")
    instances, err := suite.appsClient.ListInstances("1")
    suite.Equal(len(instancesBefore)+1, len(instances), "The number of instances was not correctly updated")

    logger.Info(appInstance)

    // Undeploy the application
    err = suite.orchestrator.Undeploy("1", *appInstance)
    suite.Nil(err, "There is an error undeploying")

    // Check if was correctly undeployed. This means, is not there any longer
    instances, err = suite.appsClient.ListInstances("1")
    suite.Equal(len(instancesBefore), len(instances), "The instance was not correctly removed")

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBasicOrchestrator(t *testing.T) {
    log.SetLevel(log.DebugLevel)
    suite.Run(t, new(BasicOrchestratorTestSuite))
}

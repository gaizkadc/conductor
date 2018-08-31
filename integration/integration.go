//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This is an integration class used in order to test the integration of the different elements running
// in daisho with the virtualized environment.
//
// This code is not right now part of a testing environment because we cannot assure the availability of
// the necessary elements to deploy a Daisho platform.
//

package integration

import (
    smServer "github.com/daishogroup/system-model/server"
    smClient "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
    "github.com/nalej/conductor/server"
    log "github.com/sirupsen/logrus"
    "os"
    "strings"
    "strconv"
    "errors"
    entitiesConductor "github.com/nalej/conductor/entities"
    clientConductor "github.com/nalej/conductor/client"
)

const (
    SmPort             = 8800
    SmAddress          = "http://localhost:8800"
    ConductorPort      = 8900
    ConductorAddress   = "localhost"
    TestNetwork        = "network"
    TestDescription    = "description"
    TestAdminName      = "admin"
    TestAdminPhone     = "phone"
    TestEmail          = "email@email.com"
    TestClusterName    = "cluster"
    TestLocation       = "location"
    TestingUsername    = "user"
    TestingPass        = "pass"
    TestingSSH         = "key"
    TestNodeName       = "node"
    TestDescriptorName = "influxdb"
    TestService        = "influxdb"
    TestVersion        = "0.2.1"
    TestPort           = 30088
    NodeIP1            = "172.28.128.4"
    NodeIP2            = "172.28.128.5"
    NodeIP3            = "172.28.128.6"
    AppMgrIP           = "172.28.128.4"
    AppMgrPort         = 30088
    NumColonyNodes     = 3
)

// Logger for the package manager.
var logger = log.WithField("package", "integration").WithField("file", "integration.go")

// The integration structure contains the clients to be used during the integration.
type Integration struct {
    systemModel      smServer.Service
    conductorManager server.Service
    data             *IntegrationData
    configuration    Config
}

// This function creates a new integration object.
// params:
//   configuration set of variables to be used for the configuration
// return:
//   ready to use integration object
func NewIntegration(configuration Config) Integration {
    smConfiguration := smServer.Config{Port: SmPort, UseInMemoryProviders: true}
    conductorConfiguration := server.Config{
        Port:               ConductorPort,
        SystemModelAddress: SmAddress}

    return Integration{
        smServer.Service{Configuration: smConfiguration},
        server.Service{Configuration: conductorConfiguration},
        nil, configuration,
    }
}

// Structure storing entity values for testing
type IntegrationData struct {
    network entities.Network

    cluster1 entities.Cluster
    cluster2 entities.Cluster
    cluster3 entities.Cluster

    node11 entities.Node
    node21 entities.Node
    node31 entities.Node

    descriptor entities.AppDescriptor
    instance   entities.AppInstance
}

// Deploy the services needed for this integration
// return:
//   true whether all services were deployed
func (i *Integration) deployDaishoServices() bool {
    logger.Debug("Deploy daisho services...")
    // deploy system model
    err := i.systemModel.Run()
    if err != nil {
        logger.Error(err)
        return false
    }
    // deploy conductor
    err = i.conductorManager.Run()
    if err != nil {
        logger.Error(err)
        return false
    }
    logger.Info("Done")
    return true
}

// Finalize the services needed for this integration
// return:
//   true whether all services were deployed
func (i *Integration) finalizeServices() {
    logger.Debug("Finalize services...")
    // finalize system model
    i.systemModel.Finalize(true)
    // finalize conductor
    i.conductorManager.Finalize(false)
    logger.Debug("Done")
}

// Initialize the system model with some initial data.
// return:
//    true whether the content was correctly initialized or not
func (i *Integration) initializeContent() bool {
    logger.Debug("Initialize content...")
    networks := smClient.NewNetworkRest(SmAddress)
    clusters := smClient.NewClusterRest(SmAddress)
    nodes := smClient.NewNodeRest(SmAddress)
    instances := smClient.NewApplicationRest(SmAddress)

    // New network
    network, err := networks.Add(*entities.NewAddNetworkRequest(TestNetwork, TestDescription, TestAdminName, TestAdminPhone, TestEmail))
    if err != nil {
        logger.Error(err)
        return false
    }

    // New Clusters
    cluster1, err := clusters.Add(network.ID,
        *entities.NewAddClusterRequest(TestClusterName, TestDescription, entities.CloudType, TestLocation, TestEmail))
    if err != nil {
        logger.Error(err)
        return false
    }
    cluster1, err = clusters.Update(network.ID, cluster1.ID,
        *entities.NewUpdateClusterRequest().WithClusterStatus(entities.ClusterInstalled))
    if err != nil {
        logger.Error(err)
        return false
    }

    cluster2, err := clusters.Add(network.ID,
        *entities.NewAddClusterRequest(TestClusterName, TestDescription, entities.GatewayType, TestLocation, TestEmail))
    if err != nil {
        logger.Error(err)
        return false
    }
    cluster2, err = clusters.Update(network.ID, cluster2.ID,
        *entities.NewUpdateClusterRequest().WithClusterStatus(entities.ClusterInstalled))
    if err != nil {
        logger.Error(err)
        return false
    }

    cluster3, err := clusters.Add(network.ID,
        *entities.NewAddClusterRequest(TestClusterName, TestDescription, entities.EdgeType, TestLocation, TestEmail))
    if err != nil {
        logger.Error(err)
        return false
    }
    cluster3, err = clusters.Update(network.ID, cluster3.ID,
        *entities.NewUpdateClusterRequest().WithClusterStatus(entities.ClusterInstalled))
    if err != nil {
        logger.Error(err)
        return false
    }

    labels := []string{"master"}
    // New nodes
    node11, err := nodes.Add(network.ID, cluster1.ID,
        *entities.NewAddNodeRequest(TestNodeName, TestDescription, labels, NodeIP1,
            "0.0.0.0", true, TestingUsername, TestingPass, TestingSSH))
    if err != nil {
        logger.Error(err)
        return false
    }
    node21, err := nodes.Add(network.ID, cluster2.ID,
        *entities.NewAddNodeRequest(TestNodeName, TestDescription, labels, NodeIP2,
            "0.0.0.0", true, TestingUsername, TestingPass, TestingSSH))

    if err != nil {
        logger.Error(err)
        return false
    }
    node31, err := nodes.Add(network.ID, cluster3.ID,
        *entities.NewAddNodeRequest(TestNodeName, TestDescription, labels, NodeIP3,
            "0.0.0.0", true, TestingUsername, TestingPass, TestingSSH))

    if err != nil {
        logger.Error(err)
        return false
    }

    nodes.Update(network.ID, cluster1.ID, node11.ID, *entities.NewUpdateNodeRequest().WithStatus(entities.NodeReadyToInstall))
    nodes.Update(network.ID, cluster2.ID, node21.ID, *entities.NewUpdateNodeRequest().WithStatus(entities.NodeReadyToInstall))
    nodes.Update(network.ID, cluster3.ID, node31.ID, *entities.NewUpdateNodeRequest().WithStatus(entities.NodeReadyToInstall))

    // A testing descriptor
    descriptor, err := instances.AddApplicationDescriptor(network.ID,
        *entities.NewAddAppDescriptorRequest(TestDescriptorName, TestDescription,
            TestService, TestVersion, string(entities.GatewayType), TestPort, []string {"image:version"}))

    if err != nil {
        logger.Error(err)
        return false
    }

    i.data = &IntegrationData{
        descriptor: *descriptor,
        network:    *network,
        cluster1:   *cluster1,
        cluster2:   *cluster2,
        cluster3:   *cluster3,
        node11:     *node11,
        node21:     *node21,
        node31:     *node31,
        instance:   entities.AppInstance{},
    }

    return true
}

// Upload a given descriptor using an existing appmgr.
// return:
//    true whether the descriptor was correctly uploaded or not.
func (i *Integration) uploadDescriptor() bool {
    logger.Debug("Upload descriptor using appmgr...")
    // TODO this should use a programmatically designed appmgr
    descriptorPath := i.configuration.AppdevKitPath + "/bazel-out/local-fastbuild/bin/packages/influxdb/influxdb-asm-package.tar.gz"
    output, err := i.runASMCommand("package", "upload", descriptorPath)
    logger.Debugf("%s\n", output)
    if err != nil {
        os.Stderr.WriteString(err.Error())
        return false
    }
    return true
}

// Execute a deployment action
// return:
//    true if the deployment was correct, false otherwise
func (i *Integration) deploy() bool {
    client := clientConductor.NewConductorRest(ConductorAddress,8900)
    instance, err := client.Deploy(i.data.network.ID, entitiesConductor.NewDeployAppRequest(TestDescriptorName,
        i.data.descriptor.ID, TestDescription, string(entities.GatewayType), make(map[string]string, 0), "",
        "1Gb", entities.AppStorageDefault))

    logger.Debugf("The returned instance after deployment is \n%s", instance)

    if err != nil {
        logger.Error(err)
        return false
    }
    // update the instance
    i.data.instance = *instance
    return true
}

// Stop previous applications if they were previously running
// return:
//    true if the running application was stopped, false otherwise
func (i *Integration) stopRunningApplication() bool {
    output, err := i.runASMCommand("app", "stop", TestDescriptorName)

    if err != nil {
        logger.Error(err)
        return false
    }
    logger.Debug(string(output))
    return true
}

// Execute the undeployment
// return:
//    true if the application was correctly undeployed, false otherwise
func (i *Integration) undeploy() bool {
    logger.Debugf("Called undeploy %s", i.data.instance.DeployedID)
    client := clientConductor.NewConductorRest(ConductorAddress, 8900)
    err := client.Undeploy(i.data.network.ID, i.data.instance.DeployedID)
    if err != nil {
        logger.Errorf("Error undeploying application: %s", err)
        return false
    }
    logger.Debug("Successfully undeployed")
    return true
}

// This functions checks using the appmgr the running applications and print the output
// return:
//   false if there was any error calling the appmgr
func (i *Integration) checkRunningApps() bool {
    logger.Debug("----- Checking running apps with asmcli says:")
    output, err := i.runASMCommand("app", "list")
    if err != nil {
        logger.Error(err)
        return false
    }
    logger.Debug(string(output))
    return true
}

// ---- Colony-related functions

// Deploy a colony environment for testing purposes.
// return:
//    true if there was no error, false, otherwise
func (i *Integration) deployColonyEnvironment() bool {
    logger.Debug("Deploy colony environment...")
    err := i.runColonyCommandStdout("-nodes", strconv.Itoa(NumColonyNodes), "provision", "noc",
        i.configuration.FactoryTargetPath)

    if err != nil {
        logger.Error(err)
        return false
    }

    return true
}

// Destroy colonize environment
// return:
//    true if there was no error, false, otherwise
func (i *Integration) destroyColonyEnvironment() bool {
    logger.Debug("Destroy colony environment...")
    err := i.runColonyCommandStdout("destroy")

    if err != nil {
        logger.Error(err)
        return false
    }

    return true
}

func (i *Integration) checkDeployedColony() bool {
    logger.Debug("Check deployed colony...")
    output, err := i.runColonyCommand("show")
    if err != nil {
        logger.Error(err)
        return false
    }
    logger.Debug(string(output))
    if !strings.Contains(string(output), "nodes") {
        return false
    }
    lines := strings.Split(string(output), "\n")
    words := strings.Split(lines[1], ":")
    availableColonyNodes, err := strconv.Atoi(strings.Trim(strings.Trim(words[1], "\t"), " "))

    if err != nil {
        logger.Error(err.Error())
        return false
    }

    if availableColonyNodes == NumColonyNodes {
        return true
    } else {
        logger.Errorf("We found colony has deployed %s", words[1])
    }

    return false
}

// Run the integration for a given configuration
// params:
//   configuration set of variables to be used during the integration test
// return:
//   error if any
func (i *Integration) RunIntegration() error {

    logger.Debug("Started i test")

    if !i.checkDeployedColony() {
        logger.Error("Colonize has not deployed the correct number of instances")
        logger.Info("We are going to retry the deployment using colonize")
        if i.deployColonyEnvironment() {
            logger.Info("A new colony environment was deployed!")
        } else {
            logger.Error("Something went wrong deploying colonize")
            return errors.New("impossible to deploy infrastructure with colonize")
        }
    }

    ok := i.deployDaishoServices()
    if !ok {
        logger.Error("There was an error initializing services. We finalize them.")
        i.finalizeServices()
        return errors.New("there was an error initializing services. We finalize them")
    }

    ok = i.initializeContent()
    if !ok {
        logger.Error("There was an error initializing content.")
        i.finalizeServices()
        return errors.New("there was an error initializing content")
    }

    if !i.uploadDescriptor() {
        i.finalizeServices()
        return errors.New("impossible to upload a descriptor")
    }

    /*
    if !i.stopRunningApplication(){
        i.finalizeServices()
        return errors.New("impossible to stop running applications")
    }
    */

    if !i.deploy() {
        i.finalizeServices()
        return errors.New("impossible to deploy with an app")
    }

    // with the current scenario there should be one running instance
    i.checkRunningApps()

    if !i.undeploy() {
        i.finalizeServices()
        return errors.New("impossible to undeploy")
    }

    i.finalizeServices()
    i.stopRunningApplication()
    logger.Debug("finally, after running the whole process the running applications are:")
    // with the current scenario there should not be any running application
    i.checkRunningApps()

    i.destroyColonyEnvironment()
    logger.Debug("finished i test")
    return nil

}

/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package plandesigner

/*
import (
    "context"
    "fmt"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/utils"

    pbApplication "github.com/nalej/grpc-application-go"
    pbOrganization "github.com/nalej/grpc-organization-go"
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    "google.golang.org/grpc"
    "os"
)

var _ = ginkgo.Describe("Simple replica plan designer", func() {

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

        _, err := pool.AddConnection(systemModelHost, int(utils.SYSTEM_MODEL_PORT))
        gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
        // clients
        appClient = pbApplication.NewApplicationsClient(connSM)
        orgClient = pbOrganization.NewOrganizationsClient(connSM)

        localPlanDesigner = NewSimpleReplicaPlanDesigner(connHelper)
    })


    ginkgo.Context("One group with multiple replicas, one group with single replication",func(){

        ginkgo.It("Run the planner",func(){
            if !isReady {
                ginkgo.Skip("run integration test not configured")
            }

            // Create a test instance
            singleCombinedAppInstance := createCombinedReplicaApp()


            score := entities.DeploymentScore{NumEvaluatedClusters: 2,
                DeploymentsScore: []entities.ClusterDeploymentScore{
                    {
                        Scores:map[string]float32{
                            "g1": 0.99,
                            "g2": 0.99,
                            "g1g2": 0.99,
                        },
                        ClusterId: "cluster1",
                    },
                    {
                        Scores:map[string]float32{
                            "g1": 0.99,
                            "g2": 0.99,
                            "g1g2": 0.99,
                        },
                        ClusterId: "cluster2",
                    },
                },
            }
            req := entities.DeploymentRequest{
                InstanceId: singleCombinedAppInstance.AppInstanceId,
                OrganizationId: singleCombinedAppInstance.OrganizationId,
                NumRetries:0,
                TimeRetry: nil,
                RequestId: "1",
                ApplicationId: singleCombinedAppInstance.AppDescriptorId,
            }

            resultingPlan, err := localPlanDesigner.DesignPlan(singleCombinedAppInstance, score, req)

            fmt.Printf("%#v",resultingPlan)
            gomega.Expect(err).To(nil)
        })

    })
})


func createCombinedReplicaApp(appClient pbApplication.ApplicationsClient, orgClient pbOrganization.OrganizationsClient) entities.AppInstance {

    // create organization
    // add an organization
    orgRequest := pbOrganization.AddOrganizationRequest{Name: "org-001"}
    resp, err := orgClient.AddOrganization(context.Background(),&orgRequest)
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred())



    singleGroup := entities.ServiceGroupInstance{
        Name: "g1",
        OrganizationId: resp.OrganizationId,
        AppInstanceId: "app-instance-001",
        AppDescriptorId: "app-001",
        ServiceGroupInstanceId: "single-group-instance-001",
        ServiceGroupId: "single-group-001",
        ServiceInstances: []entities.ServiceInstance{
            {
                Name: "singlereplica",
                ServiceGroupId: "single-group-001",
                ServiceGroupInstanceId: "single-group-instance-001",
                AppDescriptorId: "app-001",
                AppInstanceId: "app-instance-001",
                OrganizationId: resp.OrganizationId,
                ServiceId: "singlereplica-001",
                ServiceInstanceId: "singlereplica-instance-001",
            },
        },
        Specs: entities.ServiceGroupDeploymentSpecs{
            Replicas: 1,
        },
    }

    replicaGroup := entities.ServiceGroupInstance{
        Name: "g2",
        OrganizationId: resp.OrganizationId,
        AppInstanceId: "app-instance-001",
        AppDescriptorId: "app-001",
        ServiceGroupInstanceId: "replica-group-instance-001",
        ServiceGroupId: "replica-group-001",
        ServiceInstances: []entities.ServiceInstance{
            {
                Name: "multireplica",
                ServiceGroupId: "single-group-001",
                ServiceGroupInstanceId: "single-group-instance-001",
                AppDescriptorId: "app-001",
                AppInstanceId: "app-instance-001",
                OrganizationId: resp.OrganizationId,
                ServiceId: "multireplica-001",
                ServiceInstanceId: "multireplicareplica-instance-001",
            },
        },
        Specs: entities.ServiceGroupDeploymentSpecs{
            MultiClusterReplica: true,
        },
    }

    app := entities.AppInstance{
        AppDescriptorId: "app-001",
        AppInstanceId: "app-instance-001",
        OrganizationId: resp.OrganizationId,
        Name: "combinedApp",
        Groups: []entities.ServiceGroupInstance{singleGroup, replicaGroup},
    }

    return app
}
*/

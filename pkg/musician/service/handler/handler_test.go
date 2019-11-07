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

package handler

import (
	"context"
	pbConductor "github.com/nalej/grpc-conductor-go"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/nalej/conductor/pkg/musician/scorer"
	"github.com/nalej/conductor/pkg/musician/statuscollector"
	"github.com/nalej/grpc-utils/pkg/test"
)

var _ = ginkgo.Describe("Deployment server API", func() {
	// grpc server
	var server *grpc.Server
	// conductor object
	var mgr *Manager
	// grpc test listener
	var listener *bufconn.Listener

	ginkgo.BeforeEach(func() {
		collector := statuscollector.NewFakeCollector()
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		scorerMethod := scorer.NewSimpleScorer(collector)
		mgr = NewManager(&collector, scorerMethod)
		test.LaunchServer(server, listener)
	})

	ginkgo.Context("A new score requests arrives", func() {
		var request pbConductor.ClusterScoreRequest
		var response pbConductor.ClusterScoreResponse
		var client pbConductor.MusicianClient

		ginkgo.BeforeEach(func() {
			// Register the service.
			pbConductor.RegisterMusicianServer(server, NewHandler(mgr))

			request = pbConductor.ClusterScoreRequest{
				RequestId: "myrequestId",
				Requirements: []*pbConductor.Requirement{
					{AppInstanceId: "myappinstanceid",
						Storage:                0.0,
						GroupServiceInstanceId: "mygroupserviceinstanceid",
						Replicas:               1,
						RequestId:              "myrequestid",
						Memory:                 1.0,
						Cpu:                    1.0,
					},
				},
			}
			response = pbConductor.ClusterScoreResponse{RequestId: "myrequestId", Score: []*pbConductor.DeploymentScore{
				{Score: 0.1, AppInstanceId: "myappinstanceid", GroupServiceInstances: []string{"mygroupserviceinstanceid"}},
			}}

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

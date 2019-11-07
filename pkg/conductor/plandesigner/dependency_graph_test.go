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

import (
	"github.com/nalej/conductor/internal/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Check graph dependencies", func() {

	var services []entities.Service

	ginkgo.Context("common graph", func() {
		/*
		   Create a set of services with dependencies following this graph
		   0 -> 1
		   1
		   2 -> 3, 4
		   4 -> 0
		*/
		var s0, s1, s2, s3, s4 entities.Service
		ginkgo.BeforeEach(func() {
			s0 = entities.Service{
				ServiceId:   "serv0",
				DeployAfter: []string{"serv1"},
			}
			s1 = entities.Service{
				ServiceId: "serv1",
			}
			s2 = entities.Service{
				ServiceId:   "serv2",
				DeployAfter: []string{"serv3", "serv4"},
			}
			s3 = entities.Service{
				ServiceId: "serv3",
			}
			s4 = entities.Service{
				ServiceId:   "serv4",
				DeployAfter: []string{"serv0"},
			}
			services = []entities.Service{s0, s1, s2, s3, s4}
		})

		ginkgo.It("Build the graph", func() {
			g := NewDependencyGraph(services)
			gomega.Expect(g).NotTo(gomega.BeNil())
			gomega.Expect(g.NumDependencies()).To(gomega.Equal(1))
			gomega.Expect(g.NumServices()).To(gomega.Equal(5))
		})
		/*
		   ginkgo.It("Compute group topological order", func(){
		       g := NewDependencyGraph(services)
		       order, err := g.GetDependencyOrderByGroups()
		       gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		       log.Info().Msgf("order %v",order)
		       expected := [][]entities.Service{
		           []entities.Service{s1,s3},
		           []entities.Service{s0},
		           []entities.Service{s4},
		           []entities.Service{s2}}
		       // Check everything is ok
		       for i,_ := range expected {
		           //gomega.Expect(len(expected[i])).To(gomega.Equal(len(order[i])))
		           for j, _ := range expected[i] {
		               gomega.Expect(expected[i][j].Name).To(gomega.Equal(order[i][j].Name))
		           }
		       }
		       gomega.Expect(order).To(gomega.Equal(expected))
		   })
		*/
	})

})

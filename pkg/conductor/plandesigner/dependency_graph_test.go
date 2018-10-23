/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package plandesigner

import (
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    "github.com/nalej/conductor/internal/entities"
    "github.com/rs/zerolog/log"
)

var  _ = ginkgo.Describe("Check graph dependencies" , func(){

    var services [] entities.Service

     ginkgo.Context("common graph", func() {
         /*
         Create a set of services with dependencies following this graph
         0 -> 1
         1
         2 -> 3, 4
         4 -> 0
         */
         ginkgo.BeforeEach(func(){
             s0 := entities.Service{
                 ServiceId: "serv0",
                 DeployAfter: []string{"serv1"},
             }
             s1 := entities.Service{
                 ServiceId: "serv1",
             }
             s2 := entities.Service{
                 ServiceId: "serv2",
                 DeployAfter: []string{"serv3","serv4"},
             }
             s3 := entities.Service{
                 ServiceId: "serv3",
             }
             s4 := entities.Service{
                 ServiceId: "serv4",
                 DeployAfter: []string{"serv0"},
             }
             services = [] entities.Service{s0, s1, s2, s3, s4}
         })

         ginkgo.It("Build the graph", func(){
            g := NewDependencyGraph(services)
            gomega.Expect(g).NotTo(gomega.BeNil())
            gomega.Expect(g.NumDependencies()).To(gomega.Equal(4))
            gomega.Expect(g.NumServices()).To(gomega.Equal(5))
         })

         ginkgo.It("Compute group topological order", func(){
             g := NewDependencyGraph(services)
             order, err := g.GetDependencyOrderByGroups()
             gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
             log.Info().Msgf("order %v",order)
             expected := [][]string{
                 []string{"serv1","serv3"},
                 []string{"serv0"},
                 []string{"serv4"},
                 []string{"serv2"}}
             gomega.Expect(order).To(gomega.Equal(expected))

         })
     })


})



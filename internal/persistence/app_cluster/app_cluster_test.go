/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package app_cluster

import (
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/provider"
    "github.com/nalej/conductor/pkg/provider/kv"
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    "os"
)

var _ = ginkgo.Describe("application cluster data persistence test", func(){

    var db *AppClusterDB
    var localDB provider.KeyValueProvider
    dbPath := "/tmp/app_cluster_persistence_test.db"

    ginkgo.BeforeEach(func(){
        // create a kv provider
        aux, errDB := kv.NewLocalDB(dbPath)
        gomega.Expect(errDB).ToNot(gomega.HaveOccurred())

        localDB = aux
        db = NewAppClusterDB(localDB)
    })

    ginkgo.AfterEach(func(){
        errClose := localDB.Close()
        gomega.Expect(errClose).ToNot(gomega.HaveOccurred())

        err := os.Remove(dbPath)
        gomega.Expect(err).ToNot(gomega.HaveOccurred())
    })


    ginkgo.It("add and retrieve a deployment fragment", func(){
        toAdd := entities.DeploymentFragment{
            ClusterId: "cluster1",
            OrganizationId: "someorg",
            AppName: "testApp",
        }
        // Add it
        errAdd := db.AddDeploymentFragment(&toAdd)
        gomega.Expect(errAdd).ToNot(gomega.HaveOccurred())
        // Get it
        retrieved, errGet := db.GetDeploymentFragment(toAdd.ClusterId)
        gomega.Expect(errGet).ToNot(gomega.HaveOccurred())
        gomega.Expect(retrieved).To(gomega.Equal(retrieved))
    })
})

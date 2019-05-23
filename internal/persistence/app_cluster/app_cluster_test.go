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
            DeploymentId: "deployment1",
            AppInstanceId: "myappinstance1",
            OrganizationId: "someorg",
            AppName: "testApp",
            FragmentId: "fragment1",

        }
        // Add it
        errAdd := db.AddDeploymentFragment(&toAdd)
        gomega.Expect(errAdd).ToNot(gomega.HaveOccurred())
        // Get it
        retrieved, errGet := db.GetDeploymentFragment(toAdd.ClusterId,toAdd.AppInstanceId)
        gomega.Expect(errGet).ToNot(gomega.HaveOccurred())
        gomega.Expect(retrieved).To(gomega.Equal(retrieved))
    })

    ginkgo.It("add, delete and try to retrieve a fragment", func(){
        toAdd := entities.DeploymentFragment{
            ClusterId: "cluster1",
            DeploymentId: "deployment1",
            AppInstanceId: "myappinstance1",
            OrganizationId: "someorg",
            AppName: "testApp",
            FragmentId: "fragmen1",

        }
        // Add it
        errAdd := db.AddDeploymentFragment(&toAdd)
        gomega.Expect(errAdd).ToNot(gomega.HaveOccurred())
        // Delete it
        errDel := db.DeleteDeploymentFragment(toAdd.ClusterId, toAdd.AppInstanceId)
        gomega.Expect(errDel).ToNot(gomega.HaveOccurred())
        // retrieve it and it must not be there
        retrieved, errGet := db.GetDeploymentFragment(toAdd.ClusterId,toAdd.AppInstanceId)
        gomega.Expect(errGet).ToNot(gomega.HaveOccurred())
        gomega.Expect(retrieved).To(gomega.BeNil())
    })

    ginkgo.It("get all the entries stored for a cluster", func(){
        toAdd1 := entities.DeploymentFragment{
            ClusterId: "cluster1",
            DeploymentId: "deployment1",
            AppInstanceId: "myappinstance1",
            OrganizationId: "someorg",
            AppName: "testApp",
            FragmentId: "fragment1",
        }
        // Add it
        errAdd := db.AddDeploymentFragment(&toAdd1)
        gomega.Expect(errAdd).ToNot(gomega.HaveOccurred())

        toAdd2 := entities.DeploymentFragment{
            ClusterId: "cluster1",
            DeploymentId: "deployment2",
            AppInstanceId: "myappinstance2",
            OrganizationId: "someorg",
            AppName: "testApp",
            FragmentId: "fragment2",
        }
        // Add it
        errAdd = db.AddDeploymentFragment(&toAdd2)
        gomega.Expect(errAdd).ToNot(gomega.HaveOccurred())

        // Get all the bucket data
        pairs, bErr := db.GetFragmentsInCluster(toAdd1.ClusterId)
        gomega.Expect(bErr).ToNot(gomega.HaveOccurred())
        gomega.Expect(len(pairs)).To(gomega.Equal(2))
    })

})

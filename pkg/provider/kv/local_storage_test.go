/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package kv

import (
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    "os"
)

var _ = ginkgo.Describe("test local database using bolt", func(){

    dbPath := "/tmp/localstorage.db"

    var db *LocalDB

    ginkgo.BeforeEach(func(){
        // Get a file path
        var err error
        db, err  = NewLocalDB(dbPath)
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
    })

    ginkgo.AfterEach(func(){
        // close database
        err := db.Close()
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        errR := os.Remove(dbPath)
        gomega.Expect(errR).ToNot(gomega.HaveOccurred())
    })

    ginkgo.It("put and get a new value", func(){
        newKey := []byte("testkey")
        newValue := []byte("newValue")
        bucket := []byte("testbucket")

        errPut := db.Put(bucket, newKey, newValue)
        gomega.Expect(errPut).NotTo(gomega.HaveOccurred())

        retrieved, errGet := db.Get(bucket, newKey)
        gomega.Expect(errGet).NotTo(gomega.HaveOccurred())
        gomega.Expect(retrieved).Should(gomega.Equal(newValue))
    })

    ginkgo.It("fails when getting a non-existing bucket", func(){
        retrieved, errGet := db.Get([]byte("bucket"),[]byte("notthere"))
        gomega.Expect(errGet).To(gomega.HaveOccurred())
        gomega.Expect(retrieved).Should(gomega.BeNil())

    })

    ginkgo.It("returns nil when requesting a non-existing key", func(){
        newKey := []byte("testkey")
        newValue := []byte("newValue")
        bucket := []byte("testbucket")

        errPut := db.Put(bucket, newKey, newValue)
        gomega.Expect(errPut).NotTo(gomega.HaveOccurred())

        retrieved, errGet := db.Get(bucket, []byte("notexisting"))
        gomega.Expect(errGet).ToNot(gomega.HaveOccurred())
        gomega.Expect(retrieved).Should(gomega.BeNil())

    })
})

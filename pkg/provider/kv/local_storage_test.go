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

package kv

import (
	"github.com/nalej/conductor/pkg/provider"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"os"
)

var _ = ginkgo.Describe("test local database using bolt", func() {

	dbPath := "/tmp/localstorage.db"

	var db provider.KeyValueProvider

	ginkgo.BeforeEach(func() {
		// Get a file path
		aux, err := NewLocalDB(dbPath)
		db = aux
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterEach(func() {
		// close database
		err := db.Close()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		errR := os.Remove(dbPath)
		gomega.Expect(errR).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("put and get a new value", func() {
		newKey := []byte("testkey")
		newValue := []byte("newValue")
		bucket := []byte("testbucket")

		errPut := db.Put(bucket, newKey, newValue)
		gomega.Expect(errPut).NotTo(gomega.HaveOccurred())

		retrieved, errGet := db.Get(bucket, newKey)
		gomega.Expect(errGet).NotTo(gomega.HaveOccurred())
		gomega.Expect(retrieved).Should(gomega.Equal(newValue))
	})

	ginkgo.It("fails when getting a non-existing bucket", func() {
		retrieved, errGet := db.Get([]byte("bucket"), []byte("notthere"))
		gomega.Expect(errGet).To(gomega.HaveOccurred())
		gomega.Expect(retrieved).Should(gomega.BeNil())

	})

	ginkgo.It("returns nil when requesting a non-existing key", func() {
		newKey := []byte("testkey")
		newValue := []byte("newValue")
		bucket := []byte("testbucket")

		errPut := db.Put(bucket, newKey, newValue)
		gomega.Expect(errPut).NotTo(gomega.HaveOccurred())

		retrieved, errGet := db.Get(bucket, []byte("notexisting"))
		gomega.Expect(errGet).ToNot(gomega.HaveOccurred())
		gomega.Expect(retrieved).Should(gomega.BeNil())

	})

	ginkgo.It("returns the list of bucket names", func() {
		k := []byte("key")
		v := []byte("value")
		buckets := [][]byte{[]byte("bucket1"), []byte("bucket2"), []byte("bucket3")}
		for _, b := range buckets {
			errPut := db.Put(b, k, v)
			gomega.Expect(errPut).NotTo(gomega.HaveOccurred())
		}
		// get the list of buckets
		retrievedBuckets := db.GetBuckets()
		gomega.Expect(len(retrievedBuckets)).To(gomega.Equal(len(buckets)))

	})
})

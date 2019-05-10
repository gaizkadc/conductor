/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package app_cluster

import (
    "testing"

    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
)

func TestAppClusterTest(t *testing.T) {
    gomega.RegisterFailHandler(ginkgo.Fail)
    ginkgo.RunSpecs(t, "Conductor app cluster storage Suite")
}
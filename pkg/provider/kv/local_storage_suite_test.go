/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package kv

import (
    "testing"

    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
)

func TestApiTest(t *testing.T) {
    gomega.RegisterFailHandler(ginkgo.Fail)
    ginkgo.RunSpecs(t, "Conductor ApiTest Suite")
}
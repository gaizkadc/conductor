/*
 * Copyright (C) 2018 Nalej Group -All Rights Reserved
 */


package scorer

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestApiTest(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Conductor ApiTest Suite")
}

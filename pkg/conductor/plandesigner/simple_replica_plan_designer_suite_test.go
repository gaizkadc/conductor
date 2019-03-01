/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package plandesigner


import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestApiTest(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Simple replica plan designer suite")
}

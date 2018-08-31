load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(
    name = "gazelle",
    prefix = "github.com/daishogroup/dhttp",
)

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "client_mockup.go",
        "client_sling.go",
        "client_test_handler.go",
        "client_test_helper.go",
        "client_testsuite.go",
        "config.go",
        "operation.go",
        "utils.go",
    ],
    importpath = "github.com/daishogroup/dhttp",
    visibility = ["//visibility:public"],
    deps = [
        "//errors:go_default_library",
        "//vendor/github.com/daishogroup/derrors:go_default_library",
        "//vendor/github.com/dghubble/sling:go_default_library",
        "//vendor/github.com/gorilla/mux:go_default_library",
        "//vendor/github.com/sirupsen/logrus:go_default_library",
        "//vendor/github.com/stretchr/testify/suite:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "client_mockup_test.go",
        "client_sling_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/daishogroup/derrors:go_default_library",
        "//vendor/github.com/stretchr/testify/suite:go_default_library",
    ],
)

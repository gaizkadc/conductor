# gazelle:ignore

load("@bazel_gazelle//:def.bzl", "gazelle")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

gazelle(
    name = "gazelle",
    build_tags = [
        "darwin",
    ],
    command = "fix",
    external = "vendored",
    prefix = "github.com/daishogroup/system-model",
)


go_binary(
    name = "system-model",
    embed = ["//cmd/system-model:go_default_library"],
    visibility = ["//visibility:public"],
    pure="on",
)

go_binary(
    name = "system-model-cli",
    embed = ["//cmd/system_model_cli:go_default_library"],
    visibility = ["//visibility:public"],
    pure="on",
)

pkg_tar(
    name = "core-system-model",
    srcs = [
        "//components/system-model:component.yaml",
        "//components/system-model:image.tar",
    ],
    extension = "tar.gz",
    mode = "0644",
    package_dir = "/",
)

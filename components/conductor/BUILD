package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build", "docker_bundle")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "conductor",  #keep
    library = "//cmd/conductor:go_default_library",
    # Docker image doesn't have some of the needed libraries, make an
    # all-static binary for in the image.
    pure = "on",
    visibility = ["//visibility:private"],
)

docker_build(
    # TODO: Tag with version somehow
    name = "conductor-image",
    base = "@debian-tall//image",
    directory = "/daisho",
    entrypoint = [
        "/usr/bin/dumb-init",
        "--",
        "/daisho/conductor",
        "api",
    ],
    env = {
        # debian-tall has no /tmp
        "TMPDIR": "/var/run",
    },
    files = [
        ":conductor",
    ],
    mode = "0755",
)

docker_bundle(
    name = "image",
    images = {
        "daisho/conductor:v1.0.1": ":conductor-image",
    },
)

# Export this to pick it up in the final packaging step
exports_files(["component.yaml"])

package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build", "docker_bundle")

docker_build(
    # TODO: Tag with version somehow
    name = "system-model-image",
    base = "@debian-tall//image",
    directory = "/daisho",
    entrypoint = [
        "/usr/bin/dumb-init",
        "--",
        "/daisho/system-model",
        "api",
    ],
    files = [
        "//:system-model",
    ],
    mode = "0755",
)

docker_bundle(
    name = "image",
    images = {
        "daisho/system-model:v1.1.0": ":system-model-image",
    },
)

# Export this to pick it up in the final packaging step
exports_files(["component.yaml"])

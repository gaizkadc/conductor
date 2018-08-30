http_archive(
    name = "bazel_gazelle",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.9/bazel-gazelle-0.9.tar.gz",
    sha256 = "0103991d994db55b3b5d7b06336f8ae355739635e0c2379dea16b8213ea5a223",
)

git_repository(
    name = "io_bazel_rules_go",
    tag = "0.8.1",
    remote = "https://github.com/bazelbuild/rules_go.git",
)

load("@io_bazel_rules_go//go:def.bzl", "go_repositories")

go_repositories()

git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.3.0",
)

load(
  "@io_bazel_rules_docker//docker:docker.bzl",
  "docker_repositories", "docker_pull"
)
 
docker_repositories()
 
# We will use this as a base image (Busybox-based, very small,
# perfect for stand-alone binaries). See
# https://github.com/gravitational/docker-debian
docker_pull(
    name = "debian-tall",
    registry = "quay.io",
    repository = "gravitational/debian-tall",
    tag = "0.0.1",
)


load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()
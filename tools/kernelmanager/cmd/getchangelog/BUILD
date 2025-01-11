load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "getchangelog_lib",
    srcs = ["main.go"],
    importpath = "go.resf.org/peridot/tools/kernelmanager/cmd/getchangelog",
    visibility = ["//visibility:private"],
    deps = [
        "//tools/kernelmanager/packager/v1:packager",
        "//vendor/github.com/pkg/errors",
    ],
)

go_binary(
    name = "getchangelog",
    embed = [":getchangelog_lib"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "getchangelog_test",
    srcs = ["main_test.go"],
    data = glob(["testdata/**"]),
    embed = [":getchangelog_lib"],
    deps = ["//vendor/github.com/stretchr/testify/require"],
)

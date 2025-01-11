load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "updatekernel_lib",
    srcs = ["main.go"],
    importpath = "go.resf.org/peridot/tools/kernelmanager/cmd/updatekernel",
    visibility = ["//visibility:private"],
    deps = [
        "//tools/kernelmanager/packager",
        "//tools/kernelmanager/packager/kernelorg",
        "//tools/kernelmanager/packager/v1:packager",
        "//vendor/github.com/go-git/go-billy/v5/osfs",
    ],
)

go_binary(
    name = "updatekernel",
    embed = [":updatekernel_lib"],
    visibility = ["//visibility:public"],
)

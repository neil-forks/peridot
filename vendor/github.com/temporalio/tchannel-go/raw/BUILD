load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "raw",
    srcs = [
        "call.go",
        "handler.go",
    ],
    importmap = "go.resf.org/peridot/vendor/github.com/temporalio/tchannel-go/raw",
    importpath = "github.com/temporalio/tchannel-go/raw",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/temporalio/tchannel-go",
        "@org_golang_x_net//context",
    ],
)

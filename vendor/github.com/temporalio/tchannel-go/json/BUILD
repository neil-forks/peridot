load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "json",
    srcs = [
        "call.go",
        "context.go",
        "handler.go",
    ],
    importmap = "go.resf.org/peridot/vendor/github.com/temporalio/tchannel-go/json",
    importpath = "github.com/temporalio/tchannel-go/json",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/opentracing/opentracing-go",
        "//vendor/github.com/temporalio/tchannel-go",
        "@org_golang_x_net//context",
    ],
)

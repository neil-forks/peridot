load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "relay",
    srcs = ["relay.go"],
    importmap = "go.resf.org/peridot/vendor/github.com/temporalio/tchannel-go/relay",
    importpath = "github.com/temporalio/tchannel-go/relay",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/temporalio/tchannel-go/thrift/arg2"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "arg2",
    srcs = ["kv_iterator.go"],
    importmap = "go.resf.org/peridot/vendor/github.com/temporalio/tchannel-go/thrift/arg2",
    importpath = "github.com/temporalio/tchannel-go/thrift/arg2",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/temporalio/tchannel-go/typed"],
)

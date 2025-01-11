load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "meta",
    srcs = [
        "GoUnusedProtection__.go",
        "meta.go",
        "meta-consts.go",
    ],
    importmap = "go.resf.org/peridot/vendor/github.com/temporalio/tchannel-go/thrift/gen-go/meta",
    importpath = "github.com/temporalio/tchannel-go/thrift/gen-go/meta",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/apache/thrift/lib/go/thrift"],
)

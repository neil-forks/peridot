load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "thrift",
    srcs = [
        "client.go",
        "context.go",
        "doc.go",
        "headers.go",
        "interfaces.go",
        "meta.go",
        "options.go",
        "server.go",
        "struct.go",
        "tchan-meta.go",
        "transport.go",
    ],
    importmap = "go.resf.org/peridot/vendor/github.com/temporalio/tchannel-go/thrift",
    importpath = "github.com/temporalio/tchannel-go/thrift",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/apache/thrift/lib/go/thrift",
        "//vendor/github.com/temporalio/tchannel-go",
        "//vendor/github.com/temporalio/tchannel-go/internal/argreader",
        "//vendor/github.com/temporalio/tchannel-go/thrift/gen-go/meta",
        "//vendor/github.com/temporalio/tchannel-go/typed",
        "@org_golang_x_net//context",
    ],
)

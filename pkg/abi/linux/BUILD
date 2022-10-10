load("//tools:defs.bzl", "arch_genrule", "go_library", "go_test")
load("//tools/nogo:defs.bzl", "nogo_facts")

# Package linux contains the constants and types needed to interface with a
# Linux kernel. It should be used instead of syscall or golang.org/x/sys/unix
# when the host OS may not be Linux.

package(licenses = ["notice"])

nogo_facts(
    name = "iouring_offsets",
    srcs = ["iouring.go"],
    output = "iouring_offsets_impl.go",
    template = "iouring_offsets.tmpl",
)

arch_genrule(
    name = "iouring_offsets_arch",
    src = ":iouring_offsets",
    template = "iouring_offsets_%s.go",
)

go_library(
    name = "linux",
    srcs = [
        "audit.go",
        "bpf.go",
        "futex.go",
        "prctl.go",
        "seccomp.go",
    ],
    marshal = True,
    visibility = ["//visibility:public"],
)

go_test(
    name = "linux_test",
    size = "small",
    srcs = [
        "netfilter_test.go",
    ],
    library = ":linux",
)

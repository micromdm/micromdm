# Contributing
For now, reference the contributing [README](https://github.com/micromdm/micromdm/blob/master/CONTRIBUTING.md).

# Troubleshooting
This section will assist in troubleshooting issues that arise with your environment or Go tooling in the course of building micromdm.

## Common Errors

### Incorrect Go Version
```
$ go build
# github.com/micromdm/micromdm
./serve.go:232: unknown "net/http".Server field 'ReadHeaderTimeout' in struct literal
./serve.go:233: unknown "net/http".Server field 'IdleTimeout' in struct literal
./serve.go:275: srvURL.Hostname undefined (type *url.URL has no field or method Hostname)
./serve.go:333: undefined: tls.X25519
./serve.go:339: undefined: tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305
./serve.go:340: undefined: tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305
```

**FIX:** Install Go `1.8` from the [Go downloads page](https://golang.org/dl/).

### protobuf development

MicroMDM uses protobufs, or [protocol buffers](https://developers.google.com/protocol-buffers/) for a number of internal communication channels. To be able to work with them you need both the `protoc` binary and the protobuf go packages and binaries. Note you should not need to install the protobuf tools _just_ to compile and run the source, only for actually working with the protobuf definitions.

Perhaps the easiest way to get these tools on a Mac is to have [Homebrew](https://brew.sh/) install the dependencies for you, then build protobuf from source:

Install proto3 from source:

```shel
brew install autoconf automake libtool
git clone https://github.com/google/protobuf && cd protobuf
./autogen.sh && ./configure && make && make install
```

Update protoc Go bindings via:

```shell
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

(Above borrowed from [go-kit](https://raw.githubusercontent.com/go-kit/kit/master/examples/addsvc/pb/compile.sh))

Once the tools are installed then you use the `go generate` command in the directory with the `*.proto` files (which should subsequently (re-)generate the `*.pb.go` files).

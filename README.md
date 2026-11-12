# go-pfcp

[![CI](https://github.com/gomaja-labs/go-pfcp/actions/workflows/ci.yml/badge.svg)](https://github.com/gomaja-labs/go-pfcp/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/gomaja-labs/go-pfcp.svg)](https://pkg.go.dev/github.com/gomaja-labs/go-pfcp)

A Go codec for PFCP (Packet Forwarding Control Protocol), using 3GPP TS 29.244
V20.1.0. It encodes and decodes messages and Information Elements; applications
provide networking, association/session state and procedural validation.

## Install

Requires Go 1.26 or later. Consumers track `main`; there are no tagged releases
and APIs may change without compatibility aliases.

```sh
go get github.com/gomaja-labs/go-pfcp@main
```

## Use

```go
request := message.NewHeartbeatRequest(1, ie.NewRecoveryTimeStamp(started))
wire, err := request.Marshal()
if err != nil {
    return err
}
decoded, err := message.ParseHeartbeatRequest(wire)
if err != nil {
    return err
}
_ = decoded
```

Here `started` is the PFCP entity's start time (`time.Time`); import
`github.com/gomaja-labs/go-pfcp/ie` and `github.com/gomaja-labs/go-pfcp/message`.

See the [package overview](https://pkg.go.dev/github.com/gomaja-labs/go-pfcp),
[IE guide and examples](https://pkg.go.dev/github.com/gomaja-labs/go-pfcp/ie), and
[message guide and examples](https://pkg.go.dev/github.com/gomaja-labs/go-pfcp/message)
for grouped fields, errors, null-length IEs, version handling and bundles.
The [heartbeat commands](examples/heartbeat) demonstrate UDP transport.
For GTP, see [go-gtp](https://github.com/gomaja-labs/go-gtp).

## Development

The module supports Go 1.26 and later. Format and generate with the current
stable Go release, as CI does: gofmt column alignment changed in Go 1.27, so
an older toolchain reformats some files. The checks CI runs:

```sh
go install honnef.co/go/tools/cmd/staticcheck@v0.8.2
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0

gofmt -l .
go mod tidy -diff
go generate ./... && git diff --exit-code
go vet ./...
GOARCH=386 go vet ./...
go test -race -count=1 ./...
staticcheck ./...
golangci-lint run ./...
```

Generation must leave a committed tree unchanged. Tests include executable
examples and an exported-identifier documentation check. See
[internal/spec/README.md](internal/spec/README.md) for importing a new
specification version and regenerating its catalogue and codecs.

## Author(s)

The go-pfcp authors. See [LICENSE](LICENSE) for the copyright holders.

## License

[MIT](LICENSE).

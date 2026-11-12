# Heartbeat over UDP

Run the server from the repository root:

```sh
go run ./examples/heartbeat/hb-server -s 127.0.0.1:8805
```

In a second terminal, send one request:

```sh
go run ./examples/heartbeat/hb-client -s 127.0.0.1:8805
```

The client checks the response type, sequence and Recovery Time Stamp, and has
a three-second receive deadline. The server keeps running until interrupted.
Both commands accept `-s`; its default is `127.0.0.2:8805`.

The server reports its start time in each response and copies the request's
sequence number (TS 29.244 V20.1.0 clauses 8.2.65 and 6.4). The client sends one
request without retries. These commands demonstrate application-owned UDP
transport; they do not establish an association or manage PFCP sessions.

# go-fastcounters

This demonstrates a very simple Go HTTP server that can:

- process 200k TPS at ~300microseconds latency
- with each message being persisted to disk before being processed
- all transactions processed serially and therefore atomically

## Run the benchmark

```bash
go run .
wrk http://127.0.0.1:8080/ --latency -t8 -c64 -d60s
```

## The Code

The "framework" code (i.e. the code that writes to disk and orders all messages) is in `wal/wal.go`

The other files are:

- `main.go` - starts up the WAL and the http server
- `counter_state.go` - the processing of messages to update state
- `counter_routes.go` - the http endpoints

Similarly the state and routes files are available for the chess server as well.

## The approach

By bringing all messages back to a single thread, they can be totally ordered and persisted to disk in batches.
Then they are sent in this order to the message handlers that update state (also on a single thread).

Inspired by the (LMAX Architecture)[https://martinfowler.com/articles/lmax.html]

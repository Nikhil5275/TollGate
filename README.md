# TollGate

A Go HTTP rate-limiter prototype with consistent-hash routing across five Docker nodes. Each client belongs to one server, which tracks a synchronized token bucket at **100 requests/second** with a **150-request burst**.

## Run

```sh
docker compose up --build
curl -i -H 'X-Client-ID: alice' http://localhost:8081/
```

The five entry points are localhost ports 8081–8085. Any entry point routes a client to the same owner. `/healthz` bypasses rate limiting. Accepted requests return HTTP 200; exhausted buckets return HTTP 429 with a whole-second `Retry-After` header.

For a single local instance with Go 1.22 or newer:

```sh
go run ./cmd/tollgate
```

## Design

```text
Client → any HTTP node → consistent hash (128 virtual nodes/server)
                      → owning node → synchronized token bucket → 200 / 429
```

`internal/ring` provides deterministic ownership, `internal/limiter` provides counters and token refill, and `cmd/tollgate` handles routing and HTTP responses. All servers must share the same `PEERS` list and use their matching `NODE_URL`.

## Tests and benchmarks

```sh
make test
make bench
sh scripts/benchmark.sh
```

The tests cover burst acceptance, rejection, independent clients, concurrent access, and stable ownership. The wrk script exercises rotating clients and a fixed client under concurrent load. Install wrk separately before running it.

**No measured throughput or latency results are included.** A 5,000+ requests/sec throughput and under 2 ms p99 added latency are benchmark targets, not verified results. To measure added latency, compare equivalent baseline and limiter deployments under the same workload; wrk's end-to-end latency alone is not added limiter latency. Capture hardware, concurrency, response codes, and client-key distribution with any published results.

## Prototype scope

State lives in memory and resets on restart. There is no replication, automatic membership discovery, eviction of inactive client buckets, or failover. Peer failure returns a proxy error. Client IDs are caller supplied and are intended for trusted/demo traffic; a production gateway should derive keys from authenticated identities. The handler is a demonstration endpoint, rather than an application upstream proxy. Cluster members must have identical configuration to avoid routing loops.

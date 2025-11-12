# driftq-cli

> **Public** — developer CLI for DriftQ.

A small, dependency-free Go CLI to manage the DriftQ broker and do quick produce/consume tests.  
Initial implementation calls **stubs** so you can wire the real broker/admin APIs later without breaking the UX.

## Install (local build)
```bash
git clone https://github.com/BehnamAxo/driftq-cli.git
cd driftq-cli
make install    # installs 'driftq' into your $GOBIN
```

## Usage
```
driftq health [--addr HOST:PORT]
driftq topic create <name> [--partitions N] [--compacted] [--addr HOST:PORT]
driftq topic list [--addr HOST:PORT]
driftq produce <topic> [--key K] [--value V|--file PATH] [--addr HOST:PORT]
driftq consume <topic> --group G [--from latest|earliest] [--addr HOST:PORT]
driftq lag <topic> [--group G] [--addr HOST:PORT]
```

Examples:
```bash
driftq health
driftq topic create orders --partitions 3
driftq topic list
driftq produce orders --key id-1 --value 'hello'
driftq consume orders --group demo --from latest
driftq lag orders --group demo
```

## Design notes
- Zero third‑party deps; easy to vendor.
- Real broker integration lands by swapping methods in `pkg/driftqclient` with gRPC/HTTP calls.
- Flags are stable; help text won’t change when internals do.

## Roadmap
- Wire to admin/produce/consume RPCs
- Add `driftq offsets` and `driftq group describe`
- Output JSON via `--json` when needed

## License
Apache 2.0

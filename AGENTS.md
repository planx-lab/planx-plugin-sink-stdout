# planx-plugin-sink-stdout

## Collaboration and references

Use [workspace guidance](../AGENTS.md) for task completion, authorization and
proportionate checks. Read relevant clauses of the [canonical contract](../planx-spec/AI_CONTRACT.md),
[architecture](../planx-spec/planx-architecture.md) and accepted
[ADR-017](../planx-spec/adr/017-builtin-typed-data-integration.md) when their
subject changes. Use [repo.lock](repo.lock) for source ownership. Do not reread
the entire specification for an unrelated edit.

## Plugin contract

- Implement business logic through `planx-sdk-go/sdk` SPI. Keep this protected external regression implementation under `internal/plugin/`.
- One self-describing binary declares components and calls `sdk.Serve` from
  `cmd/plugin/main.go`; no YAML manifest or runtime logic in the entry point.
- No plugin-owned gRPC server, session, flow control, concurrency or backpressure.
  Do not import Engine, Proto directly or SDK internals.
- Keep business operations synchronous, deterministic and batch-oriented. Use
  approved public typed-data/SPI APIs; transport encoding belongs to the SDK.
- The stdout sink’s existing business output is intentional; preserve its tested formatting.
- Runtime work belongs in its owning SDK/Engine layer. If the task authorizes
  that cross-repository work, make the coordinated change there; otherwise pause
  that addition and complete the plugin work already in scope.
- Use `go test ./...` and `go build ./...` for plugin changes, and the protected
  external regression gate when a shared protocol or SPI changes.

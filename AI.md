# AI RULES — PLANX GO PLUGIN (v4)

## Authority Documents

1. [planx-architecture.md](../planx-architecture.md)
2. [planx-ai-guardrails.md](../planx-ai-guardrails.md)
3. [AI_CONTRACT.md](../AI_CONTRACT.md)
4. [planx-sdk-go/AI.md](../planx-sdk-go/AI.md) — SDK semantics

---

## SCOPE
This repository implements ONLY plugin business logic.

---

## PLUGIN HARD RULES

AI MUST NOT:
- Implement gRPC servers
- Manage sessions or flow control
- Start goroutines for concurrency
- Import planx-engine
- Import planx-proto directly
- Read from STDIN or write to STDOUT (except logging)

AI MUST:
- Implement EXACTLY ONE plugin type (Source OR Sink OR Processor)
- Implement SPI interfaces from planx-sdk-go/sdk
- Keep logic synchronous and deterministic
- Place all business logic under `internal/plugin/`
- Treat Batch as opaque bytes
- Ensure `main.go` only calls one SDK function: `sdk.ServeSource`, `sdk.ServeSink`, or `sdk.ServeProcessor`

If a requirement seems to need runtime logic:
STOP. That belongs to SDK.

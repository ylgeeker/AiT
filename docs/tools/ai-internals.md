# ai internals (historical)

> The standalone `ai` binary and `internal/opsai` were **removed**. Agent admin and Loop now live under **`ait`**.

## Current call chain

```text
cmd/ait/main.go                → minimal main; delegates to ait.Run
internal/ait/router            → bare ait / ait agent → Loop; ait agent <admin> → agent
internal/ait/loop              → local Agent Loop (agentcore.RunLoop)
internal/ait/dispatch          → top-level plugin exec (tooldispatch)
internal/ait/agent             → admin CLI: connect, model, provider, ollama, config
internal/pkg/agentcore         → Loop kernel (risk, approve, runner, RunLoop)
internal/pkg/agentops          → Ops toolbox DispatchToolbox + context snapshot
internal/pkg/agentwire         → LocalLoopConfig / RemoteSessionConfig SSOT
internal/pkg/tooldispatch      → plugin Lookup+ExecPath (ait dispatch + agentops)
internal/pkg/ai                → Provider routing / chat facade
internal/pkg/config            → ai.* keys (ValidateAI via ait agent config; SSH keys via aissh config)
```

## Layering

- `internal/ait/agent` imports `agentcore` / `agentops` / `agentui` / `ai` / `agentwire` (**must not** import `internal/ssh` for execution).
- `internal/ait` must not import product app packages (`internal/ssh`, `internal/log`, …) — plugins run via `dispatch` / `tooldispatch`.
- `agentcore` must not import product packages — `TestAgentcore_NoProductImports` (AC-V4-1).
- `aissh` must not import `agentcore` / `internal/ait/agent` (OA24).

## Interactive Loop

| Entry | Behavior |
|-------|----------|
| `ait` / `ait agent` | Multi-turn Agent Loop (`agentcore.RunLoop`); TTY default TUI |
| `ait agent connect` / bare destination after `agent` | Remote Loop via `RunAgentSession` |
| `ai` (retired) | Use `ait` |

## Remote connect

`internal/ait/agent/connect` — see [`ait-internals.md`](ait-internals.md).

## Tests

Former `opsai` tests now live under `internal/ait/agent` (`oa_test.go`, `connect/run_test.go`, …). See change-acceptance **OA\*** / **AG\*** / **AT49+**.

## User docs

- Migration: [`ai.md`](ai.md)
- Current guide: [`ait.md`](ait.md)
- Implementation: [`ait-internals.md`](ait-internals.md)

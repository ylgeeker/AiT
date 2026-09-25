# ait internals

## Package layering

```text
cmd/ait → internal/ait (package ait: Run, Execute, help)
         → internal/ait/router     (top-level argv: Loop, agent admin, plugins, shell)
         → internal/ait/loop       (local Agent Loop)
         → internal/ait/dispatch   (plugin exec via tooldispatch)
         → internal/ait/shell      (toolbox REPL: install/update/upgrade)
         → internal/ait/agent      (admin CLI: connect, model, provider, ollama, config)
         → internal/pkg/aitplugin  (manifest schema + Discover)
         → internal/pkg/aitrelease (GitHub Release fetch + apply)
         → internal/pkg/agentwire  (Loop config SSOT)
         → internal/pkg/tooldispatch (shared plugin exec)
         → internal/pkg/pluginbin  (ExecPath + Runner)
```

- `internal/ait` and subpackages must not import `internal/ssh` / `ailog` / `aiimg` / `aihost` for execution (exec only via `dispatch` / `tooldispatch`).
- `aitplugin` / `pluginbin` must not import `internal/ait`.
- `agentcore` must not import product packages.

## Top-level routing (`internal/ait/router`)

| argv | Route |
|------|--------|
| (empty) | `loop.Run` — Agent Loop |
| `agent …` | admin verb → `agent.RunAdmin`; else Loop with optional `--plain` |
| `shell` | toolbox REPL (`shell.RunREPL`) |
| `aissh`, `log`, … | `dispatch.RunPlugin` (aliases via `ResolvePluginName`) |
| `help` / `-h` | `WriteRootHelp` |

## Key types

| Piece | Role |
|-------|------|
| `aitplugin.Manifest` | v1 JSON schema (`apiVersion`, `name`, `binary`, `flags` for help only) |
| `aitplugin.Discover` | `sync.Once` scan of `InstallDir/plugins/<id>/` + PATH builtin fill |
| `tooldispatch.Exec` | Lookup plugin + `pluginbin.ExecPath` (stdin/stdout/stderr) |
| `shell.ExecuteLine` | one REPL line; swallows bare `*exec.ExitError` (AT14) |
| `aitrelease.Execute` | install/update/upgrade from GitHub Release |
| `agentwire.LocalLoopConfig` | TUI/plain, timeout, aissh CLI runner for Loop |
| `agent.RunAdmin` | Cobra tree for connect/model/provider/ollama/config |

## Config scope

- **`ait agent config`** — `ai.*` keys (`config.ValidateAI`, `EnsureAIDefaults`)
- **`aissh config`** — SSH/session keys only; rejects `ai.*` (`config.ValidateSSH`)

## Release install (`aitrelease`)

- Install root: `aitplugin.PrimaryInstallDir()`.
- Assets: `ait-<tag>-<goos>-<arch>.tar.gz` (Windows: `.zip`).
- Tests inject `Options.BaseURL` or env `AIT_RELEASE_TEST_BASE_URL`.
- Loopback fixture: [`mock/testrelease`](../../mock/testrelease); long-running daemon: `go run ./mock/aitrelease-mock`.

## Manifest hook

Each shipped tool's `ExecuteWriters` handles `ait-plugin-manifest` before default routing and writes embedded JSON via `aitplugin.WriteManifest`.

## Further reading

- User guide: [`ait.md`](ait.md)
- Agent admin (former `ai`): [`ai-internals.md`](ai-internals.md)
- Loop kernel: `internal/pkg/agentcore`

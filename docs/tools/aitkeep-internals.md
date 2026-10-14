# aitkeep internals

## Package layering

```text
cmd/keep/main.go
  → internal/pkg/keepalive (package keepalive: Run, ExecuteWriters, Serve, Supervisor)
       → platform.DataDir (AppRoot)
       → helpdoc / errhint / buildinfo / toolname
```

There is no `internal/keep` application package. The daemon and the CLI share `keepalive` because that package already owned ping, supervision, and the unit file, and nothing else imported it. `cmd/keep` only calls `keepalive.Run`.

Import direction: `cmd/keep` → `internal/pkg/keepalive` → other `internal/pkg/*`. `keepalive` must not import `internal/ssh`, `internal/ait`, `internal/host`, or `internal/log`.

The plugin command name is `keep` (`toolname.Keep`). The binary name is `aitkeep` (`toolname.AITKeep`). Official discovery uses the embedded manifest at `internal/pkg/aitplugin/manifests/aitkeep/plugin.json`.

## Routing

| argv | Route |
|------|-------|
| empty / `-h` / `--help` / `help` | `WriteRootHelp` |
| `-V` / `--version` / `version` | `Version()` |
| `ait-plugin-manifest` | embedded manifest JSON |
| `add`, `reload`, `remove`, `serve`, `start`, `status`, `stop` | Cobra commands from `commandOrder` |
| else | errhint + Try |

## Control plane

`Serve` takes the parent context, acquires `keepalive.lock`, listens for `GET /ping`, and opens the admin endpoint (`keepalive.sock`, or a Windows pipe named from a hash of AppRoot). `reload` and `start` use that parent context.

Each unit slot has a generation. `reconcile` leaves a running unit alone when `reflect.DeepEqual` says the spec is unchanged, leaves a stopped unit stopped when the spec is unchanged, stops a running unit that disappeared, and restarts a unit whose spec changed. A `Watch` goroutine writes state only when its generation is still current.

## Limits

Sampling sums the process tree (CPU, RSS, and IO bytes). Linux joins the child to a temporary cgroup before `exec` and removes that directory when `runProcess` returns. Hard caps are `cpu.max`, `memory.max`, and `io.max` on Linux, and Job Object CPU rate plus job memory on Windows. `memory.high` is not a hard cap. A Windows CPU rate is refused when `cpu_percent / NumCPU` is above 100 percent of the machine.

`Decide` does not share one sustain clock:

- CPU and process bytes/sec: pause inside `SampleEvery`; pause time is not added to the duty clock.
- Memory and device `%util`: no pause; the clock is wall time. Device `%util` kills only when this process had read or write bytes in the window.
- Disk space bytes or filesystem percent: terminate and return `ErrDiskFailed` (no restart). A walk capped at `DiskWalkBudget` is the same failure.

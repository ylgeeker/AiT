# aitkeep

Supervise processes from `keepalive.yaml` and answer `GET /ping` with `pong`.

The shipped binary is `aitkeep`. The toolbox plugin command is `keep` (`ait keep …`).

## Build

```bash
make build   # → bin/aitkeep
./bin/aitkeep --help
```

## Usage

```text
aitkeep [command]
aitkeep serve [--listen addr]
aitkeep status
aitkeep start NAME
aitkeep stop NAME
aitkeep reload
aitkeep add --name NAME [--cwd DIR] [--cpu P] [--mem BYTES] [--disk-bytes N] [--disk-percent P] [--read-bps N] [--write-bps N] [--io-percent P] -- COMMAND...
aitkeep remove NAME
aitkeep version
```

No arguments print help. `serve` is the daemon.

Config, the single-instance lock, and the admin endpoint live under AppRoot (`AIT_HOME`, otherwise the platform data directory):

| File | Role |
|------|------|
| `keepalive.yaml` | Listen address and units (mode 0600) |
| `keepalive.lock` | One `serve` per AppRoot |
| `keepalive.sock` | Unix admin socket |

On Windows the admin endpoint is a named pipe derived from the AppRoot path, so two homes do not share one pipe.

`GET /ping` returns `200` `text/plain` and the body `pong\n`. Any other method or path is `404`. The default listen address is `127.0.0.1:8757` (`--listen` or yaml `listen`). A listen address written into the yaml takes effect the next time `serve` starts.

## Units

```yaml
listen: 127.0.0.1:8757
units:
  - name: demo
    command: ["sleep", "3600"]
    cwd: /var/lib/demo
    limits:
      cpu_percent: 20
      memory_bytes: 268435456
      disk_space_bytes: 1073741824
      disk_space_percent: 90
      disk_read_bps: 1048576
      disk_write_bps: 1048576
      disk_io_percent: 80
```

A zero limit turns that dimension off. When both an absolute limit and a percent are set, either one over the ceiling triggers that dimension's action.

| Limit | Meaning | When exceeded |
|-------|---------|----------------|
| `cpu_percent` | Percent of one core (`150` is 1.5 cores) | Hard cap when the kernel accepts it; otherwise pause inside each 1s sample, then kill after 10s of running-while-over |
| `memory_bytes` | Process-tree RSS | Linux `memory.max` or a Windows job memory limit; otherwise kill after 10s of wall time. The process is not paused |
| `disk_read_bps` / `disk_write_bps` | Process-tree bytes per second | Linux `io.max` when available; otherwise the same pause-then-kill path as CPU |
| `disk_io_percent` | Block-device busy percent, including other processes | Kill after 10s only when this process also read or wrote. Idle processes are left alone. No device stats → `io=unsupported` |
| `disk_space_bytes` | Regular-file byte sum under `cwd` | Terminate and do not restart |
| `disk_space_percent` | Used-block percent of the filesystem holding `cwd` | Terminate and do not restart |

An empty `cwd` skips both disk-space percent and device busy percent. A directory walk that hits the entry budget, or a filesystem stat that fails, counts as over for the disk-space limit that is configured.

Abnormal exits restart with backoff from 1s to 30s. A run that stays up 30s resets the backoff and the restart count. Five abnormal restarts inside 10 minutes mark the unit `failed`. Exit code 0 does not restart. A disk-space failure does not restart until a later `start` sees the usage back under the limit.

`reload` restarts only units whose spec changed. A unit stopped with `stop` stays stopped across `reload` when its spec is unchanged. `add` replaces a unit of the same name. `remove` deletes the yaml entry even when `serve` is not running.

`status` prints each unit's state and the current sample for every nonzero limit.

## Examples

```bash
aitkeep serve
aitkeep add --name demo --cpu 20 -- sleep 3600
aitkeep status
aitkeep stop demo
aitkeep remove demo
curl -fsS http://127.0.0.1:8757/ping
```

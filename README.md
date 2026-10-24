# AiT

**AiT** is a cross-platform (Linux / macOS / Windows) AI-enhanced command-line **toolbox** for operations engineers—terminal and SSH workflows today, more tools over time.

**Tools:**

- [`ait`](cmd/ait) — Agent interaction entry (Loop by default; `ait agent` admin; `ait shell` toolbox REPL). Guide: [ait](docs/tools/ait.md).
- [`aitssh`](cmd/ssh) — SSH client (product docs: [aissh](docs/tools/aissh.md); plugin command `ssh`; Go package `aissh`; **no AI**).
- [`aitimg`](cmd/img) — image metadata extractor (EXIF, GPS, reverse-geocoded address). Guide: [aiimg](docs/tools/aiimg.md) (plugin `img`).
- [`aitlog`](cmd/log) — system log viewer and search (Linux/macOS/Windows). Guide: [ailog](docs/tools/ailog.md) (plugin `log`).
- [`aithost`](cmd/host) — live host CPU/mem/disk/net/process monitor. Guide: [aihost](docs/tools/aihost.md) (plugin `host`).
- [`aitkeep`](cmd/keep) — process supervisor (`keep` plugin command; `GET /ping` → `pong`). Guide: [aitkeep](docs/tools/aitkeep.md).

## What you get

- OpenSSH-compatible connect (`-J`, config Host aliases, `ssh -G` resolution)
- Session bookmarks: `--alias`, `session --load`, `session --switch`
- In-session `!ait` meta commands (suspend, chain hop, transfer)
- Default AiT-managed remote prompt (`AiT user@host:…>`); set `AIT_REMOTE_SYSTEM_PROMPT=1` for the native shell prompt
- Ops AI Agent via **`ait`** / **`ait agent`** (Agent Loop, connect, model, provider, ollama, config); `aitssh` has no AI subtree
- `!ait get` / `!ait put` on the live SSH session, standalone `aitssh get` / `put` via SFTP
- Optional `--tmux` to keep remote shells alive after disconnect
- Local `aitssh shell` REPL with history and UTF-8 line editing

## Examples

```bash
# 1) Bookmark a host and resume later (OpenSSH-style flags still work)
aitssh deploy@prod --alias prod
aitssh session --list
aitssh session --load prod

# 2) Jump host in one command (same as ssh -J; Host aliases resolve via ssh -G)
aitssh -J bastion app.internal

# 3) Agent Loop and remote connect via ait
ait agent config init
ait agent model list
ait                              # local Agent Loop
ait agent connect user@host      # remote Agent

# 4) Transfer over the current session (single-arg defaults to ./basename)
!ait get ./build/app-linux-amd64
!ait put ./patch.tar.gz

# 5) tmux keep-alive: remote work survives disconnect
aitssh devhost --tmux
# Later: aitssh session --load devhost --tmux

# 6) Interactive Agent Loop (on by default; re-enable only if you turned it off)
ait agent "summarize this log pattern: connection reset"
```

## Installation

**Prerequisites:** Go **1.25+** if building from source; **OpenSSH client** (`ssh`) on PATH for `aitssh` connect/`load` to resolve `~/.ssh/config`.

### Install `ait`

```bash
go install github.com/os-artificer/AiT/cmd/ait@latest
ait --help
```

Ensure `$(go env GOPATH)/bin` is on your PATH.

### Install official plugins (`ssh`, `log`, `img`, `host`, `keep`)

After a [GitHub Release](https://github.com/os-artificer/AiT/releases) is available, open the toolbox REPL with `ait shell` (bare `ait` starts the Agent Loop, or exits with the not-configured hint when AI is unset):

```bash
ait shell
```

```text
ait> install              # all official plugins
ait> install ssh keep     # one or more by plugin command name
ait> update               # update installed plugins
ait> upgrade              # update ait + installed plugins
ait> help install
```

Legacy install aliases such as `aissh` → `ssh` still resolve. Binaries land next to the running `ait` executable (same layout as `make install`). Details: [docs/tools/ait.md](docs/tools/ait.md#install--update--upgrade).

### Build from source (development)

```bash
git clone https://github.com/os-artificer/AiT.git
cd AiT
make build
./bin/ait --help
```

For local development, `make build` also produces `bin/aitssh`, `bin/aitimg`, `bin/aitlog`, `bin/aithost`, and `bin/aitkeep`. To install plugins into `bin/` without a Release, run `make install`. Once a Release exists, run `ait shell` and then `install` at the `ait>` prompt. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Documentation

| Document | Path |
|----------|------|
| Project structure | [docs/architecture.md](docs/architecture.md) |
| aissh internals | [docs/tools/aissh-internals.md](docs/tools/aissh-internals.md) |
| aissh guide | [docs/tools/aissh.md](docs/tools/aissh.md) |
| ait internals | [docs/tools/ait-internals.md](docs/tools/ait-internals.md) |
| ait guide | [docs/tools/ait.md](docs/tools/ait.md) |
| aiimg internals | [docs/tools/aiimg-internals.md](docs/tools/aiimg-internals.md) |
| aiimg guide | [docs/tools/aiimg.md](docs/tools/aiimg.md) |
| ailog internals | [docs/tools/ailog-internals.md](docs/tools/ailog-internals.md) |
| ailog guide | [docs/tools/ailog.md](docs/tools/ailog.md) |
| aihost internals | [docs/tools/aihost-internals.md](docs/tools/aihost-internals.md) |
| aihost guide | [docs/tools/aihost.md](docs/tools/aihost.md) |
| aitkeep internals | [docs/tools/aitkeep-internals.md](docs/tools/aitkeep-internals.md) |
| aitkeep guide | [docs/tools/aitkeep.md](docs/tools/aitkeep.md) |
| Python E2E tests | [tests/README.md](tests/README.md) |
| In-repo SSH/HTTP mock | [docs/testenv-mock.md](docs/testenv-mock.md) |
| Contributing | [CONTRIBUTING.md](CONTRIBUTING.md) |
| Code layout | [AGENTS.md](AGENTS.md) |
| Code of conduct | [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) |
| Changelog | [CHANGELOG.md](CHANGELOG.md) |

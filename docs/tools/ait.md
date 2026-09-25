# ait — Agent interaction entry

**`ait`** is the unified Agent interaction entry for AiT. Bare **`ait`** (or **`ait agent`**) starts the interactive Agent Loop. Plugin tools run as top-level dispatch (`ait aissh`, `ait log`, …). The toolbox REPL for install/update/upgrade is **`ait shell`**.

## Usage

```bash
ait                              # Agent Loop (default)
ait agent [goal...]              # same Loop; optional initial goal
ait agent --plain [goal...]      # plain UI
ait agent connect user@host      # remote Agent Loop
ait agent model list             # list configured models
ait agent config set ai.foo bar  # AI-scoped config (not aissh config)
ait aissh user@host.example.com  # plugin dispatch (no REPL)
ait log search error             # alias → ailog
ait shell                        # toolbox REPL (install/update/upgrade)
ait -h                           # Agent entry help
ait shell -h                     # toolbox REPL help
```

Natural-language text after flags is treated as an Agent goal (e.g. `ait list files`).

## Toolbox shell (`ait shell`)

```text
ait> help
ait> install
ait> update
ait> upgrade
ait> aissh user@host.example.com
ait> ailog search error
ait> exit
```

Reserved builtins in the shell REPL: `help`, `exit`/`quit`, `version`, `history`, `clear`, `install`, `update`, `upgrade`. Agent Loop is **not** started from the shell REPL — use bare `ait` or `ait agent`.

## Plugin aliases (top-level)

| Alias | Plugin |
|-------|--------|
| `log` | `ailog` |
| `host` | `aihost` |
| `img` | `aiimg` |
| `ssh` | `aissh` |

## install / update / upgrade

From **`ait shell`**, download official binaries from the latest [GitHub Release](https://github.com/os-artificer/AiT/releases). Official plugins: `aissh`, `ailog`, `aiimg`, `aihost` (no standalone `ai` binary).

See [ait-internals.md](ait-internals.md) for routing and package layout.

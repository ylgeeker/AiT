# ait — toolbox entry and Ops AI Agent

**`ait`** is the AiT toolbox entry. Bare **`ait`** (or **`ait shell`**) opens the toolbox REPL (prompt `ait>`). The interactive Agent Loop starts only with **`ait agent`**. When AI is not configured, `ait agent` exits with `AI not configured` and a `Try:` block for `ait agent config set`. Plugin tools run as top-level dispatch (`ait ssh`, `ait log`, …). Unknown top-level tokens print an errhint with `Try: ait agent`.

## Usage

```bash
ait                              # toolbox shell (ait>)
ait shell                        # same toolbox REPL
ait agent [goal...]              # Agent Loop; optional initial goal
ait agent --plain [goal...]      # plain UI
ait agent connect user@host      # remote Agent Loop
ait agent model list             # list configured models
ait agent config set ai.foo bar  # AI-scoped config (not aitssh config)
ait ssh user@host.example.com    # plugin dispatch (no REPL)
ait log search error             # plugin log → binary aitlog
ait -h                           # toolbox + Agent root help
ait shell -h                     # toolbox REPL help
```

Natural-language goals belong after `ait agent` (for example `ait agent list files`). A bare free-form first token (for example `ait list files`) is an unknown command.

## Toolbox shell (`ait` / `ait shell`)

```text
ait> help
ait> install
ait> update
ait> upgrade
ait> ssh user@host.example.com
ait> log search error
ait> exit
```

Reserved builtins in the shell REPL: `help`, `exit`/`quit`, `version`, `history`, `clear`, `install`, `update`, `upgrade`. Agent Loop is **not** started from the shell REPL — use `ait agent`.

## Plugin command → binary

| Command | Binary |
|---------|--------|
| `ssh` | `aitssh` |
| `log` | `aitlog` |
| `img` | `aitimg` |
| `host` | `aithost` |
| `keep` | `aitkeep` |

Legacy install names such as `aissh` / `ailog` still map to these plugin commands.

## install / update / upgrade

From the toolbox shell (`ait` or `ait shell`), download official binaries from the latest [GitHub Release](https://github.com/os-artificer/AiT/releases). Official plugins: `ssh`, `log`, `img`, `host`, `keep`.

See [ait-internals.md](ait-internals.md) for routing and package layout.

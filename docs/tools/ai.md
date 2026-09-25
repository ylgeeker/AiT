# ai — migrated to `ait`

> **The standalone `ai` binary and plugin were removed.** Use **`ait`** as the Agent interaction entry.

| Former | Use instead |
|--------|-------------|
| `ai` (no args) | `ait` or `ait agent` — Agent Loop |
| `ai connect user@host` | `ait agent connect user@host` or `ait agent user@host` |
| `ai model list` | `ait agent model list` |
| `ai provider …` / `ai ollama …` | `ait agent provider …` / `ait agent ollama …` |
| `aissh config set ai.*` | `ait agent config set ai.*` |
| Toolbox REPL (`ait>` install/update) | `ait shell` |

**`aissh` has zero AI** — no `aissh ai`, no `!ait ai` / `!!` AI entry.

Implementation: [`ait-internals.md`](ait-internals.md). User guide: [`ait.md`](ait.md). Historical internals: [`ai-internals.md`](ai-internals.md).

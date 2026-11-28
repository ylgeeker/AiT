# Unreleased

## ait

- 【Changed】Bare `ait` (and `ait --ai-experimental` with no remaining argv) opens the toolbox shell (`ait>`), same as `ait shell`. Agent Loop starts only with `ait agent`. Top-level free-form text (for example `ait list files`) and a leading `--plain` without `agent` return an unknown-command hint with `Try: ait agent`.

"""Local Agent admin CLI smoke tests (`ait agent`)."""

from __future__ import annotations

import pytest

from tests.lib.cli import AITRunner, assert_exit
from tests.lib.hints import assert_hint

pytestmark = [pytest.mark.smoke]


def test_ai_experimental_disabled_by_default(runner, ai_bin):
    runner.run_ok(["config", "init"])
    off = AITRunner(ai_bin, runner.ait_home, base_env=runner.base_env)
    proc = off.run(["model", "list"], timeout=5)
    assert_exit(1, proc)
    assert_hint(
        proc,
        sentinel="AI experimental features are disabled",
        try_commands=(
            "ai --ai-experimental",
            "AIT_AI_EXPERIMENTAL=1 ai",
        ),
    )
    combined = (proc.stdout or "") + (proc.stderr or "")
    assert "aitssh" not in combined.lower()


def test_ai_help_without_experimental(ai_bin, runner):
    off = AITRunner(ai_bin, runner.ait_home, base_env=runner.base_env)
    proc = off.run_ok(["--help"])
    assert "Usage:" in proc.stdout
    assert "experimental" in proc.stdout.lower() or "connect" in proc.stdout


def test_ai_model_list_non_interactive(runner, ai_runner):
    runner.run_ok(["config", "init"])
    proc = ai_runner.run(["model", "list"], timeout=5)
    # may fail without keys but must not hang
    assert proc.returncode in (0, 1)


def test_ai_ask_removed(ai_runner):
    proc = ai_runner.run(["ask", "ping"], timeout=5)
    assert_exit(1, proc)
    assert_hint(proc, try_commands=("ait agent",))


def test_ai_agent_removed(ai_runner):
    proc = ai_runner.run(["agent", "ping"], timeout=5)
    assert_exit(1, proc)
    assert_hint(proc, try_commands=("ait agent",))


def test_ai_connect_help(ai_runner):
    proc = ai_runner.run_ok(["connect", "--help"])
    assert "Usage:" in proc.stdout

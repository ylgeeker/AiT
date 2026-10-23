"""Local Agent admin CLI smoke tests (`ait agent`)."""

from __future__ import annotations

import pytest

from tests.lib.cli import assert_exit
from tests.lib.hints import assert_hint

pytestmark = [pytest.mark.smoke]


def test_ai_help_without_experimental(ai_runner):
    proc = ai_runner.run_ok(["--help"])
    assert "Usage:" in proc.stdout
    assert "connect" in proc.stdout.lower()


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

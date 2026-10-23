"""Remote Ops AI via `ait agent connect`. Legacy !! / !ait ai are not on aitssh."""

from __future__ import annotations

import pytest

pytestmark = [pytest.mark.ssh, pytest.mark.ai]


def test_ai_connect_help(ai_runner):
    proc = ai_runner.run_ok(["connect", "--help"])
    assert "Usage:" in proc.stdout
    assert "connect" in proc.stdout.lower()


@pytest.mark.skip(reason="v4 remote Agent uses ait agent connect; !! / !ait ai removed from aitssh")
def test_remote_ai_legacy_meta_removed():
    pass

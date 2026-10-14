//go:build unix

/**
 * Copyright 2026 AiT authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/

package keepalive

import (
	"context"
	"os/exec"
	"syscall"
	"time"
)

func applyProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func stopProcess(pid int) {
	if pid <= 0 {
		return
	}
	_ = syscall.Kill(-pid, syscall.SIGTERM)
	time.Sleep(StopGrace)
	_ = syscall.Kill(-pid, syscall.SIGKILL)
}

func pauseFor(ctx context.Context, pid int, d time.Duration) {
	if pid <= 0 || d <= 0 {
		return
	}
	_ = syscall.Kill(-pid, syscall.SIGSTOP)
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
	_ = syscall.Kill(-pid, syscall.SIGCONT)
}

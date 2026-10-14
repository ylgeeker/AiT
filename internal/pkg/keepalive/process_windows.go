//go:build windows

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
	"time"

	"golang.org/x/sys/windows"
)

func applyProcAttr(cmd *exec.Cmd) {}

func stopProcess(pid int) {
	if pid <= 0 {
		return
	}
	if stopJob(pid) {
		return
	}
	proc, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return
	}
	defer windows.CloseHandle(proc)
	_ = windows.TerminateProcess(proc, 1)
}

func pauseFor(ctx context.Context, pid int, d time.Duration) {
	if pid <= 0 || d <= 0 {
		return
	}
	suspendProcess(pid)
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
	resumeProcess(pid)
}

func suspendProcess(pid int) {
	proc, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
	if err != nil {
		return
	}
	defer windows.CloseHandle(proc)
	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	suspend := ntdll.NewProc("NtSuspendProcess")
	_, _, _ = suspend.Call(uintptr(proc))
}

func resumeProcess(pid int) {
	proc, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
	if err != nil {
		return
	}
	defer windows.CloseHandle(proc)
	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	resume := ntdll.NewProc("NtResumeProcess")
	_, _, _ = resume.Call(uintptr(proc))
}

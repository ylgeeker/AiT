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
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func command(unit Unit, cgroupDir string) *exec.Cmd {
	var cmd *exec.Cmd
	if cgroupDir != "" {
		procs := filepath.Join(cgroupDir, "cgroup.procs")
		argv := append([]string{"-c", cgroupJoinScript, shellArgv0, procs}, unit.Command...)
		cmd = exec.Command(shellBin, argv...)
	} else {
		cmd = exec.Command(unit.Command[0], unit.Command[1:]...)
	}
	cmd.Dir = unit.Cwd
	if len(unit.Env) > 0 {
		cmd.Env = append(os.Environ(), unit.Env...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	applyProcAttr(cmd)
	return cmd
}

func runProcess(ctx context.Context, unit Unit, onSample func(Usage)) (int, error) {
	if len(unit.Command) == 0 {
		return 1, ErrNoCommand
	}
	dir, caps, release := openLimits(unit)
	defer release()
	cmd := command(unit, dir)
	if err := cmd.Start(); err != nil {
		return 1, err
	}
	if dir == "" {
		caps = attachLimits(cmd.Process.Pid, unit)
	}

	done := make(chan struct{})
	defer close(done)
	diskFail := make(chan struct{}, 1)
	go monitor(ctx, cmd.Process.Pid, unit, caps, done, diskFail, onSample)

	wait := make(chan error, 1)
	go func() { wait <- cmd.Wait() }()

	select {
	case <-ctx.Done():
		stopProcess(cmd.Process.Pid)
		<-wait
		return 0, ErrOperatorStop
	case err := <-wait:
		select {
		case <-diskFail:
			return 1, ErrDiskFailed
		default:
		}
		if err == nil {
			return 0, nil
		}
		if ee, ok := err.(*exec.ExitError); ok {
			return ee.ExitCode(), nil
		}
		return 1, err
	}
}

func monitor(ctx context.Context, pid int, unit Unit, caps HardCaps, done <-chan struct{}, diskFail chan<- struct{}, onSample func(Usage)) {
	var held Hold
	var prev ReadWrite
	wait := SampleEvery
	timer := time.NewTimer(wait)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		case <-timer.C:
		}
		usage := sampleUsage(pid, &prev)
		fillDisk(&usage, unit)
		if onSample != nil {
			onSample(usage)
		}
		held.Duty = accrue(held.Duty, dutyOver(usage, unit.Limits, caps), wait)
		held.Memory = accrue(held.Memory, memoryOver(usage, unit.Limits, caps), wait)
		held.IO = accrue(held.IO, ioUtilHot(usage, unit.Limits), wait)
		pause := time.Duration(0)
		switch Decide(usage, unit.Limits, caps, held, Sustain) {
		case ActionDiskFail:
			select {
			case diskFail <- struct{}{}:
			default:
			}
			stopProcess(pid)
			return
		case ActionKill:
			stopProcess(pid)
			return
		case ActionThrottle:
			pause = throttlePause(usage, unit.Limits, caps)
			pauseFor(ctx, pid, pause)
			held.Memory = accrue(held.Memory, memoryOver(usage, unit.Limits, caps), pause)
			held.IO = accrue(held.IO, ioUtilHot(usage, unit.Limits), pause)
		case ActionNone:
		}
		next := SampleEvery - pause
		if next < time.Millisecond {
			next = time.Millisecond
		}
		wait = next
		timer.Reset(wait)
	}
}

func fillDisk(usage *Usage, unit Unit) {
	if unit.Cwd == "" {
		return
	}
	n, err := DirUsage(unit.Cwd, DiskWalkBudget)
	if err != nil {
		usage.DiskWalkFailed = true
	} else {
		usage.DiskSpaceBytes = n
	}
	if unit.Limits.DiskSpacePercent == 0 && unit.Limits.DiskIOPercent == 0 {
		return
	}
	vol, err := volumeSampler(unit.Cwd)
	if err != nil {
		usage.VolumeFailed = true
		return
	}
	usage.DiskSpacePercent = vol.UsedPercent
	usage.DiskIOPercent = vol.IOPercent
	usage.IOKnown = vol.IOKnown
}

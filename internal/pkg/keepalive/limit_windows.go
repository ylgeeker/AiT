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
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

var jobs = jobTable{jobs: map[int]windows.Handle{}}

const (
	cpuRateControlEnable  = 0x1
	cpuRateControlHardCap = 0x4
)

type cpuRateInfo struct {
	ControlFlags uint32
	CPURate      uint32
}

type jobTable struct {
	mu   sync.Mutex
	jobs map[int]windows.Handle
}

func applyLimits(pid int, unit Unit) HardCaps {
	var caps HardCaps
	if pid <= 0 {
		return caps
	}
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return caps
	}
	proc, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		windows.CloseHandle(job)
		return caps
	}
	defer windows.CloseHandle(proc)
	if err := windows.AssignProcessToJobObject(job, proc); err != nil {
		windows.CloseHandle(job)
		return caps
	}
	caps.CPU = setCPURate(job, unit.Limits.CPUPercent)
	caps.Memory = setJobMemory(job, unit.Limits.MemoryBytes)
	jobs.mu.Lock()
	jobs.jobs[pid] = job
	jobs.mu.Unlock()
	return caps
}

func setCPURate(job windows.Handle, percent float64) bool {
	if percent <= 0 {
		return false
	}
	rate, ok := cpuRateFromCores(percent, runtime.NumCPU())
	if !ok {
		return false
	}
	info := cpuRateInfo{ControlFlags: cpuRateControlEnable | cpuRateControlHardCap, CPURate: rate}
	ret, err := windows.SetInformationJobObject(
		job, windows.JobObjectCpuRateControlInformation,
		uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info)),
	)
	return err == nil && ret != 0
}

func setJobMemory(job windows.Handle, bytes uint64) bool {
	if bytes == 0 {
		return false
	}
	var info windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_JOB_MEMORY
	info.JobMemoryLimit = uintptr(bytes)
	ret, err := windows.SetInformationJobObject(
		job, windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info)),
	)
	return err == nil && ret != 0
}

func openLimits(Unit) (string, HardCaps, func()) {
	return "", HardCaps{}, func() {}
}

func attachLimits(pid int, unit Unit) HardCaps {
	return applyLimits(pid, unit)
}

func stopJob(pid int) bool {
	jobs.mu.Lock()
	job, ok := jobs.jobs[pid]
	if ok {
		delete(jobs.jobs, pid)
	}
	jobs.mu.Unlock()
	if !ok {
		return false
	}
	_ = windows.TerminateJobObject(job, 1)
	_ = windows.CloseHandle(job)
	return true
}

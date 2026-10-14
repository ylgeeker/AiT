//go:build linux

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

import "os"

func openLimits(unit Unit) (dir string, caps HardCaps, release func()) {
	release = func() {}
	if !wantsCgroup(unit) {
		return "", caps, release
	}
	var err error
	dir, err = os.MkdirTemp("/sys/fs/cgroup", "ait-keepalive-"+unit.Name+"-")
	if err != nil {
		return "", caps, release
	}
	made := dir
	release = func() { _ = os.RemoveAll(made) }
	caps.CPU = wroteCPUMax(dir, unit.Limits.CPUPercent)
	caps.Memory = wroteMemoryMax(dir, unit.Limits.MemoryBytes)
	caps.DiskIO = wroteIOMax(dir, unit)
	return dir, caps, release
}

func wroteCPUMax(dir string, percent float64) bool {
	if percent <= 0 {
		return false
	}
	return WriteCPUMax(dir, percent) == nil
}

func wroteMemoryMax(dir string, bytes uint64) bool {
	if bytes == 0 {
		return false
	}
	return WriteMemoryMax(dir, bytes) == nil
}

func wroteIOMax(dir string, unit Unit) bool {
	if unit.Cwd == "" || (unit.Limits.DiskReadBPS == 0 && unit.Limits.DiskWriteBPS == 0) {
		return false
	}
	major, minor, err := blockDevice(unit.Cwd)
	if err != nil {
		return false
	}
	return WriteIOMax(dir, major, minor, unit.Limits.DiskReadBPS, unit.Limits.DiskWriteBPS) == nil
}

func attachLimits(int, Unit) HardCaps { return HardCaps{} }

func wantsCgroup(unit Unit) bool {
	if unit.Limits.CPUPercent > 0 || unit.Limits.MemoryBytes > 0 {
		return true
	}
	return unit.Limits.DiskReadBPS > 0 || unit.Limits.DiskWriteBPS > 0
}

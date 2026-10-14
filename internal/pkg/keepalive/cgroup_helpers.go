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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// WriteCPUMax writes a cgroup v2 cpu.max file for percent of one CPU.
// dir is the cgroup directory. percent may exceed 100 so the quota can cover more than one core.
// The quota uses CPUPeriodUS. Returns an error when percent is not positive or the write fails.
func WriteCPUMax(dir string, percent float64) error {
	if percent <= 0 {
		return fmt.Errorf("cpu percent %v out of range", percent)
	}
	quota := int(percent * float64(CPUPeriodUS) / 100)
	if quota < 1 {
		quota = 1
	}
	body := fmt.Sprintf("%d %d\n", quota, CPUPeriodUS)
	return os.WriteFile(filepath.Join(dir, "cpu.max"), []byte(body), 0o644)
}

// WriteMemoryMax writes cgroup v2 memory.max.
// dir is the cgroup directory. bytes is the hard ceiling. Returns an error when bytes is zero or the write fails.
func WriteMemoryMax(dir string, bytes uint64) error {
	if bytes == 0 {
		return fmt.Errorf("memory max is zero")
	}
	return os.WriteFile(filepath.Join(dir, "memory.max"), []byte(strconv.FormatUint(bytes, 10)+"\n"), 0o644)
}

// WriteIOMax writes cgroup v2 io.max for one block device.
// A zero bps direction is recorded as max. Returns an error when both rates are zero or the write fails.
func WriteIOMax(dir string, major, minor uint32, readBPS, writeBPS uint64) error {
	if readBPS == 0 && writeBPS == 0 {
		return fmt.Errorf("io max is zero")
	}
	body := fmt.Sprintf("%d:%d rbps=%s wbps=%s\n", major, minor, bpsToken(readBPS), bpsToken(writeBPS))
	return os.WriteFile(filepath.Join(dir, "io.max"), []byte(body), 0o644)
}

func bpsToken(v uint64) string {
	if v == 0 {
		return "max"
	}
	return strconv.FormatUint(v, 10)
}

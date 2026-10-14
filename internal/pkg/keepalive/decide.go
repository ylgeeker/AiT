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

import "time"

// Action is what the supervisor should do with one sample.
type Action int

const (
	// ActionNone means the sample is inside the configured limits.
	ActionNone Action = iota
	// ActionThrottle means a duty-cycled limit is exceeded and the process should pause inside this sample.
	ActionThrottle
	// ActionKill means a limit stayed over its threshold and the tree should be terminated, then restarted.
	ActionKill
	// ActionDiskFail means disk space is over the ceiling. The tree is terminated and not restarted.
	ActionDiskFail
)

// Hold is how long each class of limit has stayed exceeded.
// Duty counts only time the process was running. Memory and IO count wall time, including pauses.
type Hold struct {
	Duty   time.Duration
	Memory time.Duration
	IO     time.Duration
}

// Decide chooses throttle, kill, or disk-fail for one sample.
// usage and limits describe the current tree. caps marks dimensions the kernel already hard-caps.
// held is how long each class has stayed over. sustain is the kill threshold.
// Disk space cannot be throttled. Memory is not paused. Device busy-percent is not a hard cap.
func Decide(usage Usage, limits Limits, caps HardCaps, held Hold, sustain time.Duration) Action {
	if diskSpaceOver(usage, limits) {
		return ActionDiskFail
	}
	if memoryOver(usage, limits, caps) && held.Memory >= sustain {
		return ActionKill
	}
	if ioUtilHot(usage, limits) && held.IO >= sustain {
		return ActionKill
	}
	if dutyOver(usage, limits, caps) && held.Duty >= sustain {
		return ActionKill
	}
	if dutyOver(usage, limits, caps) {
		return ActionThrottle
	}
	return ActionNone
}

// DutyPause returns how long to keep a process stopped inside window so usage can fall toward limit.
// usage and limit are the same unit (percent or bytes per second). A zero limit or usage within limit returns 0.
func DutyPause(usage, limit float64, window time.Duration) time.Duration {
	if limit <= 0 || usage <= limit || window <= 0 {
		return 0
	}
	frac := 1 - (limit / usage)
	if frac < 0 {
		return 0
	}
	if frac > 1 {
		frac = 1
	}
	return time.Duration(float64(window) * frac)
}

// throttlePause is the pause inside one SampleEvery for CPU and process byte rates.
// Memory and device busy-percent are not paused. The result is at most SampleEvery.
func throttlePause(usage Usage, limits Limits, caps HardCaps) time.Duration {
	var pause time.Duration
	if !caps.CPU {
		pause = longer(pause, DutyPause(usage.CPUPercent, limits.CPUPercent, SampleEvery))
	}
	if !caps.DiskIO {
		pause = longer(pause, DutyPause(float64(usage.ReadBPS), float64(limits.DiskReadBPS), SampleEvery))
		pause = longer(pause, DutyPause(float64(usage.WriteBPS), float64(limits.DiskWriteBPS), SampleEvery))
	}
	if pause > SampleEvery {
		return SampleEvery
	}
	return pause
}

func longer(cur, next time.Duration) time.Duration {
	if next > cur {
		return next
	}
	return cur
}

func accrue(acc time.Duration, over bool, d time.Duration) time.Duration {
	if !over || d <= 0 {
		return 0
	}
	return acc + d
}

func diskSpaceOver(usage Usage, limits Limits) bool {
	if limits.DiskSpaceBytes > 0 && (usage.DiskWalkFailed || usage.DiskSpaceBytes > limits.DiskSpaceBytes) {
		return true
	}
	if limits.DiskSpacePercent > 0 && (usage.VolumeFailed || usage.DiskSpacePercent > limits.DiskSpacePercent) {
		return true
	}
	return false
}

func dutyOver(usage Usage, limits Limits, caps HardCaps) bool {
	if limits.CPUPercent > 0 && !caps.CPU && usage.CPUPercent > limits.CPUPercent {
		return true
	}
	if limits.DiskReadBPS > 0 && !caps.DiskIO && usage.ReadBPS > limits.DiskReadBPS {
		return true
	}
	if limits.DiskWriteBPS > 0 && !caps.DiskIO && usage.WriteBPS > limits.DiskWriteBPS {
		return true
	}
	return false
}

func memoryOver(usage Usage, limits Limits, caps HardCaps) bool {
	return limits.MemoryBytes > 0 && !caps.Memory && usage.MemoryBytes > limits.MemoryBytes
}

func ioUtilHot(usage Usage, limits Limits) bool {
	if limits.DiskIOPercent <= 0 || !usage.IOKnown || usage.DiskIOPercent <= limits.DiskIOPercent {
		return false
	}
	return usage.ReadBPS > 0 || usage.WriteBPS > 0
}

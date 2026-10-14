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
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

// sampleHook replaces gopsutil sampling. Tests set it and restore the previous value.
var sampleHook func(pid int, prev *ReadWrite) Usage

// ReadWrite is the last cumulative IO sample used to compute bytes per second.
type ReadWrite struct {
	Read  uint64
	Write uint64
	At    time.Time
}

func sampleUsage(pid int, prev *ReadWrite) Usage {
	if sampleHook != nil {
		return sampleHook(pid, prev)
	}
	procs := processTree(pid)
	if len(procs) == 0 {
		return Usage{}
	}
	var usage Usage
	var read, write uint64
	for _, proc := range procs {
		if pct, err := proc.Percent(0); err == nil {
			usage.CPUPercent += pct
		}
		if mem, err := proc.MemoryInfo(); err == nil && mem != nil {
			usage.MemoryBytes += mem.RSS
		}
		io, err := proc.IOCounters()
		if err != nil || io == nil {
			continue
		}
		read += io.ReadBytes
		write += io.WriteBytes
	}
	fillRates(&usage, prev, read, write)
	return usage
}

func fillRates(usage *Usage, prev *ReadWrite, read, write uint64) {
	if prev == nil {
		return
	}
	now := time.Now()
	sec := now.Sub(prev.At).Seconds()
	usage.ReadBPS = bytesPerSec(read, prev.Read, sec, prev.At)
	usage.WriteBPS = bytesPerSec(write, prev.Write, sec, prev.At)
	prev.Read = read
	prev.Write = write
	prev.At = now
}

func bytesPerSec(cur, prev uint64, sec float64, at time.Time) uint64 {
	if at.IsZero() || sec <= 0 || cur < prev {
		return 0
	}
	return uint64(float64(cur-prev) / sec)
}

func processTree(pid int) []*process.Process {
	root, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil
	}
	out := []*process.Process{root}
	stack := []*process.Process{root}
	seen := map[int32]struct{}{int32(pid): {}}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		out, stack = appendChildren(out, stack, seen, cur)
	}
	return out
}

func appendChildren(out, stack []*process.Process, seen map[int32]struct{}, cur *process.Process) ([]*process.Process, []*process.Process) {
	kids, err := cur.Children()
	if err != nil {
		return out, stack
	}
	for _, kid := range kids {
		out, stack = appendChild(out, stack, seen, kid)
	}
	return out, stack
}

func appendChild(out, stack []*process.Process, seen map[int32]struct{}, kid *process.Process) ([]*process.Process, []*process.Process) {
	if kid == nil {
		return out, stack
	}
	if _, ok := seen[kid.Pid]; ok {
		return out, stack
	}
	seen[kid.Pid] = struct{}{}
	out = append(out, kid)
	stack = append(stack, kid)
	return out, stack
}

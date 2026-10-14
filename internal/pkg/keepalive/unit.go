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

// Limits is the resource ceiling for one hosted process tree.
// Zero values mean that dimension is not enforced.
type Limits struct {
	CPUPercent       float64 `yaml:"cpu_percent,omitempty"`
	MemoryBytes      uint64  `yaml:"memory_bytes,omitempty"`
	DiskSpaceBytes   uint64  `yaml:"disk_space_bytes,omitempty"`
	DiskSpacePercent float64 `yaml:"disk_space_percent,omitempty"`
	DiskReadBPS      uint64  `yaml:"disk_read_bps,omitempty"`
	DiskWriteBPS     uint64  `yaml:"disk_write_bps,omitempty"`
	DiskIOPercent    float64 `yaml:"disk_io_percent,omitempty"`
}

// Unit is one hosted command.
type Unit struct {
	Name    string   `yaml:"name"`
	Command []string `yaml:"command"`
	Cwd     string   `yaml:"cwd,omitempty"`
	Env     []string `yaml:"env,omitempty"`
	Limits  Limits   `yaml:"limits,omitempty"`
}

// File is the on-disk keepalive.yaml document.
type File struct {
	Listen string `yaml:"listen,omitempty"`
	Units  []Unit `yaml:"units,omitempty"`
}

// Usage is one sample of a hosted process tree.
type Usage struct {
	CPUPercent       float64
	MemoryBytes      uint64
	ReadBPS          uint64
	WriteBPS         uint64
	DiskSpaceBytes   uint64
	DiskSpacePercent float64
	DiskIOPercent    float64
	IOKnown          bool
	DiskWalkFailed   bool
	VolumeFailed     bool
}

// HardCaps reports which dimensions are already enforced by the kernel.
// A hard cap is not sustain-killed; the kernel throttles it.
type HardCaps struct {
	CPU    bool
	Memory bool
	// DiskIO is a kernel cap on this process tree's read/write bytes per second.
	// Device busy-percent is not a hard cap.
	DiskIO bool
}

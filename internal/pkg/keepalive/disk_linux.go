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

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/disk"

	"golang.org/x/sys/unix"
)

type ioPoint struct {
	at   time.Time
	busy uint64
}

var deviceIO = struct {
	mu sync.Mutex
	m  map[string]ioPoint
}{m: map[string]ioPoint{}}

func attachDeviceIO(cwd string, vol *Volume) {
	name, ok := blockDevName(cwd)
	if !ok {
		return
	}
	counters, err := disk.IOCounters(name)
	if err != nil {
		return
	}
	stat, ok := counters[name]
	if !ok {
		return
	}
	now := time.Now()
	deviceIO.mu.Lock()
	prev := deviceIO.m[name]
	deviceIO.m[name] = ioPoint{at: now, busy: stat.IoTime}
	deviceIO.mu.Unlock()
	vol.IOKnown = true
	if prev.at.IsZero() {
		return
	}
	wall := now.Sub(prev.at).Milliseconds()
	if wall <= 0 || stat.IoTime < prev.busy {
		return
	}
	pct := float64(stat.IoTime-prev.busy) / float64(wall) * 100
	if pct > 100 {
		pct = 100
	}
	vol.IOPercent = pct
}

func blockDevName(cwd string) (string, bool) {
	var st unix.Stat_t
	if err := unix.Stat(cwd, &st); err != nil {
		return "", false
	}
	path := "/sys/dev/block/" + strconv.FormatUint(uint64(unix.Major(st.Dev)), 10) + ":" + strconv.FormatUint(uint64(unix.Minor(st.Dev)), 10) + "/uevent"
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "DEVNAME=") {
			return strings.TrimPrefix(line, "DEVNAME="), true
		}
	}
	return "", false
}

func blockDevice(cwd string) (major, minor uint32, err error) {
	var st unix.Stat_t
	if err = unix.Stat(cwd, &st); err != nil {
		return 0, 0, err
	}
	return uint32(unix.Major(st.Dev)), uint32(unix.Minor(st.Dev)), nil
}

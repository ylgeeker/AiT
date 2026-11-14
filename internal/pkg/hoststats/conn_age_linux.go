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

package hoststats

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	nfConntrackPath = "/proc/net/nf_conntrack"
	deltaTimeKey    = "delta-time="
)

type connTuple struct {
	proto        string
	lip, rip     string
	lport, rport uint32
}

// enrichConnAge sets AgeSeconds from /proc/<pid>/fd/<fd> Ctim when Pid > 0.
// On Lstat failure, tries nf_conntrack delta-time for the 4-tuple.
func enrichConnAge(conns []ConnStat) {
	if len(conns) == 0 {
		return
	}
	now := time.Now()
	var track map[connTuple]int64
	for i := range conns {
		age := ageFromFD(conns[i].PID, conns[i].Fd, now)
		if age == nil {
			if track == nil {
				track = loadConntrackAges()
			}
			if track != nil {
				age = ageFromConntrack(track, conns[i])
			}
		}
		conns[i].AgeSeconds = age
	}
}

func ageFromFD(pid int32, fd uint32, now time.Time) *int64 {
	if pid <= 0 {
		return nil
	}
	path := fmt.Sprintf("/proc/%d/fd/%d", pid, fd)
	fi, err := os.Lstat(path)
	if err != nil {
		return nil
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return ageFromModTime(fi.ModTime(), now)
	}
	sec, nsec := st.Ctim.Unix()
	if sec == 0 && nsec == 0 {
		sec, nsec = st.Mtim.Unix()
	}
	created := time.Unix(sec, nsec)
	return ageFromModTime(created, now)
}

func ageFromModTime(created, now time.Time) *int64 {
	if created.IsZero() || created.After(now) {
		return nil
	}
	sec := int64(now.Sub(created).Seconds())
	if sec < 0 {
		sec = 0
	}
	return &sec
}

func ageFromConntrack(track map[connTuple]int64, c ConnStat) *int64 {
	if len(track) == 0 {
		return nil
	}
	key := connTuple{
		proto: c.Proto,
		lip:   c.LocalIP,
		rip:   c.RemoteIP,
		lport: c.LocalPort,
		rport: c.RemotePort,
	}
	if age, ok := track[key]; ok {
		v := age
		return &v
	}
	// Try swapped direction (conntrack may store reply first).
	key.lip, key.rip = key.rip, key.lip
	key.lport, key.rport = key.rport, key.lport
	if age, ok := track[key]; ok {
		v := age
		return &v
	}
	return nil
}

func loadConntrackAges() map[connTuple]int64 {
	f, err := os.Open(nfConntrackPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	out := map[connTuple]int64{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		t, age, ok := parseConntrackLine(sc.Text())
		if !ok {
			continue
		}
		out[t] = age
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// parseConntrackLine extracts a 4-tuple and delta-time when present.
// Exported for tests via the same package.
func parseConntrackLine(line string) (connTuple, int64, bool) {
	if !strings.Contains(line, deltaTimeKey) {
		return connTuple{}, 0, false
	}
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return connTuple{}, 0, false
	}
	proto := strings.ToLower(fields[0])
	if proto != ConnProtoTCP && proto != ConnProtoUDP {
		// Some kernels prefix with ipv4/ipv6; proto is then fields[2].
		if len(fields) > 2 {
			proto = strings.ToLower(fields[2])
		}
	}
	if proto != ConnProtoTCP && proto != ConnProtoUDP {
		return connTuple{}, 0, false
	}

	var src, dst string
	var sport, dport uint32
	var age int64
	haveAge := false
	for _, f := range fields {
		switch {
		case strings.HasPrefix(f, "src="):
			if src == "" {
				src = strings.TrimPrefix(f, "src=")
			}
		case strings.HasPrefix(f, "dst="):
			if dst == "" {
				dst = strings.TrimPrefix(f, "dst=")
			}
		case strings.HasPrefix(f, "sport="):
			if sport == 0 {
				sport = parseUint32Field(f, "sport=")
			}
		case strings.HasPrefix(f, "dport="):
			if dport == 0 {
				dport = parseUint32Field(f, "dport=")
			}
		case strings.HasPrefix(f, deltaTimeKey):
			v, err := strconv.ParseInt(strings.TrimPrefix(f, deltaTimeKey), 10, 64)
			if err == nil {
				age = v
				haveAge = true
			}
		}
	}
	if !haveAge || src == "" || dst == "" {
		return connTuple{}, 0, false
	}
	return connTuple{proto: proto, lip: src, rip: dst, lport: sport, rport: dport}, age, true
}

func parseUint32Field(f, prefix string) uint32 {
	v, err := strconv.ParseUint(strings.TrimPrefix(f, prefix), 10, 32)
	if err != nil {
		return 0
	}
	return uint32(v)
}

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
	"strings"
)

func formatStatusLine(name, state string, unit Unit, usage Usage) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s", name, state)
	lim := unit.Limits
	if lim.CPUPercent > 0 {
		fmt.Fprintf(&b, " cpu=%.1f%%", usage.CPUPercent)
	}
	if lim.MemoryBytes > 0 {
		fmt.Fprintf(&b, " mem=%d", usage.MemoryBytes)
	}
	if lim.DiskSpaceBytes > 0 {
		fmt.Fprintf(&b, " disk_bytes=%d", usage.DiskSpaceBytes)
	}
	if lim.DiskSpacePercent > 0 {
		fmt.Fprintf(&b, " disk_space=%.1f%%", usage.DiskSpacePercent)
	}
	if lim.DiskReadBPS > 0 {
		fmt.Fprintf(&b, " read_bps=%d", usage.ReadBPS)
	}
	if lim.DiskWriteBPS > 0 {
		fmt.Fprintf(&b, " write_bps=%d", usage.WriteBPS)
	}
	if lim.DiskIOPercent > 0 {
		if usage.IOKnown {
			fmt.Fprintf(&b, " io=%.1f%%", usage.DiskIOPercent)
		} else {
			b.WriteString(" io=unsupported")
		}
	}
	b.WriteByte('\n')
	return b.String()
}

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

// cpuRateFromCores converts a single-core percentage into a Windows job CPU rate.
// percent is percent of one core. ncpu is the machine's logical CPU count.
// The rate is hundredths of a percent of the whole machine, from 1 to 10000.
// The boolean is false when the requested share is above 100% of the machine.
func cpuRateFromCores(percent float64, ncpu int) (uint32, bool) {
	if percent <= 0 || ncpu < 1 {
		return 0, false
	}
	system := percent / float64(ncpu)
	if system > 100 {
		return 0, false
	}
	rate := uint32(system * 100)
	if rate < 1 {
		rate = 1
	}
	if rate > 10000 {
		return 0, false
	}
	return rate, true
}

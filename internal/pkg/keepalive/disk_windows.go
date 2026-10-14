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
	"path/filepath"

	"golang.org/x/sys/windows"
)

func sampleVolume(cwd string) (Volume, error) {
	if cwd == "" {
		return Volume{}, nil
	}
	path, err := windows.UTF16PtrFromString(filepath.Clean(cwd))
	if err != nil {
		return Volume{}, err
	}
	var free, total, avail uint64
	if err := windows.GetDiskFreeSpaceEx(path, &avail, &total, &free); err != nil {
		return Volume{}, err
	}
	if total == 0 {
		return Volume{}, nil
	}
	used := float64(total-free) / float64(total) * 100
	return Volume{UsedPercent: used}, nil
}

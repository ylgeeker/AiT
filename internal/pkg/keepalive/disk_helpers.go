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
	"io/fs"
	"os"
	"path/filepath"
)

// DirUsage sums regular-file sizes under root without following symlinks.
// budget is the maximum number of directory entries visited. Past that, the sum so far is returned
// with ErrWalkCapped so the check fails closed. A missing root returns 0, nil.
func DirUsage(root string, budget int) (uint64, error) {
	if root == "" {
		return 0, nil
	}
	info, err := os.Lstat(root)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return uint64(info.Size()), nil
	}

	var total uint64
	var seen int
	walkErr := filepath.WalkDir(root, func(_ string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		seen++
		if budget > 0 && seen > budget {
			return ErrWalkCapped
		}
		size, skip, ok := regularFileSize(d)
		if skip != nil {
			return skip
		}
		if ok {
			total += size
		}
		return nil
	})
	if walkErr != nil {
		return total, walkErr
	}
	return total, nil
}

func regularFileSize(d fs.DirEntry) (size uint64, walkCtrl error, ok bool) {
	symlink := d.Type()&os.ModeSymlink != 0
	if symlink && d.IsDir() {
		return 0, fs.SkipDir, false
	}
	if symlink || d.IsDir() {
		return 0, nil, false
	}
	st, err := d.Info()
	if err != nil {
		return 0, nil, false
	}
	return uint64(st.Size()), nil, true
}

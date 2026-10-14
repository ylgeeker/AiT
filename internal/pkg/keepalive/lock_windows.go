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
	"os"

	"golang.org/x/sys/windows"
)

// ReleaseFunc unlocks a single-instance lock.
type ReleaseFunc func()

// Acquire takes an exclusive non-blocking lock on path.
// A second caller receives ErrAlreadyRunning. release closes the lock. The operating system
// drops it when the process exits. Returns a path error when the file cannot be created.
func Acquire(path string) (ReleaseFunc, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	handle := windows.Handle(f.Fd())
	err = windows.LockFileEx(handle, windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &windows.Overlapped{})
	if err != nil {
		f.Close()
		return nil, ErrAlreadyRunning
	}
	return func() {
		_ = f.Close()
	}, nil
}

func removeStaleSocket(path string) {}

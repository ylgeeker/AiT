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

import "errors"

var (
	// ErrDaemonDown means the admin endpoint is not accepting connections.
	ErrDaemonDown = errors.New("keepalive serve is not running")

	// ErrAlreadyRunning means another serve holds the AppRoot lock.
	ErrAlreadyRunning = errors.New("keepalive serve is already running")

	// ErrBadName means a unit name is not a single path segment.
	ErrBadName = errors.New("invalid unit name")

	// ErrUnitExists means add was asked to create a name that is already in the file.
	ErrUnitExists = errors.New("unit already exists")

	// ErrUnitMissing means the named unit is not in the file.
	ErrUnitMissing = errors.New("unit not found")

	// ErrNoCommand means a unit has an empty argv.
	ErrNoCommand = errors.New("unit command is empty")

	// ErrDiskFailed means cwd disk usage is still over the limit after the process stopped.
	ErrDiskFailed = errors.New("unit disk space still over limit")

	// ErrOperatorStop means stop or daemon shutdown ended the process.
	ErrOperatorStop = errors.New("unit stopped by operator")

	// ErrWalkCapped means a directory walk hit DiskWalkBudget before finishing.
	ErrWalkCapped = errors.New("directory walk hit entry budget")

	// ErrBadLimit means a percent limit other than cpu_percent is outside 0..100.
	ErrBadLimit = errors.New("limit out of range")
)

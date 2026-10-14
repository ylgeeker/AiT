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

import "time"

const (
	// DefaultListen is the loopback HTTP address used when config and flags omit one.
	DefaultListen = "127.0.0.1:8757"

	// PingPath is the only HTTP route. Other paths respond 404.
	PingPath = "/ping"

	// PongBody is the plain-text body of a successful ping.
	PongBody = "pong\n"

	// ConfigFileName is the unit file under AppRoot.
	ConfigFileName = "keepalive.yaml"

	// SockFileName is the Unix admin socket under AppRoot.
	SockFileName = "keepalive.sock"

	// LockFileName is the single-instance lock under AppRoot.
	LockFileName = "keepalive.lock"

	// pipeNamePrefix is the Windows admin pipe stem. pipeName appends an AppRoot hash.
	pipeNamePrefix = `\\.\pipe\ait-keepalive-`

	// cgroupJoinScript runs in the child so the user command is already in the cgroup at exec.
	cgroupJoinScript = `echo $$ > "$1" || exit 1; shift; exec "$@"`

	// BackoffMin is the first restart delay.
	BackoffMin = time.Second

	// BackoffMax caps the restart delay.
	BackoffMax = 30 * time.Second

	// StableRun resets backoff and the restart budget after a process stays up this long.
	StableRun = 30 * time.Second

	// RestartWindow is the period that counts toward RestartBudget.
	RestartWindow = 10 * time.Minute

	// RestartBudget is how many abnormal restarts are allowed inside RestartWindow.
	RestartBudget = 5

	// Sustain is how long a soft limit may stay exceeded before the process is killed.
	Sustain = 10 * time.Second

	// SampleEvery is the soft-limit sample period.
	SampleEvery = time.Second

	// DutyWindow is the period used to turn a usage ratio into a pause.
	DutyWindow = time.Second

	// StopGrace is how long Unix stop waits after SIGTERM before SIGKILL.
	StopGrace = 3 * time.Second

	// DiskWalkBudget is the maximum directory entries visited for a disk-space check.
	DiskWalkBudget = 10000

	// CPUPeriodUS is the cgroup v2 cpu.max period in microseconds.
	CPUPeriodUS = 100000

	// stateRunning is a unit slot that has a live Watch.
	stateRunning = "running"
	// stateStopped is a unit slot stopped by the operator or a clean exit.
	stateStopped = "stopped"
	// stateFailed is a unit slot that exhausted its restart budget or hit a disk ceiling.
	stateFailed = "failed"

	// shellBin is the interpreter that joins a cgroup before exec.
	shellBin = "/bin/sh"
	// shellArgv0 is the $0 passed to cgroupJoinScript.
	shellArgv0 = "sh"

	// pipeBufSize is the named-pipe read and write buffer in bytes.
	pipeBufSize = 4096

	pingMethod = "GET"

	opStatus = "status"
	opReload = "reload"
	opStart  = "start"
	opStop   = "stop"
)

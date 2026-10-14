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
	"context"
	"time"
)

// WatchResult is the outcome of supervising one unit until it stays stopped.
type WatchResult struct {
	Starts int
	Failed bool
	Sleeps []time.Duration
}

// Supervisor restarts a unit after abnormal exits.
// Run, Now, Sleep, DiskUsed, and Volume may be replaced in tests. Nil fields use the process runner and the wall clock.
type Supervisor struct {
	Run      func(ctx context.Context, unit Unit) (code int, err error)
	Now      func() time.Time
	Sleep    func(time.Duration)
	DiskUsed func(cwd string) (uint64, error)
	Volume   func(cwd string) (Volume, error)
	OnSample func(name string, usage Usage)
}

// Watch runs unit until a clean exit, an operator stop, disk failure, or the restart budget is spent.
// ctx cancel is an operator stop and does not set Failed. A nil Run uses the operating-system process runner.
// The error return is reserved for future runner failures; the current implementation always returns nil.
func (s *Supervisor) Watch(ctx context.Context, unit Unit) (WatchResult, error) {
	var result WatchResult
	backoff := BackoffMin
	var count int
	var windowStart time.Time

	for {
		if ctx.Err() != nil {
			return result, nil
		}
		if s.overDisk(unit) {
			result.Failed = true
			return result, nil
		}

		started := s.now()
		result.Starts++
		code, err := s.run(ctx, unit)
		if ctx.Err() != nil || err == ErrOperatorStop {
			return result, nil
		}
		if err == ErrDiskFailed {
			result.Failed = true
			return result, nil
		}
		if s.now().Sub(started) >= StableRun {
			backoff = BackoffMin
			count = 0
			windowStart = time.Time{}
		}
		if err == nil && code == 0 {
			return result, nil
		}
		if s.overDisk(unit) {
			result.Failed = true
			return result, nil
		}

		now := s.now()
		if windowStart.IsZero() || now.Sub(windowStart) > RestartWindow {
			windowStart = now
			count = 0
		}
		if !AllowRestart(count) {
			result.Failed = true
			return result, nil
		}
		count++
		result.Sleeps = append(result.Sleeps, backoff)
		if err := s.sleep(ctx, backoff); err != nil {
			return result, nil
		}
		backoff = NextBackoff(backoff)
	}
}

func (s *Supervisor) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func (s *Supervisor) sleep(ctx context.Context, d time.Duration) error {
	if s.Sleep != nil {
		s.Sleep(d)
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (s *Supervisor) run(ctx context.Context, unit Unit) (int, error) {
	if s.Run != nil {
		return s.Run(ctx, unit)
	}
	return runProcess(ctx, unit, func(usage Usage) {
		if s.OnSample != nil {
			s.OnSample(unit.Name, usage)
		}
	})
}

func (s *Supervisor) overDisk(unit Unit) bool {
	if unit.Cwd == "" {
		return false
	}
	if unit.Limits.DiskSpaceBytes > 0 {
		used, err := s.disk(unit.Cwd)
		if err != nil || used > unit.Limits.DiskSpaceBytes {
			return true
		}
	}
	if unit.Limits.DiskSpacePercent > 0 {
		vol, err := s.volume(unit.Cwd)
		if err != nil || vol.UsedPercent > unit.Limits.DiskSpacePercent {
			return true
		}
	}
	return false
}

func (s *Supervisor) disk(cwd string) (uint64, error) {
	if s.DiskUsed != nil {
		return s.DiskUsed(cwd)
	}
	return DirUsage(cwd, DiskWalkBudget)
}

func (s *Supervisor) volume(cwd string) (Volume, error) {
	if s.Volume != nil {
		return s.Volume(cwd)
	}
	return volumeSampler(cwd)
}

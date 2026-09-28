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

package zmodem

import (
	"context"
	"fmt"
	"time"
)

func waitDataPhaseSelect(
	s *Session,
	ctx context.Context,
	ticker *time.Ticker,
	handshakeDeadline, overallDeadline time.Time,
) (dataDone bool, teardownDeadline, teardownStarted time.Time, err error) {
	select {
	case <-ctx.Done():
		s.Abort()
		return false, time.Time{}, time.Time{}, ctx.Err()
	case err := <-s.done:
		if err != nil {
			s.Abort()
			return false, time.Time{}, time.Time{}, err
		}
		now := time.Now()
		return true, now.Add(s.teardownLimit()), now, nil
	case <-ticker.C:
		if err := checkHandshakeAndOverallDeadlines(s, handshakeDeadline, overallDeadline); err != nil {
			s.Abort()
			return false, time.Time{}, time.Time{}, err
		}
		return false, time.Time{}, time.Time{}, nil
	}
}

func checkHandshakeAndOverallDeadlines(s *Session, handshakeDeadline, overallDeadline time.Time) error {
	now := time.Now()
	if !s.started.Load() && now.After(handshakeDeadline) {
		return fmt.Errorf("%w", ErrHandshakeTimeout)
	}
	if !now.After(overallDeadline) {
		return nil
	}
	if !s.started.Load() {
		return fmt.Errorf("%w", ErrNoProtocolResponse)
	}
	return fmt.Errorf("zmodem: transfer timed out")
}

func waitTeardownPhaseSelect(
	s *Session,
	ctx context.Context,
	ticker *time.Ticker,
	teardownDeadline, teardownStarted, overallDeadline time.Time,
) (continueLoop bool, err error) {
	select {
	case <-ctx.Done():
		s.Abort()
		return false, ctx.Err()
	case err := <-s.done:
		if err != nil {
			s.Abort()
			return false, err
		}
		return true, nil
	case <-ticker.C:
		return onTeardownTick(s, teardownDeadline, teardownStarted, overallDeadline)
	}
}

func onTeardownTick(
	s *Session,
	teardownDeadline, teardownStarted, overallDeadline time.Time,
) (continueLoop bool, err error) {
	if s.engine.HasUnread() {
		s.engine.PumpUnread()
	}

	if s.isSend && s.engine.SendFileEOF() {
		if s.engine.TryCompleteSendTeardown() {
			return true, nil
		}
		if shouldForceReleaseSendTeardown(teardownStarted) {
			s.engine.ForceRelease()
			return true, nil
		}
	}

	if time.Now().After(teardownDeadline) {
		s.Abort()
		return false, fmt.Errorf("%w", ErrTeardownTimeout)
	}
	if time.Now().After(overallDeadline) {
		s.Abort()
		return false, fmt.Errorf("zmodem: transfer timed out")
	}
	return true, nil
}

func shouldForceReleaseSendTeardown(teardownStarted time.Time) bool {
	if teardownStarted.IsZero() {
		return false
	}
	return time.Since(teardownStarted) >= sendTeardownForceReleaseDuration
}

func advanceWaitDataPhase(
	s *Session,
	ctx context.Context,
	ticker *time.Ticker,
	dataDone *bool,
	teardownDeadline, teardownStarted *time.Time,
	handshakeDeadline, overallDeadline time.Time,
) error {
	if *dataDone {
		return nil
	}

	done, td, ts, err := waitDataPhaseSelect(s, ctx, ticker, handshakeDeadline, overallDeadline)
	if err != nil {
		return err
	}
	if !done {
		return nil
	}

	*dataDone = true
	*teardownDeadline = td
	*teardownStarted = ts
	return nil
}

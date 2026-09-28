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

package pager

import "context"

// runStep handles one select iteration during the main paging loop. stop means
// the session should finish (quit, cancel, or closed key stream).
func (s *session) runStep(ctx context.Context) (stop bool) {
	select {
	case <-ctx.Done():
		return true
	case line, ok := <-s.lines:
		s.consume(line, ok, true)
		return false
	case key, ok := <-s.keys.stream():
		return s.handleRunKey(key, ok)
	}
}

func (s *session) handleRunKey(key byte, ok bool) bool {
	if !ok {
		return true
	}
	if s.handleKey(key) {
		return true
	}

	s.redraw()
	return false
}

// prefillStep waits for one line, key, or cancel during initial buffering.
// quit is true when the user asked to exit before the first frame.
func (s *session) prefillStep(ctx context.Context) (quit bool) {
	select {
	case <-ctx.Done():
		return false
	case line, ok := <-s.lines:
		s.consume(line, ok, false)
		return false
	case key, ok := <-s.keys.stream():
		return !ok || isQuitKey(key)
	}
}

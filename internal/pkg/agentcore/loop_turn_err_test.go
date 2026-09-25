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

package agentcore

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/ai/provider"
)

func TestTurnErrorEndsSession(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		err  error
		ends bool
	}{
		{name: "nil", err: nil, ends: false},
		{name: "malformed", err: formatMalformedReplyErr(), ends: false},
		{name: "generic", err: errors.New("provider 429"), ends: false},
		{name: "canceled", err: context.Canceled, ends: true},
		{name: "deadline", err: context.DeadlineExceeded, ends: true},
		{name: "session ended", err: ErrSessionEnded, ends: true},
		{name: "not configured", err: provider.NotConfiguredErr(), ends: true},
	}
	for _, tc := range cases {
		if got := turnErrorEndsSession(tc.err); got != tc.ends {
			t.Fatalf("%s: turnErrorEndsSession = %v, want %v", tc.name, got, tc.ends)
		}
	}
}

func TestAbsorbTurnError_MalformedEmitsWarning(t *testing.T) {
	t.Parallel()
	var events []Event
	var errOut bytes.Buffer
	cfg := Config{
		ErrOut:  &errOut,
		Emitter: FuncEmitter(func(e Event) { events = append(events, e) }),
	}

	err := absorbTurnError(cfg, formatMalformedReplyErr())
	if err != nil {
		t.Fatalf("absorbTurnError returned %v, want nil", err)
	}
	if len(events) != 1 || events[0].Kind != EventWarning {
		t.Fatalf("events=%v", events)
	}
	if !strings.Contains(events[0].Detail, "could not parse agent response") {
		t.Fatalf("warning detail=%q", events[0].Detail)
	}
	if !strings.Contains(errOut.String(), "could not parse agent response") {
		t.Fatalf("plain UI should print hint, got %q", errOut.String())
	}
}

func TestAbsorbTurnError_ManagedUISkipsStderr(t *testing.T) {
	t.Parallel()
	var errOut bytes.Buffer
	cfg := Config{
		ErrOut:   &errOut,
		ReadLine: func(string) (string, error) { return "", nil },
		Emitter:  NopEmitter,
	}
	if err := absorbTurnError(cfg, formatMalformedReplyErr()); err != nil {
		t.Fatal(err)
	}
	if errOut.Len() != 0 {
		t.Fatalf("managed UI must not print stderr, got %q", errOut.String())
	}
}

func TestAbsorbTurnError_FatalPropagates(t *testing.T) {
	t.Parallel()
	cfg := Config{ErrOut: &bytes.Buffer{}, Emitter: NopEmitter}
	if err := absorbTurnError(cfg, context.Canceled); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
}

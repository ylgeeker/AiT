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

package session

import (
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/ssh/openssh"
)

func TestParseTargetTmuxPassword_LeadingTarget(t *testing.T) {
	hints := SessionFlagHints{
		PasswordRequiresValue: func() error { return errSentinel("pw") },
		UnknownFlag:           func(f string) error { return errSentinel("unk:" + f) },
		UnexpectedExtra:       func(a string) error { return errSentinel("extra:" + a) },
	}
	target, pw, tmux, err := ParseTargetTmuxPassword([]string{"prod", openssh.FlagTmux, openssh.FlagPassword, "secret"}, true, hints)
	if err != nil {
		t.Fatal(err)
	}
	if target != "prod" || pw != "secret" || !tmux {
		t.Fatalf("got target=%q pw=%q tmux=%v", target, pw, tmux)
	}
}

func TestParseTargetTmuxPassword_FlagsFirst(t *testing.T) {
	hints := SessionFlagHints{
		PasswordRequiresValue: func() error { return errSentinel("pw") },
		UnknownFlag:           func(f string) error { return errSentinel("unk:" + f) },
		UnexpectedExtra:       func(a string) error { return errSentinel("extra:" + a) },
	}
	target, pw, tmux, err := ParseTargetTmuxPassword([]string{openssh.FlagTmux, "staging"}, false, hints)
	if err != nil {
		t.Fatal(err)
	}
	if target != "staging" || pw != "" || !tmux {
		t.Fatalf("got target=%q pw=%q tmux=%v", target, pw, tmux)
	}
}

func TestParseTargetTmuxPassword_ExtraPositional(t *testing.T) {
	hints := SessionFlagHints{
		PasswordRequiresValue: func() error { return errSentinel("pw") },
		UnknownFlag:           func(f string) error { return errSentinel("unk:" + f) },
		UnexpectedExtra:       func(a string) error { return errSentinel("extra:" + a) },
	}
	_, _, _, err := ParseTargetTmuxPassword([]string{"a", "b"}, true, hints)
	if err == nil || !strings.Contains(err.Error(), "extra:b") {
		t.Fatalf("err=%v", err)
	}
}

type errSentinel string

func (e errSentinel) Error() string { return string(e) }

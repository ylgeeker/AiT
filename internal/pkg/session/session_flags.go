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

	"github.com/os-artificer/AiT/internal/pkg/ssh/openssh"
)

// SessionFlagHints supplies surface-specific errors for --tmux/--password parsing.
type SessionFlagHints struct {
	PasswordRequiresValue func() error
	UnknownFlag           func(flag string) error
	UnexpectedExtra       func(arg string) error
}

// ParseTargetTmuxPassword walks argv for an optional target id/alias plus --tmux / --password.
// When requireLeadingTarget is true, argv[0] must be the target (shell `load`); otherwise
// the target may appear anywhere among the flags (chain load/switch).
// Errors come from hints (password value required, unknown flag, unexpected extra arg).
func ParseTargetTmuxPassword(argv []string, requireLeadingTarget bool, hints SessionFlagHints) (target, password string, tmux bool, err error) {
	if requireLeadingTarget {
		if len(argv) < 1 {
			return "", "", false, hints.UnexpectedExtra("")
		}
		if strings.HasPrefix(argv[0], "-") {
			return "", "", false, hints.UnknownFlag(argv[0])
		}
		target = argv[0]
		return consumeSessionFlags(argv[1:], target, true, hints)
	}
	return consumeSessionFlags(argv, "", false, hints)
}

func consumeSessionFlags(argv []string, target string, flagsOnly bool, hints SessionFlagHints) (string, string, bool, error) {
	var password string
	var tmux bool

	for i := 0; i < len(argv); i++ {
		arg := argv[i]
		switch arg {
		case openssh.FlagTmux:
			tmux = true
			continue
		case openssh.FlagPassword:
			next, pw, err := takePasswordValue(argv, i, hints)
			if err != nil {
				return "", "", false, err
			}
			password = pw
			i = next
			continue
		}

		newTarget, err := takePositionalTarget(arg, target, flagsOnly, hints)
		if err != nil {
			return "", "", false, err
		}
		target = newTarget
	}
	return target, password, tmux, nil
}

func takePasswordValue(argv []string, i int, hints SessionFlagHints) (next int, password string, err error) {
	next = i + 1
	if next >= len(argv) {
		return 0, "", hints.PasswordRequiresValue()
	}
	return next, argv[next], nil
}

func takePositionalTarget(arg, target string, flagsOnly bool, hints SessionFlagHints) (string, error) {
	if strings.HasPrefix(arg, "-") {
		return "", hints.UnknownFlag(arg)
	}
	if flagsOnly || target != "" {
		return "", hints.UnexpectedExtra(arg)
	}
	return arg, nil
}

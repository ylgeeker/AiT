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

package helpdoc

import (
	"slices"
	"testing"
)

func TestHelpRequested(t *testing.T) {
	if !HelpRequested([]string{"-h"}) || !HelpRequested([]string{"--help"}) || !HelpRequested([]string{"help"}) {
		t.Error("expected help requested")
	}
	if HelpRequested([]string{"host"}) || HelpRequested([]string{"-h", "host"}) {
		t.Error("unexpected help requested")
	}
}

func TestVersionRequested(t *testing.T) {
	if !VersionRequested([]string{"version"}) || !VersionRequested([]string{"-V"}) || !VersionRequested([]string{"--version"}) {
		t.Error("expected version requested")
	}
	if VersionRequested([]string{"host"}) || VersionRequested([]string{"-V", "host"}) {
		t.Error("unexpected version requested")
	}
}

func TestTakePlainFlag(t *testing.T) {
	plain, rest := TakePlainFlag(nil)
	if plain || len(rest) != 0 {
		t.Fatalf("empty: plain=%v rest=%v", plain, rest)
	}

	plain, rest = TakePlainFlag([]string{"--plain"})
	if !plain || len(rest) != 0 {
		t.Fatalf("only flag: plain=%v rest=%v", plain, rest)
	}

	plain, rest = TakePlainFlag([]string{"--plain", "list", "files"})
	if !plain || !slices.Equal(rest, []string{"list", "files"}) {
		t.Fatalf("with goal: plain=%v rest=%v", plain, rest)
	}

	orig := []string{"list", "files"}
	plain, rest = TakePlainFlag(orig)
	if plain || !slices.Equal(rest, orig) {
		t.Fatalf("no flag: plain=%v rest=%v", plain, rest)
	}

	orig = []string{"list", "--plain"}
	plain, rest = TakePlainFlag(orig)
	if plain || !slices.Equal(rest, orig) {
		t.Fatalf("not first: plain=%v rest=%v", plain, rest)
	}
}

func TestHasHelpArg(t *testing.T) {
	if !HasHelpArg([]string{"help"}) {
		t.Error("expected bare help positional")
	}
	if HasHelpArg([]string{"--help"}) || HasHelpArg([]string{"ping"}) {
		t.Error("unexpected has help arg")
	}
}

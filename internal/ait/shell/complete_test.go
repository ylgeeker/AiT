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

package shell

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/aitplugin"
)

func testNamesOfficialPlugins() []string {
	return []string{
		aitplugin.NameHost,
		aitplugin.NameImg,
		aitplugin.NameLog,
		aitplugin.NameSSH,
		aitplugin.CmdAgent,
		aitplugin.CmdClear,
		aitplugin.CmdExit,
		aitplugin.CmdHelp,
		aitplugin.CmdHistory,
		aitplugin.CmdQuit,
		aitplugin.CmdVersion,
	}
}

func TestCompleteFirstToken_UniquePlugin(t *testing.T) {
	// AT30a
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), "ss", 2)
	if display != "" {
		t.Fatalf("display = %q, want empty", display)
	}
	if line != "ssh " || cur != 4 {
		t.Fatalf("got line=%q cur=%d, want ssh /4", line, cur)
	}
}

func TestCompleteFirstToken_AgentPrefix(t *testing.T) {
	// OA10 / AT30b: agent prefix completes to agent + space
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), "ag", 2)
	if display != "" {
		t.Fatalf("display = %q, want empty", display)
	}
	if line != "agent " || cur != 6 {
		t.Fatalf("line should be agent + space, got %q cur=%d", line, cur)
	}
}

func TestCompleteFirstToken_HoToHost(t *testing.T) {
	// OA10: ho → host
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), "ho", 2)
	if display != "" {
		t.Fatalf("display = %q", display)
	}
	if line != "host " || cur != 5 {
		t.Fatalf("got line=%q cur=%d, want host /5", line, cur)
	}
}

func TestCompleteFirstToken_SecondTokenNoop(t *testing.T) {
	// AT30d
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), "ssh ", 4)
	if line != "ssh " || cur != 4 || display != "" {
		t.Fatalf("want noop, got line=%q cur=%d display=%q", line, cur, display)
	}
}

func TestCompleteFirstToken_UniqueHelp(t *testing.T) {
	// AT30e
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), "he", 2)
	if display != "" {
		t.Fatalf("display = %q", display)
	}
	if line != "help " || cur != 5 {
		t.Fatalf("got line=%q cur=%d, want help /5", line, cur)
	}
}

func TestCompleteFirstToken_HListsHelpHistory(t *testing.T) {
	// AT30f
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), "h", 1)
	if line != "h" || cur != 1 {
		t.Fatalf("line unchanged want, got %q cur=%d", line, cur)
	}
	want := "help\nhistory\nhost"
	if display != want {
		t.Fatalf("display = %q, want %q", display, want)
	}
}

func TestCompleteFirstToken_LeadingSpacesPreserved(t *testing.T) {
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), "  ss", 4)
	if display != "" {
		t.Fatalf("display = %q", display)
	}
	if line != "  ssh " {
		t.Fatalf("line = %q, want \"  ssh \"", line)
	}
	if cur != 6 {
		t.Fatalf("cur = %d, want 6", cur)
	}
}

func TestCompleteFirstToken_MidTokenWithArgs(t *testing.T) {
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), "ss foo", 2)
	if display != "" {
		t.Fatalf("display = %q", display)
	}
	if line != "ssh foo" {
		t.Fatalf("line = %q, want ssh foo", line)
	}
	if cur != 3 {
		t.Fatalf("cur = %d, want 3", cur)
	}
}

func TestCompleteFirstToken_QuotedNoop(t *testing.T) {
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), `"ss`, 3)
	if line != `"ss` || cur != 3 || display != "" {
		t.Fatalf("want noop, got line=%q cur=%d display=%q", line, cur, display)
	}
}

func TestCompleteFirstToken_ZeroMatches(t *testing.T) {
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), "xyz", 3)
	if line != "xyz" || cur != 3 || display != "" {
		t.Fatalf("want noop, got line=%q cur=%d display=%q", line, cur, display)
	}
}

func TestCompleteFirstToken_EmptyListsAll(t *testing.T) {
	_, _, display := completeFirstToken(testNamesOfficialPlugins(), "", 0)
	if display == "" {
		t.Fatal("empty line should list candidates")
	}
	if !strings.Contains(display, "ssh") || !strings.Contains(display, "help") {
		t.Fatalf("display missing expected names: %q", display)
	}
}

func TestCompleteFirstToken_InjectedThirdName(t *testing.T) {
	names := append(testNamesOfficialPlugins(), "aicustom", "aixtra")
	_, _, display := completeFirstToken(names, "ai", 2)
	if !strings.Contains(display, "aixtra") || !strings.Contains(display, "aicustom") {
		t.Fatalf("expected third-party ai-prefix names in display: %q", display)
	}
}

func TestCompleteFirstToken_MidCursorPrefix(t *testing.T) {
	line, cur, display := completeFirstToken(testNamesOfficialPlugins(), "ssh", 2)
	if display != "" {
		t.Fatalf("display = %q", display)
	}
	if line != "ssh " || cur != 4 {
		t.Fatalf("got line=%q cur=%d", line, cur)
	}
}

func TestCompleteREPLLine_EmptyPluginsStillReserved(t *testing.T) {
	aitplugin.ResetDiscoverForTest()
	aitplugin.RestoreDiscoverHooksForTest()
	t.Cleanup(func() {
		aitplugin.ResetDiscoverForTest()
		aitplugin.RestoreDiscoverHooksForTest()
	})
	aitplugin.SetLookPathForTest(func(string) (string, error) {
		return "", os.ErrNotExist
	})
	aitplugin.SetWarnWriterForTest(io.Discard)
	aitplugin.SetInstallDirsForTest([]string{t.TempDir()})
	aitplugin.ResetDiscoverForTest()

	line, cur, display := completeREPLLine("ai", 2)
	if line != "ai" || cur != 2 {
		t.Fatalf("line unchanged want, got %q cur=%d", line, cur)
	}
	if display != "" {
		t.Fatalf("display = %q, want empty (no reserved starts with ai)", display)
	}

	line, cur, display = completeREPLLine("he", 2)
	if line != "help " || cur != 5 || display != "" {
		t.Fatalf("reserved help: line=%q cur=%d display=%q", line, cur, display)
	}
}

func TestReservedCommandNames_Sorted(t *testing.T) {
	names := aitplugin.ReservedCommandNames()
	if len(names) != 11 {
		t.Fatalf("len = %d, want 11", len(names))
	}
	for i := 1; i < len(names); i++ {
		if names[i-1] >= names[i] {
			t.Fatalf("not sorted: %v", names)
		}
	}
}

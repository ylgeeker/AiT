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

package tui

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

func TestWrapWithPrefixes_LongAnswerWraps(t *testing.T) {
	text := "I can see you have 3 SSH sessions to ylgeeker.com (one to dev.ylgeeker.com). 2 are suspended and 1 is active, all showing as expired. How can I help you with these sessions? For example: reconnect to one, terminate one, or list them again."
	lines := wrapWithPrefixes("● ", "  ", text, 40)
	if len(lines) < 2 {
		t.Fatalf("expected multiple lines, got %d: %v", len(lines), lines)
	}
	joined := strings.Join(lines, "\n")
	if strings.Contains(joined, "…") {
		t.Fatalf("should not truncate with ellipsis: %q", joined)
	}
	if !strings.Contains(joined, "terminate") {
		t.Fatalf("missing wrapped content: %q", joined)
	}
	if !strings.HasPrefix(lines[0], "● ") {
		t.Fatalf("first line prefix: %q", lines[0])
	}
	for i := 1; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "  ") {
			t.Fatalf("continuation prefix line %d: %q", i, lines[i])
		}
	}
}

func TestRenderBlocks_FinalAnswerWraps(t *testing.T) {
	long := "I can see you have 3 SSH sessions to ylgeeker.com. How can I help you with these sessions? For example: reconnect to one, terminate one, or list them again."
	blocks := []transcriptBlock{
		{kind: blockFinal, raw: long, preview: finalPreviewLines},
	}
	out := renderBlocks(blocks, 50, true)
	if strings.Contains(out, "…") && strings.Contains(out, "term") && !strings.Contains(out, "terminate") {
		t.Fatalf("truncated instead of wrapped: %q", out)
	}
	if !strings.Contains(out, "terminate") {
		t.Fatalf("missing content: %q", out)
	}
	if strings.Count(out, "\n") < 1 {
		t.Fatalf("expected wrapped lines: %q", out)
	}
}

func TestTakeHardBreak_LongToken(t *testing.T) {
	part, rest := takeHardBreak("ylgeeker.com", 5)
	if part == "" || rest == "" {
		t.Fatalf("part=%q rest=%q", part, rest)
	}
	if runewidth.StringWidth(part) > 5 {
		t.Fatalf("part too wide: %q", part)
	}
}

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

package aidisplay_test

import (
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/aidisplay"
	"github.com/os-artificer/AiT/internal/pkg/config"

	_ "github.com/os-artificer/AiT/internal/pkg/ai/display/registry/all"
)

func TestRender_MF1_BoldNotLiteral(t *testing.T) {
	out := aidisplay.Render(aidisplay.ContentKindMarkdown, config.ProviderDeepseek, "**主机** 状态", aidisplay.FormatOpts{Width: 40})
	if strings.Contains(out, "**") {
		t.Fatalf("expected rendered bold, got %q", out)
	}
	if !strings.Contains(out, "主机") {
		t.Fatalf("missing text: %q", out)
	}
}

func TestRender_MF2_ListReadable(t *testing.T) {
	src := "- alpha\n\n- beta"
	out := aidisplay.Render(aidisplay.ContentKindMarkdown, "", src, aidisplay.FormatOpts{Width: 40})
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Fatalf("list items missing: %q", out)
	}
}

func TestRender_MF4_UnregisteredProviderUsesDefault(t *testing.T) {
	out := aidisplay.Render(aidisplay.ContentKindMarkdown, "unknown-vendor", "*item*", aidisplay.FormatOpts{Width: 40})
	if strings.Contains(out, "*item*") {
		t.Fatalf("expected markdown render, got %q", out)
	}
}

func TestRender_MF6_PlainPreservesShellOutput(t *testing.T) {
	src := "total 8\n# not-a-heading\n* not-emphasis"
	out := aidisplay.Render(aidisplay.ContentKindPlain, config.ProviderMinimax, src, aidisplay.FormatOpts{Width: 80})
	if out != src {
		t.Fatalf("plain changed bytes:\n got %q\nwant %q", out, src)
	}
}

func TestRender_MF9_PlainSkipsNormalizer(t *testing.T) {
	src := "<minimax:tool_call></minimax:tool_call>\nls"
	out := aidisplay.Render(aidisplay.ContentKindPlain, config.ProviderMinimax, src, aidisplay.FormatOpts{Width: 80})
	if out != src {
		t.Fatalf("plain must not run minimax normalizer: %q", out)
	}
}

type stubRenderer struct {
	called bool
}

func (s *stubRenderer) RenderMarkdown(text string, width int, thinking bool) (string, error) {
	s.called = true
	return "STUB:" + text, nil
}

func TestRender_MF5_PipelineUsesInjectedRenderer(t *testing.T) {
	stub := &stubRenderer{}
	aidisplay.SetMarkdownRenderer(stub)
	t.Cleanup(func() { aidisplay.SetMarkdownRenderer(nil) })

	out := aidisplay.Render(aidisplay.ContentKindMarkdown, "", "hello", aidisplay.FormatOpts{Width: 40})
	if !stub.called {
		t.Fatal("expected stub renderer")
	}
	if out != "STUB:hello" {
		t.Fatalf("out=%q", out)
	}
}

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

	"github.com/os-artificer/AiT/internal/pkg/agentcore"
)

func TestRenderBlocks_GoalAndTool(t *testing.T) {
	blocks := []transcriptBlock{
		{kind: blockGoal, title: "ping"},
		{kind: blockTool, title: "run ls", lines: []string{"a", "b", "c", "d", "e", "done"}, preview: toolPreviewLines},
	}
	out := renderBlocks(blocks, 80, false)
	if !strings.Contains(out, "> ping") {
		t.Fatalf("missing goal: %q", out)
	}
	if !strings.Contains(out, "● run ls") {
		t.Fatalf("missing tool: %q", out)
	}
	if !strings.Contains(out, "more lines") {
		t.Fatalf("expected fold marker: %q", out)
	}
	expanded := renderBlocks(blocks, 80, true)
	if strings.Contains(expanded, "more lines") {
		t.Fatalf("expanded should not fold: %q", expanded)
	}
}

func TestRenderBlocks_MF1_MarkdownFinal(t *testing.T) {
	blocks := []transcriptBlock{
		{kind: blockFinal, raw: "**主机** 摘要", providerID: "deepseek", preview: finalPreviewLines},
	}
	out := renderBlocks(blocks, 80, true)
	if strings.Contains(out, "**") {
		t.Fatalf("markdown should render bold: %q", out)
	}
	if !strings.Contains(out, "主机") {
		t.Fatalf("missing content: %q", out)
	}
}

func TestApplyEvent_MF8_StoresEventProvider(t *testing.T) {
	m := newModel(newLineBridge(), "v0", "status", "/tmp")
	m.applyEvent(agentcore.Event{Kind: agentcore.EventMessageFinal, Detail: "ok", Provider: "minimax"})
	if len(m.blocks) != 1 || m.blocks[0].providerID != "minimax" {
		t.Fatalf("blocks=%+v", m.blocks)
	}
}

func TestActivity_PlanDoesNotAppendBlock(t *testing.T) {
	m := newModel(newLineBridge(), "v0", "status", "/tmp")
	m.applyEvent(agentcore.Event{Kind: agentcore.EventPlanUpdated, Detail: "step 1/12"})
	if len(m.blocks) != 0 {
		t.Fatalf("plan should not append transcript blocks, got %d", len(m.blocks))
	}
	if !m.live.busy {
		t.Fatal("expected live status busy")
	}
}

func TestWarning_ClearsLiveStatus(t *testing.T) {
	m := newModel(newLineBridge(), "v0", "status", "/tmp")
	m.live.set("Consulting…", "step 1/16")
	m.applyEvent(agentcore.Event{Kind: agentcore.EventWarning, Detail: "could not parse agent response"})
	if m.live.busy {
		t.Fatal("warning should clear live status")
	}
	if len(m.blocks) != 1 || m.blocks[0].kind != blockWarn {
		t.Fatalf("blocks=%v", m.blocks)
	}
}

func TestToolFinished_CommitsBlock(t *testing.T) {
	m := newModel(newLineBridge(), "v0", "status", "/tmp")
	m.applyEvent(agentcore.Event{Kind: agentcore.EventToolStarted, Tool: "run ls"})
	m.applyEvent(agentcore.Event{Kind: agentcore.EventToolOutput, Output: "a\nb"})
	m.applyEvent(agentcore.Event{Kind: agentcore.EventToolFinished})
	if len(m.blocks) != 1 || m.blocks[0].kind != blockTool {
		t.Fatalf("blocks=%v", m.blocks)
	}
	if m.live.busy {
		t.Fatal("live status should clear")
	}
}

func TestApplySize_Adaptive(t *testing.T) {
	m := newModel(newLineBridge(), "v0", "status", "/tmp")
	m.applySize(100, 40)
	if m.viewport.Width != 100 {
		t.Fatalf("viewport width=%d", m.viewport.Width)
	}
	h1 := m.viewport.Height
	m.applySize(100, 20)
	if m.viewport.Height >= h1 {
		t.Fatalf("expected shorter mid on smaller height: %d vs %d", m.viewport.Height, h1)
	}
	view := m.composeView()
	if lineCount(view) != m.height {
		t.Fatalf("view lines=%d want %d", lineCount(view), m.height)
	}
}

func TestConfirm_HighRiskOmitsSession(t *testing.T) {
	var c confirmState
	_ = c.start("⚠ High-risk command (local)\n  rm -rf /tmp/x\n[y/N] ")
	if c.allowSession {
		t.Fatal("high risk should omit session option")
	}
	if c.form == nil {
		t.Fatal("expected form")
	}
}

func TestConfirm_ResultTokens(t *testing.T) {
	var c confirmState
	c.choice = confirmYes
	if c.resultToken() != "y" {
		t.Fatal(c.resultToken())
	}
	c.choice = confirmSession
	if c.resultToken() != "s" {
		t.Fatal(c.resultToken())
	}
}

func TestHeader_RecentActivity(t *testing.T) {
	out := renderHeader(80, "v0.1.0", "minimax", "/tmp", []string{"ping", "/help"}, false)
	if !strings.Contains(out, "Recent activity") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "ping") {
		t.Fatal(out)
	}
	if strings.Contains(out, "Activity\n") && strings.Contains(out, "Ready") {
		t.Fatal("should not have old Activity panel Ready")
	}
}

func TestFoldLines(t *testing.T) {
	lines := []string{"1", "2", "3", "4", "5"}
	got := foldLines(lines, 2, false)
	if len(got) != 3 || !strings.Contains(got[2], "Ctrl+O") {
		t.Fatalf("%v", got)
	}
}

func TestParseStepAndPlanned(t *testing.T) {
	if _, ok := parseStepLabel("step 1/16"); !ok {
		t.Fatal("step")
	}
	got, ok := parsePlannedLabel("planned: run ls")
	if !ok || got != "run ls" {
		t.Fatalf("%q %v", got, ok)
	}
}

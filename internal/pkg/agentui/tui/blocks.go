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
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

const (
	thinkPreviewLines = 2
	toolPreviewLines  = 4
	finalPreviewLines = 12
)

type blockKind int

const (
	blockGoal blockKind = iota
	blockThinking
	blockTool
	blockFinal
	blockWarn
	blockMeta
)

type transcriptBlock struct {
	kind       blockKind
	title      string
	raw        string
	lines      []string
	providerID string
	folded     bool
	preview    int
}

type liveStatus struct {
	busy      bool
	label     string
	detail    string
	startedAt time.Time
	spinner   spinner.Model
}

func newLiveStatus() liveStatus {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = processStyle
	return liveStatus{spinner: s}
}

func (l *liveStatus) set(label, detail string) {
	wasBusy := l.busy
	l.busy = true
	l.label = label
	l.detail = detail
	if !wasBusy {
		l.startedAt = time.Now()
	}
}

func (l *liveStatus) clear() {
	l.busy = false
	l.label = ""
	l.detail = ""
}

func (l liveStatus) render() string {
	if !l.busy {
		return ""
	}
	elapsed := time.Since(l.startedAt).Round(time.Second)
	text := l.label
	if l.detail != "" {
		text += ": " + l.detail
	}
	return processStyle.Render(fmt.Sprintf("%s %s (%s)", l.spinner.View(), text, elapsed))
}

func formatToolLabel(tool string) string {
	tool = strings.TrimSpace(tool)
	if strings.HasPrefix(tool, "run ") {
		return "run " + strings.TrimSpace(strings.TrimPrefix(tool, "run "))
	}
	return tool
}

func foldLines(lines []string, preview int, expanded bool) []string {
	if expanded || preview <= 0 || len(lines) <= preview {
		return lines
	}
	out := append([]string{}, lines[:preview]...)
	out = append(out, fmt.Sprintf("… %d more lines (Ctrl+O to expand)", len(lines)-preview))
	return out
}

func writeWrappedLines(b *strings.Builder, style lipgloss.Style, firstPrefix, contPrefix, text string, width int) {
	for i, line := range wrapWithPrefixes(firstPrefix, contPrefix, text, width) {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(style.Render(line))
	}
}

func renderBlocks(blocks []transcriptBlock, width int, expanded bool) string {
	if width < 8 {
		width = 8
	}
	var b strings.Builder
	for i, blk := range blocks {
		if i > 0 {
			b.WriteString("\n")
		}
		switch blk.kind {
		case blockGoal:
			writeWrappedLines(&b, goalStyle.Width(width), "> ", "  ", blk.title, width)
		case blockThinking:
			lines := foldLines(renderMarkdownBlock(blk.raw, blk.providerID, width, true), thinkPreviewLines, expanded)
			writePrefixedLines(&b, thinkingStyle, false, "  ", "  ", lines)
		case blockTool:
			writeWrappedLines(&b, toolHeaderStyle, "● ", "  ", blk.title, width)
			b.WriteString("\n")
			writeToolFoldedLines(&b, renderPlainBlock(blk.lines, width), expanded)
		case blockFinal:
			lines := foldLines(renderMarkdownBlock(blk.raw, blk.providerID, width, false), finalPreviewLines, expanded)
			writePrefixedLines(&b, answerStyle, false, "● ", "  ", lines)
		case blockWarn:
			writeWrappedLines(&b, warnStyle, "! ", "  ", blk.title, width)
		case blockMeta:
			writeWrappedLines(&b, statusStyle, "", "  ", blk.title, width)
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func parseStepLabel(detail string) (label string, ok bool) {
	detail = strings.TrimSpace(detail)
	if !strings.HasPrefix(detail, "step ") {
		return "", false
	}
	return "Consulting…", true
}

func parsePlannedLabel(detail string) (planned string, ok bool) {
	detail = strings.TrimSpace(detail)
	if !strings.HasPrefix(detail, "planned: ") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(detail, "planned: ")), true
}

func splitNonEmptyLines(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	raw := strings.Split(text, "\n")
	out := make([]string, 0, len(raw))
	for _, ln := range raw {
		ln = strings.TrimRight(ln, "\r")
		if strings.TrimSpace(ln) == "" {
			continue
		}
		out = append(out, ln)
	}
	return out
}

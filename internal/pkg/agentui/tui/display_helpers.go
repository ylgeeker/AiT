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

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

	"github.com/os-artificer/AiT/internal/pkg/aidisplay"
)

func blockContentWidth(totalWidth int, firstPrefix, contPrefix string) int {
	overhead := runewidth.StringWidth(firstPrefix)
	if w := runewidth.StringWidth(contPrefix); w > overhead {
		overhead = w
	}
	w := totalWidth - overhead
	if w < 8 {
		return 8
	}
	return w
}

func splitRenderedLines(text string) []string {
	text = strings.TrimRight(text, "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

func renderMarkdownBlock(raw, providerID string, width int, thinking bool) []string {
	firstPrefix := "● "
	contPrefix := "  "
	if thinking {
		firstPrefix = "  "
	}
	contentW := blockContentWidth(width, firstPrefix, contPrefix)
	rendered := aidisplay.Render(aidisplay.ContentKindMarkdown, providerID, raw, aidisplay.FormatOpts{
		Width:    contentW,
		Thinking: thinking,
	})
	lines := splitRenderedLines(rendered)
	if len(lines) == 0 && rendered != "" {
		return []string{rendered}
	}
	return lines
}

func renderPlainBlock(lines []string, width int) []string {
	if len(lines) == 0 {
		return nil
	}
	text := strings.Join(lines, "\n")
	contentW := blockContentWidth(width, "└──── ", "      ")
	rendered := aidisplay.Render(aidisplay.ContentKindPlain, "", text, aidisplay.FormatOpts{Width: contentW})
	out := splitRenderedLines(rendered)
	if len(out) == 0 && rendered != "" {
		return []string{rendered}
	}
	return out
}

func writePrefixedLines(b *strings.Builder, style lipgloss.Style, applyStyle bool, firstPrefix, contPrefix string, lines []string) {
	for j, ln := range lines {
		if j > 0 {
			b.WriteString("\n")
		}
		prefix := contPrefix
		if j == 0 {
			prefix = firstPrefix
		}
		line := prefix + ln
		if applyStyle {
			line = style.Render(line)
		}
		b.WriteString(line)
	}
}

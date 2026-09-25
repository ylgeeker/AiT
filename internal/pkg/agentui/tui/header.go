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

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

const (
	minTerminalWidth  = 40
	minTerminalHeight = 12
	maxRecentItems    = 5
)

func clipTextWidth(s string, width int) string {
	if width <= 0 {
		return s
	}
	if runewidth.StringWidth(s) <= width {
		return s
	}
	return runewidth.Truncate(s, width, "…")
}

func innerWidth(outer int) int {
	w := outer - 4
	if w < 8 {
		return 8
	}
	return w
}

func lineCount(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

func splitLines(s string) []string {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func renderHeader(width int, version, status, cwd string, recent []string, compact bool) string {
	title := headerTitle.Render(fmt.Sprintf("AiT Agent %s", version))

	var tipsLines []string
	tipsLines = append(tipsLines, headerSection.Render("Tips for getting started"))
	if !compact {
		tipsLines = append(tipsLines,
			"  Press / for commands; type natural language goals",
			"  Run /help for the full command list",
			"  Shell commands require approval",
		)
	} else {
		tipsLines = append(tipsLines, "  Press / for commands")
	}

	tipsLines = append(tipsLines, "", headerSection.Render("Recent activity"))
	if len(recent) == 0 {
		tipsLines = append(tipsLines, "  No recent activity")
	} else {
		limit := maxRecentItems
		if compact {
			limit = 2
		}
		for i, item := range recent {
			if i >= limit {
				break
			}
			tipsLines = append(tipsLines, "  "+clipTextWidth(item, innerWidth(width)-2))
		}
	}

	logo := logoStyle.Render(logoArt)
	right := lipgloss.JoinVertical(lipgloss.Left, title, "", strings.Join(tipsLines, "\n"))
	top := lipgloss.JoinHorizontal(lipgloss.Top, logo, "  ", right)

	footer := statusStyle.Render(clipTextWidth(strings.TrimSpace(status+" · "+cwd), innerWidth(width)))
	inner := lipgloss.JoinVertical(lipgloss.Left, top, "", footer)

	box := headerBox
	if width > 4 {
		box = box.Width(width - 2)
	}
	return box.Render(inner)
}

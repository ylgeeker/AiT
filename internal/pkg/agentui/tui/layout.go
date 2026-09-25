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

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) applySize(width, height int) {
	if width < minTerminalWidth {
		width = minTerminalWidth
	}
	if height < minTerminalHeight {
		height = minTerminalHeight
	}
	atBottom := m.viewport.AtBottom() || !m.ready
	m.width = width
	m.height = height
	m.compact = height < 20 || width < 60

	header := renderHeader(m.width, m.version, m.status, m.cwd, m.recent, m.compact)
	m.headerCache = header
	headerH := lineCount(header)

	statusH := 0
	if m.live.busy {
		statusH = 1
	}

	footerH := 3 // separator + input + hint
	if m.confirm.active && m.confirm.form != nil {
		footerH += 6
	}
	if menu := m.renderSlashMenu(); menu != "" {
		footerH += lineCount(menu) + 1
	}

	midH := height - headerH - 1 - statusH - 1 - footerH
	if midH < 1 {
		midH = 1
	}

	m.viewport.Width = width
	m.viewport.Height = midH
	m.input.Width = max(1, width-2)
	m.refreshViewportContent()
	if atBottom {
		m.viewport.GotoBottom()
	}
}

func (m *model) refreshViewportContent() {
	content := renderBlocks(m.blocks, m.width, m.expanded)
	m.viewport.SetContent(content)
}

func (m model) composeView() string {
	if m.width < 1 {
		m.width = 80
	}
	if m.height < 1 {
		m.height = 24
	}

	header := m.headerCache
	if header == "" {
		header = renderHeader(m.width, m.version, m.status, m.cwd, m.recent, m.compact)
	}

	var mid strings.Builder
	mid.WriteString(m.viewport.View())

	statusLine := m.live.render()
	sep := separatorStyle.Render(strings.Repeat("─", max(1, m.width)))

	footer := m.renderFooter()

	parts := []string{header, "", mid.String()}
	if statusLine != "" {
		parts = append(parts, statusLine)
	}
	parts = append(parts, sep, footer)

	view := strings.Join(parts, "\n")
	lines := splitLines(view)
	if len(lines) > m.height {
		// Prefer keeping footer: drop from header/mid top if needed.
		lines = lines[len(lines)-m.height:]
	}
	for len(lines) < m.height {
		// Pad above footer (before last footer block) — simple pad at mid.
		insert := len(lines) - lineCount(footer) - 1
		if insert < 0 {
			insert = 0
		}
		lines = append(lines[:insert], append([]string{""}, lines[insert:]...)...)
	}
	return strings.Join(lines[:m.height], "\n")
}

func (m model) renderFooter() string {
	var b strings.Builder
	if m.confirm.active && m.confirm.form != nil {
		b.WriteString(m.confirm.form.View())
		b.WriteString("\n")
	}
	b.WriteString(m.input.View())
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("? for shortcuts"))
	if menu := m.renderSlashMenu(); menu != "" {
		b.WriteString("\n")
		b.WriteString(menu)
	}
	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func newViewport(w, h int) viewport.Model {
	vp := viewport.New(w, h)
	vp.Style = lipgloss.NewStyle()
	return vp
}

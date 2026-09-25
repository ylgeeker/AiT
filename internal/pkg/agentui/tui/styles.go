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

import "github.com/charmbracelet/lipgloss"

const colorCyan = "#00D7FF"

var (
	headerBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(colorCyan)).
			Padding(0, 1)

	headerTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCyan)).
			Bold(true)

	headerSection = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCyan)).
			Bold(true)

	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCyan))

	statusStyle = lipgloss.NewStyle().Faint(true)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCyan)).
			Bold(true)

	hintStyle = lipgloss.NewStyle().Faint(true)

	goalStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A2F1F")).
			Foreground(lipgloss.Color("#FFFFFF"))

	thinkingStyle = lipgloss.NewStyle().
			Faint(true).
			Italic(true)

	toolHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7CFC9A"))

	toolOutputStyle = lipgloss.NewStyle().Faint(true)

	answerStyle = lipgloss.NewStyle()

	warnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC66"))

	processStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCyan))

	menuName = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCyan))

	menuDesc = lipgloss.NewStyle().Faint(true)

	confirmSelected = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCyan)).
			Bold(true)

	confirmOption = lipgloss.NewStyle()

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCyan)).
			Faint(true)
)

const logoArt = ` ▄▀▀▀▄
 █▀█▀█
 ▀ ▀ ▀`

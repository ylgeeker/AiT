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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// confirmChoice maps CodeBuddy-style menu picks to approve.ParseChoice tokens.
type confirmChoice string

const (
	confirmYes     confirmChoice = "y"
	confirmSession confirmChoice = "s"
	confirmNo      confirmChoice = "n"
)

type confirmState struct {
	active       bool
	highRisk     bool
	allowSession bool
	cmd          string
	prompt       string
	form         *huh.Form
	choice       confirmChoice
}

func (c *confirmState) start(prompt string) tea.Cmd {
	c.active = true
	c.prompt = prompt
	c.highRisk = strings.Contains(prompt, "High-risk")
	c.allowSession = !c.highRisk
	c.cmd = extractConfirmCmd(prompt)
	c.choice = confirmNo
	c.form = buildConfirmForm(c)
	return c.form.Init()
}

const confirmDescWidth = 72

func buildConfirmForm(c *confirmState) *huh.Form {
	opts := []huh.Option[confirmChoice]{
		huh.NewOption("Yes", confirmYes),
	}
	if c.allowSession {
		opts = append(opts, huh.NewOption("Yes, and don't ask again for session", confirmSession))
	}
	opts = append(opts, huh.NewOption("No", confirmNo))

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[confirmChoice]().
				Title("Do you want to proceed?").
				Description(wrapPlainText(c.cmd, confirmDescWidth)).
				Options(opts...).
				Value(&c.choice),
		),
	).WithTheme(huh.ThemeBase())
}

func (c *confirmState) clear() {
	c.active = false
	c.form = nil
	c.prompt = ""
	c.cmd = ""
	c.choice = confirmNo
}

func (c *confirmState) resultToken() string {
	if c.choice == "" {
		return string(confirmNo)
	}
	return string(c.choice)
}

func extractConfirmCmd(prompt string) string {
	lines := strings.Split(prompt, "\n")
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.HasPrefix(ln, "Run command?") || strings.HasPrefix(ln, "⚠") || strings.HasPrefix(ln, "[y/N]") {
			continue
		}
		return ln
	}
	return strings.TrimSpace(prompt)
}

func isConfirmPrompt(text string) bool {
	t := strings.TrimSpace(text)
	return strings.Contains(t, "Run command?") || strings.Contains(t, "High-risk command") || strings.HasPrefix(t, "[y/N]")
}

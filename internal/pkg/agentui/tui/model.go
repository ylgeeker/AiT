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
	"os"
	"strings"

	"github.com/os-artificer/AiT/internal/pkg/agentcore"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type model struct {
	bridge      *lineBridge
	version     string
	status      string
	cwd         string
	width       int
	height      int
	ready       bool
	compact     bool
	expanded    bool
	headerCache string
	recent      []string
	blocks      []transcriptBlock
	pendingTool string
	pendingOut  []string
	live        liveStatus
	confirm     confirmState
	viewport    viewport.Model
	input       textinput.Model
}

func newModel(b *lineBridge, version, status, cwd string) model {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.PromptStyle = promptStyle
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 40

	return model{
		bridge:   b,
		version:  version,
		status:   status,
		cwd:      cwd,
		live:     newLiveStatus(),
		viewport: newViewport(40, 10),
		input:    ti,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.waitBridge(), m.live.spinner.Tick, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.applySize(msg.Width, msg.Height)
		m.ready = true
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.live.spinner, cmd = m.live.spinner.Update(msg)
		if m.live.busy {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case bridgeMsg:
		if msg.kind == bridgeDone {
			return m, tea.Quit
		}
		if cmd := m.applyBridge(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, m.waitBridge())
		if m.live.busy {
			cmds = append(cmds, m.live.spinner.Tick)
		}
		m.applySize(m.width, m.height)
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		if m.confirm.active && m.confirm.form != nil {
			return m.handleConfirmKey(msg)
		}
		return m.handleKey(msg)

	default:
		return m.handleConfirmFormUpdate(msg)
	}
}

func (m model) handleConfirmFormUpdate(msg tea.Msg) (model, tea.Cmd) {
	if !m.confirm.active || m.confirm.form == nil {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	form, fcmd := m.confirm.form.Update(msg)
	f, ok := form.(*huh.Form)
	if !ok {
		return m, fcmd
	}
	m.confirm.form = f
	if f.State != huh.StateCompleted {
		return m, fcmd
	}

	token := m.confirm.resultToken()
	m.confirm.clear()
	m.live.clear()
	m.bridge.sendLine(token)
	m.applySize(m.width, m.height)
	return m, tea.Batch(fcmd, m.waitBridge())
}

func (m model) View() string {
	if !m.ready {
		return "Loading...\n"
	}
	return m.composeView()
}

func (m model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "ctrl+d":
		m.bridge.close()
		return m, tea.Quit
	case "esc":
		m.confirm.choice = confirmNo
		m.confirm.clear()
		m.live.clear()
		m.bridge.sendLine(string(confirmNo))
		m.applySize(m.width, m.height)
		return m, m.waitBridge()
	case "1":
		m.confirm.choice = confirmYes
		token := m.confirm.resultToken()
		m.confirm.clear()
		m.live.clear()
		m.bridge.sendLine(token)
		m.applySize(m.width, m.height)
		return m, m.waitBridge()
	case "2":
		if m.confirm.allowSession {
			m.confirm.choice = confirmSession
		} else {
			m.confirm.choice = confirmNo
		}
		token := m.confirm.resultToken()
		m.confirm.clear()
		m.live.clear()
		m.bridge.sendLine(token)
		m.applySize(m.width, m.height)
		return m, m.waitBridge()
	case "3":
		m.confirm.choice = confirmNo
		token := m.confirm.resultToken()
		m.confirm.clear()
		m.live.clear()
		m.bridge.sendLine(token)
		m.applySize(m.width, m.height)
		return m, m.waitBridge()
	}

	var fcmd tea.Cmd
	var form tea.Model
	form, fcmd = m.confirm.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.confirm.form = f
		if f.State == huh.StateCompleted {
			token := m.confirm.resultToken()
			m.confirm.clear()
			m.live.clear()
			m.bridge.sendLine(token)
			m.applySize(m.width, m.height)
			return m, tea.Batch(fcmd, m.waitBridge())
		}
	}
	return m, fcmd
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "ctrl+d":
		m.bridge.close()
		return m, tea.Quit
	case "ctrl+o":
		m.expanded = !m.expanded
		m.refreshViewportContent()
		return m, nil
	case "enter":
		line := strings.TrimSpace(m.input.Value())
		m.input.SetValue("")
		if line == "" {
			return m, nil
		}
		m.addRecent(line)
		m.appendGoal(line)
		m.bridge.sendLine(line)
		m.refreshViewportContent()
		m.viewport.GotoBottom()
		return m, nil
	case "esc":
		m.input.SetValue("")
		return m, nil
	case "up", "down", "pgup", "pgdown":
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *model) applyBridge(msg bridgeMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.kind {
	case bridgeTranscript:
		cmd = m.handleTranscript(msg.text)
	case bridgeEvent:
		m.applyEvent(msg.event)
	}
	m.refreshViewportContent()
	m.viewport.GotoBottom()
	return cmd
}

func (m *model) handleTranscript(text string) tea.Cmd {
	if isConfirmPrompt(text) {
		cmd := m.confirm.start(text)
		m.live.set("Awaiting approval", extractConfirmCmd(text))
		return cmd
	}
	for _, ln := range splitNonEmptyLines(text) {
		switch {
		case strings.HasPrefix(ln, "Slash commands:"), strings.HasPrefix(ln, "  /"), strings.HasPrefix(ln, "Type a goal"):
			m.blocks = append(m.blocks, transcriptBlock{kind: blockMeta, title: ln})
		case strings.HasPrefix(ln, "Conversation cleared"):
			m.blocks = append(m.blocks, transcriptBlock{kind: blockMeta, title: "└ " + ln})
		case strings.HasPrefix(ln, "⊘"):
			m.blocks = append(m.blocks, transcriptBlock{kind: blockMeta, title: "└ " + strings.TrimPrefix(ln, "⊘ ")})
		default:
			m.blocks = append(m.blocks, transcriptBlock{kind: blockMeta, title: ln})
		}
	}
	return nil
}

func (m *model) applyEvent(e agentcore.Event) {
	switch e.Kind {
	case agentcore.EventPlanUpdated:
		if _, ok := parseStepLabel(e.Detail); ok {
			m.live.set("Consulting…", strings.TrimSpace(e.Detail))
			return
		}
		if planned, ok := parsePlannedLabel(e.Detail); ok {
			m.live.set("Planned", planned)
		}
	case agentcore.EventThinking:
		raw := strings.TrimSpace(e.Detail)
		if raw == "" {
			return
		}
		m.blocks = append(m.blocks, transcriptBlock{
			kind:       blockThinking,
			raw:        raw,
			providerID: e.Provider,
			preview:    thinkPreviewLines,
		})
		m.live.set("Consulting…", "")
	case agentcore.EventToolStarted:
		m.pendingTool = formatToolLabel(e.Tool)
		m.pendingOut = nil
		m.live.set("Running", m.pendingTool)
	case agentcore.EventToolOutput:
		for _, ln := range splitNonEmptyLines(e.Output) {
			m.pendingOut = append(m.pendingOut, ln)
		}
	case agentcore.EventToolFinished:
		result := "done"
		if e.Detail == "skipped" {
			result = "skipped"
		}
		lines := append([]string{}, m.pendingOut...)
		lines = append(lines, result)
		m.blocks = append(m.blocks, transcriptBlock{
			kind:    blockTool,
			title:   m.pendingTool,
			lines:   lines,
			preview: toolPreviewLines,
		})
		m.pendingTool = ""
		m.pendingOut = nil
		m.live.clear()
	case agentcore.EventApprovalRequired:
		m.live.set("Awaiting approval", formatToolLabel(e.Tool))
	case agentcore.EventApprovalResolved:
		if e.Detail == "auto-read" {
			return
		}
		if !e.Allowed {
			m.blocks = append(m.blocks, transcriptBlock{
				kind:  blockTool,
				title: formatToolLabel(e.Tool),
				lines: []string{"denied"},
			})
			m.live.clear()
		}
	case agentcore.EventMessageFinal:
		raw := strings.TrimSpace(e.Detail)
		if raw == "" {
			raw = "(done)"
		}
		m.blocks = append(m.blocks, transcriptBlock{
			kind:       blockFinal,
			raw:        raw,
			providerID: e.Provider,
			preview:    finalPreviewLines,
		})
		m.live.clear()
	case agentcore.EventWarning:
		m.live.clear()
		m.blocks = append(m.blocks, transcriptBlock{kind: blockWarn, title: e.Detail})
	case agentcore.EventStepLimitReached:
		m.live.clear()
		m.blocks = append(m.blocks, transcriptBlock{
			kind:  blockWarn,
			title: fmt.Sprintf("step limit %s reached", e.Detail),
		})
	}
}

func (m *model) appendGoal(line string) {
	m.blocks = append(m.blocks, transcriptBlock{kind: blockGoal, title: line})
}

func (m *model) addRecent(line string) {
	m.recent = append([]string{line}, m.recent...)
	if len(m.recent) > maxRecentItems {
		m.recent = m.recent[:maxRecentItems]
	}
	m.headerCache = renderHeader(m.width, m.version, m.status, m.cwd, m.recent, m.compact)
}

func (m model) renderSlashMenu() string {
	val := m.input.Value()
	if !strings.HasPrefix(val, "/") {
		return ""
	}
	prefix := strings.ToLower(strings.TrimSpace(val))
	var b strings.Builder
	for _, cmd := range agentcore.SlashCommands() {
		if prefix != "/" && !strings.HasPrefix(strings.ToLower(cmd.Name), prefix) {
			continue
		}
		fmt.Fprintf(&b, "%s  %s\n", menuName.Render(cmd.Name), menuDesc.Render(cmd.Desc))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) waitBridge() tea.Cmd {
	return func() tea.Msg {
		select {
		case msg, ok := <-m.bridge.msgs:
			if !ok {
				return bridgeMsg{kind: bridgeDone}
			}
			return msg
		case <-m.bridge.closed:
			return bridgeMsg{kind: bridgeDone}
		}
	}
}

func defaultCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}

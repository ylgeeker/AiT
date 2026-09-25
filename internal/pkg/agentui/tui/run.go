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

// Package tui provides the Bubble Tea agent UI (CodeBuddy / Cursor CLI style).
package tui

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/os-artificer/AiT/internal/pkg/agentcore"
	"github.com/os-artificer/AiT/internal/pkg/ai"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

// RunLoop runs the agent task loop behind the interactive TUI on TTY.
func RunLoop(ctx context.Context, cfg agentcore.Config, initialGoal, version string) error {
	in := cfg.In
	if in == nil {
		in = os.Stdin
	}
	out := cfg.Out
	if out == nil {
		out = os.Stdout
	}
	if !shouldUseTUI(in, out) {
		return agentcore.RunLoop(ctx, cfg, initialGoal)
	}

	bridge := newLineBridge()
	status := "ops agent"
	if cfg.AI != nil {
		status = cfg.AI.StatusLine()
	}
	cwd := defaultCWD()
	if cfg.Scope == ai.ScopeRemote && cfg.HostLabel != "" {
		cwd = cfg.HostLabel
	}

	cfg.ReadLine = bridge.readLine
	cfg.ConfirmRead = bridge.confirmRead
	cfg.PlainUI = true
	cfg.Out = bridge.transcriptWriter()
	cfg.Emitter = agentcore.FuncEmitter(bridge.emit)
	cfg.SuspendTTY = func() func() { return func() {} }
	cfg.In = in

	done := make(chan error, 1)
	go func() {
		err := agentcore.RunLoop(ctx, cfg, initialGoal)
		bridge.finish()
		done <- err
	}()

	m := newModel(bridge, version, status, cwd)
	p := tea.NewProgram(
		m,
		tea.WithContext(ctx),
		tea.WithInput(in),
		tea.WithOutput(os.Stdout),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		bridge.close()
		return fmt.Errorf("agent tui: %w", err)
	}
	return <-done
}

func shouldUseTUI(in io.Reader, out io.Writer) bool {
	inFile, inOK := in.(*os.File)
	outFile, outOK := out.(*os.File)
	if !inOK || !outOK {
		return false
	}
	return term.IsTerminal(int(inFile.Fd())) && term.IsTerminal(int(outFile.Fd()))
}

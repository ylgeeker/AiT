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

// Package agent wires ait agent / ait> agent into agentcore + plain UI.
package agent

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/os-artificer/AiT/internal/ait/version"
	"github.com/os-artificer/AiT/internal/pkg/agentcore"
	"github.com/os-artificer/AiT/internal/pkg/agentops"
	"github.com/os-artificer/AiT/internal/pkg/agentui/plain"
	"github.com/os-artificer/AiT/internal/pkg/agentui/tui"
	"github.com/os-artificer/AiT/internal/pkg/agentwire"
	pkgai "github.com/os-artificer/AiT/internal/pkg/ai"
	"github.com/os-artificer/AiT/internal/pkg/config"
	"github.com/os-artificer/AiT/internal/pkg/store"

	"golang.org/x/term"
)

// Options configures an ait agent session.
type Options struct {
	Nested bool
	Goal   string
	Plain  bool
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

// RunLoop starts the v4 agent task loop with experimental gate and plain UI.
func RunLoop(ctx context.Context, opt Options) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	st, err := store.Open()
	if err != nil {
		return err
	}
	defer st.Close()

	aiClient, err := pkgai.NewClient(cfg)
	if err != nil {
		return pkgai.NotConfiguredErr()
	}

	in, out, errOut := opt.In, opt.Out, opt.ErrOut
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stderr
	}

	renderer := plain.Renderer{Out: out}
	warn := func(msg string) {
		renderer.NewEmitter().Emit(agentcore.Event{Kind: agentcore.EventWarning, Detail: msg})
	}
	envExtra := agentops.BuildContext(ctx, warn)

	maxSteps := 0
	var cmdTimeout time.Duration
	if cfg != nil {
		maxSteps = cfg.AgentMaxSteps()
		cmdTimeout = cfg.AgentCommandTimeout()
	}

	acfg := agentwire.LocalLoopConfig(agentwire.LocalOpts{
		AI:             aiClient,
		Nested:         opt.Nested,
		Plain:          opt.Plain,
		In:             in,
		Out:            out,
		ErrOut:         errOut,
		MaxSteps:       maxSteps,
		CommandTimeout: cmdTimeout,
		EnvExtra:       envExtra,
		ExecCLI:        execAisshCLI,
	})

	goal := opt.Goal
	if goal == "" && !isTTY(in) {
		return agentcore.RunLoop(ctx, acfg, "")
	}
	if isTTY(in) && !opt.Plain {
		if outFile, outOK := out.(*os.File); outOK && term.IsTerminal(int(outFile.Fd())) {
			return tui.RunLoop(ctx, acfg, goal, version.Short())
		}
	}
	return agentcore.RunLoop(ctx, acfg, goal)
}

func isTTY(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

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

package agentwire

import (
	"io"
	"os"
	"time"

	"github.com/os-artificer/AiT/internal/pkg/agentcore"
	"github.com/os-artificer/AiT/internal/pkg/agentops"
	"github.com/os-artificer/AiT/internal/pkg/agentui/plain"
	pkgai "github.com/os-artificer/AiT/internal/pkg/ai"
	"github.com/os-artificer/AiT/internal/pkg/ai/apitypes"
)

// CLIRunner runs an external CLI with argv and writers.
type CLIRunner func(argv []string, out, errOut io.Writer) error

// LocalOpts configures a local Agent Loop session.
type LocalOpts struct {
	AI             *pkgai.Client
	Nested         bool
	Plain          bool
	In             io.Reader
	Out            io.Writer
	ErrOut         io.Writer
	MaxSteps       int
	CommandTimeout time.Duration
	EnvExtra       string
	ExecCLI        CLIRunner
}

// LocalLoopConfig builds agentcore.Config for a local ops Agent Loop.
func LocalLoopConfig(opts LocalOpts) agentcore.Config {
	in, out, errOut := opts.In, opts.Out, opts.ErrOut
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
	return agentcore.Config{
		Scope:          pkgai.ScopeLocal,
		AI:             opts.AI,
		CLI:            opts.ExecCLI,
		Toolbox:        agentops.DispatchToolbox,
		In:             in,
		Out:            out,
		ErrOut:         errOut,
		MaxSteps:       opts.MaxSteps,
		CommandTimeout: opts.CommandTimeout,
		Domain:         apitypes.DomainOps,
		Emitter:        renderer.NewEmitter(),
		Nested:         opts.Nested,
		PlainUI:        opts.Plain,
		EnvExtra:       opts.EnvExtra,
	}
}

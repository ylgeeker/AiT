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
	"github.com/os-artificer/AiT/internal/pkg/agentcore/runner"
	"github.com/os-artificer/AiT/internal/pkg/agentops"
	"github.com/os-artificer/AiT/internal/pkg/agentui/plain"
	pkgai "github.com/os-artificer/AiT/internal/pkg/ai"
	"github.com/os-artificer/AiT/internal/pkg/ai/apitypes"
)

// RemoteOpts configures a remote Agent session over SSH.
type RemoteOpts struct {
	AI             *pkgai.Client
	HostLabel      string
	Scrollback     func() string
	Remote         agentcore.RemoteRunner
	Capture        runner.CapturingRunner
	ExecCLI        CLIRunner
	MaxSteps       int
	CommandTimeout time.Duration
	Out            io.Writer
	ErrOut         io.Writer
}

// RemoteSessionConfig builds agentcore.Config for ai connect remote sessions.
func RemoteSessionConfig(opts RemoteOpts) agentcore.Config {
	out, errOut := opts.Out, opts.ErrOut
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stderr
	}
	renderer := plain.Renderer{Out: out}
	return agentcore.Config{
		Scope:          pkgai.ScopeRemote,
		AI:             opts.AI,
		HostLabel:      opts.HostLabel,
		Scrollback:     opts.Scrollback,
		Remote:         opts.Remote,
		Capture:        opts.Capture,
		CLI:            opts.ExecCLI,
		Toolbox:        agentops.DispatchToolbox,
		In:             os.Stdin,
		Out:            out,
		ErrOut:         errOut,
		MaxSteps:       opts.MaxSteps,
		CommandTimeout: opts.CommandTimeout,
		Domain:         apitypes.DomainOps,
		Emitter:        renderer.NewEmitter(),
	}
}

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

package agent

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/ait/agent/runner"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

func TestWriteAIHelpSections(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteAIHelp(&buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{
		"Usage:", "Examples:", "help for " + toolname.AgentCLI(""), toolname.AgentCLI("model list"), toolname.AgentCLI("ollama status"), toolname.AgentCLI("ollama setup"),
		toolname.AgentCLI("--model minimax"), "Subcommands:", "provider enable", "ollama status", "start", "stop", "disable",
		"--model", toolname.AgentCLI("connect"), toolname.AITCLI("agent"),
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	for _, gone := range []string{
		toolname.AgentCLI(`ask "summarize this"`),
		toolname.AgentCLI(`agent "find large files in /var/log"`),
		"One-shot agent:",
	} {
		if strings.Contains(out, gone) {
			t.Errorf("help must not mention one-shot %q", gone)
		}
	}
	// Subcommands list must not advertise ask/agent entries.
	if strings.Contains(out, "  ask ") || strings.Contains(out, "  agent [message") {
		t.Errorf("help still lists ask/agent subcommands:\n%s", out)
	}
}

func TestAIHelpSections(t *testing.T) {
	TestWriteAIHelpSections(t)
}

func TestLocalAgentConfig_HasCaptureAndCLI(t *testing.T) {
	exec := func(args []string, out, errOut io.Writer) error { return nil }
	cfg := withOpsDomain(runner.LocalAgentConfig(nil, 0, exec))
	if cfg.Capture == nil {
		t.Fatal("local agent missing Capture runner")
	}
	if cfg.CLI == nil {
		t.Fatal("local agent missing CLI executor")
	}
	if cfg.Domain != promptDomainOps {
		t.Fatalf("Domain=%q want %q", cfg.Domain, promptDomainOps)
	}
	if cfg.Toolbox == nil {
		t.Fatal("ops domain missing Toolbox")
	}
}

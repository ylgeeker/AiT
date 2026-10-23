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
	"fmt"
	"io"
	"strings"

	ollamasub "github.com/os-artificer/AiT/internal/ait/agent/ollama/subregistry"
	providersub "github.com/os-artificer/AiT/internal/ait/agent/provider/subregistry"
	"github.com/os-artificer/AiT/internal/ait/agent/subregistry"

	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
	pkgsubregistry "github.com/os-artificer/AiT/internal/pkg/subregistry"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

var aiNote = "Experimental (on by default): disable with ai.experimental.enabled=false " +
	"or AIT_AI_EXPERIMENTAL=0. " +
	"Requires ai.<provider>.api_key for cloud providers. " +
	"With ai.tiered.enabled, simple prompts use local Ollama. " +
	"Interactive agent loop: " + toolname.AITCLI("agent") + "."

var aiUsageLine = toolname.AgentCLI("[subcommand]")

// WriteAIHelp prints ait agent --help to w.
func WriteAIHelp(w io.Writer) error {
	b := pkghelpdoc.NewBody(w)
	if err := b.WriteIntro(pkghelpdoc.AIIntroLocal, aiNote); err != nil {
		return err
	}
	if err := b.WriteUsage(aiUsageLine); err != nil {
		return err
	}
	if err := b.WriteSection(pkghelpdoc.SectionFlags, pkghelpdoc.AIHelpFlagsBody); err != nil {
		return err
	}
	if err := b.WriteSection(pkghelpdoc.SectionSubcommands, formatAISubcommands()); err != nil {
		return err
	}
	if err := b.WriteSection(pkghelpdoc.SectionTieredRouting, pkghelpdoc.AIHelpTieredRoutingBody); err != nil {
		return err
	}
	if err := b.WriteSection("Agent loop", "Multi-turn agent tasks run in the ait toolbox: "+toolname.AITCLI("agent")+".\n"+
		"There is no one-shot ask/agent subcommand on "+toolname.Agent+"; use the loop.\n"); err != nil {
		return err
	}
	if err := b.WriteSection("Remote", "Use "+toolname.AgentCLI("connect <destination>")+
		" (or "+toolname.AgentCLI("<destination>")+") for a remote Agent session over SSH.\n"); err != nil {
		return err
	}
	if err := b.WriteExamples(
		toolname.AITCLI("agent"),
		toolname.AgentCLI("model list"),
		toolname.AgentCLI("connect user@host"),
		pkghelpdoc.ExampleAIOllamaStatus,
		pkghelpdoc.ExampleAIProviderEnableDeepseek,
		pkghelpdoc.ExampleAIOllamaSetupYes,
		toolname.AgentCLI("--model minimax"),
	); err != nil {
		return err
	}
	return b.Finish(toolname.AgentCLI(""))
}

func formatAISubcommands() string {
	providerChildren := pkgsubregistry.JoinChildNames(providersub.Entries(), "|")
	ollamaChildren := pkgsubregistry.JoinChildNames(ollamasub.Entries(), "|")

	var b strings.Builder
	for _, e := range subregistry.Entries() {
		switch e.Name {
		case "provider":
			fmt.Fprintf(&b, "  %s %s <id>\n", e.Name, providerChildren)
		case "ollama":
			fmt.Fprintf(&b, "  %s %s\n", e.Name, ollamaChildren)
		case "model":
			fmt.Fprintf(&b, "  %s [list]           %s\n", e.Name, e.Short)
		default:
			fmt.Fprintf(&b, "  %s   %s\n", e.Name, e.Short)
		}
	}
	return b.String()
}

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

package agentcore

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/os-artificer/AiT/internal/pkg/ai"

	"golang.org/x/term"
)

const (
	metaHelp = "help"
	metaQuit = "quit"

	promptHighRiskMarker = "High-risk"
)

// normalizeLoopLine maps Cursor-style slash commands to loop meta verbs.
func normalizeLoopLine(line string) string {
	line = strings.TrimSpace(line)
	if line == "" || !strings.HasPrefix(line, "/") {
		return line
	}

	rest := strings.TrimSpace(strings.TrimPrefix(line, "/"))
	if rest == "" {
		return line
	}

	word, tail := firstWord(rest)
	switch word {
	case metaExit, metaStop, metaQuit:
		return metaExit
	case metaHelp:
		if tail == "" {
			return metaHelp
		}
		return metaHelp + " " + tail
	case metaClear, metaHistory, metaModel:
		if tail == "" {
			return word
		}
		return word + " " + tail
	default:
		return line
	}
}

// SlashCommand describes one interactive loop slash command for UI menus.
type SlashCommand struct {
	Name string
	Desc string
}

// SlashCommands returns slash commands for the agent loop UI menu.
func SlashCommands() []SlashCommand {
	return []SlashCommand{
		{Name: "/help", Desc: "Show slash commands and tips"},
		{Name: "/exit", Desc: "End the agent session"},
		{Name: "/quit", Desc: "End the agent session"},
		{Name: "/clear", Desc: "Clear conversation history"},
		{Name: "/history", Desc: "Show or edit input history"},
		{Name: "/model", Desc: "Show current model and provider"},
		{Name: "/model list", Desc: "List available models"},
	}
}

// writeAgentLoopHelp prints slash-command help for the interactive loop.
func writeAgentLoopHelp(w io.Writer) {
	body := strings.Join([]string{
		"Slash commands:",
		"  /help          show this help",
		"  /exit, /quit   end the session",
		"  /clear         clear conversation history",
		"  /history       input history",
		"  /model         show model status",
		"  /model list    list available models",
		"",
		"Type a goal in natural language. Shell commands require approval [y/N].",
	}, "\r\n") + "\r\n"
	agentWrite(w, body)
}

func printSessionWelcome(cfg Config) {
	if cfg.PlainUI || cfg.ReadLine != nil || !sessionIsInteractive(cfg) {
		return
	}
	scope, host, _, _, _ := resolveDynamic(cfg)
	if scope == ai.ScopeRemote && host != "" {
		agentWrite(cfg.Out, fmt.Sprintf(msgWelcomeRemoteFmt, host))
		return
	}
	agentWrite(cfg.Out, msgWelcomeBanner)
}

func sessionIsInteractive(cfg Config) bool {
	return cfg.In == os.Stdin && isTerminalReader(cfg.In)
}

func isTerminalReader(in io.Reader) bool {
	f, ok := in.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

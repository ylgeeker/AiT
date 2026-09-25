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

package ait

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/os-artificer/AiT/internal/pkg/aitplugin"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

// WriteRootHelp prints ait --help to w (Agent entry + plugin dispatch).
func WriteRootHelp(w io.Writer) error {
	plugins := aitplugin.Discover(context.Background())

	b := pkghelpdoc.NewBody(w)
	if err := b.WriteIntro(
		"AiT Agent interaction entry.",
		"Run without arguments to start the Agent Loop. Plugin tools run as ait <name> (e.g. ait aissh).",
	); err != nil {
		return err
	}
	if err := b.WriteUsage(toolname.AIT + " [goal]"); err != nil {
		return err
	}
	cmdSection := formatAgentRootCommands(plugins)
	if err := b.WriteSection(pkghelpdoc.SectionCommands, cmdSection); err != nil {
		return err
	}

	if len(plugins) == 0 {
		if _, err := fmt.Fprintf(w, "%s\n\n%s\n\n", noPluginsMsg, errhint.FormatTryBlock(
			"make build",
			"make install",
			toolname.AITCLI("shell")+" (then) install",
			toolname.AIT,
		)); err != nil {
			return err
		}
	}

	if err := b.WriteExamples(
		toolname.AIT,
		toolname.AITCLI("agent connect user@host.example.com"),
		toolname.AITCLI("agent model list"),
		toolname.AIT+" ssh user@host.example.com",
		toolname.AITCLI("shell"),
	); err != nil {
		return err
	}
	return b.Finish(toolname.AIT)
}

func formatAgentRootCommands(plugins []aitplugin.Plugin) string {
	var b strings.Builder
	fmt.Fprintf(&b, "  %-12s %s\n", aitplugin.CmdAgent, "Agent Loop (default) and admin subcommands")
	for _, p := range plugins {
		fmt.Fprintf(&b, "  %-12s %s\n", p.Manifest.Name, p.Manifest.Short)
	}
	fmt.Fprintf(&b, "  %-12s %s\n", aitplugin.CmdShell, "toolbox REPL (install/update/upgrade)")
	fmt.Fprintf(&b, "  %-12s %s\n", aitplugin.CmdVersion, "print ait version")
	return b.String()
}

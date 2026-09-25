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

package shell

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

// WriteShellHelp prints ait shell --help to w (toolbox REPL for install/update).
func WriteShellHelp(w io.Writer) error {
	plugins := aitplugin.Discover(context.Background())

	b := pkghelpdoc.NewBody(w)
	if err := b.WriteIntro(
		"AiT toolbox shell (install/update/upgrade).",
		"Opens the interactive REPL for plugin management. Agent Loop: run ait with no arguments.",
	); err != nil {
		return err
	}
	if err := b.WriteUsage(toolname.AITCLI("shell")); err != nil {
		return err
	}
	if err := b.WriteSection(pkghelpdoc.SectionCommands, formatShellCommands(plugins)); err != nil {
		return err
	}
	if err := b.WriteExamples(
		toolname.AITCLI("shell"),
		toolname.AITCLI("shell")+" (then) install",
		toolname.AITCLI("shell")+" (then) update",
		toolname.AITCLI("shell")+" (then) upgrade",
	); err != nil {
		return err
	}
	return b.Finish(toolname.AITCLI("shell"))
}

func formatShellCommands(plugins []aitplugin.Plugin) string {
	var b strings.Builder
	for _, p := range plugins {
		fmt.Fprintf(&b, "  %-12s %s\n", p.Manifest.Name, p.Manifest.Short)
	}
	for _, e := range reservedHelpEntries() {
		fmt.Fprintf(&b, "  %-12s %s\n", e.name, e.desc)
	}
	fmt.Fprint(&b, formatManageCommands())
	fmt.Fprintf(&b, "  %-12s %s\n", cmdExit+"/"+cmdQuit, "leave the REPL")
	return b.String()
}

func writePluginHelp(w io.Writer, name string) error {
	if _, ok := manageHelpTopics[name]; ok {
		return writeManageHelp(w, name)
	}

	p, ok := aitplugin.Lookup(name)
	if !ok {
		return errhint.WithTryMsg(
			helpPluginNotFound+": "+name,
			"Use help with no arguments to list discovered plugins.",
			toolname.AIT,
			cmdHelp,
		)
	}
	return aitplugin.WriteHelp(w, p.Manifest)
}

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

package keepalive

import (
	"fmt"
	"io"
	"strings"

	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

// WriteRootHelp prints aitkeep --help to w.
// w receives the help body. Returns a write error from the underlying writer.
func WriteRootHelp(w io.Writer) error {
	b := pkghelpdoc.NewBody(w)
	if err := b.WriteIntro("Supervise processes from keepalive.yaml and answer GET /ping with pong."); err != nil {
		return err
	}
	if err := b.WriteUsage(toolname.KeepCLI("[command]")); err != nil {
		return err
	}
	if err := b.WriteSection(pkghelpdoc.SectionCommands, formatCommands()); err != nil {
		return err
	}
	if err := b.WriteExamples(
		toolname.KeepCLI("serve"),
		toolname.KeepCLI("status"),
		toolname.KeepCLI("add --name demo -- sleep 30"),
	); err != nil {
		return err
	}
	return b.Finish(toolname.KeepCLI(""))
}

func writeCommandHelp(w io.Writer, name, usage, flags string, examples ...string) error {
	b := pkghelpdoc.NewBody(w)
	if err := b.WriteIntro(commandShort[name]); err != nil {
		return err
	}
	if err := b.WriteUsage(toolname.KeepCLI(usage)); err != nil {
		return err
	}
	if err := b.WriteSection(pkghelpdoc.SectionFlags, flags); err != nil {
		return err
	}
	if err := b.WriteExamples(examples...); err != nil {
		return err
	}
	return b.Finish(name)
}

func formatCommands() string {
	var b strings.Builder
	for _, name := range commandOrder {
		fmt.Fprintf(&b, "  %-8s %s\n", name, commandShort[name])
	}
	return b.String()
}

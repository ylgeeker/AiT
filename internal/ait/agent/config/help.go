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

package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/os-artificer/AiT/internal/ait/agent/config/shared"
	"github.com/os-artificer/AiT/internal/ait/agent/config/subregistry"

	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

const configIntro = "Read and write AI and Agent YAML configuration (providers, API keys, tiered routing)."

// WriteConfigHelp prints ait agent config --help to w.
func WriteConfigHelp(w io.Writer) error {
	b := pkghelpdoc.NewBody(w)
	if err := b.WriteIntro(configIntro); err != nil {
		return err
	}
	if err := b.WriteUsage(toolname.AITCLI("agent config")); err != nil {
		return err
	}
	if err := b.WriteSection(pkghelpdoc.SectionSubcommands, formatConfigSubcommands()); err != nil {
		return err
	}
	if err := b.WriteSection(pkghelpdoc.SectionKeys, shared.KeysHelp); err != nil {
		return err
	}
	if err := b.WriteExamples(
		toolname.AITCLI("agent config init"),
		toolname.AITCLI("agent config view"),
		toolname.AITCLI("agent config validate"),
		toolname.AITCLI("agent config get ai.default"),
		toolname.AITCLI("agent config set ai.providers.deepseek.model deepseek-chat"),
	); err != nil {
		return err
	}
	return b.Finish("config")
}

func formatConfigSubcommands() string {
	var b strings.Builder
	for _, e := range subregistry.Entries() {
		switch e.Name {
		case "set":
			fmt.Fprintf(&b, "  %s <key> <value>   %s\n", e.Name, e.Short)
		case "get":
			fmt.Fprintf(&b, "  %s <key>           %s\n", e.Name, e.Short)
		default:
			fmt.Fprintf(&b, "  %-18s %s\n", e.Name, e.Short)
		}
	}
	return b.String()
}

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

// Package tooldispatch runs discovered toolbox plugins and allowed PATH binaries.
package tooldispatch

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/os-artificer/AiT/internal/pkg/aitplugin"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
	"github.com/os-artificer/AiT/internal/pkg/pluginbin"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

const toolboxAllowedHint = toolname.Host + ", " + toolname.Log + ", " + toolname.SSH + ", " + toolname.Img

var toolboxBins = map[string]struct{}{
	toolname.Host:    {},
	toolname.Log:     {},
	toolname.SSH:     {},
	toolname.Img:     {},
	toolname.AITHost: {},
	toolname.AITLog:  {},
	toolname.AITSSH:  {},
	toolname.AITImg:  {},
}

// Exec runs argv via plugin discovery (install dir) using the resolved binary path.
func Exec(ctx context.Context, argv []string, in io.Reader, out, errOut io.Writer) error {
	if len(argv) == 0 {
		return errhint.WithTryMsg("plugin name required", "Use a discovered plugin name.", toolname.AITCLI("help"))
	}
	p, ok := aitplugin.Lookup(argv[0])
	if !ok {
		return errhint.WithTryMsg("unknown plugin: "+argv[0], "Use a discovered plugin name.", toolname.AITCLI("help"))
	}
	if in == nil {
		in = os.Stdin
	}
	execArgv := append([]string{p.Manifest.Binary}, argv[1:]...)
	return pluginbin.ExecPath(ctx, p.Path, execArgv, p.Manifest.Name, in, out, errOut)
}

// ExecToolbox runs argv[0] when it is an allowed toolbox binary (discover path, then PATH).
func ExecToolbox(ctx context.Context, argv []string, out, errOut io.Writer) error {
	if len(argv) == 0 {
		return errhint.WithTryMsg("toolbox requires a command",
			"Pass a full argv starting with "+toolboxAllowedHint+".",
			toolname.AITCLI("agent \"check disk\""),
		)
	}
	name := argv[0]
	if _, ok := toolboxBins[name]; !ok {
		return errhint.WithTryMsg("unknown toolbox binary: "+name,
			"Allowed: "+toolboxAllowedHint+".",
			toolname.HostCLI("snapshot --json"),
			toolname.LogCLI("search error"),
		)
	}
	lookup := name
	if canon, ok := aitplugin.LegacyPluginNames[name]; ok {
		lookup = canon
	}
	if p, ok := aitplugin.Lookup(lookup); ok {
		execArgv := append([]string{p.Manifest.Binary}, argv[1:]...)
		return pluginbin.ExecPath(ctx, p.Path, execArgv, p.Manifest.Name, os.Stdin, out, errOut)
	}
	bin, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("%s not found on PATH: %w", name, err)
	}
	cmd := exec.CommandContext(ctx, bin, argv[1:]...)
	cmd.Stdout = out
	cmd.Stderr = errOut
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

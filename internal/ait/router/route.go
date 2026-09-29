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

// Package router implements top-level ait argv routing.
package router

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/os-artificer/AiT/internal/ait/agent"
	"github.com/os-artificer/AiT/internal/ait/dispatch"
	"github.com/os-artificer/AiT/internal/ait/loop"
	"github.com/os-artificer/AiT/internal/pkg/aitplugin"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

// RouteDeps holds callbacks from package ait for help, shell REPL, and version.
type RouteDeps struct {
	WriteRootHelp  func(io.Writer) error
	WriteShellHelp func(io.Writer) error
	RunShellREPL   func(context.Context, io.Reader, io.Writer, io.Writer) error
	VersionLine    func() string
	RunPlugin      func(context.Context, []string, io.Writer, io.Writer) error
}

// Route runs the top-level ait CLI for args after experimental-flag stripping.
func Route(ctx context.Context, args []string, stdout, stderr io.Writer, d RouteDeps) error {
	if len(args) == 0 {
		return runLoop(ctx, "", false, stdout, stderr)
	}

	if pkghelpdoc.HelpRequested(args) {
		return d.WriteRootHelp(stdout)
	}

	if pkghelpdoc.VersionRequested(args) {
		_, err := fmt.Fprintln(stdout, d.VersionLine())
		return err
	}

	if args[0] == aitplugin.CmdShell {
		return routeShell(ctx, args[1:], stdout, stderr, d)
	}

	if args[0] == aitplugin.CmdAgent {
		return routeAgent(ctx, args[1:], stdout, stderr)
	}

	if name, ok := ResolvePluginName(args[0]); ok {
		argv := append([]string{name}, args[1:]...)
		if d.RunPlugin != nil {
			return d.RunPlugin(ctx, argv, stdout, stderr)
		}
		return execPlugin(ctx, argv, stdout, stderr)
	}

	return runLoop(ctx, strings.Join(args, " "), false, stdout, stderr)
}

func routeShell(ctx context.Context, args []string, stdout, stderr io.Writer, d RouteDeps) error {
	if len(args) == 0 {
		return d.RunShellREPL(ctx, os.Stdin, stdout, stderr)
	}
	if pkghelpdoc.HelpRequested(args) {
		return d.WriteShellHelp(stdout)
	}
	return errhint.WithTryMsg(
		"unknown ait shell command",
		"ait shell runs the toolbox REPL for install/update/upgrade.",
		toolname.AITCLI("shell"),
	)
}

func routeAgent(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	plain, rest := pkghelpdoc.TakePlainFlag(args)
	if len(rest) == 0 {
		return runLoop(ctx, "", plain, stdout, stderr)
	}

	if agent.IsAdminVerb(rest[0]) {
		return agent.RunAdmin(rest, stdout, stderr)
	}
	if agent.LooksLikeConnectDestination(rest[0]) {
		return agent.RunAdmin(append([]string{"connect"}, rest...), stdout, stderr)
	}

	return runLoop(ctx, strings.Join(rest, " "), plain, stdout, stderr)
}

func runLoop(ctx context.Context, goal string, plain bool, stdout, stderr io.Writer) error {
	return loop.Run(ctx, goal, plain, stdout, stderr)
}

func execPlugin(ctx context.Context, argv []string, stdout, stderr io.Writer) error {
	return dispatch.RunPlugin(ctx, argv, stdout, stderr)
}

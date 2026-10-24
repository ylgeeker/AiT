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

package ait_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/os-artificer/AiT/internal/ait"
	"github.com/os-artificer/AiT/internal/pkg/aitplugin"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
	"github.com/os-artificer/AiT/internal/pkg/pluginbin"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

func TestExecuteWriters_Help(t *testing.T) {
	resetPlugins(t)
	var out bytes.Buffer
	if err := ait.ExecuteWriters([]string{"-h"}, &out, io.Discard); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	if !strings.Contains(s, "Usage:") || !strings.Contains(s, toolname.AIT) {
		t.Fatalf("help = %q", s)
	}
	if strings.Contains(s, "OpenSSH") {
		t.Fatal("root help must not include SSH connect section")
	}
}

func TestExecuteWriters_Version(t *testing.T) {
	var out bytes.Buffer
	if err := ait.ExecuteWriters([]string{"-V"}, &out, io.Discard); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out.String()) == "" {
		t.Fatal("empty version")
	}
}

func TestExecuteWriters_PluginDispatch(t *testing.T) {
	resetPlugins(t)
	installBuiltinPlugin(t, "aitssh")

	t.Cleanup(func() { pluginbin.SetRunner(nil) })
	var gotArgv []string
	pluginbin.SetRunner(func(_ context.Context, _ string, argv []string, _ string) error {
		gotArgv = append([]string(nil), argv...)
		return nil
	})

	err := ait.ExecuteWriters([]string{"ssh", "--help"}, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(gotArgv) < 1 || gotArgv[0] != "aitssh" {
		t.Fatalf("argv = %#v", gotArgv)
	}
}

func TestRunREPL_Exit(t *testing.T) {
	resetPlugins(t)
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		_, _ = io.WriteString(w, "exit\n")
		_ = w.Close()
	}()
	var out bytes.Buffer
	if err := ait.RunREPL(context.Background(), r, &out, io.Discard); err != nil {
		t.Fatal(err)
	}
}

func TestExecuteLine_VersionBuiltin(t *testing.T) {
	var out bytes.Buffer
	exit, err := ait.ExecuteLine(context.Background(), "version", nil, &out, io.Discard)
	if err != nil || exit {
		t.Fatalf("exit=%v err=%v", exit, err)
	}
	if strings.TrimSpace(out.String()) == "" {
		t.Fatal("empty version")
	}
}

func TestExecuteLine_HelpPlugin(t *testing.T) {
	resetPlugins(t)
	installBuiltinPlugin(t, "aitssh")

	var out bytes.Buffer
	exit, err := ait.ExecuteLine(context.Background(), "help ssh", nil, &out, io.Discard)
	if err != nil || exit {
		t.Fatalf("exit=%v err=%v", exit, err)
	}
	if !strings.Contains(out.String(), "--save") {
		t.Fatalf("help = %q", out.String())
	}
}

func TestExecuteLine_HelpUnknown(t *testing.T) {
	resetPlugins(t)
	_, err := ait.ExecuteLine(context.Background(), "help nosuch", nil, io.Discard, io.Discard)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown help topic") {
		t.Fatalf("err = %v", err)
	}
}

func TestExecuteLine_UnknownCommand(t *testing.T) {
	resetPlugins(t)
	_, err := ait.ExecuteLine(context.Background(), "nosuch", nil, io.Discard, io.Discard)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("err = %v", err)
	}
}

func TestExecuteLine_RunPluginSwallowsExitError(t *testing.T) {
	resetPlugins(t)
	installBuiltinPlugin(t, "aitssh")

	t.Cleanup(func() { pluginbin.SetRunner(nil) })
	pluginbin.SetRunner(func(context.Context, string, []string, string) error {
		return &exec.ExitError{}
	})

	exit, err := ait.ExecuteLine(context.Background(), "ssh --help", nil, io.Discard, io.Discard)
	if exit || err != nil {
		t.Fatalf("exit=%v err=%v", exit, err)
	}
}

func TestExecuteLine_JoinedExitErrorNotSwallowed(t *testing.T) {
	resetPlugins(t)
	installBuiltinPlugin(t, "aitssh")

	t.Cleanup(func() { pluginbin.SetRunner(nil) })
	pluginbin.SetRunner(func(context.Context, string, []string, string) error {
		return errors.Join(fmt.Errorf("terminal: make raw"), &exec.ExitError{})
	})

	_, err := ait.ExecuteLine(context.Background(), "ssh --help", nil, io.Discard, io.Discard)
	if err == nil {
		t.Fatal("expected MakeRaw+ExitError to surface")
	}
	if !strings.Contains(err.Error(), "terminal") {
		t.Fatalf("err = %v", err)
	}
}

func TestExecuteLine_UnclosedQuoteHint(t *testing.T) {
	_, err := ait.ExecuteLine(context.Background(), "ssh 'open", nil, io.Discard, io.Discard)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unclosed quote") || !strings.Contains(err.Error(), "Try:") {
		t.Fatalf("err = %v", err)
	}
}

func TestExecuteLine_DoubleQuotedArg(t *testing.T) {
	resetPlugins(t)
	installBuiltinPlugin(t, "aitssh")

	t.Cleanup(func() { pluginbin.SetRunner(nil) })
	var gotArgv []string
	pluginbin.SetRunner(func(_ context.Context, _ string, argv []string, _ string) error {
		gotArgv = append([]string(nil), argv...)
		return nil
	})

	_, err := ait.ExecuteLine(context.Background(), `ssh "user@host.example.com"`, nil, io.Discard, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if len(gotArgv) < 2 || gotArgv[1] != "user@host.example.com" {
		t.Fatalf("argv = %#v", gotArgv)
	}
}

func TestExecuteLine_WrappedExitErrorNotSwallowed(t *testing.T) {
	resetPlugins(t)
	installBuiltinPlugin(t, "aitssh")

	t.Cleanup(func() { pluginbin.SetRunner(nil) })
	pluginbin.SetRunner(func(context.Context, string, []string, string) error {
		return fmt.Errorf("wrap: %w", &exec.ExitError{})
	})

	_, err := ait.ExecuteLine(context.Background(), "ssh --help", nil, io.Discard, io.Discard)
	if err == nil {
		t.Fatal("expected wrapped ExitError to surface")
	}
}

func TestWriteRootHelp_NoPlugins(t *testing.T) {
	resetPlugins(t)
	var out bytes.Buffer
	if err := ait.WriteRootHelp(&out); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	if !strings.Contains(s, "no plugins found") || !strings.Contains(s, "Try:") {
		t.Fatalf("help = %q", s)
	}
	if !strings.Contains(s, "make install") {
		t.Fatalf("missing make install Try: %q", s)
	}
	if !strings.Contains(s, "agent") || !strings.Contains(s, "shell") {
		t.Fatalf("missing agent/shell commands: %q", s)
	}
}

func TestWriteRootHelp_PluginsBeforeBuiltins(t *testing.T) {
	resetPlugins(t)
	installBuiltinPlugin(t, "aitssh")

	var out bytes.Buffer
	if err := ait.WriteRootHelp(&out); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	cmdSSH := strings.Index(s, "  ssh ")
	cmdShell := strings.Index(s, "  shell ")
	if cmdSSH < 0 || cmdShell < 0 {
		t.Fatalf("missing command rows in: %q", s)
	}
	if cmdSSH > cmdShell {
		t.Fatalf("plugins must list before shell: ssh@%d shell@%d\n%s", cmdSSH, cmdShell, s)
	}
}

func TestExecuteLine_RunPluginUsesPath(t *testing.T) {
	resetPlugins(t)
	path := installBuiltinPlugin(t, "aitlog")

	t.Cleanup(func() { pluginbin.SetRunner(nil) })
	var gotPath string
	var gotPluginCMD string
	pluginbin.SetRunner(func(_ context.Context, p string, argv []string, pluginCMD string) error {
		gotPath = p
		gotPluginCMD = pluginCMD
		if argv[0] != "aitlog" {
			t.Fatalf("argv[0]=%q", argv[0])
		}
		return nil
	})

	_, err := ait.ExecuteLine(context.Background(), "log sources", nil, io.Discard, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != path {
		t.Fatalf("path = %q want %q", gotPath, path)
	}
	if gotPluginCMD != "log" {
		t.Fatalf("pluginCMD=%q want log", gotPluginCMD)
	}
}

func TestHistoryPath_NotReplHistory(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("AIT_HOME", dir)
	var out bytes.Buffer
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		_, _ = io.WriteString(w, "exit\n")
		_ = w.Close()
	}()
	if err := ait.RunREPL(context.Background(), r, &out, io.Discard); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "repl_history")); err == nil {
		t.Fatal("must not use repl_history")
	}
}

func TestToolnameAIT(t *testing.T) {
	if toolname.AIT != "ait" {
		t.Fatalf("AIT = %q", toolname.AIT)
	}
}

func TestErrhintPreservesIs(t *testing.T) {
	cause := errors.New("interactive only")
	err := errhint.WithTry(cause, "x", "ait")
	if !errors.Is(err, cause) {
		t.Fatal("errors.Is broken")
	}
}

func resetPlugins(t *testing.T) {
	t.Helper()
	aitplugin.ResetDiscoverForTest()
	aitplugin.RestoreDiscoverHooksForTest()
	t.Cleanup(func() {
		aitplugin.ResetDiscoverForTest()
		aitplugin.RestoreDiscoverHooksForTest()
	})
	aitplugin.SetLookPathForTest(func(string) (string, error) {
		return "", os.ErrNotExist
	})
	aitplugin.SetWarnWriterForTest(io.Discard)
	aitplugin.SetInstallDirsForTest([]string{t.TempDir()})
}

func installBuiltinPlugin(t *testing.T, binary string) string {
	t.Helper()
	root := t.TempDir()
	cmd := aitplugin.PluginCommandName(binary)
	exe := aitplugin.PluginBinary(cmd)
	plugDir := filepath.Join(root, "plugins", cmd)
	if err := os.MkdirAll(plugDir, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := aitplugin.BuiltinJSON(binary)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(plugDir, "plugin.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	name := exe
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(plugDir, name)
	mode := os.FileMode(0o755)
	if runtime.GOOS == "windows" {
		mode = 0o666
	}
	if err := os.WriteFile(path, []byte("x"), mode); err != nil {
		t.Fatal(err)
	}
	aitplugin.SetInstallDirsForTest([]string{root})
	aitplugin.ResetDiscoverForTest()
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

// TestExecuteWriters_BareAitRoutesLoop (AT1/AT49) empty argv starts Agent Loop route (not shell REPL).
func TestExecuteWriters_BareAitRoutesLoop(t *testing.T) {
	err := ait.ExecuteWriters(nil, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "experimental") {
		t.Fatalf("bare ait must not be experimental-gated: %v", err)
	}
}

// TestExecuteWriters_AgentAliasRoutesLoop (AT49) ait agent is a silent alias for the Loop.
func TestExecuteWriters_AgentAliasRoutesLoop(t *testing.T) {
	err := ait.ExecuteWriters([]string{"agent"}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "experimental") {
		t.Fatalf("ait agent must not be experimental-gated: %v", err)
	}
}

// TestExecuteWriters_AgentModelList (AT50/OA14) admin subcommand via ait agent.
func TestExecuteWriters_AgentModelList(t *testing.T) {
	done := make(chan error, 1)
	go func() {
		done <- ait.ExecuteWriters([]string{"agent", "model", "list"}, io.Discard, io.Discard)
	}()
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("model list hung >5s")
	case err := <-done:
		if err != nil && strings.Contains(err.Error(), `unknown command "model"`) {
			t.Fatalf("model list must reach model subcommand, got %v", err)
		}
	}
}

// TestExecuteWriters_AgentHelpWithoutConfig prints ait agent --help when AI is not configured.
func TestExecuteWriters_AgentHelpWithoutConfig(t *testing.T) {
	t.Setenv("AIT_HOME", t.TempDir())
	for _, arg := range []string{"--help", "-h", "help"} {
		var out, errOut bytes.Buffer
		if err := ait.ExecuteWriters([]string{"agent", arg}, &out, &errOut); err != nil {
			t.Fatalf("%s: %v\nstderr: %s", arg, err, errOut.String())
		}
		s := out.String()
		for _, sub := range []string{"Usage:", "provider", "ollama", "model", "connect"} {
			if !strings.Contains(s, sub) {
				t.Fatalf("%s help missing %q: %q", arg, sub, s)
			}
		}
		if strings.Contains(errOut.String(), "AI not configured") {
			t.Fatalf("%s help required AI config: %s", arg, errOut.String())
		}
	}
}

// TestExecuteWriters_ShellHelp (AT51) ait shell --help is toolbox REPL help, not Agent root help.
func TestExecuteWriters_ShellHelp(t *testing.T) {
	resetPlugins(t)
	var out bytes.Buffer
	if err := ait.ExecuteWriters([]string{"shell", "--help"}, &out, io.Discard); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	if !strings.Contains(s, "install") || !strings.Contains(s, "toolbox") {
		t.Fatalf("shell help = %q", s)
	}
	if strings.Contains(s, "Agent interaction entry") {
		t.Fatal("shell help must not duplicate Agent root intro")
	}
}

// TestExecuteWriters_PluginAliasLog (AT56) ait log routes to ailog plugin.
func TestExecuteWriters_PluginAliasLog(t *testing.T) {
	resetPlugins(t)
	installBuiltinPlugin(t, "aitlog")

	t.Cleanup(func() { pluginbin.SetRunner(nil) })
	var gotArgv []string
	pluginbin.SetRunner(func(_ context.Context, _ string, argv []string, _ string) error {
		gotArgv = append([]string(nil), argv...)
		return nil
	})

	if err := ait.ExecuteWriters([]string{"log", "version"}, io.Discard, io.Discard); err != nil {
		t.Fatal(err)
	}
	if len(gotArgv) < 1 || gotArgv[0] != "aitlog" {
		t.Fatalf("argv = %#v", gotArgv)
	}
}

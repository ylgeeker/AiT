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

package tooldispatch_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/aitplugin"
	"github.com/os-artificer/AiT/internal/pkg/pluginbin"
	"github.com/os-artificer/AiT/internal/pkg/tooldispatch"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

type discard struct{}

func (discard) Write(p []byte) (int, error) { return len(p), nil }

// TestExecToolbox_DiscoverPath (PKG-5) toolbox dispatch prefers install-dir plugin path.
func TestExecToolbox_DiscoverPath(t *testing.T) {
	root := t.TempDir()
	plugDir := filepath.Join(root, "plugins", "host")
	if err := os.MkdirAll(plugDir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"apiVersion":1,"name":"host","binary":"aithost","short":"host"}`
	if err := os.WriteFile(filepath.Join(plugDir, "plugin.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(plugDir, "aithost"), []byte("#!/bin/sh\necho ok\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	aitplugin.SetInstallDirsForTest([]string{root})
	aitplugin.ResetDiscoverForTest()
	t.Cleanup(func() {
		aitplugin.SetInstallDirsForTest(nil)
		aitplugin.ResetDiscoverForTest()
	})

	t.Cleanup(func() { pluginbin.SetRunner(nil) })
	var gotPath string
	var gotPluginCMD string
	pluginbin.SetRunner(func(_ context.Context, path string, _ []string, pluginCMD string) error {
		gotPath = path
		gotPluginCMD = pluginCMD
		return nil
	})

	if err := tooldispatch.ExecToolbox(context.Background(), []string{toolname.Host, "snapshot"}, discard{}, discard{}); err != nil {
		t.Fatal(err)
	}
	if gotPath == "" {
		t.Fatal("expected plugin path exec")
	}
	if gotPluginCMD != toolname.Host {
		t.Fatalf("pluginCMD=%q want %q", gotPluginCMD, toolname.Host)
	}
}

func TestExec_ForwardsPluginCMD(t *testing.T) {
	root := t.TempDir()
	plugDir := filepath.Join(root, "plugins", "log")
	if err := os.MkdirAll(plugDir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"apiVersion":1,"name":"log","binary":"aitlog","short":"log"}`
	if err := os.WriteFile(filepath.Join(plugDir, "plugin.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(plugDir, "aitlog"), []byte("#!/bin/sh\necho ok\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	aitplugin.SetInstallDirsForTest([]string{root})
	aitplugin.ResetDiscoverForTest()
	t.Cleanup(func() {
		aitplugin.SetInstallDirsForTest(nil)
		aitplugin.ResetDiscoverForTest()
	})

	t.Cleanup(func() { pluginbin.SetRunner(nil) })
	var gotPluginCMD string
	pluginbin.SetRunner(func(_ context.Context, _ string, _ []string, pluginCMD string) error {
		gotPluginCMD = pluginCMD
		return nil
	})

	if err := tooldispatch.Exec(context.Background(), []string{toolname.Log, "--help"}, nil, discard{}, discard{}); err != nil {
		t.Fatal(err)
	}
	if gotPluginCMD != toolname.Log {
		t.Fatalf("pluginCMD=%q want %q", gotPluginCMD, toolname.Log)
	}
}

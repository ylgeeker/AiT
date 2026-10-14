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
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/aitplugin"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

func TestCLI_NoArgsHelp(t *testing.T) {
	var out, errBuf bytes.Buffer
	if err := ExecuteWriters(nil, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Usage:") || !strings.Contains(out.String(), "serve") {
		t.Fatalf("help=%s", out.String())
	}
}

func TestCLI_ServeHelp(t *testing.T) {
	var out, errBuf bytes.Buffer
	if err := ExecuteWriters([]string{"serve", "--help"}, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	text := out.String()
	if !strings.Contains(text, "Usage:") || !strings.Contains(text, "Flags:") || !strings.Contains(text, "--listen") {
		t.Fatalf("help=%s", text)
	}
	if !strings.Contains(text, "Examples:") {
		t.Fatalf("help=%s", text)
	}
}

func TestCLI_AddHelpFlags(t *testing.T) {
	var out, errBuf bytes.Buffer
	if err := ExecuteWriters([]string{"add", "--help"}, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	for _, flag := range []string{"--name", "--cpu", "--mem", "--disk-bytes", "--disk-percent", "--read-bps", "--write-bps", "--io-percent", "--cwd"} {
		if !strings.Contains(out.String(), flag) {
			t.Fatalf("missing %s in %s", flag, out.String())
		}
	}
}

func TestCLI_Manifest(t *testing.T) {
	var out, errBuf bytes.Buffer
	if err := ExecuteWriters([]string{aitplugin.ManifestArg}, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"binary": "aitkeep"`) {
		t.Fatalf("manifest=%s", out.String())
	}
}

func TestCLI_StatusDaemonDown(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("AIT_HOME", dir)
	t.Setenv("AIT_DATA_DIR", "")
	var out, errBuf bytes.Buffer
	err := ExecuteWriters([]string{"status"}, &out, &errBuf)
	if err == nil || !strings.Contains(err.Error(), "not running") {
		t.Fatalf("err=%v", err)
	}
}

func TestCLI_AddRemove(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("AIT_HOME", dir)
	t.Setenv("AIT_DATA_DIR", "")
	var out, errBuf bytes.Buffer
	err := ExecuteWriters([]string{"add", "--name", "demo", "--cpu", "20", "--", "sleep", "1"}, &out, &errBuf)
	if err != nil {
		t.Fatal(err)
	}
	cfg, _, _ := Layout(dir)
	doc, err := Load(cfg)
	if err != nil || len(doc.Units) != 1 || doc.Units[0].Limits.CPUPercent != 20 {
		t.Fatalf("doc=%+v err=%v", doc, err)
	}
	if err := ExecuteWriters([]string{"remove", "demo"}, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	doc, err = Load(cfg)
	if err != nil || len(doc.Units) != 0 {
		t.Fatalf("after remove=%+v err=%v", doc, err)
	}
	if toolname.Keep != "keep" || toolname.AITKeep != "aitkeep" {
		t.Fatalf("names %s %s", toolname.Keep, toolname.AITKeep)
	}
}

func TestCLI_AddBadNameHint(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("AIT_HOME", dir)
	t.Setenv("AIT_DATA_DIR", "")
	var out, errBuf bytes.Buffer
	err := ExecuteWriters([]string{"add", "--name", "bad/name", "--", "sleep", "1"}, &out, &errBuf)
	if err == nil || !errors.Is(err, ErrBadName) {
		t.Fatalf("err=%v", err)
	}
	text := err.Error()
	if !strings.Contains(text, "Try:") || !strings.Contains(text, "add --name demo") {
		t.Fatalf("hint=%s", text)
	}
}

func TestCLI_RemoveMissingHint(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("AIT_HOME", dir)
	t.Setenv("AIT_DATA_DIR", "")
	var out, errBuf bytes.Buffer
	err := ExecuteWriters([]string{"remove", "nosuch"}, &out, &errBuf)
	if err == nil || !errors.Is(err, ErrUnitMissing) {
		t.Fatalf("err=%v", err)
	}
	if !strings.Contains(err.Error(), "Try:") {
		t.Fatalf("hint=%s", err.Error())
	}
}

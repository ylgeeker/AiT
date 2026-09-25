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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestAit_NoDirectAppImports (PKG-4) ait must not import product app packages directly.
func TestAit_NoDirectAppImports(t *testing.T) {
	root := moduleRoot(t)
	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}: {{join .Imports \"\\n\"}}", "./internal/ait/...")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	forbidden := []string{
		"github.com/os-artificer/AiT/internal/ssh",
		"github.com/os-artificer/AiT/internal/log",
		"github.com/os-artificer/AiT/internal/host",
		"github.com/os-artificer/AiT/internal/img",
	}
	for _, line := range strings.Split(string(out), "\n") {
		for _, imp := range strings.Fields(line) {
			for _, f := range forbidden {
				if strings.HasPrefix(imp, f) {
					t.Fatalf("forbidden import %q in %s", imp, line)
				}
			}
		}
	}
}

// TestBuild_NoAIBinary (AT57) standalone ai binary and cmd/ai are retired.
func TestBuild_NoAIBinary(t *testing.T) {
	root := moduleRoot(t)
	if _, err := os.Stat(filepath.Join(root, "cmd", "ai")); err == nil {
		t.Fatal("cmd/ai must not exist after ai binary retirement")
	}
	makefile, err := os.ReadFile(filepath.Join(root, "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(makefile), "\n") {
		if strings.HasPrefix(line, "TOOLS") && strings.Contains(line, " ai ") {
			t.Fatalf("Makefile TOOLS must not include standalone ai: %q", strings.TrimSpace(line))
		}
	}
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

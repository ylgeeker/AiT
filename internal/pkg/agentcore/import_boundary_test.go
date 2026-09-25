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

package agentcore_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestAgentcore_NoProductImports (AC-V4-1) agentcore must not import product CLI packages.
func TestAgentcore_NoProductImports(t *testing.T) {
	root := moduleRoot(t)
	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}: {{join .Imports \"\\n\"}}", "./internal/pkg/agentcore/...")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		for _, imp := range strings.Fields(line) {
			if strings.HasPrefix(imp, "github.com/os-artificer/AiT/internal/ssh") ||
				strings.HasPrefix(imp, "github.com/os-artificer/AiT/internal/ait") ||
				strings.HasPrefix(imp, "github.com/os-artificer/AiT/internal/ait/agent") {
				t.Fatalf("forbidden import %q in %s", imp, line)
			}
		}
	}
}

// TestRepo_NoAimodeImport (AC-V4-6) aimode package was removed in v4.
func TestRepo_NoAimodeImport(t *testing.T) {
	root := moduleRoot(t)
	var offenders []string
	_ = filepath.Walk(filepath.Join(root, "internal"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if strings.Contains(path, filepath.Join("internal", "pkg", "aimode")) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		if strings.Contains(string(data), "internal/pkg/aimode") {
			offenders = append(offenders, path)
		}
		return nil
	})
	if len(offenders) > 0 {
		t.Fatalf("aimode imports remain: %v", offenders)
	}
	if _, err := os.Stat(filepath.Join(root, "internal/pkg/aimode")); err == nil {
		t.Fatal("internal/pkg/aimode directory still exists")
	}
}

// TestRepo_NoAichatOrRunOnce (RO-2) one-shot aichat and agentcore.RunOnce were removed.
func TestRepo_NoAichatOrRunOnce(t *testing.T) {
	root := moduleRoot(t)
	if _, err := os.Stat(filepath.Join(root, "internal/pkg/aichat")); err == nil {
		t.Fatal("internal/pkg/aichat directory still exists")
	}
	loopPath := filepath.Join(root, "internal/pkg/agentcore/loop.go")
	data, err := os.ReadFile(loopPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "func RunOnce(") {
		t.Fatal("agentcore.RunOnce must be removed")
	}
}

// TestRepo_NoOpsaiPackage (PKG-1) internal/opsai was merged into internal/ait/agent.
func TestRepo_NoOpsaiPackage(t *testing.T) {
	root := moduleRoot(t)
	if _, err := os.Stat(filepath.Join(root, "internal", "opsai")); err == nil {
		t.Fatal("internal/opsai directory still exists")
	}
	cmd := exec.Command("go", "list", "./...")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "internal/opsai") {
		t.Fatalf("go list still includes internal/opsai:\n%s", out)
	}
}

// TestRepo_NoCmdAI (PKG-1) standalone ai binary was removed; use ait agent.
func TestRepo_NoCmdAI(t *testing.T) {
	root := moduleRoot(t)
	if _, err := os.Stat(filepath.Join(root, "cmd", "ai")); err == nil {
		t.Fatal("cmd/ai directory still exists")
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

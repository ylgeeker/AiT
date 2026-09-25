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

package shellext_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/shellext"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

func TestRunWithOpts_SetsPluginCMDEnv(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("printenv helper script is Unix-only")
	}

	dir := t.TempDir()
	script := filepath.Join(dir, "printenv.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nprintenv AIT_PLUGIN_CMD\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	var out strings.Builder
	err := shellext.RunWithOpts(context.Background(), script, []string{"printenv", "-"}, toolname.Log, nil, &out, nil)
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			t.Fatalf("run: %v", err)
		}
	}
	if got := strings.TrimSpace(out.String()); got != toolname.Log {
		t.Fatalf("AIT_PLUGIN_CMD=%q want %q", got, toolname.Log)
	}
}

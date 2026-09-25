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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEnsureSSHDefaults_Creates locks create semantics: session keys are merged, and
// Save still materializes ai: via Load defaults()/normalizeAI (pre-existing behavior).
func TestEnsureSSHDefaults_Creates(t *testing.T) {
	root := isolateConfigHome(t)

	res, err := EnsureSSHDefaults(false)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Created {
		t.Fatal("want created")
	}

	b, err := os.ReadFile(filepath.Join(root, "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(b)
	if !strings.Contains(text, "session:") {
		t.Fatalf("missing session:\n%s", text)
	}
	if !strings.Contains(text, "auth_ttl:") {
		t.Fatalf("missing auth_ttl:\n%s", text)
	}
	if !strings.Contains(text, "ai:") {
		t.Fatalf("current create path still writes ai: via defaults(); missing:\n%s", text)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Session.AuthTTL == "" {
		t.Fatal("empty session.auth_ttl")
	}
	if cfg.Session.RemotePrompt.Path == "" {
		t.Fatal("empty session.remote_prompt.path")
	}
}

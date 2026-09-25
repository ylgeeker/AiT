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

package configcli

import (
	"errors"
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

func TestRejectAIConfigKey(t *testing.T) {
	if err := RejectAIConfigKey("session.auth_ttl"); err != nil {
		t.Fatalf("session key: %v", err)
	}

	err := RejectAIConfigKey("ai.default")
	if err == nil {
		t.Fatal("want reject")
	}
	if !errors.Is(err, ErrAIKeyOnSSHConfig) {
		t.Fatalf("errors.Is: %v", err)
	}
	s := err.Error()
	if !strings.Contains(s, "Try:") {
		t.Fatalf("missing Try: %q", s)
	}
	if !strings.Contains(s, toolname.AgentCLI("config init")) {
		t.Fatalf("missing agent try: %q", s)
	}
}

func TestRejectSSHConfigKey(t *testing.T) {
	if err := RejectSSHConfigKey("ai.default"); err != nil {
		t.Fatalf("ai key: %v", err)
	}

	err := RejectSSHConfigKey("session.auth_ttl")
	if err == nil {
		t.Fatal("want reject")
	}
	s := err.Error()
	if !strings.Contains(s, "key belongs to aissh config") {
		t.Fatalf("summary: %q", s)
	}
	if !strings.Contains(s, "Try:") {
		t.Fatalf("missing Try: %q", s)
	}
	if !strings.Contains(s, toolname.CLI("config get session.auth_ttl")) {
		t.Fatalf("missing aissh try: %q", s)
	}
}

func TestSSH_AI_Options(t *testing.T) {
	ssh := SSH()
	if ssh.HandleAPIKeys || ssh.Scope != ScopeSSH {
		t.Fatalf("SSH opts: %+v", ssh)
	}
	ai := AI()
	if !ai.HandleAPIKeys || ai.Scope != ScopeAI {
		t.Fatalf("AI opts: %+v", ai)
	}
}

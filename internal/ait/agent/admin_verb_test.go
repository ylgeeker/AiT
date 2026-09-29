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

package agent

import (
	"testing"

	"github.com/os-artificer/AiT/internal/ait/agent/subregistry"
)

func TestIsAdminVerb_MatchesRegistry(t *testing.T) {
	entries := subregistry.Entries()
	want := map[string]struct{}{
		"connect":  {},
		"model":    {},
		"provider": {},
		"ollama":   {},
		"config":   {},
	}
	if len(entries) != len(want) {
		t.Fatalf("registry len=%d want %d", len(entries), len(want))
	}
	for _, entry := range entries {
		if _, ok := want[entry.Name]; !ok {
			t.Fatalf("unexpected admin verb %q", entry.Name)
		}
		if !IsAdminVerb(entry.Name) {
			t.Fatalf("IsAdminVerb(%q) = false", entry.Name)
		}
	}
	if IsAdminVerb("") || IsAdminVerb("not-a-verb") {
		t.Fatal("unknown name reported as admin verb")
	}
}

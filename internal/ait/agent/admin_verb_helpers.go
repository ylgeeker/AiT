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

import "github.com/os-artificer/AiT/internal/ait/agent/subregistry"

// IsAdminVerb reports whether name is a registered agent management subcommand.
// name is the first token after "ait agent" (for example "model" or "connect").
// Returns false when name is empty or not registered. The set comes from the
// Cobra subregistry, which is filled by the blank import in cli_root.go.
func IsAdminVerb(name string) bool {
	if name == "" {
		return false
	}
	for _, entry := range subregistry.Entries() {
		if entry.Name == name {
			return true
		}
	}
	return false
}

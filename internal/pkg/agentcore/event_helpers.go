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

package agentcore

// eventProvider returns the provider ID for display formatting on assistant events.
func eventProvider(cfg Config) string {
	if cfg.AI == nil {
		return ""
	}
	if p := cfg.AI.LastCompletionProvider(); p != "" {
		return string(p)
	}
	if p := cfg.AI.SessionProvider(); p != "" {
		return string(p)
	}
	return string(cfg.AI.DefaultProviderID())
}

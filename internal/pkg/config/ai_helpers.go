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

// providerSpec is one built-in AI provider. baseURL depends on api style because
// MiniMax publishes separate OpenAI and Anthropic endpoints.
type providerSpec struct {
	id      string
	cloud   bool
	model   string
	baseURL func(apiStyle string) string
}

// CloudProviderIDs returns cloud providers used for ai.default, fallback, and --model.
func CloudProviderIDs() []string {
	var ids []string
	for _, spec := range builtinProviderSpecs() {
		if spec.cloud {
			ids = append(ids, spec.id)
		}
	}
	return ids
}

// AllProviderIDs returns every built-in provider ID including local Ollama.
func AllProviderIDs() []string {
	specs := builtinProviderSpecs()
	ids := make([]string, len(specs))
	for i, spec := range specs {
		ids[i] = spec.id
	}
	return ids
}

// IsCloudProvider reports whether id may be ai.default or session --model override.
func IsCloudProvider(id string) bool {
	spec, ok := lookupProviderSpec(id)
	return ok && spec.cloud
}

// IsKnownProvider reports whether id is a built-in provider (including ollama).
func IsKnownProvider(id string) bool {
	_, ok := lookupProviderSpec(id)
	return ok
}

func builtinProviderSpecs() []providerSpec {
	return []providerSpec{
		{
			id:      ProviderDeepseek,
			cloud:   true,
			model:   DefaultDeepseekModel,
			baseURL: func(string) string { return DefaultDeepseekBase },
		},
		{
			id:    ProviderMinimax,
			cloud: true,
			model: DefaultMinimaxModel,
			baseURL: func(apiStyle string) string {
				if apiStyle == APIStyleAnthropic {
					return DefaultMinimaxBaseAnthropic
				}
				return DefaultMinimaxBaseOpenAI
			},
		},
		{
			id:      ProviderOllama,
			cloud:   false,
			model:   DefaultOllamaModel,
			baseURL: func(string) string { return DefaultOllamaBase },
		},
	}
}

func lookupProviderSpec(id string) (providerSpec, bool) {
	for _, spec := range builtinProviderSpecs() {
		if spec.id == id {
			return spec, true
		}
	}
	return providerSpec{}, false
}

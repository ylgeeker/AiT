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

package ollama

import (
	"github.com/os-artificer/AiT/internal/pkg/config"
)

func ollamaEndpointsFromConfig(cfg *config.Config, baseURL, model string) (string, string) {
	if cfg == nil {
		return baseURL, model
	}
	pc, ok := cfg.Provider(config.ProviderOllama)
	if !ok {
		return baseURL, model
	}
	if pc.BaseURL != "" {
		baseURL = pc.BaseURL
	}
	if pc.Model != "" {
		model = pc.Model
	}
	return baseURL, model
}

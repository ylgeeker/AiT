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

package runner

import (
	pkgai "github.com/os-artificer/AiT/internal/pkg/ai"
	"github.com/os-artificer/AiT/internal/pkg/config"
)

// NewAIClient builds an AI client from cfg, optionally overriding the cloud provider.
func NewAIClient(cfg *config.Config, modelOverride string) (*pkgai.Client, error) {
	var opts []pkgai.ClientOption
	if modelOverride != "" {
		opts = append(opts, pkgai.WithProviderOverride(modelOverride))
	}
	return pkgai.NewClient(cfg, opts...)
}

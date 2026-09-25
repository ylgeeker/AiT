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
	pkgagent "github.com/os-artificer/AiT/internal/pkg/ai/agent"
)

// EnvelopeToAction maps a pkg/agent AgentEnvelope to the aissh Action DTO.
// Fields are copied verbatim; callers must ensure e passed ValidateEnvelope.
func EnvelopeToAction(e pkgagent.AgentEnvelope) Action {
	return Action{
		Action:  e.Action,
		Command: e.Command,
		Args:    e.Args,
		Text:    e.Text,
	}
}

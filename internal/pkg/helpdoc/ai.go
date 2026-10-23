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

package helpdoc

import (
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

// AI help strings — single source for ait agent, ollama, and in-REPL help.
var (
	ExampleAIProviderEnableDeepseek  = toolname.AgentCLI("provider enable deepseek")
	ExampleAIProviderEnableOllamaYes = toolname.AgentCLI("provider enable ollama --yes")
	ExampleAIProviderDisableMinimax  = toolname.AgentCLI("provider disable minimax")
	ExampleAIOllamaStatus            = toolname.AgentCLI("ollama status")
	ExampleAIOllamaSetupYes          = toolname.AgentCLI("ollama setup --yes")
	ExampleAIOllamaStartWait         = toolname.AgentCLI("ollama start --wait")
	ExampleAIOllamaDisableStop       = toolname.AgentCLI("ollama disable --stop")

	ExampleMetaSuspend = toolname.MetaLine("suspend")
	ExampleMetaLocal   = toolname.MetaLine("local")

	AIHelpFlagsBody = "  --model   Override AI provider for this session (uses provider model/base_url from config)\n"

	AIHelpTieredRoutingBody = "When ai.tiered.enabled in config.yaml, short first-turn prompts use local Ollama; " +
		"longer or complex prompts use cloud providers. Setup: " + toolname.AgentCLI("ollama setup") +
		". List: " + toolname.AgentCLI("model list") + "\n"

	AIIntroLocal = "Ops AI Agent plugin (experimental, on by default): model list, " +
		"provider/ollama, and remote connect. Interactive Agent Loop: " + toolname.AIT + " agent."
)

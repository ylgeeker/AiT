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

package provider

import (
	"fmt"

	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

const (
	hintSetCloudKeysOrOllama   = "Set a cloud provider API key or enable tiered Ollama routing."
	hintSetCloudAPIKey         = "Set an API key for a cloud provider."
	hintStoreAPIKeyInKeyring   = "Store the provider API key in the OS keyring."
	hintProvideChatMessage     = "Provide a goal on the command line or in the Agent Loop."
	hintEnableProviderOverride = "Enable the provider before using it as --model override."
)

var (
	hintSetDeepseekAPIKey = toolname.AgentCLI("config set ai.deepseek.api_key <key>")
	hintSetMinimaxAPIKey  = toolname.AgentCLI("config set ai.minimax.api_key <key>")
	hintEnableOllama      = toolname.AgentCLI("provider enable ollama")
	hintAIHelp            = toolname.AgentCLI("--help")
	hintOllamaSetup       = toolname.AgentCLI("ollama setup")
	hintAitAgent          = toolname.AITCLI("agent")
	hintAitAgentGoal      = toolname.AITCLI(`agent "check disk usage"`)

	// HintAIHelp is the standard ait agent --help Try line (shared with client).
	HintAIHelp = hintAIHelp
	// HintOllamaSetup is the ait agent ollama setup Try line.
	HintOllamaSetup = hintOllamaSetup
	// HintStoreAPIKeyInKeyring explains keyring storage in API key errors.
	HintStoreAPIKeyInKeyring = hintStoreAPIKeyInKeyring
	// HintSetCloudAPIKey is the one-line cloud API key setup hint.
	HintSetCloudAPIKey = hintSetCloudAPIKey
)

// ConfigSetAPIKeyCmd returns the ait agent config set command for a provider API key.
func ConfigSetAPIKeyCmd(id string) string {
	return toolname.AgentCLI(fmt.Sprintf("config set ai.%s.api_key <key>", id))
}

// ConfigGetAPIKeyCmd returns the ait agent config get command for a provider API key.
func ConfigGetAPIKeyCmd(id string) string {
	return toolname.AgentCLI(fmt.Sprintf("config get ai.%s.api_key", id))
}

// ProviderEnableCmd returns the ait agent provider enable command for id.
func ProviderEnableCmd(id string) string {
	return toolname.AgentCLI("provider enable " + id)
}

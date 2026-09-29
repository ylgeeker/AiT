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

import "encoding/json"

// Public action names emitted in the Agent-mode JSON envelope.
const (
	ToolRunShell = "run_shell"
	ToolAISshCLI = "aish_cli"
	ToolAitCLI   = "ait_cli" // legacy alias for ToolAISshCLI
	ToolToolbox  = "toolbox"
	ToolAnswer   = "answer"
)

// FCCall is a provider-neutral function-call payload (name + JSON arguments).
type FCCall struct {
	Name      string
	Arguments string
}

type toolArgField struct {
	jsonKey   string
	outField  string
	actionOut string
}

func toolCallConverter(name string) (toolArgField, bool) {
	for _, def := range ToolDefinitions() {
		if def.Name != name || def.CanonicalAction == "" || def.ParamName == "" {
			continue
		}
		return toolArgField{
			jsonKey:   def.ParamName,
			outField:  def.ParamName,
			actionOut: def.CanonicalAction,
		}, true
	}
	return toolArgField{}, false
}

// ToolCallToJSON converts a function call into the AiT Agent-mode JSON action envelope.
func ToolCallToJSON(tc FCCall) (string, error) {
	if tc.Name == "" {
		return "", ErrUnsupportedToolCall
	}
	conv, ok := toolCallConverter(tc.Name)
	if !ok {
		return "", ErrUnknownTool(tc.Name)
	}

	var raw map[string]string
	if err := json.Unmarshal([]byte(NormalizeES6UnicodeJSON(tc.Arguments)), &raw); err != nil {
		return "", err
	}
	b, err := json.Marshal(map[string]string{
		JSONFieldAction: conv.actionOut,
		conv.outField:   raw[conv.jsonKey],
	})
	if err != nil {
		return "", err
	}
	return string(b), nil
}

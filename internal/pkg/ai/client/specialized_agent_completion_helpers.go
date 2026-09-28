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

package client

func appendSystemBlock(system, content string) string {
	if system == "" {
		return content
	}
	return system + "\n\n" + content
}

func anthropicBlockTextParts(block anthropicContentBlock) (reply, thinking string) {
	if block.Type == anthropicBlockThinking {
		content := block.Thinking
		if content == "" {
			content = block.Text
		}
		if content == "" {
			return "", ""
		}
		return "", content
	}
	if block.Type == anthropicBlockRedactedThinking {
		return "", ""
	}
	if block.Text == "" {
		return "", ""
	}
	return block.Text, ""
}

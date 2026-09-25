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

package minimax

import (
	"regexp"
	"strings"

	"github.com/os-artificer/AiT/internal/pkg/aidisplay"
	"github.com/os-artificer/AiT/internal/pkg/config"
)

const (
	thinkOpen  = "\x3Cthink\x3E"
	thinkClose = "\x3C/think\x3E"
)

var (
	xmlToolCallRe = regexp.MustCompile(`(?s)<minimax:tool_call>.*?</minimax:tool_call>`)
	xmlInvokeRe   = regexp.MustCompile(`(?s)<invoke\b[^>]*>.*?</invoke>`)
	thinkBlockRe  = regexp.MustCompile(`(?s)` + thinkOpen + `.*?` + thinkClose)
	orphanCloseRe = regexp.MustCompile(`(?s)</(?:minimax:tool_call|invoke|parameter)>`)
)

func init() {
	aidisplay.RegisterNormalizer(config.ProviderMinimax, normalizer{})
}

type normalizer struct{}

// Normalize strips display-layer minimax XML/thinking residue that escaped protocol parsing.
func (normalizer) Normalize(text string) string {
	cleaned := thinkBlockRe.ReplaceAllString(text, "")
	cleaned = xmlToolCallRe.ReplaceAllString(cleaned, "")
	cleaned = xmlInvokeRe.ReplaceAllString(cleaned, "")
	cleaned = orphanCloseRe.ReplaceAllString(cleaned, "")
	return strings.TrimSpace(cleaned)
}

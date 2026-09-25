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

package minimax_test

import (
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/aidisplay"
	"github.com/os-artificer/AiT/internal/pkg/config"

	_ "github.com/os-artificer/AiT/internal/pkg/ai/display/minimax"
)

func TestNormalizer_MF3_StripsXMLFallback(t *testing.T) {
	raw := "answer <minimax:tool_call><invoke name=\"run\"></invoke></minimax:tool_call>"
	out := aidisplay.Render(aidisplay.ContentKindMarkdown, config.ProviderMinimax, raw, aidisplay.FormatOpts{Width: 60})
	if strings.Contains(out, "minimax:tool_call") || strings.Contains(out, "<invoke") {
		t.Fatalf("xml leaked into display: %q", out)
	}
	if !strings.Contains(out, "answer") {
		t.Fatalf("missing answer text: %q", out)
	}
}

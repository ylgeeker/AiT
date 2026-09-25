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

package plain_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/agentcore"
	"github.com/os-artificer/AiT/internal/pkg/agentui/plain"

	_ "github.com/os-artificer/AiT/internal/pkg/ai/display/registry/all"
)

func TestRenderer_MF6_PlainToolOutputUnchanged(t *testing.T) {
	var buf bytes.Buffer
	r := plain.Renderer{Out: &buf}
	emit := r.NewEmitter()
	emit.Emit(agentcore.Event{Kind: agentcore.EventToolOutput, Output: "total 8\n# line"})
	if !strings.Contains(buf.String(), "total 8") || !strings.Contains(buf.String(), "# line") {
		t.Fatalf("output=%q", buf.String())
	}
}

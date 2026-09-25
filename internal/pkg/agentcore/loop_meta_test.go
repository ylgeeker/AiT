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

package agentcore

import "testing"

func TestNormalizeLoopLine_SlashCommands(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, want string
	}{
		{"/exit", metaExit},
		{"/quit", metaExit},
		{"/help", metaHelp},
		{"/clear", "clear"},
		{"/history -c", "history -c"},
		{"/model list", "model list"},
		{"list files", "list files"},
		{"/unknown", "/unknown"},
	}
	for _, tc := range cases {
		if got := normalizeLoopLine(tc.in); got != tc.want {
			t.Fatalf("normalizeLoopLine(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

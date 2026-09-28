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

package shellwords

import (
	"strings"
	"testing"
)

func TestAppend_EmptySkipped(t *testing.T) {
	t.Parallel()

	var words []string
	var cur strings.Builder
	Append(&words, &cur)
	if len(words) != 0 {
		t.Fatalf("words = %v, want empty", words)
	}
}

func TestAppend_NonEmptyThenReset(t *testing.T) {
	t.Parallel()

	var words []string
	var cur strings.Builder
	cur.WriteString("hello")
	Append(&words, &cur)

	if len(words) != 1 || words[0] != "hello" {
		t.Fatalf("words = %v, want [hello]", words)
	}
	if cur.Len() != 0 {
		t.Fatalf("cur.Len() = %d, want 0 after Append", cur.Len())
	}

	cur.WriteString("world")
	Append(&words, &cur)
	if len(words) != 2 || words[1] != "world" {
		t.Fatalf("words = %v, want [hello world]", words)
	}
}

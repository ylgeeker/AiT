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

package terminal

import (
	"bytes"
	"os"
	"testing"
)

func TestIsTerminalReader_NonFile(t *testing.T) {
	if IsTerminalReader(bytes.NewReader(nil)) {
		t.Fatal("bytes.Reader must not be a TTY")
	}
}

func TestIsTerminalFile_Nil(t *testing.T) {
	if IsTerminalFile(nil) {
		t.Fatal("nil file")
	}
}

func TestIsTerminalFile_StdinStat(t *testing.T) {
	// Non-TTY CI: may be false; only assert no panic and consistency with Reader.
	got := IsTerminalFile(os.Stdin)
	if IsTerminalReader(os.Stdin) != got {
		t.Fatalf("File vs Reader mismatch: file=%v", got)
	}
}

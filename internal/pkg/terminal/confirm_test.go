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
	"io"
	"strings"
	"testing"
)

func TestDefaultConfirm_Yes(t *testing.T) {
	in := strings.NewReader("y\n")
	var out strings.Builder
	ok, err := DefaultConfirm(in, &out)("Overwrite? [y/N] ")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected yes")
	}
	if !strings.Contains(out.String(), "[y/N]") {
		t.Fatalf("prompt not written: %q", out.String())
	}
}

func TestDefaultConfirm_NoBlank(t *testing.T) {
	in := strings.NewReader("\n")
	ok, err := DefaultConfirm(in, io.Discard)("p ")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected no on blank")
	}
}

func TestDefaultConfirm_EOFIsNo(t *testing.T) {
	in := strings.NewReader("")
	ok, err := DefaultConfirm(in, io.Discard)("p ")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected no on EOF")
	}
}

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
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ConfirmFunc prompts the user and returns whether they chose yes (y).
type ConfirmFunc func(prompt string) (bool, error)

// DefaultConfirm prompts on out and reads from in (TTY-style [y/N]).
// Blank line, EOF, and any answer other than "y" are treated as no.
func DefaultConfirm(in io.Reader, out io.Writer) ConfirmFunc {
	return func(prompt string) (bool, error) {
		_, _ = fmt.Fprint(out, prompt)
		line, err := ReadConfirmLine(in)
		if err != nil && err != io.EOF {
			return false, err
		}
		return strings.EqualFold(strings.TrimSpace(line), ConfirmResponseYes), nil
	}
}

// ReadConfirmLine reads one line from r (including trailing newline when present).
// EOF with partial data returns the data and nil when the line is non-empty; bare EOF returns "", io.EOF.
func ReadConfirmLine(r io.Reader) (string, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	line, err := br.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	if err == io.EOF && line == "" {
		return "", io.EOF
	}
	return line, nil
}

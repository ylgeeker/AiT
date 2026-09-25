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

package helpdoc

import (
	"io"
)

const sectionTTYKeysCommon = `  Ctrl+A, Home     move cursor to start of line
  Ctrl+E, End      move cursor to end of line
  Left, Right      move cursor one character (UTF-8 aware)
  Up, Down         recall previous / next line from loaded input history (incl. prior sessions)
  Ctrl+R           reverse incremental search through input history
  Backspace, Del   delete character before cursor
  Ctrl+U           delete entire line
  Ctrl+K           delete from cursor to end of line
  Ctrl+W           delete word before cursor
  Ctrl+L           clear screen and redraw prompt (line kept)
  Ctrl+C           cancel current line (REPL continues)
  Ctrl+D           exit REPL when line is empty (EOF)
  Tab              insert indentation
`

const sectionTTYKeysShellNote = "  clear            clear terminal screen (command)\n"

// WriteTTYKeysHelp writes the Keys section for local TTY REPL input (shell).
func WriteTTYKeysHelp(w io.Writer) error {
	b := NewBody(w)
	return b.WriteSection(SectionKeys, sectionTTYKeysCommon+sectionTTYKeysShellNote)
}

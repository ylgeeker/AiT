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

package localmode

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func (l *Shell) finishCopyOutOnError(err error) {
	if l.closing.Load() || errors.Is(err, io.EOF) {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, "%sAiT: local shell: %v%s", noticeEOL, err, noticeEOL)
}

func (l *Shell) teardownCopyOutAfterRead() {
	l.mu.Lock()
	onExit := l.onExit
	l.pty = nil
	l.cmd = nil
	l.onExit = nil
	l.mu.Unlock()
	if onExit != nil {
		onExit()
	}
}

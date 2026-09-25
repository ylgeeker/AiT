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

package ait

import (
	"context"
	"io"

	"github.com/os-artificer/AiT/internal/ait/shell"
	"github.com/os-artificer/AiT/internal/pkg/terminal"
)

// RunREPL runs the ait> toolbox REPL (install/update/upgrade).
func RunREPL(ctx context.Context, in io.Reader, out, errOut io.Writer) error {
	return shell.RunREPL(ctx, in, out, errOut)
}

// ExecuteLine runs one toolbox REPL line (tests and shell package).
func ExecuteLine(ctx context.Context, line string, hist *terminal.LineHistory, out, errOut io.Writer) (exitREPL bool, err error) {
	return shell.ExecuteLine(ctx, line, hist, out, errOut)
}

// WriteShellHelp prints ait shell --help.
func WriteShellHelp(w io.Writer) error {
	return shell.WriteShellHelp(w)
}

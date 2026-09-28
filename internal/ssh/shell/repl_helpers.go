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

package shell

import (
	"context"
	"io"
	"os"

	"github.com/os-artificer/AiT/internal/pkg/shellext"
	"github.com/os-artificer/AiT/internal/pkg/terminal"

	"golang.org/x/term"
)

// setShellTTYIO configures CRLF stdout and silenced stderr for the duration of aissh shell.
func setShellTTYIO(enabled bool) {
	if enabled {
		shellCLIOut = terminal.NewCRLFWriter(os.Stdout)
		shellCLIErr = io.Discard
		return
	}

	shellCLIOut = os.Stdout
	shellCLIErr = os.Stderr
}

// runExternalWithTTYHandoff returns an externalRunner that restores the cooked TTY
// for the child, then re-enters raw mode afterward. fd and old must match the
// MakeRaw state owned by the REPL. Restore errors return as-is; MakeRaw failure
// is returned only when runErr is nil (runErr takes precedence when both fail).
func runExternalWithTTYHandoff(fd int, old *term.State) func(context.Context, string, []string) error {
	return func(ctx context.Context, path string, argv []string) error {
		if restoreErr := term.Restore(fd, old); restoreErr != nil {
			return restoreErr
		}

		runErr := shellext.Run(ctx, path, argv, os.Stdin, os.Stdout, os.Stderr)

		if _, rawErr := term.MakeRaw(fd); rawErr != nil && runErr == nil {
			return rawErr
		}
		return runErr
	}
}

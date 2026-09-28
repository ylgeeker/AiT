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
	"errors"
	"fmt"
	"os"

	"github.com/os-artificer/AiT/internal/pkg/pluginbin"
	"github.com/os-artificer/AiT/internal/pkg/shellext"

	"golang.org/x/term"
)

// runPluginWithTTYHandoff returns a Runner that restores the cooked TTY for the
// child process, then re-enters raw mode afterward. fd and old must match the
// MakeRaw state owned by the REPL. Returns Restore errors as-is; MakeRaw failure
// wraps as "terminal: …" and Joins with runErr when both fail.
func runPluginWithTTYHandoff(fd int, old *term.State) pluginbin.Runner {
	return func(ctx context.Context, path string, argv []string, pluginCMD string) error {
		if restoreErr := term.Restore(fd, old); restoreErr != nil {
			return restoreErr
		}

		runErr := shellext.RunWithOpts(ctx, path, argv, pluginCMD, os.Stdin, os.Stdout, os.Stderr)

		_, rawErr := term.MakeRaw(fd)
		if rawErr == nil {
			return runErr
		}
		termErr := fmt.Errorf("terminal: %w", rawErr)
		if runErr == nil {
			return termErr
		}
		return errors.Join(termErr, runErr)
	}
}

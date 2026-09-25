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
	"os"

	"github.com/os-artificer/AiT/internal/ait/dispatch"
	"github.com/os-artificer/AiT/internal/ait/router"
	"github.com/os-artificer/AiT/internal/ait/version"
	"github.com/os-artificer/AiT/internal/pkg/config"
)

// Execute runs the top-level ait CLI with args on os.Stdout and os.Stderr.
func Execute(args []string) error {
	return ExecuteWriters(args, os.Stdout, os.Stderr)
}

// ExecuteWriters runs the top-level ait CLI with custom stdout and stderr.
func ExecuteWriters(args []string, stdout, stderr io.Writer) error {
	args, saw := config.StripAIExperimentalFlag(args)
	if saw {
		config.SetAIExperimentalCLI(true)
	}

	return router.Route(context.Background(), args, stdout, stderr, router.RouteDeps{
		WriteRootHelp:  WriteRootHelp,
		WriteShellHelp: WriteShellHelp,
		RunShellREPL:   RunREPL,
		VersionLine:    version.Version,
		RunPlugin:      dispatch.RunPlugin,
	})
}

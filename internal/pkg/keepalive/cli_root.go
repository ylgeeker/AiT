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

package keepalive

import (
	"fmt"
	"io"
	"os"

	"github.com/os-artificer/AiT/internal/pkg/aitplugin"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
	"github.com/os-artificer/AiT/internal/pkg/platform"
	"github.com/os-artificer/AiT/internal/pkg/toolname"

	"github.com/spf13/cobra"
)

// Execute runs the keepalive CLI on os.Stdout and os.Stderr.
// args are the arguments after the program name. Empty args print help.
// Returns an unknown-command hint, a config or serve error, or a write error from help or version.
func Execute(args []string) error {
	return ExecuteWriters(args, os.Stdout, os.Stderr)
}

// ExecuteWriters runs the keepalive CLI with custom stdout and stderr.
// args are the arguments after the program name. Empty args print help. stdout and stderr receive help and command output.
// Returns an unknown-command hint, a config or serve error, or a write error from help or version.
func ExecuteWriters(args []string, stdout, stderr io.Writer) error {
	if len(args) == 1 && args[0] == aitplugin.ManifestArg {
		return aitplugin.WriteManifest(stdout, toolname.AITKeep)
	}
	if len(args) > 0 && pkghelpdoc.HelpRequested(args) {
		return WriteRootHelp(stdout)
	}
	if pkghelpdoc.VersionRequested(args) {
		_, err := fmt.Fprintln(stdout, Version())
		return err
	}
	if len(args) == 0 {
		return WriteRootHelp(stdout)
	}
	root := newRootCmd()
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs(args)
	if registered(args[0]) {
		return root.Execute()
	}
	return errhint.WithTryMsg(
		fmt.Sprintf("unknown command %q", args[0]),
		"Unrecognized keepalive command.",
		toolname.KeepCLI("--help"),
		toolname.KeepCLI("serve"),
	)
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           toolname.KeepCLI(""),
		Short:         "Keep hosted processes alive and within resource limits",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	for _, name := range commandOrder {
		root.AddCommand(newKeepCommand(name))
	}
	root.SetHelpFunc(func(c *cobra.Command, _ []string) {
		_ = WriteRootHelp(c.OutOrStdout())
	})
	return root
}

func appRoot() (string, error) {
	return platform.DataDir()
}

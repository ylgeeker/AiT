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

package configcli

import (
	"io"

	pkgconfig "github.com/os-artificer/AiT/internal/pkg/config"
	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
	"github.com/os-artificer/AiT/internal/pkg/platform"

	"github.com/spf13/cobra"
)

// NewSet builds `config set <key> <value>` for the given scope.
// writeHelp is bound via pkghelpdoc.Bind for --help.
// RunE errors: RejectKey; keyring SetProviderAPIKey; Load/SetKey/Save failures.
func NewSet(o Options, writeHelp func(io.Writer) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "set <key> <value>",
		Short:         "Set config value",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSet(cmd, o, args[0], args[1])
		},
	}
	pkghelpdoc.Bind(cmd, writeHelp)
	return cmd
}

// NewGet builds `config get <key>` for the given scope.
// writeHelp is bound via pkghelpdoc.Bind for --help.
// RunE errors: RejectKey; keyring GetProviderAPIKey; Load/GetKey failures.
func NewGet(o Options, writeHelp func(io.Writer) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get <key>",
		Short:         "Get config value",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd, o, args[0])
		},
	}
	pkghelpdoc.Bind(cmd, writeHelp)
	return cmd
}

// NewInit builds `config init` for the given scope.
// writeHelp is bound via pkghelpdoc.Bind for --help.
// RunE errors: o.Ensure failure (path, load, or save).
func NewInit(o Options, writeHelp func(io.Writer) error) *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:           "init",
		Short:         o.InitShort,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			res, err := o.Ensure(dryRun)
			if err != nil {
				return err
			}
			printEnsureResult(cmd.OutOrStdout(), res, dryRun)
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print changes without writing config.yaml")
	pkghelpdoc.Bind(cmd, writeHelp)
	return cmd
}

// NewView builds `config view` for the given scope.
// writeHelp is bound via pkghelpdoc.Bind for --help.
// RunE errors: EnsureConfigFileIfMissing; ReadRawFile.
func NewView(o Options, writeHelp func(io.Writer) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "view",
		Short:         "Print config.yaml contents",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := o.EnsureConfigFileIfMissing(); err != nil {
				return err
			}

			_, content, err := pkgconfig.ReadRawFile()
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write(content)
			return err
		},
	}
	pkghelpdoc.Bind(cmd, writeHelp)
	return cmd
}

// NewEdit builds `config edit` for the given scope.
// writeHelp is bound via pkghelpdoc.Bind for --help.
// RunE errors: EnsureConfigFileIfMissing; ConfigPath; OpenEditorFn (via WrapEditorErr); FinishConfigValidation.
func NewEdit(o Options, writeHelp func(io.Writer) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "edit",
		Short:         "Edit config.yaml in $VISUAL or $EDITOR",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := o.EnsureConfigFileIfMissing(); err != nil {
				return err
			}

			path, err := platform.ConfigPath()
			if err != nil {
				return err
			}

			if err := OpenEditorFn(path); err != nil {
				return WrapEditorErr(err)
			}

			return o.FinishConfigValidation(cmd, path)
		},
	}
	pkghelpdoc.Bind(cmd, writeHelp)
	return cmd
}

// NewValidate builds `config validate` for the given scope.
// writeHelp is bound via pkghelpdoc.Bind for --help.
// RunE errors: o.Validate I/O; ErrValidationFailed when issues are present.
func NewValidate(o Options, writeHelp func(io.Writer) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "validate",
		Short:         o.ValidateShort,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			report, err := o.Validate("")
			if err != nil {
				return err
			}
			return EmitValidationResult(cmd, report)
		},
	}
	pkghelpdoc.Bind(cmd, writeHelp)
	return cmd
}

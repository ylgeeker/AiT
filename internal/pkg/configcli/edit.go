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
	"errors"
	"os"
	"runtime"

	pkgconfig "github.com/os-artificer/AiT/internal/pkg/config"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
	"github.com/os-artificer/AiT/internal/pkg/platform"

	"github.com/spf13/cobra"
)

// OpenEditorFn is the editor launcher; tests override this to skip launching $EDITOR.
var OpenEditorFn = platform.OpenEditor

// EnsureConfigFileIfMissing creates config.yaml with scoped defaults when missing.
// Uses o.Ensure(false). Errors: platform.ConfigPath failure; Stat failure; Ensure failure.
func (o Options) EnsureConfigFileIfMissing() error {
	path, err := platform.ConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	_, err = o.Ensure(false)
	return err
}

// FinishConfigValidation validates YAML after edit and prints the report.
// path is the config file path passed to o.Validate. Errors: Validate I/O or parse errors;
// ErrValidationFailed when the report has non-warning issues.
func (o Options) FinishConfigValidation(cmd *cobra.Command, path string) error {
	report, err := o.Validate(path)
	if err != nil {
		return err
	}
	return EmitValidationResult(cmd, report)
}

// EmitValidationResult writes validation output to cmd writers.
// Returns nil when report.OK(); otherwise ErrValidationFailed after writing issues to stderr.
func EmitValidationResult(cmd *cobra.Command, report *pkgconfig.ValidationReport) error {
	if report.OK() {
		pkgconfig.FormatValidation(cmd.OutOrStdout(), report)
		return nil
	}

	pkgconfig.FormatValidation(cmd.ErrOrStderr(), report)
	return ErrValidationFailed
}

// WrapEditorErr adds hints when no editor is configured.
// Non-ErrNoEditor errors are returned unchanged. ErrNoEditor is wrapped with OS-specific Try lines.
func WrapEditorErr(err error) error {
	if !errors.Is(err, platform.ErrNoEditor) {
		return err
	}
	if runtime.GOOS == "windows" {
		return errhint.WithTryMsg(
			platform.ErrNoEditor.Error(),
			"Set EDITOR or use a GUI editor.",
			"set EDITOR=notepad",
			"set EDITOR=code --wait",
		)
	}
	return errhint.WithTryMsg(
		platform.ErrNoEditor.Error(),
		"Set EDITOR or install vi.",
		"export EDITOR=vi",
		"export EDITOR=\"code --wait\"",
	)
}

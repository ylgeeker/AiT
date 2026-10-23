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

package agent

import (
	"context"
	"io"

	"github.com/os-artificer/AiT/internal/ait/agent/connect"
	"github.com/os-artificer/AiT/internal/ait/agent/deps"
	"github.com/os-artificer/AiT/internal/ait/agent/subregistry"
	pkgai "github.com/os-artificer/AiT/internal/pkg/ai"
	"github.com/os-artificer/AiT/internal/pkg/clihints"
	"github.com/os-artificer/AiT/internal/pkg/config"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
	"github.com/os-artificer/AiT/internal/pkg/toolname"

	"github.com/spf13/cobra"
)

const adminRootName = "_agent"

// RunAdmin runs the agent management Cobra tree (connect, model, provider, ollama, config).
func RunAdmin(args []string, out, errOut io.Writer) error {
	d := deps.Deps{
		Exec:           execAisshCLI,
		RunWithMgr:     runWithConfig,
		ApplyOpsDomain: withOpsDomain,
	}
	root := newAdminCmd(d)
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetArgs(args)
	root.SilenceUsage = true
	root.SilenceErrors = true
	return wrapRemovedOneShotErr(args, root.Execute())
}

// LooksLikeConnectDestination reports whether first may be an SSH destination for bare connect routing.
func LooksLikeConnectDestination(first string) bool {
	return connect.LooksLikeDestination(first)
}

// newAdminCmd builds the hidden admin root with subcommands at the top level.
func newAdminCmd(d Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:           adminRootName,
		Hidden:        true,
		Short:         "Ops AI Agent (experimental; on by default)",
		Long:          pkghelpdoc.LongHint(toolname.Agent),
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if pkghelpdoc.HasHelpArg(args) {
				return WriteAIHelp(cmd.OutOrStdout())
			}
			if len(args) == 0 {
				if err := WriteAIHelp(cmd.OutOrStdout()); err != nil {
					return err
				}
				return errLoopMovedToAit()
			}
			if len(args) > 0 && connect.LooksLikeDestination(args[0]) {
				return d.RunWithMgr(func(ctx context.Context, cfg *config.Config, aiClient *pkgai.Client) error {
					return connect.Run(ctx, cfg, aiClient, args, d.Exec)
				})
			}
			return errhint.WithTryMsg(
				"unknown agent command",
				"Use "+toolname.AgentCLI("")+" subcommands for the task loop.",
				clihints.TryAitAgent(),
				toolname.AgentCLI("connect user@host"),
				toolname.AgentCLI("model list"),
			)
		},
	}
	cmd.PersistentFlags().String("model", "", "Override AI provider for this session (deepseek|minimax)")
	subregistry.Attach(cmd, deps.Deps{Exec: d.Exec, RunWithMgr: d.RunWithMgr, ApplyOpsDomain: d.ApplyOpsDomain})
	pkghelpdoc.Bind(cmd, WriteAIHelp)
	return cmd
}

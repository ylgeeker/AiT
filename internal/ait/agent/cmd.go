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
	"github.com/os-artificer/AiT/internal/ait/agent/deps"
	"github.com/os-artificer/AiT/internal/pkg/toolname"

	"github.com/spf13/cobra"
)

// Deps is an alias for shared ai subtree dependencies.
type Deps = deps.Deps

// CLIExecutor invokes aissh (or toolbox CLIs) with custom IO.
type CLIExecutor = deps.CLIExecutor

// RunWithMgr opens config and AI client under a signal-cancellable context.
type RunWithMgr = deps.RunWithMgr

// New builds the root `ai` cobra command tree from the admin command tree.
func New(d Deps) *cobra.Command {
	cmd := newAdminCmd(d)
	cmd.Use = toolname.Agent
	cmd.Hidden = false
	return cmd
}

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
	"strings"

	"github.com/os-artificer/AiT/internal/pkg/clihints"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

// removedOneShotCmds are former one-shot AI verbs; keep reserved from connect routing
// and map unknown-command errors to a Try that points at the Agent Loop.
var removedOneShotCmds = map[string]struct{}{
	"ask":   {},
	"agent": {},
}

func wrapRemovedOneShotErr(args []string, err error) error {
	if err == nil || len(args) == 0 {
		return err
	}
	word := strings.TrimSpace(args[0])
	if _, ok := removedOneShotCmds[word]; !ok {
		return err
	}
	if !strings.Contains(err.Error(), "unknown command") {
		return err
	}
	return errhint.WithTryMsg(
		"unknown agent command",
		"One-shot ask/agent were removed; use the interactive Agent Loop.",
		clihints.TryAitAgent(),
		toolname.AgentCLI("connect user@host"),
		toolname.AgentCLI("model list"),
	)
}

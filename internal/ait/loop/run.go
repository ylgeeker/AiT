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

// Package loop starts the interactive Agent task loop for bare ait / ait agent.
package loop

import (
	"context"
	"io"

	"github.com/os-artificer/AiT/internal/ait/agent"
)

// Run starts the v4 Agent Loop with optional initial goal and plain UI.
func Run(ctx context.Context, goal string, plain bool, out, errOut io.Writer) error {
	return agent.RunLoop(ctx, agent.Options{
		Nested: false,
		Plain:  plain,
		Goal:   goal,
		Out:    out,
		ErrOut: errOut,
	})
}

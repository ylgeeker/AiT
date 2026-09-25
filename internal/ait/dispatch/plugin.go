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

// Package dispatch runs discovered toolbox plugins for top-level ait argv.
package dispatch

import (
	"context"
	"io"
	"os"

	"github.com/os-artificer/AiT/internal/pkg/tooldispatch"
)

// RunPlugin executes argv via plugin discovery (install dir, then PATH for toolbox bins).
func RunPlugin(ctx context.Context, argv []string, out, errOut io.Writer) error {
	return tooldispatch.Exec(ctx, argv, os.Stdin, out, errOut)
}

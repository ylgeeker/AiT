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

package client

import (
	"context"
	"io"
)

// waitForwardStdinStep blocks for one stdin/ctx event.
// done is true when the forward loop should exit (returning err, which may be nil).
func (c *Client) waitForwardStdinStep(
	ctx context.Context,
	stdin io.Writer,
	stdinCh <-chan stdinByte,
	scan *metaScanner,
	opts InteractiveOpts,
) (err error, done bool) {
	select {
	case <-ctx.Done():
		return ctx.Err(), true
	case sb, ok := <-stdinCh:
		if !ok {
			return io.EOF, true
		}
		res := c.processForwardStdinByte(ctx, stdin, sb, scan, opts)
		if res.action == fwdStdinReturn {
			return res.err, true
		}
		return nil, false
	}
}

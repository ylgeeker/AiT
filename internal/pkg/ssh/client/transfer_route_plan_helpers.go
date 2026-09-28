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
	"io"

	"github.com/os-artificer/AiT/internal/pkg/transfer/route"
)

func (c *Client) planTransferNonChain(stdin io.Writer, plan route.Plan) (activeHost string, out route.Plan, earlyReturn bool) {
	out = plan
	if c.scroll == nil {
		out.Backend = route.BackendSFTP
		return "", out, true
	}
	activeHost, out = c.planTransferFromActiveShell(stdin, out)
	if out.Backend == route.BackendSFTP && activeHost == "" {
		return activeHost, out, true
	}
	return activeHost, out, false
}

func (c *Client) planTransferFromActiveShell(stdin io.Writer, plan route.Plan) (activeHost string, out route.Plan) {
	out = plan
	host, probeErr := c.probeActiveShellHost(stdin)
	if probeErr != nil {
		out.Backend = route.BackendSFTP
		return "", out
	}
	activeHost = host
	if !activeShellNeedsChain(c.opts.Host, c.opts.LoginDestination, activeHost) {
		return activeHost, out
	}
	c.SetChainMode(true)
	out.ChainOnly = true
	out.AutoChain = true
	autoChainNotice(activeHost, c.opts.Host)
	return activeHost, out
}

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
	"os"
	"strings"

	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

func (c *Client) handleMetaLineDefault(ctx context.Context, stdin io.Writer, line string, opts InteractiveOpts) error {
	if chain, ok := ParseMetaChainLine(line, opts.AliasLookup); ok {
		if err := c.handleSuspendMeta(stdin, opts); err != nil {
			return err
		}
		_ = c.Close()
		return chain
	}

	if line == metaLineBare || line == metaLineHelp {
		return c.runMetaCLI(stdin, opts, nil)
	}
	if cmd, args, isXfer := transferArgsFromMetaLine(line); isXfer {
		switch cmd {
		case "get":
			return c.handleTransferGet(ctx, stdin, args)
		case "put":
			return c.handleTransferPut(ctx, stdin, args)
		}
	}
	if strings.HasPrefix(line, metaPrefixAitSpace) {
		args := strings.Fields(strings.TrimPrefix(line, metaPrefixAitSpace))
		return c.runMetaCLI(stdin, opts, args)
	}
	termPrintln(os.Stdout, "AiT: unknown meta command (try "+toolname.MetaLine("help")+")")
	c.finishMetaInteraction(stdin)
	return nil
}

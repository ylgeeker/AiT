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

	"github.com/os-artificer/AiT/internal/pkg/tmux"

	"golang.org/x/crypto/ssh"
)

func (c *Client) runTmuxPreAttach(tmuxName string, initial []string) (postAttach []string, err error) {
	setup, after := partitionInitialCommands(initial, c.opts.User, c.opts.Host)
	if err := c.tmuxPreAttachSetup(tmuxName, setup); err != nil {
		return initial, wrapTmuxPreAttachErr(tmuxName, err)
	}
	return after, nil
}

func (c *Client) runInteractiveTmuxAttach(sess *ssh.Session, stdin io.Writer, opts InteractiveOpts, setupNeeded bool, postAttach *([]string)) error {
	if err := c.tmuxEnsure(stdin, opts.TmuxName); err != nil {
		sess.Close()
		return err
	}
	if err := c.ensureTmuxSessionReady(stdin, opts.TmuxName); err != nil {
		sess.Close()
		return err
	}

	if setupNeeded {
		after, err := c.runTmuxPreAttach(opts.TmuxName, opts.InitialCommands)
		if err != nil {
			sess.Close()
			return err
		}
		*postAttach = after
	}

	if err := sess.Start(tmux.AttachCommand(opts.TmuxName)); err != nil {
		sess.Close()
		return err
	}
	return nil
}

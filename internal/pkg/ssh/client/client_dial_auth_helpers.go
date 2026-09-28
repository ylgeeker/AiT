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
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/os-artificer/AiT/internal/pkg/ssh/client/dial"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func appendDefaultSSHKeys(methods []ssh.AuthMethod, opts Options) []ssh.AuthMethod {
	home, _ := os.UserHomeDir()
	for _, name := range []string{"id_ed25519", "id_rsa"} {
		p := filepath.Join(home, ".ssh", name)
		if opts.IdentityPath != "" && p == opts.IdentityPath {
			continue
		}
		if _, err := os.Stat(p); err != nil {
			continue
		}
		signer, err := dial.LoadSigner(p, opts.Password)
		if err != nil {
			continue
		}
		return append(methods, ssh.PublicKeys(signer))
	}
	return methods
}

func appendIdentityAuth(methods []ssh.AuthMethod, opts Options) ([]ssh.AuthMethod, error) {
	if opts.IdentityPath == "" {
		return appendDefaultSSHKeys(methods, opts), nil
	}
	if _, err := os.Stat(opts.IdentityPath); err != nil {
		if opts.IdentityExplicit {
			return methods, fmt.Errorf("identity file not found: %s", opts.IdentityPath)
		}
		return appendDefaultSSHKeys(methods, opts), nil
	}
	signer, err := dial.LoadSigner(opts.IdentityPath, opts.Password)
	if err == nil {
		return append(methods, ssh.PublicKeys(signer)), nil
	}
	if opts.IdentityExplicit {
		return methods, err
	}
	return appendDefaultSSHKeys(methods, opts), nil
}

func appendSSHAgentAuth(methods []ssh.AuthMethod, opts Options) ([]ssh.AuthMethod, net.Conn, error) {
	if !opts.UseAgent {
		return methods, nil, nil
	}
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return methods, nil, nil
	}
	conn, err := net.Dial("unix", sock)
	if err != nil {
		return methods, nil, nil
	}
	ag := agent.NewClient(conn)
	methods = append(methods, ssh.PublicKeysCallback(ag.Signers))
	return methods, conn, nil
}

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

package resolve

import (
	"time"

	"github.com/os-artificer/AiT/internal/pkg/platform"
	"github.com/os-artificer/AiT/internal/pkg/ssh/client"
	"github.com/os-artificer/AiT/internal/pkg/ssh/openssh"
)

// Resolved is the effective SSH connection parameters after ssh -G and CLI overrides.
type Resolved struct {
	Host         string
	Port         int
	User         string
	IdentityPath string
	// IdentityExplicit is true when the user passed -i on the CLI.
	IdentityExplicit bool
	Password         string

	ProxyJump    []string
	ProxyCommand string

	KnownHostsPath        string
	StrictHostKeyChecking string // yes, no, ask, accept-new
	ConnectTimeout        time.Duration
	BatchMode             bool

	RemoteCommand []string
	Alias         string
	// Destination is the CLI/OpenSSH destination before hostname resolution (e.g. devhost).
	Destination string
	Save        bool
	Tmux        bool
}

// Resolve runs ssh -G for spec.Destination and merges the parsed config with CLI flags in spec.
// Returns effective connection parameters for dial. Errors: platform.OpenSSHPath; ssh -G failure;
// identity file not found; could not resolve hostname; cannot combine --tmux with remote command.
func Resolve(spec openssh.Spec) (Resolved, error) {
	sshPath, err := platform.OpenSSHPath()
	if err != nil {
		return Resolved{}, err
	}

	out, err := sshGOutput(sshPath, sshGArgs(spec))
	if err != nil {
		return Resolved{}, err
	}

	r, err := resolvedFromSSHG(sshPath, spec, parseSSHGOutput(out))
	if err != nil {
		return Resolved{}, err
	}
	return validateResolved(r, spec)
}

// ClientOptions maps r to client.Options for dialing (UseAgent is always true).
// Returns a value copy suitable for client.Dial or DialChain.
func (r Resolved) ClientOptions() client.Options {
	return client.Options{
		Host:                  r.Host,
		Port:                  r.Port,
		User:                  r.User,
		LoginDestination:      r.Destination,
		IdentityPath:          r.IdentityPath,
		IdentityExplicit:      r.IdentityExplicit,
		Password:              r.Password,
		UseAgent:              true,
		ProxyJump:             append([]string(nil), r.ProxyJump...),
		ProxyCommand:          r.ProxyCommand,
		KnownHostsPath:        r.KnownHostsPath,
		StrictHostKeyChecking: r.StrictHostKeyChecking,
		ConnectTimeout:        r.ConnectTimeout,
		BatchMode:             r.BatchMode,
	}
}

// ToSSHOptions returns the Resolved fields as separate return values for legacy call sites.
// The returned slices and strings are copies or primitives from r; this method never errors.
func (r Resolved) ToSSHOptions() (
	host string,
	port int,
	user, identity, password string,
	proxyJump []string,
	proxyCommand string,
	knownHosts string,
	timeout time.Duration,
	batch bool,
) {
	return r.Host, r.Port, r.User, r.IdentityPath, r.Password,
		append([]string(nil), r.ProxyJump...), r.ProxyCommand,
		r.KnownHostsPath, r.ConnectTimeout, r.BatchMode
}

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

package hostkeys

import (
	"fmt"
	"net"
	"strings"

	"github.com/os-artificer/AiT/internal/pkg/ssh/client/dial"

	"golang.org/x/crypto/ssh"
)

func rejectChangedHostKey(host, primary string) error {
	return fmt.Errorf(
		"host key for %s has changed — possible MITM (refusing to connect); remove the old entry from %q",
		host, primary,
	)
}

func noteDisplayHost(displayHost, seen string) string {
	if seen != "" && displayHost == "host" {
		return seen
	}
	return displayHost
}

func unknownHostPort(remote net.Addr, fallback int) int {
	p := fallback
	if _, portStr, splitErr := net.SplitHostPort(remote.String()); splitErr == nil {
		if pp, perr := dial.ParsePortString(portStr); perr == nil && pp != 0 {
			p = pp
		}
	}
	if p == 0 {
		p = fallback
	}
	if p == 0 {
		p = 22
	}
	return p
}

func unknownHostMode(mode string, batch bool) string {
	if mode == "ask" && (batch || !IsInteractiveTTY()) {
		return "yes"
	}
	return mode
}

func applyUnknownHostPolicy(effective, primary, displayHost string, port int, key ssh.PublicKey) error {
	switch effective {
	case "accept-new":
		return appendKnownHost(primary, displayHost, port, key)
	case "ask":
		return promptAndAppendHostKey(primary, displayHost, port, key)
	default:
		return hostKeyUnknownError(displayHost, port, primary)
	}
}

func promptAndAppendHostKey(primary, displayHost string, port int, key ssh.PublicKey) error {
	portNote := ""
	if port != 22 {
		portNote = fmt.Sprintf(" (port %d)", port)
	}
	prompt := fmt.Sprintf(
		"The authenticity of host '%s'%s can't be established.\n"+
			"%s key fingerprint is %s.\n"+
			"Are you sure you want to continue connecting (yes/no)? ",
		displayHost, portNote, strings.ToUpper(key.Type()), ssh.FingerprintSHA256(key),
	)
	if !Confirm(prompt) {
		return fmt.Errorf("host key verification failed")
	}
	return appendKnownHost(primary, displayHost, port, key)
}

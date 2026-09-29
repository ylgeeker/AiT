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
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/os-artificer/AiT/internal/pkg/platform"
	"github.com/os-artificer/AiT/internal/pkg/ssh/openssh"
	"github.com/os-artificer/AiT/internal/pkg/ssh/stricthostkey"
)

type sshGConfig struct {
	single map[string]string
	multi  map[string][]string
}

func parseSSHGOutput(out []byte) sshGConfig {
	cfg := sshGConfig{
		single: make(map[string]string),
		multi:  make(map[string][]string),
	}
	sc := bufio.NewScanner(bytes.NewReader(out))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.ToLower(fields[0])
		val := strings.TrimSpace(strings.Join(fields[1:], " "))
		cfg.single[key] = val
		cfg.multi[key] = append(cfg.multi[key], val)
	}
	return cfg
}

func expandIdentityPaths(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if p == "" {
			continue
		}
		out = append(out, platform.ExpandSSHPath(p))
	}
	return out
}

func firstExistingIdentity(paths []string) string {
	for _, p := range paths {
		if p == "" {
			continue
		}
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func parsePort(s string) int {
	if s == "" {
		return 0
	}
	p, _ := strconv.Atoi(s)
	return p
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func resolveIdentityPath(spec openssh.Spec, identityFiles []string) (path string, explicit bool, err error) {
	if spec.Identity != "" {
		path = platform.ExpandSSHPath(spec.Identity)
		return path, true, statIdentityRequired(path)
	}
	if spec.Options["identityfile"] != "" {
		path = platform.ExpandSSHPath(spec.Options["identityfile"])
		return path, false, statIdentityRequired(path)
	}
	return firstExistingIdentity(identityFiles), false, nil
}

func statIdentityRequired(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("identity file not found: %s", path)
	}
	return nil
}

func sshGArgs(spec openssh.Spec) []string {
	args := []string{"-G"}
	if spec.ConfigFile != "" {
		args = append(args, "-F", spec.ConfigFile)
	}
	for _, j := range spec.ProxyJump {
		args = append(args, "-J", j)
	}
	if spec.Port != 0 {
		args = append(args, "-p", strconv.Itoa(spec.Port))
	}
	if spec.User != "" {
		args = append(args, "-l", spec.User)
	}
	if spec.Identity != "" {
		args = append(args, "-i", spec.Identity)
	}
	for k, v := range spec.Options {
		args = append(args, "-o", k+"="+v)
	}
	args = append(args, spec.Destination)
	return args
}

// endpointFromSSHG applies ssh -G hostname/user/port, then a non-zero CLI port and a non-empty CLI user.
// Default port 22 is filled before the override so an explicit port replaces it. ProxyJump uses the same order.
func endpointFromSSHG(hostname, fallbackHost, cfgUser string, cfgPort, overridePort int, overrideUser string) (host, user string, port int) {
	host = firstNonEmpty(hostname, fallbackHost)
	port = cfgPort
	if port == 0 {
		port = 22
	}
	if overridePort != 0 {
		port = overridePort
	}
	user = cfgUser
	if overrideUser != "" {
		user = overrideUser
	}
	return host, user, port
}

func resolvedFromSSHG(sshPath string, spec openssh.Spec, cfg sshGConfig) (Resolved, error) {
	host, user, port := endpointFromSSHG(
		cfg.single["hostname"], spec.Destination, cfg.single["user"],
		parsePort(cfg.single["port"]), spec.Port, spec.User,
	)
	r := Resolved{
		Host:                  host,
		Port:                  port,
		User:                  user,
		ProxyCommand:          cfg.single["proxycommand"],
		KnownHostsPath:        platform.ExpandSSHPath(cfg.single["userknownhostsfile"]),
		StrictHostKeyChecking: stricthostkey.NormalizeStrictHostKeyChecking(cfg.single["stricthostkeychecking"]),
		BatchMode:             strings.EqualFold(cfg.single["batchmode"], "yes"),
		RemoteCommand:         spec.RemoteCommand,
		Alias:                 spec.Alias,
		Destination:           spec.Destination,
		Save:                  spec.Save,
		Tmux:                  spec.Tmux,
		Password:              spec.Password,
	}

	identityFiles := expandIdentityPaths(cfg.multi["identityfile"])
	var idErr error
	r.IdentityPath, r.IdentityExplicit, idErr = resolveIdentityPath(spec, identityFiles)
	if idErr != nil {
		return Resolved{}, idErr
	}

	if err := applyProxyJump(&r, sshPath, spec, cfg); err != nil {
		return Resolved{}, err
	}
	applyCLIOptionOverrides(&r, spec, cfg)
	applyDestinationParse(&r, spec)
	return r, nil
}

func applyProxyJump(r *Resolved, sshPath string, spec openssh.Spec, cfg sshGConfig) error {
	if spec.ProxyJump != nil {
		r.ProxyJump = append([]string(nil), spec.ProxyJump...)
	} else if pj := cfg.single["proxyjump"]; pj != "" && !strings.EqualFold(pj, "none") {
		r.ProxyJump = appendNonEmptyProxyJumps(nil, pj)
	}

	if len(r.ProxyJump) > 0 {
		expanded, err := expandProxyJumpList(sshPath, spec, r.ProxyJump)
		if err != nil {
			return err
		}
		r.ProxyJump = expanded
	}
	if pc := cfg.single["proxycommand"]; pc != "" && !strings.EqualFold(pc, "none") {
		r.ProxyCommand = pc
	}
	return nil
}

func applyCLIOptionOverrides(r *Resolved, spec openssh.Spec, cfg sshGConfig) {
	if spec.Options["port"] != "" {
		if p, err := strconv.Atoi(spec.Options["port"]); err == nil {
			r.Port = p
		}
	}
	if spec.Options["user"] != "" {
		r.User = spec.Options["user"]
	}
	if kh := spec.Options["userknownhostsfile"]; kh != "" {
		r.KnownHostsPath = platform.ExpandSSHPath(kh)
	}
	if spec.Options["stricthostkeychecking"] != "" {
		r.StrictHostKeyChecking = stricthostkey.NormalizeStrictHostKeyChecking(spec.Options["stricthostkeychecking"])
	}
	if to := cfg.single["connecttimeout"]; to != "" && to != "none" && to != "0" {
		if sec, err := strconv.Atoi(to); err == nil && sec > 0 {
			r.ConnectTimeout = time.Duration(sec) * time.Second
		}
	}
}

func applyDestinationParse(r *Resolved, spec openssh.Spec) {
	u, h, p, err := openssh.ParseDestination(spec.Destination)
	if err != nil || h == "" {
		return
	}
	if r.Host == spec.Destination || r.Host == "" {
		r.Host = h
	}
	if u != "" && r.User == "" {
		r.User = u
	}
	if p != 0 && spec.Port == 0 {
		r.Port = p
	}
}

func validateResolved(r Resolved, spec openssh.Spec) (Resolved, error) {
	if r.Host == "" {
		return Resolved{}, fmt.Errorf("could not resolve hostname from %q", spec.Destination)
	}
	if len(spec.RemoteCommand) > 0 && spec.Tmux {
		return Resolved{}, openssh.TmuxWithRemoteCommandErr()
	}
	return r, nil
}

func appendNonEmptyProxyJumps(dst []string, csv string) []string {
	for _, j := range strings.Split(csv, ",") {
		j = strings.TrimSpace(j)
		if j == "" {
			continue
		}
		dst = append(dst, j)
	}
	return dst
}

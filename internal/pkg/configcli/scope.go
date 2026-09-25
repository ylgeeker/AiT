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

// Package configcli builds scoped Cobra commands for aissh config and ait agent config.
package configcli

import (
	"errors"

	pkgconfig "github.com/os-artificer/AiT/internal/pkg/config"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

// ErrAIKeyOnSSHConfig is returned when an AI config key is used on aissh config.
var ErrAIKeyOnSSHConfig = errhint.WithTryMsg(
	"AI config keys are not available via aissh config",
	"Configure AI providers and agent settings with ait agent config.",
	toolname.AgentCLI("config set ai.deepseek.api_key <key>"),
	toolname.AgentCLI("config init"),
)

// ErrValidationFailed indicates config validate saw at least one issue.
var ErrValidationFailed = errors.New("config validation failed")

// Scope identifies which config surface a command tree serves.
type Scope int

const (
	// ScopeSSH is aissh config (session.* only).
	ScopeSSH Scope = iota
	// ScopeAI is ait agent config (ai.* and provider API keys).
	ScopeAI
)

// Options configures a scoped config command factory.
type Options struct {
	Scope         Scope
	HandleAPIKeys bool
	InitShort     string
	ValidateShort string
	RejectKey     func(key string) error
	Ensure        func(dryRun bool) (pkgconfig.EnsureFileResult, error)
	Validate      func(path string) (*pkgconfig.ValidationReport, error)
}

// SSH returns Options for aissh config (session keys; EnsureSSHDefaults / ValidateSSH).
func SSH() Options {
	return Options{
		Scope:         ScopeSSH,
		HandleAPIKeys: false,
		InitShort:     "Create or update config.yaml with defaults for missing keys",
		ValidateShort: "Validate config.yaml syntax and values",
		RejectKey:     RejectAIConfigKey,
		Ensure:        pkgconfig.EnsureSSHDefaults,
		Validate:      pkgconfig.ValidateSSH,
	}
}

// AI returns Options for ait agent config (ai.* keys and provider API keys).
func AI() Options {
	return Options{
		Scope:         ScopeAI,
		HandleAPIKeys: true,
		InitShort:     "Create or update config.yaml with AI defaults for missing keys",
		ValidateShort: "Validate config.yaml syntax and AI values",
		RejectKey:     RejectSSHConfigKey,
		Ensure:        pkgconfig.EnsureAIDefaults,
		Validate:      pkgconfig.ValidateAI,
	}
}

// RejectAIConfigKey returns ErrAIKeyOnSSHConfig when key is an AI-domain key.
// Non-AI keys return nil.
func RejectAIConfigKey(key string) error {
	if pkgconfig.IsAIConfigKey(key) {
		return ErrAIKeyOnSSHConfig
	}
	return nil
}

// RejectSSHConfigKey returns an errhint when key belongs to aissh config (session.*).
// Non-SSH keys return nil. Errors: errhint-wrapped sentinel for cross-domain keys.
func RejectSSHConfigKey(key string) error {
	if !pkgconfig.IsSSHConfigKey(key) {
		return nil
	}
	return errhint.WithTryMsg(
		"key belongs to aissh config",
		"Session and SSH settings are managed with aissh config.",
		toolname.CLI("config get "+key),
		toolname.CLI("config set "+key+" <value>"),
	)
}

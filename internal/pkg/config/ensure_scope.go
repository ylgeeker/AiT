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

package config

import (
	"errors"
	"os"
	"reflect"

	"github.com/os-artificer/AiT/internal/pkg/platform"
)

// EnsureSSHDefaults creates or updates config.yaml with session defaults only.
func EnsureSSHDefaults(dryRun bool) (EnsureFileResult, error) {
	return ensureScopedDefaults(dryRun, mergeSSHDefaults)
}

// EnsureAIDefaults creates or updates config.yaml with AI/provider defaults only.
func EnsureAIDefaults(dryRun bool) (EnsureFileResult, error) {
	path, err := platform.ConfigPath()
	if err != nil {
		return EnsureFileResult{}, err
	}
	tieredInFile := hadTieredInFile(path)
	return ensureScopedDefaults(dryRun, func(cfg *Config) []string {
		return mergeAIDefaults(cfg, tieredInFile || cfg.AI.Tiered.Enabled)
	})
}

type mergeDefaultsFn func(cfg *Config) []string

func ensureScopedDefaults(dryRun bool, merge mergeDefaultsFn) (EnsureFileResult, error) {
	path, err := platform.ConfigPath()
	if err != nil {
		return EnsureFileResult{}, err
	}

	res := EnsureFileResult{Path: path}
	_, statErr := os.Stat(path)
	res.Created = errors.Is(statErr, os.ErrNotExist)
	if statErr != nil && !res.Created {
		return res, statErr
	}

	cfg, err := Load()
	if err != nil {
		return res, err
	}

	before := cloneConfig(cfg)
	added := merge(cfg)
	normalizeAI(cfg)

	if len(added) > 0 {
		res.AddedKeys = added
		res.Updated = !res.Created
	} else if !res.Created && !reflect.DeepEqual(before, cfg) {
		res.Updated = true
	}

	if res.Created || res.Updated {
		if dryRun {
			return res, nil
		}
		if err := Save(cfg); err != nil {
			return res, err
		}
	}
	return res, nil
}

func mergeSSHDefaults(cfg *Config) []string {
	def := defaults()
	var added []string

	if cfg.Session.AuthTTL == "" {
		cfg.Session.AuthTTL = def.Session.AuthTTL
		added = append(added, keySessionAuthTTL)
	}
	if cfg.Session.RemotePrompt.Path == "" {
		cfg.Session.RemotePrompt.Path = def.Session.RemotePrompt.Path
		added = append(added, keySessionRemotePromptPath)
	}

	return added
}

func mergeAIDefaults(cfg *Config, fillTiered bool) []string {
	def := defaults()
	var added []string

	if cfg.AI.Default == "" {
		cfg.AI.Default = def.AI.Default
		added = append(added, keyAIDefault)
	}

	ensureProviders(cfg)
	for _, id := range CloudProviderIDs() {
		added = append(added, mergeProvider(cfg, def, id)...)
	}

	if fillTiered {
		added = append(added, mergeTieredDefaults(cfg)...)
	}

	return added
}

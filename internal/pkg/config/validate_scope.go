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

	"github.com/os-artificer/AiT/internal/pkg/platform"

	"gopkg.in/yaml.v3"
)

func loadConfigAtPath(path string) (*Config, *ValidationReport, error) {
	if path == "" {
		var err error
		path, err = platform.ConfigPath()
		if err != nil {
			return nil, nil, err
		}
	}

	report := &ValidationReport{Path: path}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, nil, configNotFoundErr(path)
	} else if err != nil {
		return nil, nil, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		report.Issues = append(report.Issues, ValidationIssue{
			Field:   "yaml",
			Message: err.Error(),
		})
		return nil, report, nil
	}

	return &cfg, report, nil
}

// ValidateSSH runs semantic validation for SSH/session keys only.
func ValidateSSH(path string) (*ValidationReport, error) {
	cfg, report, err := loadConfigAtPath(path)
	if err != nil || report != nil && len(report.Issues) > 0 {
		return report, err
	}
	report.Issues = validateSession(cfg)
	return report, nil
}

// ValidateAI runs semantic validation for AI/agent keys only.
func ValidateAI(path string) (*ValidationReport, error) {
	cfg, report, err := loadConfigAtPath(path)
	if err != nil || report != nil && len(report.Issues) > 0 {
		return report, err
	}
	var issues []ValidationIssue
	issues = append(issues, validateAI(cfg)...)
	issues = append(issues, legacyDeepseekWarning(cfg)...)
	report.Issues = issues
	return report, nil
}

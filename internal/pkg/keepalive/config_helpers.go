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

package keepalive

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

var unitName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)

// Layout returns the config, admin socket, and lock paths under root.
// root is the AppRoot directory. On Windows the admin pipe is pipeName(root), not the socket path.
func Layout(root string) (cfg, sock, lock string) {
	return filepath.Join(root, ConfigFileName), filepath.Join(root, SockFileName), filepath.Join(root, LockFileName)
}

// ValidateName reports whether name is a single safe path segment.
// Returns ErrBadName when name is empty or contains characters other than letters, digits, '-' and '_'.
func ValidateName(name string) error {
	if !unitName.MatchString(name) {
		return ErrBadName
	}
	return nil
}

// Load reads keepalive.yaml. A missing file returns an empty File and a nil error.
// Returns a YAML decode error when the file exists but is not a document.
func Load(path string) (File, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return File{}, nil
	}
	if err != nil {
		return File{}, err
	}
	var doc File
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return File{}, err
	}
	return doc, nil
}

// Save writes doc to path, creating parent directories with mode 0700.
// Returns an error from mkdir or write. The file mode is 0600.
func Save(path string, doc File) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Upsert adds unit to doc. It returns ErrBadName, ErrNoCommand, or ErrUnitExists.
func Upsert(doc File, unit Unit) (File, error) {
	if err := validateUnit(unit); err != nil {
		return doc, err
	}
	for _, existing := range doc.Units {
		if existing.Name == unit.Name {
			return doc, ErrUnitExists
		}
	}
	doc.Units = append(doc.Units, unit)
	return doc, nil
}

// Replace inserts unit or overwrites the unit with the same name.
// It returns ErrBadName, ErrNoCommand, or ErrBadLimit.
func Replace(doc File, unit Unit) (File, error) {
	if err := validateUnit(unit); err != nil {
		return doc, err
	}
	for i, existing := range doc.Units {
		if existing.Name == unit.Name {
			doc.Units[i] = unit
			return doc, nil
		}
	}
	doc.Units = append(doc.Units, unit)
	return doc, nil
}

// Remove deletes the named unit. It returns ErrUnitMissing when name is absent.
func Remove(doc File, name string) (File, error) {
	out := doc
	out.Units = nil
	found := false
	for _, unit := range doc.Units {
		if unit.Name == name {
			found = true
			continue
		}
		out.Units = append(out.Units, unit)
	}
	if !found {
		return doc, ErrUnitMissing
	}
	return out, nil
}

func validateUnit(unit Unit) error {
	if err := ValidateName(unit.Name); err != nil {
		return err
	}
	if len(unit.Command) == 0 || unit.Command[0] == "" {
		return ErrNoCommand
	}
	if unit.Limits.DiskSpacePercent < 0 || unit.Limits.DiskSpacePercent > 100 {
		return ErrBadLimit
	}
	if unit.Limits.DiskIOPercent < 0 || unit.Limits.DiskIOPercent > 100 {
		return ErrBadLimit
	}
	return nil
}

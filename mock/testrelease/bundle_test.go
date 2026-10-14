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

package testrelease

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/aitplugin"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

func TestWriteBundleTree_IncludesKeep(t *testing.T) {
	dir := t.TempDir()
	if err := WriteBundleTree(dir, "aitssh-bin", "aitlog-bin"); err != nil {
		t.Fatal(err)
	}
	for _, plugin := range aitplugin.BuiltinBinaryNames {
		if _, err := os.Stat(filepath.Join(dir, plugin.Binary)); err != nil {
			t.Fatalf("%s: %v", plugin.Binary, err)
		}
		manifest := filepath.Join(dir, "plugins", plugin.Name, "plugin.json")
		data, err := os.ReadFile(manifest)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), `"binary": "`+plugin.Binary+`"`) {
			t.Fatalf("%s manifest=%s", plugin.Name, data)
		}
	}
	body, err := os.ReadFile(filepath.Join(dir, "plugins", "keep", toolname.AITKeep))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "aitkeep-bin" {
		t.Fatalf("aitkeep=%q", body)
	}
}

func TestPackZipDir_IncludesKeep(t *testing.T) {
	dir := t.TempDir()
	if err := WriteBundleTree(dir, "aitssh-bin", "aitlog-bin"); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(t.TempDir(), "bundle.zip")
	if err := PackZipDir(dir, archive); err != nil {
		t.Fatal(err)
	}
	zr, err := zip.OpenReader(archive)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	var found bool
	for _, f := range zr.File {
		if f.Name == "plugins/keep/aitkeep" || f.Name == "aitkeep" {
			found = true
		}
	}
	if !found {
		t.Fatal("zip missing aitkeep")
	}
}

func TestLoadReleaseArchive_WindowsZip(t *testing.T) {
	name, data, err := loadReleaseArchive(ServeOptions{Tag: "v9.9.9", GOOS: "windows", GOARCH: "amd64"})
	if err != nil {
		t.Fatal(err)
	}
	if name != "ait-v9.9.9-windows-amd64.zip" || len(data) == 0 {
		t.Fatalf("name=%s len=%d", name, len(data))
	}
}

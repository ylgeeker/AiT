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

package aitrelease

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
)

func extractTarGzEntry(hdr *tar.Header, tr *tar.Reader, destDir string, total *int64) error {
	target := filepath.Join(destDir, filepath.FromSlash(hdr.Name))
	if hdr.Typeflag == tar.TypeDir {
		return os.MkdirAll(target, 0o755)
	}
	if hdr.Typeflag != tar.TypeReg && hdr.Typeflag != tar.TypeRegA {
		return nil
	}
	*total += hdr.Size
	if *total > maxExtractBytes {
		return fmt.Errorf("archive exceeds size limit")
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	return writeExtractedFile(target, tr, hdr.FileInfo().Mode())
}

func extractZipEntry(f *zip.File, destDir string) error {
	target := filepath.Join(destDir, filepath.FromSlash(f.Name))
	if f.FileInfo().IsDir() {
		return os.MkdirAll(target, 0o755)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	err = writeExtractedFile(target, rc, f.Mode())
	rc.Close()
	return err
}

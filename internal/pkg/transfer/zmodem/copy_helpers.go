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

package zmodem

import (
	"io"
)

const copyBufSize = 32 * 1024

func copyWithProgress(dst io.Writer, src io.Reader, total int64, progress ProgressFunc) (int64, error) {
	buf := make([]byte, copyBufSize)
	var written int64

	for {
		done, err := copyOneProgressChunk(dst, src, buf, &written, total, progress)
		if err != nil {
			return written, err
		}
		if done {
			return written, nil
		}
	}
}

func copyOneProgressChunk(
	dst io.Writer,
	src io.Reader,
	buf []byte,
	written *int64,
	total int64,
	progress ProgressFunc,
) (done bool, err error) {
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}
	if n > 0 {
		var werr error
		*written, werr = appendCopyProgressChunk(dst, buf[:n], *written, total, progress)
		if werr != nil {
			return false, werr
		}
	}
	if err == io.EOF {
		return true, nil
	}
	return false, nil
}

func appendCopyProgressChunk(
	dst io.Writer,
	chunk []byte,
	written int64,
	total int64,
	progress ProgressFunc,
) (int64, error) {
	if _, err := dst.Write(chunk); err != nil {
		return written, err
	}
	written += int64(len(chunk))
	if progress != nil && total > 0 {
		progress(written, total)
	}
	return written, nil
}

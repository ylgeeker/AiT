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

import "os"

func collectExistingKnownHostsFiles(files []string) (existing []string, err error) {
	for _, fn := range files {
		existing, err = appendKnownHostsFileIfPresent(existing, fn)
		if err != nil {
			return nil, err
		}
	}
	return existing, nil
}

func appendKnownHostsFileIfPresent(existing []string, fn string) ([]string, error) {
	_, statErr := os.Stat(fn)
	if statErr == nil {
		return append(existing, fn), nil
	}
	if os.IsNotExist(statErr) {
		return existing, nil
	}
	return existing, statErr
}

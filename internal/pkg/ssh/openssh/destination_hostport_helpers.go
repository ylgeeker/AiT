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

package openssh

import (
	"strconv"
	"strings"
)

func parseDestinationHostPort(dest string) (hostPart string, port int, err error) {
	hostPart = dest
	if strings.HasPrefix(dest, "[") {
		return parseBracketHost(dest)
	}
	colon := strings.LastIndex(dest, ":")
	if colon <= 0 {
		return hostPart, 0, nil
	}
	p, perr := strconv.Atoi(dest[colon+1:])
	if perr != nil {
		return hostPart, 0, nil
	}
	return dest[:colon], p, nil
}

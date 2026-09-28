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

import "strings"

func applyUnknownOrGluedOFlag(s *Spec, arg string, i int, cmdLabel string) (int, error) {
	if !strings.HasPrefix(arg, "-o") || len(arg) <= 2 {
		return i, unknownFlagErr(arg, cmdLabel)
	}
	if err := mergeOption(s.Options, arg[2:]); err != nil {
		return i, err
	}
	return i, nil
}

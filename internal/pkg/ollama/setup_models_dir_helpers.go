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

package ollama

import (
	"fmt"
	"io"
	"strings"
)

func readCustomModelsDir(in io.Reader, out io.Writer) (string, error) {
	_, _ = fmt.Fprint(out, "Absolute path: ")
	custom, err := readChoiceLine(in)
	if err != nil {
		return "", err
	}
	return expandModelsPath(strings.TrimSpace(custom))
}

func dirFromNumericChoice(choice string, opts []SuggestedModelsDirOption, customN int, in io.Reader, out io.Writer) (string, error) {
	n, err := parseChoiceIndex(choice)
	if err != nil {
		return "", fmt.Errorf("invalid selection %q", choice)
	}
	if n >= 1 && n <= len(opts) {
		return opts[n-1].Path, nil
	}
	if n == customN {
		return readCustomModelsDir(in, out)
	}
	return "", fmt.Errorf("invalid selection %q", choice)
}

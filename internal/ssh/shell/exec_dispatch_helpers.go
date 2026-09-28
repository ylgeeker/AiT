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

package shell

import (
	"github.com/os-artificer/AiT/internal/pkg/cliroute"
	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
)

// finishCLIDispatch runs the aissh Cobra subtree for words. When stop is true, dispatchShellCommand
// must return err immediately (nil means CLI handled the line). When stop is false, try PATH fallback.
func finishCLIDispatch(words []string) (stop bool, err error) {
	err = runCLI(words)
	if err == nil {
		return true, nil
	}
	if !pkghelpdoc.IsUnknownCommand(err) {
		return true, err
	}
	if cliroute.IsRemovedTopLevelWord(words[0]) {
		return true, err
	}
	return false, nil
}

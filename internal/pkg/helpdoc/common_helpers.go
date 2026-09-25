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

package helpdoc

import (
	"fmt"
	"slices"

	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

// LongHint returns the standard one-line Cobra Long field for cmdName.
func LongHint(cmdName string) string {
	root := toolname.CLI("")
	if cmdName == toolname.AITSSH || cmdName == root {
		return fmt.Sprintf(longHintRoot, root)
	}
	if cmdName == toolname.Agent {
		return fmt.Sprintf(longHintRoot, toolname.AgentCLI(""))
	}
	return fmt.Sprintf(longHintSub, toolname.CLI(cmdName))
}

// HelpRequested reports whether argv is only -h, --help, or help.
func HelpRequested(argv []string) bool {
	return len(argv) == 1 && isHelpArg(argv[0])
}

// VersionRequested reports whether argv asks only for the version line.
// argv must be exactly one token: version, -V, or --version.
// Returns true for that single token, and false when argv is empty, longer, or some other command.
func VersionRequested(argv []string) bool {
	return len(argv) == 1 && isVersionArg(argv[0])
}

// TakePlainFlag removes a leading --plain token from argv.
// argv is the argument list after the agent verb.
// Returns plain true and the remaining args when the first token is --plain.
// Returns plain false and argv unchanged when it is missing or not first.
func TakePlainFlag(argv []string) (plain bool, rest []string) {
	if len(argv) > 0 && argv[0] == PlainArg {
		return true, argv[1:]
	}
	return false, argv
}

// isHelpArg returns true when arg equals -h, --help, or help.
func isHelpArg(arg string) bool {
	return arg == HelpArgH || arg == HelpArgHelp || arg == HelpArgCmd
}

// isVersionArg returns true when arg equals version, -V, or --version.
func isVersionArg(arg string) bool {
	return arg == VersionArgShort || arg == VersionArgLong || arg == VersionArgCmd
}

// HasHelpArg reports whether args contain a bare "help" positional argument
// (not a flag value like --load help or --switch help).
func HasHelpArg(args []string) bool {
	return slices.Contains(args, HelpArgCmd)
}

// HelpFooter returns the standard trailing Flags section for cmdName.
func HelpFooter(cmdName string) string {
	return fmt.Sprintf(helpFooterFlagsPrefix, cmdName)
}

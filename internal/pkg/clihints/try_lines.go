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

// Package clihints holds shared errhint Try-block example command lines for the aissh CLI.
package clihints

import "github.com/os-artificer/AiT/internal/pkg/toolname"

// Session and login Try lines (single source for internal/ssh and internal/pkg/ssh/client).

func TrySessionList() string         { return toolname.CLI("session --list") }
func TrySessionLoadAlias() string    { return toolname.CLI("session --load <alias>") }
func TrySessionLoad() string         { return toolname.CLI("session --load") }
func TrySessionLoadOptional() string { return toolname.CLI("session --load [<id|alias>]") }
func TrySessionLoadFullID() string {
	return toolname.CLI("session --list -v") + "    then " + toolname.CLI("session --load <FULL_ID>")
}
func TrySessionLoadPick() string {
	return toolname.CLI("session --load") + "       pick by # interactively"
}
func TrySessionSwitchID() string        { return toolname.CLI("session --switch <id|alias>") }
func TrySessionSwitch() string          { return toolname.CLI("session --switch") }
func TrySessionSwitchForce() string     { return toolname.CLI("session --switch <id|alias> --force") }
func TrySessionHelp() string            { return toolname.CLI("session --help") }
func TryConnectSave() string            { return toolname.CLI("--save <destination>") }
func TryConnectAlias() string           { return toolname.CLI("--alias <name> <destination>") }
func TryConnectHelp() string            { return toolname.CLI("--help") }
func TryConnectTmuxUserHost() string    { return toolname.CLI("--tmux user@host") }
func TryConnectTmuxDevhost() string     { return toolname.CLI("--tmux devhost") }
func TryConnectTmuxSaveDevhost() string { return toolname.CLI("--tmux --save devhost") }

// OpenSSH parse / transfer Try lines (concrete examples for get/put/connect argv hints).

func TryConnectExample() string      { return toolname.CLI("user@host") }
func TryConnectSaveUserHost() string { return toolname.CLI("--save user@host") }
func TryConnectRemoteCmd() string    { return toolname.CLI("user@host uptime") }
func TryGetHelp() string             { return toolname.CLI("get --help") }
func TryPutHelp() string             { return toolname.CLI("put --help") }
func TryGetExample() string          { return toolname.CLI("get user@host:/remote/path ./local") }
func TryPutExample() string          { return toolname.CLI("put ./local user@host:/remote/path") }
func TryGetSessionProd() string      { return toolname.CLI("get --session prod <paths...>") }

// Shell REPL builtin shorthand (maps to session --list).

func TryShellList() string { return toolname.CLI("list") }

// AI plugin Try lines (ai binary / env; never aissh config).

func TryAIBinary() string          { return toolname.AICLI("") }
func TryAIExperimentalCLI() string { return toolname.AICLI("--ai-experimental") }
func TryAIExperimentalEnv() string { return "AIT_AI_EXPERIMENTAL=1 " + toolname.AICLI("") }
func TryAitAgent() string          { return toolname.AITCLI("agent") }

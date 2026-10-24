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

package session

import "github.com/os-artificer/AiT/internal/pkg/clihints"

// Re-export Try-block lines from clihints (avoids import cycle with ssh/client).

func TrySessionList() string         { return clihints.TrySessionList() }
func TrySessionLoadAlias() string    { return clihints.TrySessionLoadAlias() }
func TrySessionLoad() string         { return clihints.TrySessionLoad() }
func TrySessionLoadOptional() string { return clihints.TrySessionLoadOptional() }
func TrySessionLoadFullID() string   { return clihints.TrySessionLoadFullID() }
func TrySessionLoadPick() string     { return clihints.TrySessionLoadPick() }
func TrySessionSwitchID() string     { return clihints.TrySessionSwitchID() }
func TrySessionSwitchAlias() string  { return clihints.TrySessionSwitchAlias() }
func TrySessionSwitch() string       { return clihints.TrySessionSwitch() }
func TrySessionSwitchFullID() string {
	return clihints.TrySessionSwitchFullID()
}
func TrySessionSwitchPick() string      { return clihints.TrySessionSwitchPick() }
func TrySessionSwitchForce() string     { return clihints.TrySessionSwitchForce() }
func TrySessionHelp() string            { return clihints.TrySessionHelp() }
func TryConnectSave() string            { return clihints.TryConnectSave() }
func TryConnectAlias() string           { return clihints.TryConnectAlias() }
func TryConnectHelp() string            { return clihints.TryConnectHelp() }
func TryConnectTmuxUserHost() string    { return clihints.TryConnectTmuxUserHost() }
func TryConnectTmuxDevhost() string     { return clihints.TryConnectTmuxDevhost() }
func TryConnectTmuxSaveDevhost() string { return clihints.TryConnectTmuxSaveDevhost() }

// Private aliases for in-package hint helpers.

func tryLoadByAlias() string     { return clihints.TrySessionLoadAlias() }
func tryLoadByFullID() string    { return clihints.TrySessionLoadFullID() }
func tryLoadInteractive() string { return clihints.TrySessionLoadPick() }
func trySwitchByAlias() string   { return clihints.TrySessionSwitchAlias() }
func trySwitchByFullID() string  { return clihints.TrySessionSwitchFullID() }
func trySwitchInteractive() string {
	return clihints.TrySessionSwitchPick()
}
func trySessionHelp() string   { return clihints.TrySessionHelp() }
func tryLoadFromLocal() string { return clihints.TrySessionLoadOptional() }
func trySwitchForce() string   { return clihints.TrySessionSwitchForce() }
func tryConnectSave() string   { return clihints.TryConnectSave() }
func tryConnectAlias() string  { return clihints.TryConnectAlias() }
func tryConnectHelp() string   { return clihints.TryConnectHelp() }

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

func stepParseLogin(s *Spec, argv []string, i int, remoteCommandStarted bool) (next int, newRemoteStarted bool, done bool, err error) {
	arg := argv[i]
	if arg == "--" {
		s.RemoteCommand = append(s.RemoteCommand, argv[i+1:]...)
		return i, remoteCommandStarted, true, nil
	}
	if strings.HasPrefix(arg, "-") {
		next, err = consumeLoginDashArg(s, argv, i, arg, remoteCommandStarted)
		if err != nil {
			return i, remoteCommandStarted, false, err
		}
		return next, remoteCommandStarted, false, nil
	}
	if s.Destination != "" {
		s.RemoteCommand = append(s.RemoteCommand, arg)
		return i, true, false, nil
	}
	s.Destination = arg
	return i, remoteCommandStarted, false, nil
}

func consumeLoginDashArg(s *Spec, argv []string, i int, arg string, remoteCommandStarted bool) (next int, err error) {
	if remoteCommandStarted {
		s.RemoteCommand = append(s.RemoteCommand, arg)
		return i, nil
	}
	next, err = applyFlag(s, argv, i, connectArgvLabel)
	if err != nil {
		return i, err
	}
	return next, nil
}

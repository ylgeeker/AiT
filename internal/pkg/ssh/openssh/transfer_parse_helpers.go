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
	"fmt"
	"strings"
)

func stepParseTransferArgv(s *Spec, argv []string, i int, cmdLabel, session string, positionals []string) (next int, outSession string, outPos []string, err error) {
	arg := argv[i]
	if arg == "--" {
		return i, session, positionals, fmt.Errorf("-- is not supported for aissh %s", cmdLabel)
	}
	if !strings.HasPrefix(arg, "-") {
		return i, session, append(positionals, arg), nil
	}
	if arg == FlagSession {
		next, outSession, err = applyTransferSessionFlag(argv, i, session)
		if err != nil {
			return i, session, positionals, err
		}
		return next, outSession, positionals, nil
	}
	next, err = applyFlag(s, argv, i, cmdLabel)
	if err != nil {
		return i, session, positionals, err
	}
	return next, session, positionals, nil
}

func consumeSessionFlag(argv []string, i int) (session string, next int, err error) {
	i++
	if i >= len(argv) {
		return "", i, sessionFlagRequiresValueErr()
	}
	return argv[i], i, nil
}

func applyTransferSessionFlag(argv []string, i int, _ string) (next int, session string, err error) {
	session, next, err = consumeSessionFlag(argv, i)
	return next, session, err
}

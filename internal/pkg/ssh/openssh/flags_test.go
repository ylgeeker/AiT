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

import "testing"

func TestFlagConstants_StableLiterals(t *testing.T) {
	cases := []struct {
		got, want string
	}{
		{FlagSession, "--session"},
		{FlagPassword, "--password"},
		{FlagHost, "--host"},
		{FlagUser, "--user"},
		{FlagPort, "--port"},
		{FlagIdentity, "--identity"},
		{FlagAlias, "--alias"},
		{FlagSave, "--save"},
		{FlagTmux, "--tmux"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Fatalf("%q != %q", tc.got, tc.want)
		}
	}
}

func TestApplyFlag_UsesLongFlagConstants(t *testing.T) {
	var s Spec
	s.Options = make(map[string]string)
	s.LegacyPort = -1

	argv := []string{FlagPassword, "secret", FlagHost, "h.example", FlagAlias, "prod"}
	for i := 0; i < len(argv); i++ {
		next, err := applyFlag(&s, argv, i, "")
		if err != nil {
			t.Fatalf("argv[%d]=%q: %v", i, argv[i], err)
		}
		i = next
	}
	if s.Password != "secret" || s.LegacyHost != "h.example" || s.Alias != "prod" {
		t.Fatalf("spec=%+v", s)
	}
}

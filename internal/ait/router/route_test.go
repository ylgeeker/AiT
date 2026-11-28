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

package router_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/ait/router"
	"github.com/os-artificer/AiT/internal/pkg/clihints"
	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

// TestRoute_EmptyArgsRunsShell (AT1) empty argv opens the toolbox shell REPL.
func TestRoute_EmptyArgsRunsShell(t *testing.T) {
	called := false
	err := router.Route(context.Background(), nil, io.Discard, io.Discard, router.RouteDeps{
		RunShellREPL: func(context.Context, io.Reader, io.Writer, io.Writer) error {
			called = true
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Route: %v", err)
	}
	if !called {
		t.Fatal("empty argv must call RunShellREPL")
	}
}

// TestRoute_UnknownFirstToken (AT57) non-plugin first tokens return errhint, not Loop.
func TestRoute_UnknownFirstToken(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "freeform",
			args: []string{"list", "files"},
			want: []string{"unknown command: list", "Try:", clihints.TryAitAgent(), toolname.AITCLI("agent list files"), toolname.AITCLI("shell")},
		},
		{
			name: "plain",
			args: []string{pkghelpdoc.PlainArg},
			want: []string{"unknown command: " + pkghelpdoc.PlainArg, "Try:", clihints.TryAitAgent(), toolname.AITCLI("agent --plain"), toolname.AITCLI("shell")},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := router.Route(context.Background(), tc.args, io.Discard, io.Discard, router.RouteDeps{
				RunShellREPL: func(context.Context, io.Reader, io.Writer, io.Writer) error {
					t.Fatal("unknown first token must not open shell")
					return nil
				},
			})
			if err == nil {
				t.Fatal("expected error")
			}
			s := err.Error()
			for _, sub := range tc.want {
				if !strings.Contains(s, sub) {
					t.Fatalf("err missing %q: %q", sub, s)
				}
			}
		})
	}
}

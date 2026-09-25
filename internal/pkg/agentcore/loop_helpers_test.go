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

package agentcore

import (
	"strings"
	"testing"

	"github.com/os-artificer/AiT/internal/pkg/agentcore/runner"
	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

func TestFormatPlanStep(t *testing.T) {
	t.Parallel()
	got := formatPlanStep(2, 12, Config{})
	if got != "step 2/12" {
		t.Fatalf("got %q", got)
	}
}

func TestDescribeAction_Run(t *testing.T) {
	t.Parallel()
	got := describeAction(runner.Action{Action: runner.ActionRun, Command: "ls -la"})
	if got != "run ls -la" {
		t.Fatalf("got %q", got)
	}
}

func TestDescribeAction_Ait(t *testing.T) {
	t.Parallel()
	got := describeAction(runner.Action{Action: runner.ActionAISsh, Args: "session --list"})
	want := toolname.AITSSH + " session --list"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestDescribeAction_Answer(t *testing.T) {
	t.Parallel()
	if describeAction(runner.Action{Action: runner.ActionAnswer}) != "answer" {
		t.Fatal("expected answer")
	}
}

func TestRunAgentSteps_EmitsThinking(t *testing.T) {
	t.Parallel()
	var kinds []EventKind
	var details []string
	emit := FuncEmitter(func(e Event) {
		kinds = append(kinds, e.Kind)
		details = append(details, e.Detail)
	})
	emit(Event{Kind: EventPlanUpdated, Detail: formatPlanStep(1, 12, Config{})})
	emit(Event{Kind: EventThinking, Detail: "checking network interfaces"})
	emit(Event{Kind: EventPlanUpdated, Detail: "planned: " + describeAction(runner.Action{Action: runner.ActionRun, Command: "cat /proc/net/dev"})})

	if len(kinds) != 3 {
		t.Fatalf("events: %v", kinds)
	}
	if kinds[0] != EventPlanUpdated || kinds[1] != EventThinking || kinds[2] != EventPlanUpdated {
		t.Fatalf("unexpected sequence: %v", kinds)
	}
	if !strings.Contains(details[2], "planned:") {
		t.Fatalf("planned detail missing: %q", details[2])
	}
}

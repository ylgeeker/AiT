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

package config

import (
	"errors"
	"strings"
	"testing"
)

func resetAIExperimentalCLI(t *testing.T) {
	t.Helper()
	SetAIExperimentalCLI(false)
	t.Cleanup(func() { SetAIExperimentalCLI(false) })
}

func TestAIExperimental_DefaultOn(t *testing.T) {
	resetAIExperimentalCLI(t)
	t.Setenv(EnvAIExperimental, "")
	if !AIExperimentalEnabled(nil) {
		t.Fatal("want enabled by default when cfg is nil")
	}
	cfg := defaults()
	if !AIExperimentalEnabled(cfg) {
		t.Fatal("want enabled by default from defaults()")
	}
	if err := RequireAIExperimental(cfg); err != nil {
		t.Fatalf("RequireAIExperimental: %v", err)
	}
}

func TestAIExperimental_ConfigOff(t *testing.T) {
	resetAIExperimentalCLI(t)
	t.Setenv(EnvAIExperimental, "")
	cfg := &Config{}
	cfg.AI.Experimental.Enabled = false
	if AIExperimentalEnabled(cfg) {
		t.Fatal("want config off")
	}
	err := RequireAIExperimental(cfg)
	if !errors.Is(err, ErrAIExperimentalDisabled) {
		t.Fatalf("got %v", err)
	}
	s := err.Error()
	for _, want := range []string{"Try:", keyAIExperimentalEnabled, EnvAIExperimental} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in %q", want, s)
		}
	}
	if strings.Contains(s, "aitssh") {
		t.Errorf("experimental gate hint must not mention aissh: %q", s)
	}
}

func TestAIExperimentalDisabledErr_SSOT(t *testing.T) {
	err := AIExperimentalDisabledErr()
	if !errors.Is(err, ErrAIExperimentalDisabled) {
		t.Fatalf("got %v", err)
	}
	s := err.Error()
	if !strings.Contains(s, "Try:") {
		t.Fatalf("missing Try: %q", s)
	}
	if strings.Contains(s, "aitssh") {
		t.Fatalf("disabled err must not mention aissh: %q", s)
	}
}

func TestAIExperimental_EnvOn(t *testing.T) {
	resetAIExperimentalCLI(t)
	t.Setenv(EnvAIExperimental, "1")
	if !AIExperimentalEnabled(&Config{}) {
		t.Fatal("want env on")
	}
	if err := RequireAIExperimental(&Config{}); err != nil {
		t.Fatal(err)
	}
}

func TestAIExperimental_EnvOff(t *testing.T) {
	resetAIExperimentalCLI(t)
	t.Setenv(EnvAIExperimental, "0")
	cfg := &Config{}
	cfg.AI.Experimental.Enabled = true
	if AIExperimentalEnabled(cfg) {
		t.Fatal("want env off")
	}
	if err := RequireAIExperimental(cfg); err == nil {
		t.Fatal("expected disabled")
	}
}

func TestAIExperimental_CLIFlag(t *testing.T) {
	resetAIExperimentalCLI(t)
	t.Setenv(EnvAIExperimental, "")
	SetAIExperimentalCLI(true)
	if !AIExperimentalEnabled(&Config{}) {
		t.Fatal("want CLI on")
	}
}

func TestStripAIExperimentalFlag(t *testing.T) {
	rest, saw := StripAIExperimentalFlag([]string{FlagAIExperimental, "ai", "ask", "hi"})
	if !saw {
		t.Fatal("want saw")
	}
	if len(rest) != 3 || rest[0] != "ai" || rest[1] != "ask" || rest[2] != "hi" {
		t.Fatalf("rest=%v", rest)
	}

	rest, saw = StripAIExperimentalFlag([]string{"ai", FlagAIExperimental, "ask"})
	if !saw || len(rest) != 2 || rest[0] != "ai" || rest[1] != "ask" {
		t.Fatalf("rest=%v saw=%v", rest, saw)
	}

	rest, saw = StripAIExperimentalFlag([]string{"session", "--list"})
	if saw || len(rest) != 2 {
		t.Fatalf("rest=%v saw=%v", rest, saw)
	}
}

func TestSetGetKey_AIExperimental(t *testing.T) {
	cfg := &Config{}
	if err := cfg.SetKey(keyAIExperimentalEnabled, "true"); err != nil {
		t.Fatal(err)
	}
	if !cfg.AI.Experimental.Enabled {
		t.Fatal("want enabled")
	}
	got, err := cfg.GetKey(keyAIExperimentalEnabled)
	if err != nil || got != "true" {
		t.Fatalf("GetKey=%q err=%v", got, err)
	}
	if err := cfg.SetKey(keyAIExperimentalEnabled, "false"); err != nil {
		t.Fatal(err)
	}
	got, err = cfg.GetKey(keyAIExperimentalEnabled)
	if err != nil || got != "false" {
		t.Fatalf("GetKey=%q err=%v", got, err)
	}
}

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

package toolname

import "testing"

func TestMetaPrefix_Golden(t *testing.T) {
	if MetaPrefix != "!ait" {
		t.Fatalf("MetaPrefix=%q want !ait", MetaPrefix)
	}
}

func TestMetaLine(t *testing.T) {
	if MetaLine("") != "!ait" {
		t.Fatalf("bare meta=%q", MetaLine(""))
	}
	if MetaLine("get /r ./l") != "!ait get /r ./l" {
		t.Fatalf("meta line=%q", MetaLine("get /r ./l"))
	}
}

func TestCLI(t *testing.T) {
	if CLI("") != "aitssh" {
		t.Fatalf("CLI()=%q", CLI(""))
	}
	if CLI("session --list") != "aitssh session --list" {
		t.Fatalf("CLI(sub)=%q", CLI("session --list"))
	}
}

func TestPluginCLI_NoEnvUsesBinary(t *testing.T) {
	t.Setenv(EnvPluginCMD, "")
	if got := PluginCLI(AITLog, Log, ""); got != AITLog {
		t.Fatalf("PluginCLI()=%q want %q", got, AITLog)
	}
	if got := PluginCLI(AITSSH, SSH, "session --list"); got != "aitssh session --list" {
		t.Fatalf("PluginCLI(sub)=%q", got)
	}
}

func TestPluginCLI_EnvUsesCommand(t *testing.T) {
	t.Setenv(EnvPluginCMD, Log)
	if got := LogCLI(""); got != Log {
		t.Fatalf("LogCLI()=%q want %q", got, Log)
	}
	if got := LogCLI("search error"); got != "log search error" {
		t.Fatalf("LogCLI(sub)=%q", got)
	}

	t.Setenv(EnvPluginCMD, SSH)
	if got := CLI("session --list"); got != "ssh session --list" {
		t.Fatalf("CLI(sub)=%q", got)
	}

	t.Setenv(EnvPluginCMD, Host)
	if got := HostCLI("snapshot"); got != "host snapshot" {
		t.Fatalf("HostCLI(sub)=%q", got)
	}

	t.Setenv(EnvPluginCMD, Img)
	if got := ImgCLI("extract"); got != "img extract" {
		t.Fatalf("ImgCLI(sub)=%q", got)
	}
}

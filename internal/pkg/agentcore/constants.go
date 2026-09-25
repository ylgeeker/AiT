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
	"time"

	"github.com/os-artificer/AiT/internal/pkg/errhint"
)

const (
	defaultMaxSteps = 12

	// defaultCommandTimeout bounds one tool execution when Config leaves
	// CommandTimeout unset.
	defaultCommandTimeout = 2 * time.Minute

	msgAgentStepLimitFmt = "Step limit reached (%d). Continue? [y/N] or type a new goal.\r\n"
	msgSkipped           = "⊘ skipped\r\n"
	msgInteractiveSkip   = "⊘ interactive command — run manually\r\n"
	msgConversationClear = "Conversation cleared.\r\n"
	msgCommandTimeoutFmt = "timed out after %s (agent command timeout)"
	// msgContinueGoal is the goal replayed when the user accepts the step-limit prompt.
	msgContinueGoal = "continue the previous task"

	msgWelcomeBanner    = "  plan · search · run anything\r\n\r\n"
	msgWelcomeRemoteFmt = "  plan · search · run anything  (%s)\r\n\r\n"

	promptAgent         = "> "
	promptAgentHostFmt  = "> %s "
	promptAgentContinue = "> continue? "

	metaExit     = "exit"
	metaStop     = "stop"
	metaBack     = "back"
	metaModel    = "model"
	metaContinue = "continue"
	metaClear    = "clear"
	metaHistory  = "history"
)

// ErrSessionEnded is returned when the agent loop exits normally.
var ErrSessionEnded = errhint.NewSentinel("agent session ended")

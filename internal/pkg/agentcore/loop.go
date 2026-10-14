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
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/os-artificer/AiT/internal/pkg/agentcore/approve"
	"github.com/os-artificer/AiT/internal/pkg/agentcore/runner"
	"github.com/os-artificer/AiT/internal/pkg/ai"
	pkgagent "github.com/os-artificer/AiT/internal/pkg/ai/agent"
	"github.com/os-artificer/AiT/internal/pkg/ai/client"
	"github.com/os-artificer/AiT/internal/pkg/terminal"

	"golang.org/x/term"
)

// RunLoop runs the agent task loop until exit/stop/back or ctx cancel.
// initialGoal, if non-empty, is processed as the first user turn.
func RunLoop(ctx context.Context, cfg Config, initialGoal string) error {
	if cfg.AI == nil {
		return ai.NotConfiguredErr()
	}

	cfg = fillDefaults(cfg)
	st := newSessionState(cfg)
	emit(cfg, Event{Kind: EventSessionStarted})
	printSessionWelcome(cfg)

	if g := strings.TrimSpace(initialGoal); g != "" {
		emit(cfg, Event{Kind: EventGoalReceived, Detail: g})
		if err := absorbTurnError(cfg, runAgentTurn(ctx, cfg, st, g)); err != nil {
			return err
		}
	}

	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) && cfg.In == os.Stdin && cfg.ReadLine == nil {
		return runLoopTTY(ctx, cfg, st, fd)
	}
	return runLoopPlain(ctx, cfg, st)
}

func fillDefaults(cfg Config) Config {
	if cfg.In == nil {
		cfg.In = os.Stdin
	}
	if cfg.Out == nil {
		cfg.Out = os.Stdout
	}
	if cfg.ErrOut == nil {
		cfg.ErrOut = os.Stderr
	}
	if cfg.Emitter == nil {
		cfg.Emitter = NopEmitter
	}
	if cfg.Local == nil && cfg.Scope == ai.ScopeLocal {
		cfg.Local = defaultLocalRunner{}
	}
	if cfg.Capture == nil && cfg.Scope == ai.ScopeLocal {
		cfg.Capture = runner.DefaultLocalCaptureRunner()
	}
	return cfg
}

// runAgentTurn drives one user goal to completion, continuing with a follow-up
// goal when the user asks to proceed past the step limit.
// Iterating here instead of recursing from the step-limit prompt keeps the stack
// flat no matter how many times the user continues.
func runAgentTurn(ctx context.Context, cfg Config, st *sessionState, question string) error {
	for question != "" {
		next, err := runAgentSteps(ctx, cfg, st, question)
		if err != nil {
			return err
		}
		question = next
	}
	return nil
}

// runAgentSteps runs up to MaxAgentSteps model and tool steps for question.
// Returns the next goal when the user chose to continue past the step limit, or
// "" when the goal finished or the session ended.
func runAgentSteps(ctx context.Context, cfg Config, st *sessionState, question string) (string, error) {
	if cfg.AI == nil {
		return "", ai.NotConfiguredErr()
	}

	env := assistEnv(cfg, st)
	userMsg := question
	maxSteps := cfg.MaxAgentSteps()
	parseState := newAgentParseRetryState(st.history)

	for step := range maxSteps {
		emit(cfg, Event{Kind: EventPlanUpdated, Detail: formatPlanStep(step+1, maxSteps, cfg)})
		reply, thinking, updated, chatErr := cfg.AI.ChatAgent(ctx, st.history, userMsg, env)
		if chatErr != nil {
			return "", chatErr
		}
		if strings.TrimSpace(thinking) != "" {
			emit(cfg, Event{Kind: EventThinking, Detail: thinking, Provider: eventProvider(cfg)})
		}
		if reply != "" && cfg.StreamDelta != nil {
			cfg.StreamDelta(reply)
		}
		if reply != "" {
			emit(cfg, Event{Kind: EventTokenDelta, Detail: reply})
		}

		parseResult := pkgagent.ResolveAgentReply(reply)
		var retry bool
		userMsg, retry, chatErr = retryAgentParseOrContinue(st, &parseState, userMsg, reply, parseResult)
		if chatErr != nil {
			return "", chatErr
		}
		if retry {
			continue
		}

		canonical, err := parseResult.Envelope.CanonicalJSON()
		if err != nil {
			st.history = parseState.turnStartHistory
			return "", formatMalformedReplyErr()
		}
		st.history = client.ReplaceLastAssistant(updated, canonical)

		act := envelopeToAction(parseResult.Envelope)
		if !act.IsAnswer() {
			emit(cfg, Event{Kind: EventPlanUpdated, Detail: "planned: " + describeAction(act)})
		}
		toolResult, done, err := executeAgentAction(ctx, cfg, st, act)
		if err != nil {
			return "", err
		}
		if done {
			return "", nil
		}

		userMsg = pkgagent.ToolResultPrefix + toolResult
	}

	emit(cfg, Event{Kind: EventStepLimitReached, Detail: fmt.Sprintf("%d", maxSteps)})
	if !managedUI(cfg) {
		agentWrite(cfg.Out, fmt.Sprintf(msgAgentStepLimitFmt, maxSteps))
	}
	return waitContinueOrExit(ctx, cfg, st)
}

// waitContinueOrExit prompts after the step limit and returns the goal to run
// next, or "" when the user ended the session or input was exhausted.
// Input goes through st.ReadLine so a raw terminal uses the line editor instead
// of a newline-terminated Scanner that would never see the carriage return.
func waitContinueOrExit(ctx context.Context, cfg Config, st *sessionState) (string, error) {
	for {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		line, err := st.ReadLine(promptAgentContinue)
		if err == io.EOF {
			emit(cfg, Event{Kind: EventSessionEnded})
			return "", nil
		}
		if err != nil {
			return "", err
		}

		goal, decided := parseContinueChoice(line)
		if !decided {
			continue
		}
		if goal == "" {
			emit(cfg, Event{Kind: EventSessionEnded})
			return "", nil
		}
		if goal != msgContinueGoal {
			emit(cfg, Event{Kind: EventGoalReceived, Detail: goal})
		}
		return goal, nil
	}
}

func runLoopPlain(ctx context.Context, cfg Config, st *sessionState) error {
	hist := terminal.NewReplLineHistory(cfg.ErrOut)

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		raw, err := st.ReadLine(promptAgent)
		if err == io.EOF {
			emit(cfg, Event{Kind: EventSessionEnded})
			return nil
		}
		if err != nil {
			return err
		}

		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		line = normalizeLoopLine(line)

		exit, err := handleLoopMeta(ctx, cfg, st, line, hist)
		if exit {
			emit(cfg, Event{Kind: EventSessionEnded})
			return err
		}
		if err != nil {
			agentPrintErr(cfg.ErrOut, err.Error())
			continue
		}
		if err := absorbTurnError(cfg, runAgentTurn(ctx, cfg, st, line)); err != nil {
			return err
		}
	}
}

func runLoopTTY(ctx context.Context, cfg Config, st *sessionState, fd int) error {
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return runLoopPlain(ctx, cfg, st)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	br := bufio.NewReader(cfg.In)
	hist := terminal.NewReplLineHistory(cfg.ErrOut)
	if cfg.ConfirmRead == nil {
		cfg = cfg.withTTYConfirm(br)
	}
	cfg.SuspendTTY = suspendRawTTY(fd, oldState)
	st.readLine = ttyLineReader(br, cfg.Out)

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		scope, host, _, _, _ := resolveDynamic(cfg)
		prompt := loopPrompt(scope, host)
		// ReadLine only redraws Prompt after an edit, so the first paint is ours.
		agentWrite(cfg.Out, prompt)
		line, err := terminal.ReadLine(br, cfg.Out, terminal.LineEditOptions{
			Prompt:        prompt,
			History:       hist,
			RecordHistory: true,
		})
		if err == io.EOF {
			emit(cfg, Event{Kind: EventSessionEnded})
			return nil
		}
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = normalizeLoopLine(line)

		exit, err := handleLoopMeta(ctx, cfg, st, line, hist)
		if exit {
			emit(cfg, Event{Kind: EventSessionEnded})
			return err
		}
		if err != nil {
			agentPrintErr(cfg.ErrOut, err.Error())
			continue
		}
		if err := absorbTurnError(cfg, runAgentTurn(ctx, cfg, st, line)); err != nil {
			return err
		}
	}
}

func handleLoopMeta(ctx context.Context, cfg Config, st *sessionState, line string, hist *terminal.LineHistory) (exit bool, err error) {
	_ = ctx
	word, rest := firstWord(line)
	switch word {
	case metaExit, metaStop, metaQuit:
		st.Trust.Clear()
		return true, nil
	case metaBack:
		if cfg.Nested {
			st.Trust.Clear()
			return true, nil
		}
		return false, nil
	case metaHelp:
		writeAgentLoopHelp(cfg.Out)
		return false, nil
	case metaClear:
		st.history = nil
		agentWrite(cfg.Out, msgConversationClear)
		return false, nil
	case metaHistory:
		parts := strings.Fields(line)
		var args []string
		if len(parts) > 1 {
			args = parts[1:]
		}
		return false, terminal.RunHistoryCommand(terminal.NewCRLFWriter(cfg.Out), hist, args)
	case metaModel:
		return false, showModelStatus(cfg, rest)
	default:
		return false, nil
	}
}

func showModelStatus(cfg Config, rest string) error {
	rest = strings.TrimSpace(rest)
	if rest == "list" {
		if cfg.AI == nil {
			return ai.NotConfiguredErr()
		}
		lines, err := cfg.AI.ListModels(context.Background())
		if err != nil {
			return err
		}
		if lines != nil {
			_, _ = fmt.Fprintln(cfg.Out, lines.Format())
		}
		return nil
	}
	if cfg.AI != nil {
		_, _ = fmt.Fprintln(cfg.Out, cfg.AI.StatusLine())
	}
	return nil
}

func firstWord(line string) (word, rest string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", ""
	}
	i := strings.IndexByte(line, ' ')
	if i < 0 {
		return line, ""
	}
	return line[:i], strings.TrimSpace(line[i+1:])
}

func (cfg Config) withTTYConfirm(br *bufio.Reader) Config {
	if cfg.ConfirmRead != nil {
		return cfg
	}
	out := cfg.Out
	if out == nil {
		out = os.Stdout
	}
	cfg.ConfirmRead = func(prompt string) (approve.Result, error) {
		agentWrite(out, prompt)
		line, err := terminal.ReadTTYLineNoHistory(br, out, "")
		if err != nil {
			return approve.Result{}, err
		}
		allowPattern := !strings.Contains(prompt, promptHighRiskMarker)
		return approve.ParseChoice(line, allowPattern, allowPattern), nil
	}
	return cfg
}

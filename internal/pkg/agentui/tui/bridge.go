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

package tui

import (
	"io"
	"strings"
	"sync"

	"github.com/os-artificer/AiT/internal/pkg/agentcore"
	"github.com/os-artificer/AiT/internal/pkg/agentcore/approve"
)

type bridgeKind int

const (
	bridgeTranscript bridgeKind = iota
	bridgeEvent
	bridgeDone
)

type bridgeMsg struct {
	kind  bridgeKind
	text  string
	event agentcore.Event
}

type lineBridge struct {
	mu sync.Mutex

	submit chan string
	msgs   chan bridgeMsg
	closed chan struct{}
}

func newLineBridge() *lineBridge {
	return &lineBridge{
		submit: make(chan string, 1),
		msgs:   make(chan bridgeMsg, 64),
		closed: make(chan struct{}),
	}
}

func (b *lineBridge) finish() {
	b.push(bridgeMsg{kind: bridgeDone})
	b.close()
}

func (b *lineBridge) close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	select {
	case <-b.closed:
	default:
		close(b.closed)
	}
}

func (b *lineBridge) push(msg bridgeMsg) {
	select {
	case <-b.closed:
	case b.msgs <- msg:
	default:
		// Drop oldest-style: non-blocking best effort when UI is slow.
		select {
		case <-b.msgs:
		default:
		}
		select {
		case b.msgs <- msg:
		default:
		}
	}
}

func (b *lineBridge) readLine(prompt string) (string, error) {
	p := strings.TrimSpace(prompt)
	if p != "" && p != ">" && !strings.HasPrefix(p, "> continue") {
		b.push(bridgeMsg{kind: bridgeTranscript, text: p})
	}
	return b.readSubmit()
}

func (b *lineBridge) readSubmit() (string, error) {
	select {
	case line, ok := <-b.submit:
		if !ok {
			return "", io.EOF
		}
		return line, nil
	case <-b.closed:
		return "", io.EOF
	}
}

func (b *lineBridge) confirmRead(prompt string) (approve.Result, error) {
	b.push(bridgeMsg{kind: bridgeTranscript, text: prompt})
	line, err := b.readSubmit()
	if err != nil {
		return approve.Result{}, err
	}
	allowPattern := !strings.Contains(prompt, "High-risk")
	return approve.ParseChoice(line, allowPattern, allowPattern), nil
}

func (b *lineBridge) transcriptWriter() io.Writer {
	return bridgeWriter{b: b}
}

type bridgeWriter struct {
	b *lineBridge
}

func (w bridgeWriter) Write(p []byte) (int, error) {
	text := strings.TrimSpace(string(p))
	if text == "" {
		return len(p), nil
	}
	w.b.push(bridgeMsg{kind: bridgeTranscript, text: text})
	return len(p), nil
}

func (b *lineBridge) emit(e agentcore.Event) {
	b.push(bridgeMsg{kind: bridgeEvent, event: e})
}

func (b *lineBridge) sendLine(line string) {
	select {
	case <-b.closed:
	case b.submit <- line:
	}
}

//go:build windows

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

package keepalive

import (
	"errors"
	"net"
	"sync"

	"golang.org/x/sys/windows"
)

type pipeListener struct {
	name   string
	mu     sync.Mutex
	closed bool
}

// Accept waits for one admin client on the named pipe.
// Returns net.ErrClosed after Close. CreateNamedPipe and ConnectNamedPipe errors are returned as-is.
func (l *pipeListener) Accept() (net.Conn, error) {
	l.mu.Lock()
	closed := l.closed
	l.mu.Unlock()
	if closed {
		return nil, net.ErrClosed
	}
	path, err := windows.UTF16PtrFromString(l.name)
	if err != nil {
		return nil, err
	}
	h, err := windows.CreateNamedPipe(
		path,
		windows.PIPE_ACCESS_DUPLEX,
		windows.PIPE_TYPE_BYTE|windows.PIPE_READMODE_BYTE|windows.PIPE_WAIT,
		windows.PIPE_UNLIMITED_INSTANCES,
		pipeBufSize,
		pipeBufSize,
		0,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = windows.ConnectNamedPipe(h, nil)
	if err != nil && !errors.Is(err, windows.ERROR_PIPE_CONNECTED) {
		windows.CloseHandle(h)
		return nil, err
	}
	return netPipe(h), nil
}

// Close marks the listener closed so a later Accept returns net.ErrClosed.
// Handles already returned by Accept stay open.
func (l *pipeListener) Close() error {
	l.mu.Lock()
	l.closed = true
	l.mu.Unlock()
	return nil
}

// Addr returns the pipe name as a net.Addr.
func (l *pipeListener) Addr() net.Addr { return pipeAddr{name: l.name} }

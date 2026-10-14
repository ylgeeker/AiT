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
	"net"
	"os"

	"golang.org/x/sys/windows"
)

func listenAdmin(root string) (net.Listener, error) {
	return &pipeListener{name: pipeName(root)}, nil
}

func dialAdmin(root string) (net.Conn, error) {
	path, err := windows.UTF16PtrFromString(pipeName(root))
	if err != nil {
		return nil, err
	}
	h, err := windows.CreateFile(path, windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil, windows.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, err
	}
	return netPipe(h), nil
}

func netPipe(h windows.Handle) net.Conn {
	return &fileConn{File: os.NewFile(uintptr(h), "keepalive")}
}

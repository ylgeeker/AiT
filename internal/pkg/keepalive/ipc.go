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
	"encoding/json"
	"net"
)

// Request is one admin command sent over the local socket or named pipe.
type Request struct {
	Op   string `json:"op"`
	Name string `json:"name,omitempty"`
}

// Response is the admin reply. Error is empty when OK is true.
type Response struct {
	OK     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Status string `json:"status,omitempty"`
}

// Call dials the admin endpoint under root and writes one JSON request.
// root is AppRoot. Returns ErrDaemonDown when nothing is listening.
func Call(root string, req Request) (Response, error) {
	conn, err := dialAdmin(root)
	if err != nil {
		return Response{}, ErrDaemonDown
	}
	defer conn.Close()
	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return Response{}, err
	}
	var resp Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return Response{}, err
	}
	return resp, nil
}

func serveConn(conn net.Conn, handle func(Request) Response) {
	defer conn.Close()
	var req Request
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		return
	}
	_ = json.NewEncoder(conn).Encode(handle(req))
}

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

const pipeNetwork = "pipe"

type pipeAddr struct{ name string }

// Network returns the address network name for a named pipe.
func (pipeAddr) Network() string { return pipeNetwork }

// String returns the pipe path. An empty path means the address was not set.
func (a pipeAddr) String() string { return a.name }

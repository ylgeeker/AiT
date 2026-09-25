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

package shared

// KeysHelp is the session configuration keys section for aissh config subcommand help.
// AI keys are managed via ait agent config (see toolname.AITCLI("agent config")).
const KeysHelp = `  session.auth_ttl              duration for session unlock (e.g. 168h, 7d)
  session.remote_prompt.path    short (last path segment) or absolute (full path, ~ for home)
`

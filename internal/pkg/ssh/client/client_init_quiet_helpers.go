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

package client

import "time"

func quietAfterNoPostInject(last int64, injectAt time.Time, quiet time.Duration) bool {
	if last != 0 && !remoteOutAt(last).Before(injectAt) {
		return false
	}
	return time.Since(injectAt) >= quiet
}

func quietSinceNoOutput(last int64, since time.Time, quiet time.Duration) bool {
	if last != 0 && !remoteOutAt(last).Before(since) {
		return false
	}
	return time.Since(since) >= quiet
}

func shellReadyAfterPromptQuiet(last int64) bool {
	return last != 0 && time.Since(remoteOutAt(last)) >= remoteOutputQuietPeriod
}

func shellReadyBashFallback(last int64, attachAt time.Time) bool {
	if last != 0 && !remoteOutAt(last).Before(attachAt) {
		return time.Since(remoteOutAt(last)) >= remoteOutputQuietPeriod
	}
	return time.Since(attachAt) >= remoteOutputQuietPeriod
}

func scrollLooksShellReady(scroll, user, host string) bool {
	return containsRemotePromptLine(scroll, user, host) || lastNativePromptLine(scroll) != ""
}

func (c *Client) shellReadyBeforeInjectTick(user, host string, attachAt time.Time) bool {
	scroll := c.scroll.String()
	last := c.lastRemoteOut.Load()
	if scrollLooksShellReady(scroll, user, host) {
		return shellReadyAfterPromptQuiet(last)
	}
	return shellReadyBashFallback(last, attachAt)
}

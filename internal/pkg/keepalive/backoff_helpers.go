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

import "time"

// NextBackoff doubles cur and clamps it to BackoffMin and BackoffMax.
// cur below BackoffMin returns BackoffMin. cur at or above BackoffMax stays at BackoffMax.
func NextBackoff(cur time.Duration) time.Duration {
	if cur < BackoffMin {
		return BackoffMin
	}
	next := cur * 2
	if next > BackoffMax || next < cur {
		return BackoffMax
	}
	return next
}

// AllowRestart reports whether another abnormal restart is inside the budget.
// count is how many abnormal restarts have already happened in the current window.
// The next restart is allowed while count is below RestartBudget.
func AllowRestart(count int) bool {
	return count < RestartBudget
}

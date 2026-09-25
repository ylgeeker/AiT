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
	"context"
	"errors"

	"github.com/os-artificer/AiT/internal/pkg/ai/provider"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
)

// turnErrorEndsSession reports whether err should terminate an interactive RunLoop.
// Malformed model replies and most provider/chat failures stay in-session so the
// user can retry; cancel, NotConfigured, and SessionEnded still exit.
func turnErrorEndsSession(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	if errors.Is(err, ErrSessionEnded) {
		return true
	}
	if errors.Is(err, provider.ErrNotConfigured) {
		return true
	}
	return false
}

// absorbTurnError surfaces a non-fatal turn failure to the UI and returns nil so
// the interactive loop can prompt again. Fatal errors are returned unchanged.
func absorbTurnError(cfg Config, err error) error {
	if err == nil {
		return nil
	}
	if turnErrorEndsSession(err) {
		return err
	}

	emit(cfg, Event{Kind: EventWarning, Detail: err.Error()})
	if !managedUI(cfg) {
		errhint.Print(cfg.ErrOut, errhint.SurfaceTTY, err)
	}
	return nil
}

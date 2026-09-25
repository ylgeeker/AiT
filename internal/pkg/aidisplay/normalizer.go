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

package aidisplay

import "sync"

// Normalizer optionally cleans provider-specific noise before Markdown rendering.
type Normalizer interface {
	Normalize(text string) string
}

type identityNormalizer struct{}

func (identityNormalizer) Normalize(text string) string {
	return text
}

var (
	normalizerMu sync.RWMutex
	normalizers  = map[string]Normalizer{}
)

// RegisterNormalizer registers a display normalizer for providerID (e.g. minimax).
// Duplicate registration panics.
func RegisterNormalizer(providerID string, n Normalizer) {
	if providerID == "" || n == nil {
		panic("aidisplay: invalid normalizer registration")
	}
	normalizerMu.Lock()
	defer normalizerMu.Unlock()
	if _, ok := normalizers[providerID]; ok {
		panic("aidisplay: duplicate normalizer for " + providerID)
	}
	normalizers[providerID] = n
}

func normalizerFor(providerID string) Normalizer {
	normalizerMu.RLock()
	defer normalizerMu.RUnlock()
	if n, ok := normalizers[providerID]; ok {
		return n
	}
	return identityNormalizer{}
}

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

import "strings"

// FormatOpts configures display formatting width and style hints.
type FormatOpts struct {
	Width    int
	Thinking bool
}

var markdownRenderer Renderer = defaultGlamour

// SetMarkdownRenderer replaces the default Markdown renderer (tests only).
func SetMarkdownRenderer(r Renderer) {
	if r == nil {
		markdownRenderer = defaultGlamour
		return
	}
	markdownRenderer = r
}

// Render formats text for terminal display according to kind and provider.
// Plain kind skips provider normalizers; Markdown runs Normalize then glamour.
func Render(kind ContentKind, providerID string, text string, opts FormatOpts) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	if opts.Width < 8 {
		opts.Width = 8
	}

	switch kind {
	case ContentKindPlain:
		return wrapPlainText(text, opts.Width)
	case ContentKindMarkdown:
		normalized := normalizerFor(providerID).Normalize(text)
		out, err := markdownRenderer.RenderMarkdown(normalized, opts.Width, opts.Thinking)
		if err != nil {
			return wrapPlainText(normalized, opts.Width)
		}
		return strings.TrimRight(out, "\n")
	default:
		return wrapPlainText(text, opts.Width)
	}
}

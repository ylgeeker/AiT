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

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
)

type glamourRenderer struct {
	mu    sync.Mutex
	cache map[glamourCacheKey]*glamour.TermRenderer
}

type glamourCacheKey struct {
	width    int
	thinking bool
}

var defaultGlamour = &glamourRenderer{cache: make(map[glamourCacheKey]*glamour.TermRenderer)}

func (g *glamourRenderer) RenderMarkdown(text string, width int, thinking bool) (string, error) {
	if width < 8 {
		width = 8
	}
	r, err := g.renderer(width, thinking)
	if err != nil {
		return "", err
	}
	return r.Render(text)
}

func (g *glamourRenderer) renderer(width int, thinking bool) (*glamour.TermRenderer, error) {
	key := glamourCacheKey{width: width, thinking: thinking}

	g.mu.Lock()
	defer g.mu.Unlock()

	if r, ok := g.cache[key]; ok {
		return r, nil
	}

	opts := []glamour.TermRendererOption{
		glamour.WithWordWrap(width),
		glamour.WithPreservedNewLines(),
	}
	if thinking {
		opts = append(opts, glamour.WithStyles(thinkingANSIStyles()))
	} else {
		opts = append(opts, glamour.WithStandardStyle("dark"))
	}

	r, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return nil, fmt.Errorf("aidisplay glamour: %w", err)
	}
	g.cache[key] = r
	return r, nil
}

func thinkingANSIStyles() ansi.StyleConfig {
	base := styles.DarkStyleConfig
	base.Text.Italic = boolPtr(true)
	base.Emph.Italic = boolPtr(true)
	base.Document.Margin = uintPtr(0)
	return base
}

func boolPtr(v bool) *bool {
	return &v
}

func uintPtr(v uint) *uint {
	return &v
}

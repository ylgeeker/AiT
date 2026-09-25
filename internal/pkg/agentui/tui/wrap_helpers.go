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

package tui

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// wrapWithPrefixes breaks text into display lines that fit width, prefixing the
// first line with firstPrefix and continuations with contPrefix.
func wrapWithPrefixes(firstPrefix, contPrefix, text string, width int) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{firstPrefix}
	}

	firstMax := width - runewidth.StringWidth(firstPrefix)
	contMax := width - runewidth.StringWidth(contPrefix)
	if firstMax < 1 {
		firstMax = 1
	}
	if contMax < 1 {
		contMax = 1
	}

	var out []string
	remaining := text
	useFirst := true
	for remaining != "" {
		maxW := contMax
		prefix := contPrefix
		if useFirst {
			maxW = firstMax
			prefix = firstPrefix
			useFirst = false
		}

		line, rest := takeWrapLine(remaining, maxW)
		out = append(out, prefix+line)
		remaining = strings.TrimLeft(rest, " \t")
	}
	return out
}

// wrapPlainText breaks text into lines that fit width without prefixes.
func wrapPlainText(text string, width int) string {
	lines := wrapWithPrefixes("", "", text, width)
	return strings.Join(lines, "\n")
}

func takeWrapLine(text string, width int) (line, rest string) {
	if width < 1 {
		return text, ""
	}
	if runewidth.StringWidth(text) <= width {
		return text, ""
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return "", ""
	}

	var b strings.Builder
	lineW := 0
	for i, word := range words {
		wordW := runewidth.StringWidth(word)
		if wordW > width {
			if b.Len() == 0 {
				part, left := takeHardBreak(word, width)
				restWords := words[i+1:]
				if left != "" {
					restWords = append([]string{left}, restWords...)
				}
				return part, strings.Join(restWords, " ")
			}
			return strings.TrimRight(b.String(), " "), strings.Join(words[i:], " ")
		}

		add := wordW
		if b.Len() > 0 {
			add++
		}
		if lineW+add > width {
			return strings.TrimRight(b.String(), " "), strings.Join(words[i:], " ")
		}

		if b.Len() > 0 {
			b.WriteByte(' ')
			lineW++
		}
		b.WriteString(word)
		lineW += wordW
	}
	return b.String(), ""
}

func takeHardBreak(word string, width int) (prefix, rest string) {
	runes := []rune(word)
	lineW := 0
	for i, r := range runes {
		rw := runewidth.RuneWidth(r)
		if lineW+rw > width && i > 0 {
			return string(runes[:i]), string(runes[i:])
		}
		lineW += rw
	}
	return word, ""
}

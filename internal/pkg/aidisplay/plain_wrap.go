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
	"strings"

	"github.com/mattn/go-runewidth"
)

func wrapPlainText(text string, width int) string {
	if width < 1 {
		width = 1
	}
	raw := strings.Split(text, "\n")
	out := make([]string, 0, len(raw))
	for _, ln := range raw {
		ln = strings.TrimRight(ln, "\r")
		out = append(out, wrapPlainLine(ln, width)...)
	}
	return strings.Join(out, "\n")
}

func wrapPlainLine(line string, width int) []string {
	line = strings.TrimRight(line, " \t")
	if line == "" {
		return []string{""}
	}
	if runewidth.StringWidth(line) <= width {
		return []string{line}
	}

	var out []string
	remaining := line
	for remaining != "" {
		part, rest := takePlainWrapSegment(remaining, width)
		out = append(out, part)
		remaining = strings.TrimLeft(rest, " \t")
	}
	return out
}

func takePlainWrapSegment(text string, width int) (line, rest string) {
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
			return takeOversizedWordSegment(b, words, i, word, width)
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

func takeOversizedWordSegment(b strings.Builder, words []string, i int, word string, width int) (line, rest string) {
	if b.Len() > 0 {
		return strings.TrimRight(b.String(), " "), strings.Join(words[i:], " ")
	}

	part, left := takePlainHardBreak(word, width)
	restWords := words[i+1:]
	if left != "" {
		restWords = append([]string{left}, restWords...)
	}
	return part, strings.Join(restWords, " ")
}

func takePlainHardBreak(word string, width int) (prefix, rest string) {
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

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

package lineedit

import (
	"bufio"
	"fmt"
	"io"

	"github.com/os-artificer/AiT/internal/pkg/terminal/history"
)

const (
	ctrlG = 7
	ctrlR = 18
)

func runReverseSearch(br *bufio.Reader, out io.Writer, prompt string, hist *history.LineHistory, st *editState) lineEditResult {
	if hist == nil || hist.Len() == 0 {
		return lineEditResult{}
	}

	savedRunes := append([]rune(nil), st.runes...)
	savedCursor := st.cursor

	query := ""
	matchIdx := hist.Len()
	matchLine := ""
	found := false

	matchIdx, matchLine, found = hist.ReverseSearch(query, matchIdx)
	writeReverseSearch(out, query, matchLine, found)

	for {
		r, _, err := br.ReadRune()
		if err != nil {
			return lineEditResult{err: err, done: true}
		}

		res, handled := dispatchReverseSearchKey(br, out, prompt, hist, st, savedRunes, savedCursor, &query, &matchIdx, &matchLine, &found, r)
		if handled {
			return res
		}
	}
}

func dispatchReverseSearchKey(
	br *bufio.Reader,
	out io.Writer,
	prompt string,
	hist *history.LineHistory,
	st *editState,
	saved []rune,
	savedCursor int,
	query *string,
	matchIdx *int,
	matchLine *string,
	found *bool,
	r rune,
) (lineEditResult, bool) {
	switch r {
	case ctrlR:
		*matchIdx, *matchLine, *found = reverseSearchRepeat(hist, *query, *matchIdx, *found)
		writeReverseSearch(out, *query, *matchLine, *found)
		return lineEditResult{}, false
	case ctrlG, ctrlC:
		return abortReverseSearch(out, prompt, st, saved, savedCursor), true
	case keyESC:
		return handleReverseSearchESC(br, out, prompt, st, saved, savedCursor, *matchLine, *found)
	case ctrlE:
		return acceptReverseSearchIfFound(out, prompt, st, *matchLine, *found)
	case '\r', '\n':
		return finishReverseSearchEnter(out, prompt, st, saved, savedCursor, *matchLine, *found), true
	case keyBS, keyDEL:
		if len(*query) == 0 {
			return lineEditResult{}, false
		}
		applyReverseSearchBackspace(hist, query, matchIdx, matchLine, found)
		writeReverseSearch(out, *query, *matchLine, *found)
		return lineEditResult{}, false
	default:
		if !isLineInputRune(r) || r == '\t' {
			return lineEditResult{}, false
		}
		applyReverseSearchRune(hist, query, matchIdx, matchLine, found, r)
		writeReverseSearch(out, *query, *matchLine, *found)
		return lineEditResult{}, false
	}
}

func reverseSearchRepeat(hist *history.LineHistory, query string, matchIdx int, found bool) (int, string, bool) {
	if !found {
		return hist.ReverseSearch(query, hist.Len())
	}
	return hist.ReverseSearch(query, matchIdx)
}

func handleReverseSearchESC(
	br *bufio.Reader,
	out io.Writer,
	prompt string,
	st *editState,
	saved []rune,
	savedCursor int,
	matchLine string,
	found bool,
) (lineEditResult, bool) {
	key, err := parseEscapeSeq(br)
	if err != nil {
		return lineEditResult{err: err, done: true}, true
	}
	if key == escapeNone {
		return abortReverseSearch(out, prompt, st, saved, savedCursor), true
	}
	if (key == escapeRight || key == escapeEnd) && found {
		return acceptReverseSearchEdit(out, prompt, st, matchLine), true
	}
	return lineEditResult{}, false
}

func acceptReverseSearchIfFound(out io.Writer, prompt string, st *editState, matchLine string, found bool) (lineEditResult, bool) {
	if !found {
		return lineEditResult{}, false
	}
	return acceptReverseSearchEdit(out, prompt, st, matchLine), true
}

func finishReverseSearchEnter(
	out io.Writer,
	prompt string,
	st *editState,
	saved []rune,
	savedCursor int,
	matchLine string,
	found bool,
) lineEditResult {
	if !found {
		return abortReverseSearch(out, prompt, st, saved, savedCursor)
	}
	st.setLine(matchLine)
	_, _ = io.WriteString(out, "\r\n")
	return lineEditResult{line: matchLine, done: true}
}

func applyReverseSearchBackspace(hist *history.LineHistory, query *string, matchIdx *int, matchLine *string, found *bool) {
	*query = (*query)[:len(*query)-1]
	*matchIdx, *matchLine, *found = hist.ReverseSearch(*query, hist.Len())
}

func applyReverseSearchRune(hist *history.LineHistory, query *string, matchIdx *int, matchLine *string, found *bool, r rune) {
	*query += string(r)
	*matchIdx, *matchLine, *found = hist.ReverseSearch(*query, hist.Len())
}

func writeReverseSearch(out io.Writer, query, match string, found bool) {
	display := match
	if !found {
		display = ""
	}
	_, _ = fmt.Fprintf(out, "\r\x1b[K(reverse-i-search)`%s': %s", query, display)
}

func abortReverseSearch(out io.Writer, prompt string, st *editState, saved []rune, cursor int) lineEditResult {
	st.runes = saved
	st.cursor = cursor
	st.histIdx = -1
	redrawLine(out, prompt, st.runes, st.cursor)
	return lineEditResult{}
}

func acceptReverseSearchEdit(out io.Writer, prompt string, st *editState, match string) lineEditResult {
	st.setLine(match)
	redrawLine(out, prompt, st.runes, st.cursor)
	return lineEditResult{}
}

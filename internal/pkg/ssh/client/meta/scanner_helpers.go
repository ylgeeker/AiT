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

package meta

import (
	"strings"

	"github.com/os-artificer/AiT/internal/pkg/toolname"
)

// normalizeBangAitTypo maps a common typo (!!ait → !ait) to the real meta opener.
// Does not alter !! AI lines (those require a space or newline after !!).
func normalizeBangAitTypo(s string) string {
	if strings.HasPrefix(s, bangAitTypoPrefix) {
		return "!" + s[2:]
	}
	// Partial bangAitTypoPrefix while the user is still typing.
	if len(s) >= 3 && strings.HasPrefix(s, "!!a") && strings.HasPrefix(bangAitTypoPrefix, s) {
		return "!" + s[2:]
	}
	return s
}

func isBangAitTypoPrefix(candidate string) bool {
	return strings.HasPrefix(candidate, "!!") &&
		normalizeBangAitTypo(candidate) != candidate
}

func finalizeMetaLine(candidate string) string {
	return normalizeBangAitTypo(strings.TrimRight(candidate, "\r\n"))
}

// prefixCandidate reports whether s is a strict prefix of any opener.
func prefixCandidate(s string) bool {
	// Hold "!" / "!!" until the next byte disambiguates bash history vs !!ait typo.
	if s == "!" || s == "!!" {
		return true
	}
	s = normalizeBangAitTypo(s)
	for _, p := range prefixes {
		if len(s) < len(p) && strings.HasPrefix(p, s) {
			return true
		}
	}
	return false
}

// prefixMatch returns the matched opener if s equals one exactly.
func prefixMatch(s string) (string, bool) {
	s = normalizeBangAitTypo(s)
	for _, p := range prefixes {
		if s == p {
			return p, true
		}
	}
	return "", false
}

// lineEnd reports whether b terminates a meta line.
func lineEnd(b byte) bool {
	return b == '\n' || b == '\r'
}

// stripTrailingPinyinAtCursor removes consecutive lowercase ASCII letters before
// the cursor and returns how many runes were removed.
func (s *Scanner) stripTrailingPinyinAtCursor() int {
	rs := s.runes()
	i := s.cursor
	for i > 0 {
		c := rs[i-1]
		if c >= 'a' && c <= 'z' {
			i--
		} else {
			break
		}
	}
	removed := s.cursor - i
	if removed > 0 {
		s.setRunes(append(rs[:i], rs[s.cursor:]...))
		s.cursor = i
	}
	return removed
}

func completeLine(candidate string, b byte) (line string, done bool) {
	if !lineEnd(b) {
		return "", false
	}
	return finalizeMetaLine(candidate), true
}

func (s *Scanner) feedOneEditPrefix(b byte, fwd, ech []byte) ([]byte, []byte, bool) {
	if !s.inEditState() {
		return fwd, ech, false
	}
	if fwd, ech, handled := s.feedOneEditEsc(b, fwd, ech); handled {
		return fwd, ech, true
	}
	if b == keyESC {
		s.escState = escReceived
		return fwd, ech, true
	}
	return fwd, ech, false
}

func (s *Scanner) feedOneEditEsc(b byte, fwd, ech []byte) ([]byte, []byte, bool) {
	if s.escState == escNone {
		return fwd, ech, false
	}
	echOut, handled := s.consumeEscForEdit(b)
	if !handled {
		return fwd, ech, false
	}
	return fwd, append(ech, echOut...), true
}

func (s *Scanner) feedOneStateMeta(b byte, fwd, ech []byte) ([]byte, []byte, string, bool) {
	if len(s.buf) >= maxLine {
		fwd = append(fwd, s.buf...)
		s.buf = s.buf[:0]
		s.cursor = 0
		s.enterIdle()
		s.atLineStart = false
		return s.feedOne(b, fwd, ech)
	}
	if b == '\n' || b == '\r' {
		line := finalizeMetaLine(string(s.buf))
		s.buf = s.buf[:0]
		s.cursor = 0
		s.enterIdle()
		s.atLineStart = true
		ech = append(ech, b)
		return fwd, ech, line, true
	}
	ech = s.insertMetaByte(b, ech)
	return fwd, ech, "", false
}

func (s *Scanner) feedOneStateHolding(b byte, fwd, ech []byte) ([]byte, []byte, string, bool) {
	candidate := string(s.buf) + string([]byte{b})
	if f, e, handled := s.tryLocalBangPassthrough(candidate, b, fwd, ech); handled {
		return f, e, "", false
	}
	if f, e, line, hasLine, handled := s.feedHoldingLoose(candidate, b, fwd, ech); handled {
		return feedHoldingHandledReturn(f, e, line, hasLine)
	}
	if f, e, line, hasLine, handled := s.feedHoldingStrict(candidate, b, fwd, ech); handled {
		return feedHoldingHandledReturn(f, e, line, hasLine)
	}
	return s.feedHoldingPrefixGrow(candidate, b, fwd, ech)
}

func feedHoldingHandledReturn(f, e []byte, line string, hasLine bool) ([]byte, []byte, string, bool) {
	if hasLine {
		return f, e, line, true
	}
	return f, e, "", false
}

func (s *Scanner) feedHoldingPrefixGrow(candidate string, b byte, fwd, ech []byte) ([]byte, []byte, string, bool) {
	if !prefixCandidate(candidate) || len(s.buf)+1 >= maxLook {
		ech, fwd = s.flushHoldingToForward(ech, fwd)
		s.buf = s.buf[:0]
		s.cursor = 0
		s.enterIdle()
		return s.feedOne(b, fwd, ech)
	}
	s.buf = append(s.buf, b)
	s.cursor = len(s.runes())
	return s.holdingPrefixEcho(b, fwd, ech)
}

func (s *Scanner) holdingPrefixEcho(b byte, fwd, ech []byte) ([]byte, []byte, string, bool) {
	if s.shouldDeferHoldingEcho(b) {
		return s.holdingPrefixDeferredEcho(fwd, ech)
	}
	if s.localShellActive {
		return fwd, s.holdingAppendEcho(ech, b), "", false
	}
	return fwd, s.noteHoldingRedraw(ech), "", false
}

func (s *Scanner) holdingPrefixDeferredEcho(fwd, ech []byte) ([]byte, []byte, string, bool) {
	if s.localShellActive {
		return fwd, ech, "", false
	}
	return fwd, s.noteHoldingRedraw(ech), "", false
}

func (s *Scanner) feedOneStateIdle(b byte, fwd, ech []byte) ([]byte, []byte, string, bool) {
	if b == '!' {
		s.buf = append(s.buf[:0], b)
		s.cursor = 1
		s.state = stateHolding
		s.deferLocalEcho = true
		return fwd, ech, "", false
	}
	if b == 0x15 { // Ctrl+U idle — forward to remote shell.
		s.atLineStart = true
		fwd = append(fwd, b)
		return fwd, ech, "", false
	}
	fwd = append(fwd, b)
	s.noteForwarded(b)
	return fwd, ech, "", false
}

func (s *Scanner) feedHoldingLoose(candidate string, b byte, fwd, ech []byte) (outFwd, outEch []byte, line string, hasLine, handled bool) {
	if !s.looseBangCapture(candidate, b) {
		return fwd, ech, "", false, false
	}
	if line, ok := completeLine(candidate, b); ok {
		s.buf = s.buf[:0]
		s.cursor = 0
		s.enterIdle()
		s.atLineStart = true
		return fwd, ech, line, true, true
	}
	s.buf = append(s.buf, b)
	s.cursor = len(s.runes())
	// Exact "!ait" still needs space/newline; stay in holding so "!aitssh" can flush.
	if candidate == toolname.MetaPrefix {
		fwd, ech, _, _ := s.holdingPrefixEcho(b, fwd, ech)
		return fwd, ech, "", false, true
	}
	if s.shouldDeferHoldingEcho(b) {
		return fwd, ech, "", false, true
	}
	return fwd, s.echoEnteringStateMeta(ech, b), "", false, true
}

func (s *Scanner) feedHoldingStrict(candidate string, b byte, fwd, ech []byte) (outFwd, outEch []byte, line string, hasLine, handled bool) {
	if _, ok := prefixMatch(candidate); !ok {
		return fwd, ech, "", false, false
	}
	if line, ok := completeLine(candidate, b); ok {
		s.buf = s.buf[:0]
		s.cursor = 0
		s.enterIdle()
		s.atLineStart = true
		return fwd, ech, line, true, true
	}
	s.buf = append(s.buf, b)
	s.cursor = len(s.runes())
	if s.shouldDeferHoldingEcho(b) {
		return fwd, ech, "", false, true
	}
	return fwd, s.echoEnteringStateMeta(ech, b), "", false, true
}

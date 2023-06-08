// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package boolexpr

import (
	"errors"
	"fmt"
	"strings"
)

// A basic boolean parser which will parse a string and return the boolean result with a used variable tracker
func ParseWithUsed(str string, vars map[string]bool) (val bool, used map[string]bool, err error) {
	used = make(map[string]bool)
	str, val, err = parseSet(str, vars, used)
	if err == nil && strings.TrimSpace(str) != "" {
		return val, used, fmt.Errorf("Leftover values in parse: %q", strings.TrimSpace(str))
	}
	return
}

// A basic boolean parser which will parse a string and return the boolean result.
func Parse(str string, vars map[string]bool) (val bool, err error) {
	str, val, err = parseSet(str, vars, nil)
	if err == nil && strings.TrimSpace(str) != "" {
		return val, fmt.Errorf("Leftover values in parse: %q", strings.TrimSpace(str))
	}
	return
}

const (
	opOr = iota + 1
	opAnd
)

func parseSet(s string, vars, used map[string]bool) (string, bool, error) {
	var cur bool
	var err error
	var op uint8
	for s != "" {
		// consume the next token
		switch s[0] {
		case '(':
			var result bool
			s, result, err = parseSet(s[1:], vars, used)
			if err != nil {
				return "", false, err
			}
			if len(s) > 0 && s[0] == ')' {
				s = s[1:]
				switch op & 0xf {
				case opOr:
					cur = cur || result
					op = 0
				case opAnd:
					cur = cur && result
					op = 0
				default:
					cur = result
				}
				continue
			}
			return "", false, errors.New("boolexp: no matching )")
		}

		// Consume var.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
				continue
			}
			break
		}
		if i == 0 {
			return "", false, errors.New("boolexp: missing variable in expression " + quote(s))
		}
		u := s[:i]
		s = s[i:]
		val, ok := vars[u]
		if !ok {
			return "", false, errors.New("boolexp: unknown variable " + quote(u) + " in expression " + quote(s))
		}
		if used != nil {
			used[u] = true
		}

		// Combine
		switch op & 0xf {
		case opOr:
			cur = cur || val
			op = 0
		case opAnd:
			cur = cur && val
			op = 0
		default:
			cur = val
		}
		//fmt.Println("cur=", cur)

		// Drop space
		for len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}

		if len(s) == 0 || s[0] == ')' {
			return s, cur, nil
		}

		// Consume op.
		switch s[0] {
		case '|':
			op = opOr
			s = s[1:]

			// Drop space
			for len(s) > 0 && s[0] == ' ' {
				s = s[1:]
			}
			continue
		case '&':
			op = opAnd
			s = s[1:]

			// Drop space
			for len(s) > 0 && s[0] == ' ' {
				s = s[1:]
			}
			continue
		default:
			for ; i < len(s); i++ {
				c := s[i]
				if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
					continue
				}
				break
			}
			if i == 0 {
				return "", false, errors.New("boolexp: missing logical expression " + quote(s))
			}
			switch s[:i] {
			case "or", "OR":
				op = opOr
				s = s[2:]

				// Drop space
				for len(s) > 0 && s[0] == ' ' {
					s = s[1:]
				}
				continue
			case "and", "AND":
				op = opAnd
				s = s[3:]

				// Drop space
				for len(s) > 0 && s[0] == ' ' {
					s = s[1:]
				}
				continue
			}
			return "", false, errors.New("boolexp: invalid logical expression " + quote(s))
		}
	}
	return s, cur, err
}

// These are borrowed from unicode/utf8 and strconv and replicate behavior in
// that package, since we can't take a dependency on either.
const (
	lowerhex  = "0123456789abcdef"
	runeSelf  = 0x80
	runeError = '\uFFFD'
)

func quote(s string) string {
	buf := make([]byte, 1, len(s)+2) // slice will be at least len(s) + quotes
	buf[0] = '"'
	for i, c := range s {
		if c >= runeSelf || c < ' ' {
			// This means you are asking us to parse a string with unprintable or
			// non-ASCII characters in it.  We don't expect to hit this case very
			// often. We could try to reproduce strconv.Quote's behavior with full
			// fidelity but given how rarely we expect to hit these edge cases, speed
			// and conciseness are better.
			var width int
			if c == runeError {
				width = 1
				if i+2 < len(s) && s[i:i+3] == string(runeError) {
					width = 3
				}
			} else {
				width = len(string(c))
			}
			for j := 0; j < width; j++ {
				buf = append(buf, `\x`...)
				buf = append(buf, lowerhex[s[i+j]>>4])
				buf = append(buf, lowerhex[s[i+j]&0xF])
			}
		} else {
			if c == '"' || c == '\\' {
				buf = append(buf, '\\')
			}
			buf = append(buf, string(c)...)
		}
	}
	buf = append(buf, '"')
	return string(buf)
}

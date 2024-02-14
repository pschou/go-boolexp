// Copyright 2023 github.com/pschou
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
func ParseWithUsed(str string, vars map[string][]bool, used map[string]bool) (val bool, err error) {
	var rem []rune
	rem, val, err = parseSet([]rune(str), vars, used)
	if err == nil && strings.TrimSpace(string(rem)) != "" {
		return val, fmt.Errorf("Leftover values in parse: %q", strings.TrimSpace(string(rem)))
	}
	return
}

// A basic boolean parser which will parse a string and return the boolean result.
func Parse(str string, vars map[string][]bool) (val bool, err error) {
	var rem []rune
	rem, val, err = parseSet([]rune(str), vars, nil)
	if err == nil && strings.TrimSpace(string(rem)) != "" {
		return val, fmt.Errorf("Leftover values in parse: %q", strings.TrimSpace(string(rem)))
	}
	return
}

const (
	opOr = iota + 1
	opAnd
	opXor

	aggAny = iota + 1
	aggAll
)

// Strip out the rest of the line on the same depth or lower
func chomp(s []rune) (a, b []rune, err error) {
	var paran int
	for i, c := range s {
		switch c {
		case '(':
			paran++
		case ')':
			paran--
			if paran == -1 {
				return s[:i], s[i:], nil
			}
		}
	}
	if paran == 0 {
		return s, nil, nil
	}
	return nil, nil, errors.New("Mismatched ()'s")
}

func parseSet(s []rune, vars map[string][]bool, used map[string]bool) ([]rune, bool, error) {
	// Drop space
	for len(s) > 0 && s[0] == ' ' {
		s = s[1:]
	}

	var cur, next bool
	var err error
	var op, agg uint8
	var neg, has_cur bool
	for len(s) > 0 {
		//fmt.Println("op = ", op, "s=", s, "cur=", cur, "next=", next, "has_cur=", has_cur)
		if has_cur {
			has_cur = false
		} else {
			// consume the next token
			switch s[0] {
			case '!':
				neg = !neg
				s = s[1:]
				// Drop space
				for len(s) > 0 && s[0] == ' ' {
					s = s[1:]
				}
				continue
			case '(':
				s, next, err = parseSet(s[1:], vars, used)
				if err != nil {
					return nil, false, err
				}
				if len(s) > 0 && s[0] == ')' {
					s = s[1:]

					// Flip if negative is declared
					next = next != neg
					neg = false

					switch op & 0xf {
					case opOr:
						cur = cur || next
						op = 0
					case opAnd:
						cur = cur && next
						op = 0
					case opXor:
						cur = cur != next
						op = 0
					default:
						cur = next
					}
					if neg {
						cur = !cur
						neg = false
					}
					// Drop space
					for len(s) > 0 && s[0] == ' ' {
						s = s[1:]
					}
					has_cur = true
					continue
				}
				return nil, false, errors.New("boolexp: no matching )")
			}

			// Consume var.
			i := 0
			for ; i < len(s); i++ {
				c := s[i]
				if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' || c == '_' || c == '-' || c == '.' {
					continue
				}
				break
			}
			if i == 0 {
				return nil, false, errors.New("boolexp: missing variable in expression " + quote(s))
			}
			u := s[:i]
			s = s[i:]

			// Test for logical changers
			switch string(u) {
			case "not", "NOT":
				neg = !neg

				// Drop space
				for len(s) > 0 && s[0] == ' ' {
					s = s[1:]
				}
				continue
			case "all", "ALL":
				agg = aggAll
				// Drop space
				for len(s) > 0 && s[0] == ' ' {
					s = s[1:]
				}

				i = 0
				for ; i < len(s); i++ {
					c := s[i]
					if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
						continue
					}
					break
				}
				u = s[:i]
				s = s[i:]

			case "any", "ANY":
				agg = aggAny
				// Drop space
				for len(s) > 0 && s[0] == ' ' {
					s = s[1:]
				}

				i = 0
				for ; i < len(s); i++ {
					c := s[i]
					if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
						continue
					}
					break
				}
				u = s[:i]
				s = s[i:]

			case "none", "NONE":
				agg = aggAny
				// Drop space
				for len(s) > 0 && s[0] == ' ' {
					s = s[1:]
				}

				i = 0
				for ; i < len(s); i++ {
					c := s[i]
					if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
						continue
					}
					break
				}
				neg = !neg
				u = s[:i]
				s = s[i:]
			}

			//fmt.Println("found var", u)
			// Parse var
			val, ok := vars[string(u)]
			if !ok {
				return nil, false, errors.New("boolexp: unknown variable " + quote(u) + " in expression " + quote(s))
			}
			if used != nil {
				used[string(u)] = true
			}
			switch len(val) {
			case 0:
				return nil, false, errors.New("boolexp: variable " + quote(u) + " has no values")
			case 1:
				next = val[0]
			default:
				next = val[0]
				switch agg {
				case aggAll:
					for i := 1; i < len(val); i++ {
						next = next && val[i]
					}
				case aggAny:
					for i := 1; i < len(val); i++ {
						next = next || val[i]
					}
				default:
					return nil, false, errors.New("boolexp: multiple values for " + quote(u) + " with no aggregate operator")
				}
				agg = 0
			}
			// Flip if negative is declared
			next = next != neg
			neg = false

			// Combine
			switch op & 0xf {
			case opOr:
				cur = cur || next
				op = 0
			case opAnd:
				cur = cur && next
				op = 0
			case opXor:
				cur = cur != next
				op = 0
			default:
				cur = next
			}
			//fmt.Println("cur=", cur)
		}

		// Drop space
		for len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}

		// Return if at the end of a grouping
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

			if cur {
				op = 0
				var a []rune
				a, s, err = chomp(s)
				if err != nil {
					return nil, false, err
				}
				if len(a) > 0 {
					if _, _, err = parseSet(a, vars, used); err != nil {
						return nil, false, err
					}
				}
				has_cur = true
				continue
			}
			continue
		case '&':
			op = opAnd
			s = s[1:]

			// Drop space
			for len(s) > 0 && s[0] == ' ' {
				s = s[1:]
			}

			if !cur {
				op = 0
				var a []rune
				a, s, err = chomp(s)
				if err != nil {
					return nil, false, err
				}
				if len(a) > 0 {
					if _, _, err = parseSet(a, vars, used); err != nil {
						return nil, false, err
					}
				}
				has_cur = true
				continue
			}
			continue
		case '^':
			op = opXor
			s = s[1:]

			// Drop space
			for len(s) > 0 && s[0] == ' ' {
				s = s[1:]
			}
			continue
		default:
			i := 0
			for ; i < len(s); i++ {
				c := s[i]
				if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
					continue
				}
				break
			}
			if i == 0 {
				return nil, false, errors.New("boolexp: missing logical expression " + quote(s))
			}
			switch string(s[:i]) {
			case "or", "OR":
				op = opOr
				s = s[2:]

				// Drop space
				for len(s) > 0 && s[0] == ' ' {
					s = s[1:]
				}

				if cur {
					op = 0
					var a []rune
					a, s, err = chomp(s)
					if err != nil {
						return nil, false, err
					}
					if len(a) > 0 {
						if _, _, err = parseSet(a, vars, used); err != nil {
							return nil, false, err
						}
					}
					has_cur = true
					continue
				}
				continue
			case "and", "AND":
				op = opAnd
				s = s[3:]

				// Drop space
				for len(s) > 0 && s[0] == ' ' {
					s = s[1:]
				}

				if !cur {
					op = 0
					var a []rune
					a, s, err = chomp(s)
					if err != nil {
						return nil, false, err
					}
					if len(a) > 0 {
						if _, _, err = parseSet(a, vars, used); err != nil {
							return nil, false, err
						}
					}
					has_cur = true
					continue
				}
				continue
			case "xor", "XOR":
				op = opXor
				s = s[3:]

				// Drop space
				for len(s) > 0 && s[0] == ' ' {
					s = s[1:]
				}
				continue
			}
			return nil, false, errors.New("boolexp: invalid logical expression " + quote(s))
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

func quote(s []rune) string {
	buf := make([]byte, 1, len(s)+2) // slice will be at least len(s) + quotes
	buf[0] = '"'
	for i, c := range string(s) {
		if c >= runeSelf || c < ' ' {
			// This means you are asking us to parse a string with unprintable or
			// non-ASCII characters in it.  We don't expect to hit this case very
			// often. We could try to reproduce strconv.Quote's behavior with full
			// fidelity but given how rarely we expect to hit these edge cases, speed
			// and conciseness are better.
			var width int
			if c == runeError {
				width = 1
				if i+2 < len(s) && string(s[i:i+3]) == string(runeError) {
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

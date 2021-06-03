/**
 * Copyright 2021 The cal2cat authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cat

import (
	"path"
	"regexp"
	"strings"
)

// Matcher returns true if the string matches a certain condition.
type Matcher func(string) bool

// Mapper maps strings to other strings.
type Mapper func(string) string

// NewMatcher creates a new Matcher for the pattern.
//
// The matching operation depends on the pattern:
//
// - starts with '^': use regular expression
//
// - contains any of "*?[": use globing
//
// - otherwise: prefix search
func NewMatcher(pattern string) Matcher {
	if strings.HasPrefix(pattern, "^") {
		return regexp.MustCompile(pattern).MatchString
	}

	if strings.ContainsAny(pattern, "*?[") {
		return func(s string) bool {
			ok, err := path.Match(pattern, s)
			return err == nil && ok
		}
	}

	return func(s string) bool {
		return strings.HasPrefix(s, pattern)
	}
}

// NewMapper creates a new Mapper, which checks if the input matches pattern
// and returns either str or an empty string.
func NewMapper(pattern, str string) Mapper {
	m := NewMatcher(pattern)
	return func(s string) string {
		if m(s) {
			return str
		}
		return ""
	}
}

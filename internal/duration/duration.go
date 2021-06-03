// Copyright 2021 The cal2cat authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package duration

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var timeUnitPattern = regexp.MustCompile(`(-?\d+)(d|w|cw|m|cm|y|cy)`)

// Calc adds an offset to the given time.
func Calc(t time.Time, offset string) time.Time {
	matches := timeUnitPattern.FindStringSubmatch(strings.TrimSpace(offset))
	if len(matches) == 0 {
		log.Fatalln("cannot parse duration:", offset)
	}

	n, _ := strconv.Atoi(matches[1])
	u := matches[2]
	switch u {
	case "d":
		return t.AddDate(0, 0, n)
	case "w":
		return t.AddDate(0, 0, 7*n)
	case "cw":
		return t.AddDate(0, 0, 7*n).Truncate(7 * 24 * time.Hour)
	case "m":
		return t.AddDate(0, n, 0)
	case "cm":
		return t.AddDate(0, n, -t.Day()+1).Truncate(24 * time.Hour)
	case "y":
		return t.AddDate(n, 0, 0)
	case "cy":
		return t.AddDate(n, -int(t.Month())+1, -t.Day()+1).Truncate(24 * time.Hour)
	}
	log.Fatalln("cannot parse duration:", offset)
	return time.Time{}
}

// Format returns a textual representation of the duration value formatted by
// the given layout format.
func Format(d time.Duration, f string) string {
	switch f {
	case "minutes":
		return strconv.FormatFloat(d.Minutes(), 'f', 0, 64)
	case "hours":
		return strconv.FormatFloat(d.Hours(), 'f', 2, 64)
	default:
		return time.Time{}.Add(d).Format(f)
	}
}

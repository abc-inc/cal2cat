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

package event

// present is the value indicating that a key is present in a map.
var present interface{}

// Conflict represents overlapping events.
type Conflict struct {
	Event  Wrapper
	Reason Wrapper
}

// Conflicts represents multiple overlapping events.
type Conflicts []Conflict

// Events returns the conflicting events.
func (cs Conflicts) Events() Events {
	es := make([]Wrapper, len(cs))
	for i, c := range cs {
		es[i] = c.Event
	}
	return es
}

// Conflicts returns all conflicting events with their first overlapping event.
func (es Events) Conflicts() (cs Conflicts) {
	ignored := map[Wrapper]interface{}{}
	for i, j := 0, 1; i < len(es); j++ {
		if j >= len(es) {
			// if the second index runs out of range, proceed with the next one
			i, j = i+1, i+1
			continue
		}

		cur := es[i]
		next := es[j]

		_, ignore := ignored[cur]
		switch {
		case ignore:
			// if the current event is a known conflict, proceed with the next one
			i, j = i+1, i+1
		case !overlaps(cur, next):
			// if the current event does not overlap with the next one, proceed
			i, j = i+1, i+1
		default:
			// conflict detect: record the event and the cause
			ignored[next] = present
			cs = append(cs, Conflict{next, cur})
		}
	}
	return
}

// overlaps checks if the second event starts before the first one ends.
func overlaps(fst, snd Wrapper) bool {
	return snd.StartTime().Before(fst.EndTime())
}

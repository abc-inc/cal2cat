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

package event_test

import (
	"testing"
	"time"

	. "github.com/abc-inc/cal2cat/event"
	. "github.com/stretchr/testify/require"
)

// TestEvents_Conflicts tests the conflict resolution strategy.
//
//	07:00 A
//	08:00 A B C
//	09:00   B   D
//	10:00   B
//
//	- A is reported
//	- B is skipped because of A
//	- C is skipped because of A
//	- C is NOT skipped because of B (already conflicts with A)
//	- D is reported (NOT skipped because B is conflicting with A)
func TestEvents_Conflicts(t *testing.T) {
	today := time.Now().Truncate(24 * time.Hour)
	es := Events([]Wrapper{
		NewSimpleEvent(today.Add(7*time.Hour), today.Add(9*time.Hour), "A"),
		NewSimpleEvent(today.Add(8*time.Hour), today.Add(11*time.Hour), "B"),
		NewSimpleEvent(today.Add(8*time.Hour), today.Add(9*time.Hour), "C"),
		NewSimpleEvent(today.Add(9*time.Hour), today.Add(10*time.Hour), "D"),
	})

	cs := es.Conflicts()
	Equal(t, 2, len(cs))

	Equal(t, "B", cs[0].Event.Summary())
	Equal(t, "A", cs[0].Reason.Summary())

	Equal(t, "C", cs[1].Event.Summary())
	Equal(t, "A", cs[1].Reason.Summary())

	es = es.Filter(NewEventFilter(cs.Events()).Not())
	Equal(t, 2, len(es))
	Equal(t, "A", es[0].Summary())
	Equal(t, "D", es[1].Summary())
}

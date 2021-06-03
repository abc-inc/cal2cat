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

import (
	"time"

	"github.com/emersion/go-ical"
)

// Wrapper is the interface that defines common events.
type Wrapper interface {
	StartTime() time.Time
	EndTime() time.Time
	Duration() time.Duration
	Summary() string
}

// CalEvent represents an iCalendar event.
type CalEvent struct {
	event ical.Event
}

// NewCalEvent creates a new iCalendar event.
func NewCalEvent(e ical.Event) *CalEvent {
	return &CalEvent{e}
}

// StartTime returns the start time in the system's local time zone.
func (a CalEvent) StartTime() time.Time {
	t, _ := a.event.DateTimeStart(time.Local)
	return t
}

// EndTime returns the end time in the system's local time zone.
func (a CalEvent) EndTime() time.Time {
	t, _ := a.event.DateTimeEnd(time.Local)
	return t
}

// Duration returns the event duration.
func (a CalEvent) Duration() time.Duration {
	return a.EndTime().Sub(a.StartTime())
}

// Summary returns the summary of the iCalendar event.
func (a CalEvent) Summary() string {
	return a.event.Props.Get(ical.PropSummary).Value
}

// SimpleEvent represents an event with minimal set of properties.
type SimpleEvent struct {
	startTime time.Time
	endTime   time.Time
	summary   string
}

// NewSimpleEvent creates a new custom event.
func NewSimpleEvent(startTime, endTime time.Time, summary string) *SimpleEvent {
	return &SimpleEvent{startTime, endTime, summary}
}

// StartTime returns the start time in the system's local time zone.
func (e SimpleEvent) StartTime() time.Time {
	return e.startTime
}

// EndTime returns the end time in the system's local time zone.
func (e SimpleEvent) EndTime() time.Time {
	return e.endTime
}

// Duration returns the event duration.
func (e SimpleEvent) Duration() time.Duration {
	return e.EndTime().Sub(e.StartTime())
}

// Summary returns the summary of the event.
func (e SimpleEvent) Summary() string {
	return e.summary
}

// Filter checks whether an event matches a certain criteria.
type Filter func(e Wrapper) bool

// Events represents series of events.
type Events []Wrapper

// Duration calculates the total duration of all events.
func (es Events) Duration() (d time.Duration) {
	for _, e := range es {
		d += e.Duration()
	}
	return
}

// Filter returns a subset of events matching the filter.
func (es Events) Filter(f Filter) Events {
	sub := []Wrapper{}
	for _, e := range es {
		if f(e) {
			sub = append(sub, e)
		}
	}
	return sub
}

// NewRangeFilter returns a new Filter, which checks whether an event
// begins within the given range.
func NewRangeFilter(start, end time.Time) Filter {
	return func(e Wrapper) bool {
		return e.StartTime().After(start) && e.EndTime().Before(end)
	}
}

// NewEventFilter returns a new Filter, which checks whether an event is
// contained in another set of events.
func NewEventFilter(set Events) Filter {
	m := map[Wrapper]interface{}{}
	for _, e := range set {
		m[e] = present
	}

	return func(e Wrapper) bool {
		_, found := m[e]
		return found
	}
}

// Not negates the given Filter.
func (f Filter) Not() Filter {
	return func(e Wrapper) bool {
		return !f(e)
	}
}

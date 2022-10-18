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

package cal

import (
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/abc-inc/cal2cat/event"
	"github.com/emersion/go-ical"
)

// Load loads multiple calendars and returns an ordered series of events.
func Load(paths ...string) event.Events {
	es := []ical.Event{}
	for _, p := range paths {
		cal := decode(p)
		es = append(es, cal.Events()...)
	}

	eas := make([]event.Wrapper, len(es))
	for i, e := range es {
		e.Props.Get(ical.PropDateTimeStart).Params.Set(ical.ParamTimezoneID, "Local")
		e.Props.Get(ical.PropDateTimeEnd).Params.Set(ical.ParamTimezoneID, "Local")
		eas[i] = *event.NewCalEvent(e)
	}

	sort.SliceStable(eas, func(i, j int) bool {
		aStart := eas[i].StartTime()
		bStart := eas[j].StartTime()
		return aStart.Before(bStart)
	})

	return eas
}

// decode creates a new Calendar from the given path.
func decode(path string) *ical.Calendar {
	r, err := readFrom(path)
	if err != nil {
		log.Fatalf("cannot read calender from %s: %v", path, err)
	}

	cal, err := ical.NewDecoder(r).Decode()
	r.Close()
	if err != nil {
		log.Fatalf("cannot parse calendar from %s: %v", path, err)
	}
	return cal
}

// readFrom opens a file for reading or downloads it using HTTP.
func readFrom(path string) (io.ReadCloser, error) {
	if !strings.Contains(path, "://") {
		return os.Open(path)
	}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/"))) //nolint:gosec
	c := &http.Client{Transport: t}
	res, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

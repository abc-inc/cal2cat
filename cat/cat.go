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
	"sort"

	"github.com/abc-inc/cal2cat/event"
)

// Category is a named series of events.
type Category struct {
	Name   string
	Events event.Events
}

// Map categorizes events using the given mapping.
func Map(es event.Events, ms []Mapper) []Category {
	esByCatName := map[string]event.Events{}
	for _, e := range es {
		summary := e.Summary()
		for _, m := range ms {
			if n := m(summary); n != "" {
				esByCatName[n] = append(esByCatName[n], e)
				break
			}
		}
	}

	cns := []string{}
	for cn := range esByCatName {
		cns = append(cns, cn)
	}
	sort.Strings(cns)

	cs := make([]Category, len(cns))
	for i, cn := range cns {
		cs[i] = Category{cn, esByCatName[cn]}
	}
	return cs
}

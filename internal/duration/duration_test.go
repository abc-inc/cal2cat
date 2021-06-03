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

package duration_test

import (
	"testing"
	"time"

	. "github.com/abc-inc/cal2cat/internal/duration"
	. "github.com/stretchr/testify/assert"
)

func Test_Calc(t *testing.T) {
	now := time.Date(2020, 4, 10, 15, 4, 5, 0, time.UTC)

	Equal(t, "2020-04-10 15:04:05 +0000 UTC", Calc(now, "-0d").String())
	Equal(t, "2020-04-07 15:04:05 +0000 UTC", Calc(now, "-3d").String())
	Equal(t, "2020-04-13 15:04:05 +0000 UTC", Calc(now, "3d").String())

	Equal(t, "2020-04-10 15:04:05 +0000 UTC", Calc(now, "-0w").String())
	Equal(t, "2020-03-27 15:04:05 +0000 UTC", Calc(now, "-2w").String())
	Equal(t, "2020-04-24 15:04:05 +0000 UTC", Calc(now, "2w").String())

	Equal(t, "2020-04-06 00:00:00 +0000 UTC", Calc(now, "-0cw").String())
	Equal(t, "2020-03-23 00:00:00 +0000 UTC", Calc(now, "-2cw").String())
	Equal(t, "2020-04-20 00:00:00 +0000 UTC", Calc(now, "2cw").String())

	Equal(t, "2020-04-10 15:04:05 +0000 UTC", Calc(now, "-0m").String())
	Equal(t, "2020-02-10 15:04:05 +0000 UTC", Calc(now, "-2m").String())
	Equal(t, "2020-06-10 15:04:05 +0000 UTC", Calc(now, "2m").String())

	Equal(t, "2020-04-01 00:00:00 +0000 UTC", Calc(now, "-0cm").String())
	Equal(t, "2020-02-01 00:00:00 +0000 UTC", Calc(now, "-2cm").String())
	Equal(t, "2020-06-01 00:00:00 +0000 UTC", Calc(now, "2cm").String())

	Equal(t, "2020-04-10 15:04:05 +0000 UTC", Calc(now, "-0y").String())
	Equal(t, "2018-04-10 15:04:05 +0000 UTC", Calc(now, "-2y").String())
	Equal(t, "2022-04-10 15:04:05 +0000 UTC", Calc(now, "2y").String())

	Equal(t, "2020-01-01 00:00:00 +0000 UTC", Calc(now, "-0cy").String())
	Equal(t, "2018-01-01 00:00:00 +0000 UTC", Calc(now, "-2cy").String())
	Equal(t, "2022-01-01 00:00:00 +0000 UTC", Calc(now, "2cy").String())
}

func Test_Format(t *testing.T) {
	Equal(t, "45", Format(45*time.Minute, "minutes"))
	Equal(t, "0.75", Format(45*time.Minute, "hours"))
	Equal(t, "00:45", Format(45*time.Minute, "15:04"))
	Equal(t, "00h 45m", Format(45*time.Minute, "15h 04m"))
}

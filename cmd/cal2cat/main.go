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

package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/abc-inc/cal2cat/cal"
	"github.com/abc-inc/cal2cat/cat"
	"github.com/abc-inc/cal2cat/event"
	"github.com/abc-inc/cal2cat/internal/duration"
	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

const defDurFmt = "hours"
const defTimeFmt = "2006-01-02 15:04"

//go:embed default.ini
var defIni []byte

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var rootCmd = &cobra.Command{
		Use:   "cal2cat",
		Short: "cal2cat parses iCalendar files and categorizes the events.",
		Run:   execute,

		Example: `  Time ranges can be specified as <n><unit>, where
  <n> is an integer and
  <unit> is any of:
    d   days
    w   weeks (7 days)
    cw  calendar weeks
    m   months (same day in another month)
    cm  calendar months
    y   years (same day in another year)
    cy  calendar years

  Calendar weeks, months and years are truncated to Monday, the first day of the
  month and the first day of the year, respectively.`,
	}

	rootCmd.Flags().StringP("end", "e", "-0cw", "end time")
	rootCmd.Flags().StringP("start", "s", "-1cw", "start time")
	cobra.CheckErr(rootCmd.Execute())
}

func execute(cmd *cobra.Command, args []string) {
	cfgPath, _ := xdg.ConfigFile("cal2booking/config.ini")
	cfg := readConfig(cfgPath)

	timeFmt := cfg.Section("settings").Key("timeFormat").MustString(defTimeFmt)
	durFmt := cfg.Section("settings").Key("durationFormat").MustString(defDurFmt)

	today := time.Now().UTC().Truncate(24 * time.Hour)
	fStart, _ := cmd.Flags().GetString("start")
	rangeStart := duration.Calc(today, fStart)
	fEnd, _ := cmd.Flags().GetString("end")
	rangeEnd := duration.Calc(today, fEnd)
	fmt.Printf("Categorizing events from %s until %s\n",
		rangeStart.Format(timeFmt), rangeEnd.Format(timeFmt))

	ms := []cat.Mapper{}
	sec := cfg.Section("mapping")
	for _, k := range sec.Keys() {
		ms = append(ms, cat.NewMapper(k.Name(), k.Value()))
	}

	ps := []string{}
	sec, _ = cfg.GetSection("calendars")
	for _, k := range sec.Keys() {
		ps = append(ps, k.Value())
	}

	es := cal.Load(ps...)
	es = es.Filter(event.NewRangeFilter(rangeStart, rangeEnd))
	es = es.Filter(event.NewEventFilter(es.Conflicts().Events()).Not())

	cs := cat.Map(es, ms)
	if len(cs) == 0 {
		return
	}

	for _, c := range cs {
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%s (%d events - %s)\n", c.Name, len(c.Events),
			duration.Format(c.Events.Duration(), durFmt))

		for _, e := range c.Events {
			fmt.Printf("%v %s (%s)\n",
				e.StartTime().Format(timeFmt), e.Summary(),
				duration.Format(e.Duration(), durFmt))
		}
	}
}

func readConfig(cfgPath string) *ini.File {
	cfg, err := ini.LooseLoad("cal2booking.ini", cfgPath, defIni)
	if err != nil {
		log.Fatalf("cannot read config file from %s: %v", cfgPath, err)
	}

	cfgFile, err := os.Create(cfgPath)
	if err != nil {
		log.Fatalf("cannot create config file %s: %v", cfgPath, err)
	}
	defer cfgFile.Close()
	_, _ = cfg.WriteTo(cfgFile)

	return cfg
}

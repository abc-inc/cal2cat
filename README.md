# cal2cat

*cal2cat* parses one or more iCalendar files (ics) and categorizes the events.

## Configuration

```ini
[settings]
timeFormat=2006-01-02 15:04
durationFormat=15:04 ; possible values: "minutes", "hours", "15:04", "15h 04m", etc.

[mapping]
; Info Meetings
*All Staff*=Info Meeting
Community of Practice=Info Meeting
; Projects
^(ABC|ABD)=Project ABC
*XYZ*=Project XYZ
; Other
^Vacation=Vacation
^=OTHER ; all calendar entries not matching any of the previous categories are classified as OTHER

[calendars]
work=https://outlook.office365.com/owa/calendar/.../calendar.ics
main=https://outlook.office365.com/owa/calendar/.../calendar.ics
```

### Section `mapping`

This section contains an arbitrary amount of lines in the form:
`<summary_pattern>=<category>`.
The left side matches the summary of a calendar entry and the right side is the category. Wildcard patterns or regular expressions can be used as follows:

- pattern starts with '`^`': regular expressions are enabled
- pattern contains any of "`*?`": wildcards can be used ('`*`' denotes any
number of characters, '`?`' denotes a single character)
- otherwise, the calendar summary must begin with the pattern

#### Examples

- `Conference=Training`: if a calendar entry summary begins with `Conference`,
then it is categorized as `Training`
- `*JF*=Info Meeting`: if a calendar entry summary contains `JF`,
then it is categorized as `Info Meeting`
- `^.*Problem|Incident=Troubleshooting`: if a calendar entry summary contains
`Problem` or `Incident`, then it is categorized as `Troubleshooting`

### Section `calendars`

This section can contain multiple paths or URLs to calendars.

## Usage

The following command categorizes all calender entries in the previous calendar
week (start: `-1cw`, end: `-0cw`):

```shell
$ cal2cat --start -1cw --end -0cw

Categorizing events from 24.05.2021 00:00 until 31.05.2021 00:00
--------------------------------------------------------------------------------
Info Meeting (2 events - 03:30)
24.05.2021 15:30 All Staff Meeting (00:30)
28.05.2021 10:15 Community of Practice: Security (03:00)
--------------------------------------------------------------------------------
Project ABC (3 events - 04:00)
24.05.2021 08:30 ABC: Testing-Kickoff (01:30)
24.05.2021 11:30 ABC (01:00)
28.05.2021 14:30 ABC: Review (01:30)
--------------------------------------------------------------------------------
OTHER (1 events - 04:15)
24.05.2021 14:00 Coffee Break (00:30)
--------------------------------------------------------------------------------
Vacation (3 events - 24:00)
25.05.2021 08:00 Vacation (08:00)
26.05.2021 08:00 Vacation (08:00)
27.05.2021 08:00 Vacation (08:00)
```

package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Time struct {
	Layout string
	Zone   string
}

var FirstDayMonday bool
var TimeFormats = []string{"1/2/2006", "1/2/2006 15:4:5", "2006-1-2 15:4:5", "2006-1-2 15:4", "2006-1-2", "1-2", "15:4:5", "15:4", "15", "15:4:5 Jan 2, 2006 MST", "2006-01-02 15:04:05.999999999 -0700 MST"}

type Now struct {
	time.Time
}

func NewTime(t time.Time) *Now {
	return &Now{t}
}

func BeginningOfMinute() time.Time {
	return NewTime(time.Now()).BeginningOfMinute()
}

func BeginningOfHour() time.Time {
	return NewTime(time.Now()).BeginningOfHour()
}

func BeginningOfDay() time.Time {
	return NewTime(time.Now()).BeginningOfDay()
}

func BeginningOfWeek() time.Time {
	return NewTime(time.Now()).BeginningOfWeek()
}

func BeginningOfMonth() time.Time {
	return NewTime(time.Now()).BeginningOfMonth()
}

func BeginningOfQuarter() time.Time {
	return NewTime(time.Now()).BeginningOfQuarter()
}

func BeginningOfYear() time.Time {
	return NewTime(time.Now()).BeginningOfYear()
}

func EndOfMinute() time.Time {
	return NewTime(time.Now()).EndOfMinute()
}

func EndOfHour() time.Time {
	return NewTime(time.Now()).EndOfHour()
}

func EndOfDay() time.Time {
	return NewTime(time.Now()).EndOfDay()
}

func EndOfWeek() time.Time {
	return NewTime(time.Now()).EndOfWeek()
}

func EndOfMonth() time.Time {
	return NewTime(time.Now()).EndOfMonth()
}

func EndOfQuarter() time.Time {
	return NewTime(time.Now()).EndOfQuarter()
}

func EndOfYear() time.Time {
	return NewTime(time.Now()).EndOfYear()
}

func Monday() time.Time {
	return NewTime(time.Now()).Monday()
}

func Sunday() time.Time {
	return NewTime(time.Now()).Sunday()
}

func EndOfSunday() time.Time {
	return NewTime(time.Now()).EndOfSunday()
}

func Parse(strs ...string) (time.Time, error) {
	return NewTime(time.Now()).Parse(strs...)
}

func MustParse(strs ...string) time.Time {
	return NewTime(time.Now()).MustParse(strs...)
}

func TimeBetween(time1, time2 string) bool {
	return NewTime(time.Now()).Between(time1, time2)
}
func UnixToStr(n int64, layouts ...string) string {
	if n == 0 {
		return ""
	}
	layout := "2006-01-02 15:04:05"
	if len(layouts) > 0 && layouts[0] != "" {
		layout = layouts[0]
	}
	var ss, ns string
	s := strconv.FormatInt(n, 10)
	l := len(s)
	if l > 10 {
		ss = s[:10]
		ns = s[10:]
		fillLen := 9 - len(ns)
		ns = ns + strings.Repeat("0", fillLen)
	} else {
		ss = s[:10]
	}
	si, _ := strconv.ParseInt(ss, 10, 64)
	ni, _ := strconv.ParseInt(ns, 10, 64)
	tm := time.Unix(si, ni)

	return tm.Format(layout)
}

func StrToUnix(s string) int64 {
	if s == "" {
		return 0
	}
	Layout := "2006-01-02 15:04:05"
	Zone := "+08:00"
	ti, _ := time.Parse(Layout+" "+Zone, s)
	return ti.Unix()
}
func StrToUnixNano(layout, s string) int64 {
	if s == "" {
		return 0
	}
	loc, _ := time.LoadLocation("Local") //获取本地时区
	ti, _ := time.ParseInLocation(layout, s, loc)
	return ti.UnixNano()
}
func IntToTime(ti int64, denominator int) time.Time {
	if denominator == 0 {
		denominator = 1
	}
	return time.Unix(ti/int64(denominator), 0)
}

func (now *Now) BeginningOfMinute() time.Time {
	return now.Truncate(time.Minute)
}

func (now *Now) BeginningOfHour() time.Time {
	return now.Truncate(time.Hour)
}

func (now *Now) BeginningOfDay() time.Time {
	d := time.Duration(-now.Hour()) * time.Hour
	return now.BeginningOfHour().Add(d)
}

func (now *Now) BeginningOfWeek() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if FirstDayMonday {
		if weekday == 0 {
			weekday = 7
		}
		weekday = weekday - 1
	}
	d := time.Duration(-weekday) * 24 * time.Hour
	return t.Add(d)
}

func (now *Now) BeginningOfMonth() time.Time {
	t := now.BeginningOfDay()
	d := time.Duration(-int(t.Day())+1) * 24 * time.Hour
	return t.Add(d)
}

func (now *Now) BeginningOfQuarter() time.Time {
	month := now.BeginningOfMonth()
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

func (now *Now) BeginningOfYear() time.Time {
	t := now.BeginningOfDay()
	d := time.Duration(-int(t.YearDay())+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func (now *Now) EndOfMinute() time.Time {
	return now.BeginningOfMinute().Add(time.Minute - time.Nanosecond)
}

func (now *Now) EndOfHour() time.Time {
	return now.BeginningOfHour().Add(time.Hour - time.Nanosecond)
}

func (now *Now) EndOfDay() time.Time {
	return now.BeginningOfDay().Add(24*time.Hour - time.Nanosecond)
}

func (now *Now) EndOfWeek() time.Time {
	return now.BeginningOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

func (now *Now) EndOfMonth() time.Time {
	return now.BeginningOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func (now *Now) EndOfQuarter() time.Time {
	return now.BeginningOfQuarter().AddDate(0, 3, 0).Add(-time.Nanosecond)
}

func (now *Now) EndOfYear() time.Time {
	return now.BeginningOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

func (now *Now) Monday() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := time.Duration(-weekday+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func (now *Now) Sunday() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if weekday == 0 {
		return t
	} else {
		d := time.Duration(7-weekday) * 24 * time.Hour
		return t.Truncate(time.Hour).Add(d)
	}
}

func (now *Now) EndOfSunday() time.Time {
	return now.Sunday().Add(24*time.Hour - time.Nanosecond)
}

func parseWithFormat(str string) (t time.Time, err error) {
	for _, format := range TimeFormats {
		t, err = time.Parse(format, str)
		if err == nil {
			return
		}
	}
	err = errors.New("Can't parse string as time: " + str)
	return
}

func (now *Now) Parse(strs ...string) (t time.Time, err error) {
	var setCurrentTime bool
	parseTime := []int{}
	currentTime := []int{now.Second(), now.Minute(), now.Hour(), now.Day(), int(now.Month()), now.Year()}
	currentLocation := now.Location()

	for _, str := range strs {
		onlyTime := regexp.MustCompile(`^\s*\d+(:\d+)*\s*$`).MatchString(str) // match 15:04:05, 15

		t, err = parseWithFormat(str)
		location := t.Location()
		if location.String() == "UTC" {
			location = currentLocation
		}

		if err == nil {
			parseTime = []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
			onlyTime = onlyTime && (parseTime[3] == 1) && (parseTime[4] == 1)

			for i, v := range parseTime {
				// Don't reset hour, minute, second if it is a time only string
				if onlyTime && i <= 2 {
					continue
				}

				// Fill up missed information with current time
				if v == 0 {
					if setCurrentTime {
						parseTime[i] = currentTime[i]
					}
				} else {
					setCurrentTime = true
				}

				// Default day and month is 1, fill up it if missing it
				if onlyTime {
					if i == 3 || i == 4 {
						parseTime[i] = currentTime[i]
						continue
					}
				}
			}
		}

		if len(parseTime) > 0 {
			t = time.Date(parseTime[5], time.Month(parseTime[4]), parseTime[3], parseTime[2], parseTime[1], parseTime[0], 0, location)
			currentTime = []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
		}
	}
	return
}

func (now *Now) MustParse(strs ...string) (t time.Time) {
	t, err := now.Parse(strs...)
	if err != nil {
		panic(err)
	}
	return t
}

func (now *Now) Between(time1, time2 string) bool {
	restime := now.MustParse(time1)
	restime2 := now.MustParse(time2)
	return now.After(restime) && now.Before(restime2)
}

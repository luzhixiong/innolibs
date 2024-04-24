package utils

import (
	"fmt"
	"time"
)

const (
	year      = "2006"
	month     = "01"
	day       = "02"
	hour      = "15"
	minute    = "04"
	second    = "05"
	date      = "2006-01-02"
	date1     = "20060102"
	datetime  = "2006-01-02 15:04:05"
	datetime2 = "20060102150405"
	datetime3 = "2006010215"
)

func Unixtime() uint {
	return uint(time.Now().Unix())
}

func MilliUnixtime() uint64 {
	return uint64(time.Now().UnixNano() / 1e6)
}

func NanoUnixtime() uint {
	return uint(time.Now().UnixNano())
}

func format(key string) string {
	return time.Now().Format(key)
}

func Year() string {
	return format(year)
}

func Month() string {
	return format(month)
}

func Day() string {
	return format(day)
}

func Hour() string {
	return format(hour)
}

func Minute() string {
	return format(minute)
}

func Second() string {
	return format(second)
}

func Date() string {
	return format(date)
}

func Date1() string {
	return format(date1)
}

func DateTime() string {
	return format(datetime)
}

func DateTime2() string {
	return format(datetime2)
}

func DateTime3() string {
	return format(datetime3)
}

func StartOfDay(t time.Time) int64 {
	y, m, d := t.Date()
	t1 := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	return t1.Unix()
}

func StartOfWeek(t time.Time) int64 {
	weekDay := t.Weekday()
	diffDay := weekDay - 1
	if weekDay == time.Sunday {
		diffDay = 6
	}
	zero := StartOfDay(t)
	return zero - int64(diffDay)*86400
}

func StartOfMonth(t time.Time) int64 {
	_, _, d := t.Date()
	diffDay := d - 1
	zero := StartOfDay(t)
	return zero - int64(diffDay)*86400
}

func DateFormat(timestamp int64) string {
	return DateFormatWithLoc(timestamp, nil)
}

func DateFormatWithLoc(timestamp int64, loc *time.Location) string {
	t := time.Unix(timestamp, 0)
	if loc != nil {
		t = t.In(loc)
	}
	return t.Format(date)
}

func DateTimeFormat(timestamp int64) string {
	return DateTimeFormatWithLoc(timestamp, nil)
}

func DateTimeFormatWithLoc(timestamp int64, loc *time.Location) string {
	t := time.Unix(timestamp, 0)
	if loc != nil {
		t = t.In(loc)
	}
	return t.Format(datetime)
}

func parseDate(str string) time.Time {
	t, _ := time.ParseInLocation(date, str, time.Local)
	return t
}

func parseDateTime(str string) time.Time {
	t, _ := time.ParseInLocation(datetime, str, time.Local)
	return t
}

func parseDateTimeRfc3339(str string) time.Time {
	t, _ := time.ParseInLocation(time.RFC3339, str, time.Local)
	return t
}

func ParseDate2Ts(str string) uint {
	t := parseDate(str)
	return uint(t.Unix())
}

func ParseDateTime2Ts(str string) uint {
	t := parseDateTime(str)
	return uint(t.Unix())
}

func ParseDateTime2TsRfc3339(str string) uint {
	t := parseDateTimeRfc3339(str)
	return uint(t.Unix())
}

func AddDate(str string, days int) string {
	t := parseDate(str)
	return t.AddDate(0, 0, days).Format(date)
}

func AddDate1(str string, days int) string {
	t := parseDate(str)
	return t.AddDate(0, 0, days).Format(date1)
}

// 获取指定时间的周开始时间和结束时间
func GetDayStartWeek(mytime time.Time, weekStartDay, weekInterval int) (time.Time, time.Time, error) {
	thisYear := mytime.Year()
	var s, e time.Time
	for i := 1; i < 106; i++ {
		if weekStartDay == 0 {
			s, e, _ = weekStartEnd(thisYear, i, weekInterval)
		} else {
			s, e, _ = weekStartEnd2(thisYear, i, weekInterval)
		}
		//fmt.Println(s,e)
		if mytime.Format("20060102") >= s.Format("20060102") && mytime.Format("20060102") <= e.Format("20060102") {
			return s, e, nil
		}
	}
	thisYear--
	for i := 1; i < 106; i++ {
		if weekStartDay == 0 {
			s, e, _ = weekStartEnd(thisYear, i, weekInterval)
		} else {
			s, e, _ = weekStartEnd2(thisYear, i, weekInterval)
		}
		//fmt.Println(s,e)
		if mytime.Format("20060102") >= s.Format("20060102") && mytime.Format("20060102") <= e.Format("20060102") {
			return s, e, nil
		}
	}

	return time.Time{}, time.Time{}, nil
}

// 从周日开始
func weekStartEnd(year, week, weekInterval int) (time.Time, time.Time, error) {
	if week < 1 {
		return time.Time{}, time.Time{}, fmt.Errorf("week number must be greater than 0")
	}

	// 计算年份的1月1日
	firstDayOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)

	// 计算1月1日是周几 (0: 周日, 1: 周一, ..., 6: 周六)
	dayOfWeek := int(firstDayOfYear.Weekday())

	// 计算第一周的起始日期
	var daysToFirstWeek int
	if dayOfWeek == 0 {
		daysToFirstWeek = 1
	} else {
		tmp := 4 - dayOfWeek
		if tmp == 4 {
			dayOfWeek = -3
		} else {
			daysToFirstWeek = 8 - dayOfWeek
		}
		//daysToFirstWeek = 8 - dayOfWeek
	}
	firstWeekStart := firstDayOfYear.AddDate(0, 0, daysToFirstWeek-1)

	// 计算指定周的起始和结束日期
	weekStart := firstWeekStart.AddDate(0, 0, weekInterval*7*(week-1))
	weekEnd := weekStart.AddDate(0, 0, weekInterval*7-1)

	return weekStart, weekEnd, nil
}

// 从周一开始
func weekStartEnd2(year, week, weekInterval int) (time.Time, time.Time, error) {
	if week < 1 {
		return time.Time{}, time.Time{}, fmt.Errorf("week number must be greater than 0")
	}

	// 计算年份的1月1日
	firstDayOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)

	// 计算1月1日是周几 (0: 周日, 1: 周一, ..., 6: 周六)
	dayOfWeek := int(firstDayOfYear.Weekday())

	//fmt.Println("dayOfWeek:", dayOfWeek)

	// 计算第一周的起始日期
	var daysToFirstWeek int
	if dayOfWeek == 1 {
		daysToFirstWeek = 1
	} else {
		daysToFirstWeek = 2 - dayOfWeek
		//daysToFirstWeek = dayOfWeek - 1
	}
	firstWeekStart := firstDayOfYear.AddDate(0, 0, daysToFirstWeek-1)

	// 计算指定周的起始和结束日期
	weekStart := firstWeekStart.AddDate(0, 0, weekInterval*7*(week-1))
	weekEnd := weekStart.AddDate(0, 0, weekInterval*7-1)

	return weekStart, weekEnd, nil
}

func GetDayStartWeekYear(mytime time.Time, weekStartDay, weekInterval int, statDay uint) (start time.Time, end time.Time, err error) {
	today := mytime.Format("2006-01-02")
	endtime := ParseDate2Ts(today)
	startime := endtime - 86400*statDay
	startimetmp := time.Unix(int64(startime), 0)
	sTime, _, _ := GetDayStartWeek(startimetmp, weekStartDay, weekInterval)
	for i := 1; i < 118; i++ {
		if weekStartDay == 0 {
			s, e, _ := weekStartEnd(sTime.Year(), i, weekInterval)
			if s.Unix() <= int64(endtime) && e.Unix() >= int64(endtime) {
				start = s
				end = e
				break
			}
		} else {
			s, e, _ := weekStartEnd2(sTime.Year(), i, weekInterval)
			if s.Unix() <= int64(endtime) && e.Unix() >= int64(endtime) {
				start = s
				end = e
				break
			}
		}
	}
	return
}

func GetPeriodDate(startDate, nowDate string, periodDay int) (start time.Time, end time.Time) {
	if startDate == "" {
		startDate = "2024-01-01"
	}
	if periodDay <= 0 {
		periodDay = 14
	}
	if nowDate == "" {
		nowDate = Date()
	}
	start = parseDate(startDate)
	nowTime := parseDate(nowDate)
	for {
		end = start.AddDate(0, 0, periodDay-1)
		if nowTime.Unix() >= start.Unix() && nowTime.Unix() <= end.Unix() {
			break
		}
		//fmt.Println(nowTime.Format("2006-01-02"), start.Format("2006-01-02"), end.Format("2006-01-02"))
		start = end.AddDate(0, 0, 1)
	}

	return
}

package common

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// FmtDateTimeSecondStr - Date Time format used
	// for file names and directory names
	FmtDateTimeSecondStr = "20060102150405"
	// FmtDateTimeNanoSecondStr - Custom Date Time Format
	FmtDateTimeNanoSecondStr = "2006-01-02 15:04:05.000000000"
	// FmtDateTimeSecText - Custom Date Time Format
	FmtDateTimeSecText = "2006-01-02 15:04:05"
	// FmtDateTimeEverything - Custom Date Time Format showing virtually
	// all elements of a date time string.
	FmtDateTimeEverything = "Monday January 2, 2006 15:04:05.000000000 -0700 MST"
	// MicroSecondNanoseconds - Number of Nanoseconds in a Microsecond
	MicroSecondNanoseconds = int64(1000)
	// MilliSecondNanoseconds - Number of Nanoseconds in a MilliSecond
	MilliSecondNanoseconds = int64(1000 * 1000)
	// SecondNanoseconds - Number of Nanoseconds in a Second
	SecondNanoseconds = int64(1000 * 1000 * 1000)
	// MinuteNanoSeconds - Number of Nanoseconds in a minute
	MinuteNanoSeconds = int64(1000 * 1000 * 1000 * 60)
	// HourNanoSeconds - Number of Nanoseconds in an hour
	HourNanoSeconds = int64(1000 * 1000 * 1000 * 60 * 60)
	// DayNanoSeconds - Number of Nanoseconds in a 24-hour day
	DayNanoSeconds = int64(1000 * 1000 * 1000 * 60 * 60 * 24)
	// YearNanoSeconds - Number of Nanoseconds in a 365-day year
	YearNanoSeconds = int64(1000 * 1000 * 1000 * 60 * 60 * 24 * 365)
)

// ElapsedDuration - holds elements of
// time duration
type ElapsedDuration struct {
	Years        int64
	Days         int64
	Hours        int64
	Minutes      int64
	Seconds      int64
	MilliSeconds int64
	MicroSeconds int64
	NanoSeconds  int64
	// NanosecStr - Example: 2-Days 13-Hours 26-Minutes 46-Seconds 864197832-Nanoseconds
	NanosecStr string
	// DurationStr - Example: 2-Days 13-Hours 26-Minutes 46-Seconds 864-Milliseconds 197-Microseconds 832-Nanoseconds
	DurationStr string
	// DefaultStr - Example: 61h26m46.864197832s - format provided by 'go' library
	DefaultStr string
}

// DateTimeUtility - struct used to export
// Date Time Management methods.
type DateTimeUtility struct {
	TimeIn     time.Time
	TimeOut    time.Time
	TimeStart  time.Time
	TimeEnd    time.Time
	Duration   time.Time
	TimeInStr  string
	TimeOutStr string
	TimeFmtStr string
}

// GetDateTimeStrNowLocal - Gets current
// local time and formats it as a date time
// string
func (dt DateTimeUtility) GetDateTimeStrNowLocal() string {

	return dt.GetDateTimeStr(time.Now().Local())

}

// GetDateTimeStr - Returns a date time string
// in the format 20170427211307
func (dt DateTimeUtility) GetDateTimeStr(t time.Time) string {

	// Time Format down to the second
	return t.Format(FmtDateTimeSecondStr)

}

// GetDateTimeSecText - Returns formatted
// date time with seconds for display,
// 2006-01-02 15:04:05.
func (dt DateTimeUtility) GetDateTimeSecText(t time.Time) string {
	// Time Display Format with seconds
	return t.Format(FmtDateTimeSecText)
}

// GetDateTimeNanoSecText - Returns formated
// date time string with nanoseconds
// 2006-01-02 15:04:05.000000000.
func (dt DateTimeUtility) GetDateTimeNanoSecText(t time.Time) string {
	// Time Format down to the nanosecond
	return t.Format(FmtDateTimeNanoSecondStr)
}

// GetDateTimeEverything - Receives a time value and formats as
// a date time string in the format:
//  Saturday April 29, 2017 19:54:30.123456489 -0500 CDT
func (dt DateTimeUtility) GetDateTimeEverything(t time.Time) string {
	return t.Format(FmtDateTimeEverything)
}

// GetDateTimeCustomFmt - Returns time string
// formatted according to passed in format
// string. Example format string:
// 'Monday January 2, 2006 15:04:05.000000000 -0700 MST'
func (dt DateTimeUtility) GetDateTimeCustomFmt(t time.Time, fmt string) string {
	return t.Format(fmt)
}

// GetDuration - Returns a time.Duration structure defining the duration between
// input parameters startTime and endTime
func (dt DateTimeUtility) GetDuration(startTime time.Time, endTime time.Time) (time.Duration, error) {

	def := time.Duration(0)

	if startTime.Equal(endTime) {
		return def, nil
	}

	if endTime.Before(startTime) {
		return def, errors.New("DateTimeUtility.GetDuration() Error: endTime less than startTime")
	}

	return endTime.Sub(startTime), nil
}

// GetDurationBreakDown - Receives a Duration type
// and returns a breakdown of duration by years,
// days, hours, minutes, seconds, milliseconds,
// microseconds and nanoseconds.
func (dt DateTimeUtility) GetDurationBreakDown(d time.Duration) ElapsedDuration {
	str := ""
	ed := ElapsedDuration{}
	ed.DefaultStr = fmt.Sprintf("%v", d)
	firstEle := false
	rd := int64(d)

	if rd >= YearNanoSeconds {
		ed.Years = rd / YearNanoSeconds
		rd -= YearNanoSeconds * ed.Years
	}

	if ed.Years > 0 {
		str = fmt.Sprintf("%v-Years ", ed.Years)
		firstEle = true
	}

	if rd >= DayNanoSeconds {
		ed.Days = rd / DayNanoSeconds
		rd -= DayNanoSeconds * ed.Days
	}

	if ed.Days > 0 || firstEle {
		str += fmt.Sprintf("%v-Days ", ed.Days)
		firstEle = true
	}

	if rd >= HourNanoSeconds {
		ed.Hours = rd / HourNanoSeconds
		rd -= HourNanoSeconds * ed.Hours
	}

	if ed.Hours > 0 || firstEle {
		str += fmt.Sprintf("%v-Hours ", ed.Hours)
		firstEle = true
	}

	if rd >= MinuteNanoSeconds {
		ed.Minutes = rd / MinuteNanoSeconds
		rd -= MinuteNanoSeconds * ed.Minutes
	}

	if ed.Minutes > 0 || firstEle {
		str += fmt.Sprintf("%v-Minutes ", ed.Minutes)
		firstEle = true
	}

	if rd >= SecondNanoseconds {
		ed.Seconds = rd / SecondNanoseconds
		rd -= SecondNanoseconds * ed.Seconds
	}

	if ed.Seconds > 0 || firstEle {
		str += fmt.Sprintf("%v-Seconds ", ed.Seconds)
		firstEle = true
	}

	ed.NanosecStr = str + fmt.Sprintf("%v-Nanoseconds", rd)

	if rd >= MilliSecondNanoseconds {
		ed.MilliSeconds = rd / MilliSecondNanoseconds
		rd -= MilliSecondNanoseconds * ed.MilliSeconds
	}

	if ed.MilliSeconds > 0 || firstEle {
		str += fmt.Sprintf("%v-Milliseconds ", ed.MilliSeconds)
		firstEle = true
	}

	if rd >= MicroSecondNanoseconds {
		ed.MicroSeconds = rd / MicroSecondNanoseconds
		rd -= MicroSecondNanoseconds * ed.MicroSeconds
	}

	if ed.MicroSeconds > 0 || firstEle {
		str += fmt.Sprintf("%v-Microseconds ", ed.MicroSeconds)
		firstEle = true
	}

	ed.NanoSeconds = rd

	if ed.NanoSeconds > 0 || firstEle {
		str += fmt.Sprintf("%v-Nanoseconds", ed.NanoSeconds)
		firstEle = true
	}

	ed.DurationStr = strings.TrimRight(str, " ")

	return ed

}

// GetElapsedTime - calculates the elapsed time
// between input parameters startTime and endTime.
// The result is returned in an ElapsedDuration
// structure.
func (dt DateTimeUtility) GetElapsedTime(startTime time.Time, endTime time.Time) (ElapsedDuration, error) {

	ed := ElapsedDuration{}

	dur, err := dt.GetDuration(startTime, endTime)

	if err != nil {
		s := "DateTimeUtility-GetElapsedTime Error: " + err.Error()

		return ed, errors.New(s)
	}

	return dt.GetDurationBreakDown(dur), nil

}

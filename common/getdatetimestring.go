package common

import "time"

const (
	// DateTimeFmtSecondStr - Date Time format used
	// for file names and directory names
	dateTimeFmtSecondStr   = "20060102150405"
	dateTimeFmtNanoSecText = "2006-01-02 15:04:05.000000000"
	dateTimeFmtSecText     = "2006-01-02 15:04:05"
)

// GetDateTimeStrNowLocal - Gets current
// local time and formats it as a date time
// string
func GetDateTimeStrNowLocal() string {

	return GetDateTimeStr(time.Now().Local())

}

// GetDateTimeStr - Returns a date time string
// in the format 20170427211307
func GetDateTimeStr(t time.Time) string {

	// Time Format down to the second
	return t.Format(dateTimeFmtSecondStr)

}

// GetDateTimeSecText - Returns formatted
// date time with seconds for display,
// 2006-01-02 15:04:05.
func GetDateTimeSecText(t time.Time) string {
	// Time Display Format with seconds
	return t.Format(dateTimeFmtSecText)
}

// GetDateTimeNanoSecText - Returns formated
// date time string with nanoseconds
// 2006-01-02 15:04:05.000000000.
func GetDateTimeNanoSecText(t time.Time) string {
	// Time Format down to the nanosecond
	return t.Format(dateTimeFmtNanoSecText)
}

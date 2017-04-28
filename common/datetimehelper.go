package common

import "time"

const (
	// DateTimeFmtSecondStr - Date Time format used
	// for file names and directory names
	DateTimeFmtSecondStr = "20060102150405"
)

// GetDateTimeStrNowLocal - Gets current
// local time and formats it as a date time
// string
func GetDateTimeStrNowLocal() string {

	return GetDateTimeStr(time.Now().Local())

}

// GetDateTimeStr - Returns a date time string
// in the format 20170427211307.123
func GetDateTimeStr(t time.Time) string {

	// Time Format down to the second
	return t.Format(DateTimeFmtSecondStr)

}


package common

import (
	"errors"
	"fmt"
	"time"
	"strings"
)

// NOTE: See https://golang.org/pkg/time/#LoadLocation
// and https://www.iana.org/time-zones to ensure that
// the IANA Time Zone Database is properly configured
// on your system.
const (
	// TzUsEast - USA Eastern Time Zone
	// IANA database identifier
	TzUsEast = "America/New_York"

	// TzUsCentral - USA Central Time Zone
	// IANA database identifier
	TzUsCentral = "America/Chicago"

	// TzUsMountain - USA Mountain Time Zone
	// IANA database identifier
	TzUsMountain = "America/Denver"

	// TzUsPacific - USA Pacific Time Zone
	// IANA database identifier
	TzUsPacific = "America/Los_Angeles"

	// TzUsHawaii - USA Hawaiian Time Zone
	// IANA database identifier
	TzUsHawaii = "Pacific/Honolulu"

	// tzUTC - UTC Time Zone IANA database
	// identifier
	TzUTC = "Zulu"

	NeutralDateFmt = "2006-01-02 15:04:05.000000000"
)

// TimeZoneUtility - Time Zone Data and Methods
type TimeZoneUtility struct {
	TimeIn     time.Time
	TimeInLoc  *time.Location
	TimeOut    time.Time
	TimeOutLoc *time.Location
	TimeUTC    time.Time
}


// ConvertTz - Convert Time from existing time zone to targetTZone.
// The results are stored in the TimeZoneUtility data structure.
// Input Parameters
// tIn - time.Time initial time
// tInIanaTz - The IANA time zone string associated with tIn parameter
// targetIanaTz - The IANA time zone to which tIn will be converted
// Output Values are returned in the tzu (TimeZoneUtility) data structure
func (tzu *TimeZoneUtility) ConvertTz(tIn time.Time, tInIanaTz string, targetIanaTz string) error {

	if !tzu.IsIanaTzValid(tInIanaTz) {
		return errors.New("TimeZoneUtility:ConvertTz() Error: tInIanaTz is INVALID")
	}

	if !tzu.IsIanaTzValid(targetIanaTz) {
		return errors.New("TimeZoneUtility:ConvertTz() Error: targetIanaTz is INVALID")
	}

	if tIn.IsZero() {
		return errors.New("TimeZoneUtility:ConvertTz() Error: Input parameter tIn is zero and INVALID")
	}

	tInNoLocStr := tzu.TimeWithoutTimeZone(tIn)

	tzIn, err := time.LoadLocation(tInIanaTz)

	if err != nil {
		return fmt.Errorf("TimeZoneUtility:ConvertTz() - Error Loading Input IANA Time Zone 'tInIanaTz', %v. Errors: %v ", tInIanaTz, err.Error())
	}

	tzOut, err := time.LoadLocation(targetIanaTz)

	if err != nil {
		return fmt.Errorf("TimeZoneUtility:ConvertTz() - Error Loading Target IANA Time Zone 'targetIanaTz', %v. Errors: %v ", targetIanaTz, err.Error())
	}

	tInWithTz, err := time.ParseInLocation(NeutralDateFmt, tInNoLocStr, tzIn)

	if err != nil {
		return fmt.Errorf("TimeZoneUtility:ConvertTz() - Error ParseInLocation tIn with Time Zone , %v. Errors: %v ", tInNoLocStr, err.Error())
	}

	tzu.SetTimeIn(tInWithTz)

	tzu.SetTimeOut(tInWithTz.In(tzOut))

	tzu.SetUTCTime(tInWithTz)

	return nil
}

func(tzu *TimeZoneUtility) IsIanaTzValid(tzIana string) bool {

	if tzIana == "" {
		return false
	}

	if strings.ToLower(tzIana) == "local" {
		return false
	}

	_, err := time.LoadLocation(tzIana)

	if err != nil {
		return false
	}

	return true

}

// SetTimeIn - Assigns value to field 'TimeIn'
func (tzu *TimeZoneUtility) SetTimeIn(tIn time.Time) {
	tzu.TimeIn = tIn
	tzu.TimeInLoc = tIn.Location()
}

// SetTimeOut - Assigns value to field 'TimeOut'
func (tzu *TimeZoneUtility) SetTimeOut(tOut time.Time) {
	tzu.TimeOut = tOut
	tzu.TimeOutLoc = tOut.Location()
}

// SetUTCTime - Assigns UTC Time to field 'TimeUTC'
func (tzu *TimeZoneUtility) SetUTCTime(t time.Time) {

	tzu.TimeUTC = t.UTC()
}

// TimeWithoutTimeZone - Returns a Time String containing
// NO time zone.
func (tzu *TimeZoneUtility) TimeWithoutTimeZone(t time.Time) string {
	return t.Format(NeutralDateFmt)
}

package common

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeZoneUtilConvertTz(t *testing.T) {
	tstr := "04/29/2017 19:54:30 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05 -0700 MST"
	ianaPacificTz := "America/Los_Angeles"
	ianaCentralTz := "America/Chicago"
	tIn, _ := time.Parse(fmtstr, tstr)
	tzu := TimeZoneUtility{}
	tzu.ConvertTz(tIn, ianaCentralTz, ianaPacificTz)

	exTIn := "2017-04-29 19:54:30 -0500 CDT"
	actTIn := fmt.Sprintf("%v", tzu.TimeIn)
	if actTIn != exTIn {
		t.Error(fmt.Sprintf("Expected tzu.TimeIn %v, got ", exTIn), actTIn)
	}

	exTInLoc := "America/Chicago"
	actTInLoc := fmt.Sprintf("%v", tzu.TimeInLoc)
	if actTInLoc != exTInLoc {
		t.Error(fmt.Sprintf("Expected tzu.TimeInLoc %v, got", exTInLoc), actTInLoc)
	}

	exTOut := "2017-04-29 17:54:30 -0700 PDT"
	actTOut := fmt.Sprintf("%v", tzu.TimeOut)
	if actTOut != exTOut {
		t.Error(fmt.Sprintf("Expected tzu.TimeOut %v, got", exTOut), actTOut)
	}

	exTOutLoc := "America/Los_Angeles"
	actTOutLoc := fmt.Sprintf("%v", tzu.TimeOutLoc)

	if actTOutLoc != exTOutLoc {
		t.Error(fmt.Sprintf("Expected tzu.TimeOutLoc %v, got", exTOutLoc), actTOutLoc)
	}

	exUTC := "2017-04-30 00:54:30 +0000 UTC"
	actUTC := fmt.Sprintf("%v", tzu.TimeUTC)

	if exUTC != actUTC {
		t.Error(fmt.Sprintf("Expected tzu.TimeUTC %v, got", exUTC), actUTC)
	}

}

func TestInvalidTzInConversion(t *testing.T) {
	tstr := "04/29/2017 19:54:30 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05 -0700 MST"
	ianaPacificTz := "America/Los_Angeles"
	// Invalid Central Time Zone
	ianaCentralTz := "AmericChicago"
	tIn, _ := time.Parse(fmtstr, tstr)
	tzu := TimeZoneUtility{}
	err := tzu.ConvertTz(tIn, ianaCentralTz, ianaPacificTz)

	if err == nil {
		t.Error("ConvertTz failed to detect INVALID Time In Time Zone. Got:", "err==nil")
	}

}

func TestInvalidTargetTzInConversion(t *testing.T) {
	tstr := "04/29/2017 19:54:30 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05 -0700 MST"
	// Invalid Target Iana Time Zone
	ianaPacificTz := time.Now().Location().String()
	ianaCentralTz := "America/Chicago"
	tIn, _ := time.Parse(fmtstr, tstr)
	tzu := TimeZoneUtility{}
	err := tzu.ConvertTz(tIn, ianaCentralTz, ianaPacificTz)

	if err == nil {
		t.Error("ConverTz() failed to detect INVALID Tartet Time Zone. Got: ", "err==nil")
	}

}

func TestLocalIsInvalidIanaTimeZone(t *testing.T) {
	tIn := time.Now()

	tzu := TimeZoneUtility{}

	result := tzu.IsIanaTzValid(tIn.Location().String())

	if result == true {
		t.Error("Expected Now() Location Location to be INVALID (false), got", result)
	}

}

func TestCDTIsValidIanaTimeZone(t *testing.T) {

	tzu := TimeZoneUtility{}

	result := tzu.IsIanaTzValid("America/Chicago")

	if result == false {
		t.Error("Expected 'America/Chicago' to be VALID IANA Time Zone 'true', got", result)
	}

}

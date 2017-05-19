package common

import (
	"fmt"
	"testing"
	"time"
)

func TestGetDateTimeStr(t *testing.T) {
	tstr := "04/29/2017 19:54:30 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05 -0700 MST"
	expected := "20170429195430"

	testTime, _ := time.Parse(fmtstr, tstr)

	result := DateTimeUtility{}.GetDateTimeStr(testTime)

	if result != expected {
		t.Error("Expected '20170429195430' got", result)
	}

}

func TestGetDateTimeSecText(t *testing.T) {

	tstr := "04/29/2017 19:54:30 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05 -0700 MST"
	expected := "2017-04-29 19:54:30"
	testTime, _ := time.Parse(fmtstr, tstr)
	result := DateTimeUtility{}.GetDateTimeSecText(testTime)

	if result != expected {
		t.Error("Expected '", expected, "' got", result)
	}

}

func TestGetDateTimeNanoSecText(t *testing.T) {
	tstr := "04/29/2017 19:54:30.123456489 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05.000000000 -0700 MST"
	expected := "2017-04-29 19:54:30.123456489"
	testTime, _ := time.Parse(fmtstr, tstr)
	dt := DateTimeUtility{}
	result := dt.GetDateTimeNanoSecText(testTime)

	if result != expected {
		t.Error("Expected '", expected, "' got", result)
	}

}

func TestGetDateTimeEverything(t *testing.T) {
	tstr := "04/29/2017 19:54:30.123456489 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05.000000000 -0700 MST"
	expected := "Saturday April 29, 2017 19:54:30.123456489 -0500 CDT"
	testTime, _ := time.Parse(fmtstr, tstr)
	dt := DateTimeUtility{}
	str := dt.GetDateTimeEverything(testTime)

	if str != expected {
		t.Error(fmt.Sprintf("Expected datetime: '%v', got", expected), str)
	}

}

func TestCustomDateTimeFormat(t *testing.T) {
	dt := DateTimeUtility{}
	tstr := "04/29/2017 19:54:30.123456489 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05.000000000 -0700 MST"
	expected := "Saturday April 29, 2017 19:54:30.123456489 -0500 CDT"
	testTime, _ := time.Parse(fmtstr, tstr)
	result := dt.GetDateTimeCustomFmt(testTime, FmtDateTimeEverything)

	if result != expected {
		t.Error(fmt.Sprintf("Expected: %v, got", expected), result)
	}

}

func TestSimpleDurationBreakdown(t *testing.T) {
	t1str := "04/28/2017 19:54:30 -0500 CDT"
	t2str := "04/30/2017 22:58:32 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05 -0700 MST"
	dt := DateTimeUtility{}
	t1, err := time.Parse(fmtstr, t1str)
	if err != nil {
		t.Error("Time Parse1 Error:", err.Error())
	}

	t2, err := time.Parse(fmtstr, t2str)
	if err != nil {
		t.Error("Time Parse2 Error:", err.Error())
	}

	dur, err := dt.GetDuration(t1, t2)
	if err != nil {
		t.Error("Get Duration Failed: ", err.Error())
	}

	ed := dt.GetDurationBreakDown(dur)

	ex1 := "2-Days 3-Hours 4-Minutes 2-Seconds 0-Milliseconds 0-Microseconds 0-Nanoseconds"

	if ed.DurationStr != ex1 {
		t.Error(fmt.Sprintf("Expected duration string of %v, got", ex1), ed.DurationStr)
	}
	// 2-Days 3-Hours 4-Minutes 2-Seconds 0-Milliseconds 0-Microseconds 0-Nanoseconds

	ex2 := "51h4m2s"

	if ed.DefaultStr != ex2 {
		t.Error(fmt.Sprintf("Expected default druation string: %v, got", ex2), ed.DefaultStr)
	}
}

func TestElapsedTimeBreakdown(t *testing.T) {
	tstr1 := "04/15/2017 19:54:30.123456489 -0500 CDT"
	tstr2 := "04/18/2017 09:21:16.987654321 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05.000000000 -0700 MST"
	t1, err1 := time.Parse(fmtstr, tstr1)

	if err1 != nil {
		t.Error("Error On Time Parse #1: ", err1.Error())
	}

	t2, err2 := time.Parse(fmtstr, tstr2)

	if err2 != nil {
		t.Error("Error On Time Parse #2: ", err2.Error())
	}

	dt := DateTimeUtility{}

	ed, err4 := dt.GetElapsedTime(t1, t2)
	if err4 != nil {
		t.Error("Error On GetElapsedTime: ", err4.Error())
	}

	ex1 := "2-Days 13-Hours 26-Minutes 46-Seconds 864-Milliseconds 197-Microseconds 832-Nanoseconds"

	if ed.DurationStr != ex1 {
		t.Error(fmt.Sprintf("Expected %v, got", ex1), ed.DurationStr)
	}

	ex2 := "61h26m46.864197832s"

	if ed.DefaultStr != ex2 {
		t.Error(fmt.Sprintf("Expected %v, got", ex2), ed.DefaultStr)
	}

	ex3 := "2-Days 13-Hours 26-Minutes 46-Seconds 864197832-Nanoseconds"

	if ex3 != ed.NanosecStr {
		t.Error(fmt.Sprintf("Expected %v, got", ex3), ed.NanosecStr)
	}

}

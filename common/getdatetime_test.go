package common

import (
	"testing"
	"time"
)

func TestGetDateTimeStr(t *testing.T) {
	tstr := "04/29/2017 19:54:30 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05 -0700 MST"
	expected := "20170429195430"

	testTime, _ := time.Parse(fmtstr, tstr)

	result := GetDateTimeStr(testTime)

	if result != expected {
		t.Error("Expected '20170429195430' got", result)
	}

}

func TestGetDateTimeSecText(t *testing.T) {

	tstr := "04/29/2017 19:54:30 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05 -0700 MST"
	expected := "2017-04-29 19:54:30"
	testTime, _ := time.Parse(fmtstr, tstr)
	result := GetDateTimeSecText(testTime)

	if result != expected {
		t.Error("Expected '", expected, "' got", result)
	}

}

func TestGetDateTimeNanoSecText(t *testing.T) {
	tstr := "04/29/2017 19:54:30.123456489 -0500 CDT"
	fmtstr := "01/02/2006 15:04:05.000000000 -0700 MST"
	expected := "2017-04-29 19:54:30.123456489"
	testTime, _ := time.Parse(fmtstr, tstr)
	result := GetDateTimeNanoSecText(testTime)

	if result != expected {
		t.Error("Expected '", expected, "' got", result)
	}

}

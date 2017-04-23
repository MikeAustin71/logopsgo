package common

import "testing"

func TestLogLevelString(t *testing.T) {
	var r LogLevel

	r = OPERROR

	var s string

	s = r.String()

	if s != "OPERROR" {
		t.Error("Expected string 'OPERROR' got ", s)
	}

}

func TestLogLevelValue(t *testing.T) {
	var r LogLevel

	var i int

	r = INFO

	if r != 3 {
		t.Error("Expected INFO value 3, got", i)
	}
}

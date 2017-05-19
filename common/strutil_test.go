package common

import (
	"fmt"
	"strings"
	"testing"
)

func TestStrCenterInStr(t *testing.T) {
	strToCntr := "1234567"
	fieldLen := 79
	exPadLen := 36
	exTotalLen := 43

	exStr := strings.Repeat(" ", exPadLen) + strToCntr
	su := StringUtility{}
	str, err := su.StrCenterInStr(strToCntr, fieldLen)
	if err != nil {
		t.Error("StrCenterInStr() generated error: ", err.Error())
	}

	l1 := su.StrGetRuneCnt(str)

	if l1 != exTotalLen {
		t.Error(fmt.Sprintf("Expected total str length '%v', got", exTotalLen), l1)
	}

	if str != exStr {
		t.Error(fmt.Sprintf("Strings did not match. Expected string '%v', got ", exStr), str)
	}

}

func TestStrGetRuneCnt(t *testing.T) {
	strToCnt := "1234567"
	exCnt := 7
	su := StringUtility{}
	l1 := su.StrGetRuneCnt(strToCnt)

	if l1 != exCnt {
		t.Error(fmt.Sprintf("Expected string character count of '%v', got", exCnt), l1)
	}

}

func TestSTrGetCharCnt(t *testing.T) {
	strToCnt := "1234567"
	exCnt := 7

	su := StringUtility{}
	l1 := su.StrGetCharCnt(strToCnt)

	if l1 != exCnt {
		t.Error(fmt.Sprintf("Expected string character count of '%v', got", exCnt), l1)
	}
}

func TestStrPadLeftToCenter(t *testing.T) {
	strToCntr := "1234567"
	fieldLen := 79
	exLen := 36
	su := StringUtility{}
	padStr, err := su.StrPadLeftToCenter(strToCntr, fieldLen)

	if err != nil {
		t.Error("Error on StrPadLeftToCenter(), got", err.Error())
	}

	l1 := su.StrGetRuneCnt(padStr)

	if l1 != exLen {
		t.Error(fmt.Sprintf("Expected pad length of '%v', got ", exLen), l1)
	}

}

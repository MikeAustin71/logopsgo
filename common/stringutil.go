package common

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// StringUtility - encapsulates a collection of
// methods used to manipulate strings
type StringUtility struct {
	StrIn  string
	StrOut string
}

// StrCenterInStr - returns a string which includes
// a left pad blank string plus the original string.
// The complete string will effectively center the
// original string is a field of specified length.
func (su StringUtility) StrCenterInStr(strToCenter string, fieldLen int) (string, error) {

	pad, err := su.StrPadLeftToCenter(strToCenter, fieldLen)

	if err != nil {
		return "", errors.New("StringUtility:StrCenterInStr() - " + err.Error())
	}

	return pad + strToCenter, nil

}

// StrPadLeftToCenter - Returns a blank string
// which allows centering of the target string
// in a fixed length field.
func (su StringUtility) StrPadLeftToCenter(strToCenter string, fieldLen int) (string, error) {

	sLen := su.StrGetRuneCnt(strToCenter)

	if sLen > fieldLen {
		return "", errors.New("StringUtility:StrPadLeftToCenter() - String To Center is longer than Field Length")
	}

	if sLen == fieldLen {
		return "", nil
	}

	margin := (fieldLen - sLen) / 2

	return strings.Repeat(" ", margin), nil
}

// StrGetRuneCnt - Uses utf8 Rune Count
// function to return the number of characters
// in a string.
func (su StringUtility) StrGetRuneCnt(targetStr string) int {
	return utf8.RuneCountInString(targetStr)
}

// StrGetCharCnt - Uses the 'len' method to
// return the number of characters in a
// string.
func (su StringUtility) StrGetCharCnt(targetStr string) int {
	return len([]rune(targetStr))
}
